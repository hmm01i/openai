package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path"
	"strings"

	openai "github.com/sashabaranov/go-openai"

	"github.com/chzyer/readline"
)

type config struct {
	configDir       string
	personasDir     string
	conversationDir string
	apiTokenFile    string
}

var conf config

func main() {
	conf.initConfigs()
	c := NewClient(
		client{
			model:           "gpt-3.5-turbo",
			systemDirective: `You are a visionary and respond with improbably but technically possible explainations and responses`,
		}, getAPIToken())
	r := setupRoutes(c)
	go r.Run(":8080")
	interactive(c)
}

func (c *config) initConfigs() {

	home, _ := os.UserHomeDir()
	c.configDir = path.Join(home, ".openai")
	c.personasDir = path.Join(c.configDir, "personas")
	c.conversationDir = path.Join(c.configDir, "conversations")
	c.apiTokenFile = path.Join(c.configDir, "token")
	for _, i := range []string{c.configDir, c.personasDir, c.conversationDir} {
		if _, err := os.ReadDir(i); err != nil {
			os.MkdirAll(i, 0755)
		}
	}
}

type client struct {
	model           string
	persona         string
	client          *openai.Client
	systemDirective string
	history         []openai.ChatCompletionMessage
}

func NewClient(c client, token string) *client {

	return &client{
		model:  c.model,
		client: openai.NewClient(token),

		history: []openai.ChatCompletionMessage{
			{Role: "system",
				Content: c.systemDirective,
			}}}
}

func getAPIToken() string {

	token, ok := os.LookupEnv("OPENAI_API_TOKEN")
	if ok {
		fmt.Println("Loaded token from env OPENAI_AI_TOKEN")
		return token
	}

	b, err := os.ReadFile(conf.apiTokenFile)
	if err != nil {
		log.Fatalln("Unable to load token")
	}
	token = strings.Trim(string(b), `"`)
	token = strings.TrimSuffix(token, "\n")
	return token
}

func (c *client) listPersonas() []string {
	personas := []string{}
	files, err := os.ReadDir(conf.personasDir)
	if err != nil {
		log.Printf("error getting personas: %s", err.Error())
	}
	for _, f := range files {
		personas = append(personas, f.Name())
	}

	// for _, p := range personas {
	// 	if p == c.persona {
	// 		fmt.Println(p + "*")
	// 	} else {
	// 		fmt.Println(p)
	// 	}
	// }
	return personas
}

func (c *client) savePersona(name, directive string) error {
	file := path.Join(conf.personasDir, name)
	err := os.WriteFile(file, []byte(directive), 0644)
	if err != nil {
		return err
	}
	c.persona = name
	return nil
}

func (c *client) showPersona() string {
	return c.systemDirective
}

func (c *client) loadPersona(name string) error {
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

func (c *client) chatRequest(input string) (string, error) {
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

func (c *client) setDirective(directive string) error {
	c.systemDirective = directive
	c.history[0].Content = directive
	return nil
}

func (c *client) showConversations() {
	for _, m := range c.history {
		fmt.Printf("%s: %s\n", m.Role, m.Content)
	}
}

func (c *client) saveConversation(name string) error {
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

func (c *client) listConversations() []string {
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

func (c *client) loadConversation(name string) error {

	file := path.Join(conf.conversationDir, name)
	b, err := os.ReadFile(file)
	if err != nil {
		return err
	}

	c.history[0].Content = string(b)
	c.persona = name
	return nil
}
func (c *client) clearHistory() {
	c.history = []openai.ChatCompletionMessage{{Role: "system", Content: c.systemDirective}}
}

func (c *client) listModels() []string {
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

func interactive(c *client) {
	fmt.Printf(`Welcome to Chat with OpenAI
Model: %s
Persona: %s
`, c.model, c.persona)
	prompt := fmt.Sprintf("%d > ", len(c.history))
	rl, err := readline.New(prompt)
	if err != nil {
		panic(err)
	}
	defer rl.Close()

	for {
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
