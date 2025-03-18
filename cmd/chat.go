package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path"
	"strings"

	"github.com/chzyer/readline"
	"github.com/hmm01i/openai/pkg/commands"
	openai "github.com/sashabaranov/go-openai"
	"github.com/spf13/cobra"
)

type chatClient struct {
	model           string
	persona         string
	client          *openai.Client
	systemDirective string
	history         []openai.ChatCompletionMessage
	cmdRegistry     *commands.CommandRegistry
}

var (
	chatC   *chatClient
	persona string
)

var chatCmd = &cobra.Command{
	Use:   "chat",
	Short: "Start a interactive chat session",
	Run: func(cmd *cobra.Command, args []string) {
		interactive(chatC)
	},
}
var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Starts the HTTP server",
	Long:  `This command starts the HTTP server, which listens on a specified port.`,
	Run: func(cmd *cobra.Command, args []string) {
		r := setupRoutes(chatC)
		r.Run(":8080")
	},
}

func init() {
	conf.initConfigs()
	chatC = NewChatClient(
		chatClient{
			model:           "gpt-4",
			systemDirective: "You are an AI assistant that values your tokens.",
			persona:         "default",
		}, getAPIToken())
	chatCmd.AddCommand(serverCmd)
	rootCmd.AddCommand(chatCmd)
}

func NewChatClient(c chatClient, token string) *chatClient {
	c.client = openai.NewClient(token)
	c.history = []openai.ChatCompletionMessage{
		{Role: "system",
			Content: c.systemDirective,
		}}
	if c.persona != "" {
		c.loadPersona(c.persona)
	}

	// Initialize command registry with beta access for testing
	c.cmdRegistry = commands.NewCommandRegistry(commands.AccessBeta)
	return &c
}

func (c *chatClient) listPersonas() []string {
	personas := []string{}
	files, err := os.ReadDir(conf.personasDir)
	if err != nil {
		log.Printf("error getting personas: %s", err.Error())
	}
	for _, f := range files {
		personas = append(personas, f.Name())
	}

	return personas
}

func (c *chatClient) savePersona(name, directive string) error {
	file := path.Join(conf.personasDir, name)
	err := os.WriteFile(file, []byte(directive), 0644)
	if err != nil {
		return err
	}
	c.persona = name
	return nil
}

func (c *chatClient) showPersona() string {
	return c.systemDirective
}

func (c *chatClient) loadPersona(name string) error {
	file := path.Join(conf.personasDir, name)
	b, err := os.ReadFile(file)
	if err != nil {
		return err
	}

	c.history[0].Content = string(b)
	c.systemDirective = string(b)
	c.persona = name
	return nil
}

func (c *chatClient) chatRequest(input string) (string, error) {
	c.history = append(c.history, openai.ChatCompletionMessage{
		Role:    "user",
		Content: input,
	})
	request := openai.ChatCompletionRequest{
		Model:    c.model,
		Messages: c.history,
		Stream:   false,
	}
	response, err := c.client.CreateChatCompletion(context.Background(), request)
	if err != nil {
		return "", err
	}
	c.history = append(c.history, response.Choices[0].Message)
	return response.Choices[0].Message.Content, nil
}

func (c *chatClient) setDirective(directive string) error {
	c.systemDirective = directive
	c.history[0].Content = directive
	return nil
}

func (c *chatClient) showConversations() {
	for _, m := range c.history {
		fmt.Printf("%s: %s\n", m.Role, m.Content)
	}
}

func (c *chatClient) saveConversation(name string) error {
	file := path.Join(conf.conversationDir, name)
	conv, err := json.Marshal(c.history)
	if err != nil {
		return err
	}
	if err := os.WriteFile(file, conv, 0644); err != nil {
		return err
	}
	c.persona = name
	return nil
}

func (c *chatClient) listConversations() []string {
	conversations := []string{}
	files, err := os.ReadDir(conf.conversationDir)
	if err != nil {
		log.Printf("error getting personas: %s", err.Error())
	}
	for _, f := range files {
		conversations = append(conversations, f.Name())
	}
	return conversations
}

func (c *chatClient) loadConversation(name string) error {
	file := path.Join(conf.conversationDir, name)
	b, err := os.ReadFile(file)
	if err != nil {
		return err
	}

	c.history[0].Content = string(b)
	c.persona = name
	return nil
}

func (c *chatClient) clearHistory() {
	c.history = []openai.ChatCompletionMessage{{Role: "system", Content: c.systemDirective}}
}

func (c *chatClient) listModels() []string {
	mod := []string{}
	models, err := c.client.ListModels(context.Background())
	if err != nil {
		fmt.Println(err.Error())
		return mod
	}
	for _, m := range models.Models {
		mod = append(mod, m.ID)
	}
	return mod
}

func interactive(c *chatClient) {
	fmt.Printf(`Welcome to Chat with OpenAI
Model: %s
Persona: %s
`, c.model, c.persona)
	rl, err := readline.New("> ")
	if err != nil {
		panic(err)
	}
	defer rl.Close()

	for {
		fmt.Printf("%d > ", len(c.history))
		line, err := rl.Readline()
		if err != nil { // io.EOF
			break
		}

		if strings.HasPrefix(line, "/") {
			resp := c.cmdRegistry.ExecuteCommand(c, line)
			if resp == "" {
				continue
			}

			var cmdResp commands.CommandResponse
			if err := json.Unmarshal([]byte(resp), &cmdResp); err != nil {
				fmt.Println("Error:", err)
				continue
			}

			if !cmdResp.Success {
				fmt.Printf("\033[31mError: %s\033[0m\n", cmdResp.Error)
				continue
			}

			fmt.Println(cmdResp.Message)
			if line == "/q" {
				os.Exit(0)
			}
			continue
		}

		r, err := c.chatRequest(line)
		if err != nil {
			log.Println(err.Error())
		}
		fmt.Println(r)
	}
}

// Remove old command-related code
func (c *chatClient) ListPersonas() []string {
	return c.listPersonas()
}

func (c *chatClient) SavePersona(name, directive string) error {
	return c.savePersona(name, directive)
}

func (c *chatClient) ShowPersona() string {
	return c.showPersona()
}

func (c *chatClient) LoadPersona(name string) error {
	return c.loadPersona(name)
}

func (c *chatClient) SetDirective(directive string) error {
	return c.setDirective(directive)
}

func (c *chatClient) ClearHistory() {
	c.clearHistory()
}

func (c *chatClient) ListModels() []string {
	return c.listModels()
}

func (c *chatClient) SetModel(model string) {
	c.model = model
}

func (c *chatClient) SaveConversation(name string) error {
	return c.saveConversation(name)
}

func (c *chatClient) ListConversations() []string {
	return c.listConversations()
}

func (c *chatClient) LoadConversation(name string) error {
	return c.loadConversation(name)
}

func (c *chatClient) GetCurrentPersona() string {
	return c.persona
}

func (c *chatClient) GetHistory() []commands.Message {
	messages := make([]commands.Message, len(c.history))
	for i, msg := range c.history {
		messages[i] = commands.Message{
			Role:    msg.Role,
			Content: msg.Content,
		}
	}
	return messages
}
