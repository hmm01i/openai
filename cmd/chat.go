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
	openai "github.com/sashabaranov/go-openai"
	"github.com/spf13/cobra"
)

type chatClient struct {
	model           string
	persona         string
	client          *openai.Client
	systemDirective string
	history         []openai.ChatCompletionMessage
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
		if resp, ok := c.command(line); ok {
			fmt.Println(resp)
			continue
		}
		r, err := c.chatRequest(line)

		if err != nil {
			log.Println(err.Error())

		}
		fmt.Println(r)
	}
}

func (c *chatClient) command(input string) (string, bool) {
	if !strings.HasPrefix(input, "/") {
		return "", false
	}
	inputSlice := strings.Split(input, " ")
	var cmd string
	var args []string
	cmd = inputSlice[0]
	if len(inputSlice) > 1 {
		args = inputSlice[1:]
	}

	if _, ok := cmds[cmd]; ok {
		return cmds[cmd](c, args), true
	}
	switch cmd {

	case "/listcmd":
		commands := []string{}
		return func(c *chatClient, args []string) string {
			for o, _ := range cmds {
				commands = append(commands, o)
			}
			return strings.Join(commands, "\n")
		}(c, args), true
	}
	return "Unrecognized command (try /listcmd)", true
}

var cmds = map[string]func(c *chatClient, args []string) string{
	//Personas
	"/listper": func(c *chatClient, args []string) string {
		personas := c.listPersonas()
		for i, p := range personas {
			if p == c.persona {
				personas[i] = p + "*"
			}
		}
		return strings.Join(personas, "\n")
	},
	"/saveper": func(c *chatClient, args []string) string { c.savePersona(args[0], c.systemDirective); return "ok" },
	"/showper": func(c *chatClient, args []string) string { persona := c.showPersona(); return persona },
	"/loadper": func(c *chatClient, args []string) string {
		if len(args) < 1 {
			return `ERR: No persona provided. Usage: /loadper <persona>`
		}
		if err := c.loadPersona(args[0]); err != nil {
			return "failed to load persona"
		}
		return "ok"
	},
	// system directive
	"/setdir": func(c *chatClient, args []string) string {
		if err := c.setDirective(strings.Trim(strings.Join(args, " "), `"`)); err != nil {
			return "failed to set directive"
		}
		return "ok"
	},
	// conversation history
	"/hist": func(c *chatClient, args []string) string {
		var hist []string
		for _, m := range c.history {
			hist = append(hist, fmt.Sprintf("%s: %s\n", m.Role, m.Content))
		}
		return strings.Join(hist, "\n")
	},
	"/clearhist": func(c *chatClient, args []string) string { c.clearHistory(); return "ok" },
	// models
	"/listmod": func(c *chatClient, args []string) string { mod := c.listModels(); return strings.Join(mod, "\n") },
	"/setmod":  func(c *chatClient, args []string) string { c.model = args[0]; return "model set" },
	// conversations
	"/saveconv": func(c *chatClient, args []string) string { c.saveConversation(args[0]); return "ok" },
	"/listconv": func(c *chatClient, args []string) string {
		convos := c.listConversations()
		return strings.Join(convos, "\n")
	},
	"/loadconv": func(c *chatClient, args []string) string { c.loadConversation(args[0]); return "ok" },

	// quit
	"/q": func(c *chatClient, args []string) string { os.Exit(0); return "ok" },
}
