package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"path"

	openai "github.com/sashabaranov/go-openai"

	"github.com/chzyer/readline"
)

var (
	configDir   = "openai"
	personasDir = path.Join(configDir, "personas")
)

func main() {
	home, _ := os.UserHomeDir()
	configDir = path.Join(home, ".openai")
	personasDir = path.Join(configDir, "personas")
	if _, err := os.ReadDir(configDir); err != nil {
		os.Mkdir(configDir, 0755)
	}
	if _, err := os.ReadDir(personasDir); err != nil {
		os.Mkdir(personasDir, 0755)
	}
	token := os.Getenv("OPENAI_API_TOKEN")
	c := NewClient(
		client{
			model:           "gpt-3.5-turbo",
			systemDirective: `You are a visionary and respond with improbably but technically possible explainations and responses`,
		}, token)
	interactive(c)
}

type client struct {
	model           string
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

func (c *client) listPersonas() []string {
	personas := []string{}
	files, err := os.ReadDir(personasDir)
	if err != nil {
		log.Printf("error getting personas: %s", err.Error())
	}
	for _, f := range files {
		personas = append(personas, f.Name())
	}
	return personas
}

func (c *client) savePersona(name, directive string) error {
	file := path.Join(personasDir, name)
	err := os.WriteFile(file, []byte(directive), 666)
	if err != nil {
		return err
	}
	return nil
}

func (c *client) showPersona() error {
	println(c.systemDirective)
	return nil
}

func (c *client) loadPersona(name string) error {
	file := path.Join(personasDir, name)
	b, err := os.ReadFile(file)
	if err != nil {
		return err
	}

	c.history[0].Content = string(b)
	return nil
}

func (c *client) handleChatRequest(input string) (string, error) {
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

func (c *client) showHistory() {
	for _, m := range c.history {
		fmt.Printf("%s: %s\n", m.Role, m.Content)
	}
}

func (c *client) clearHistory() {
	c.history = []openai.ChatCompletionMessage{{Role: "system", Content: c.systemDirective}}
}

func (c *client) listModels() {
	models, err := c.client.ListModels(context.Background())
	if err != nil {
		fmt.Println(err.Error())
	}
	for _, m := range models.Models {
		fmt.Println(m.ID)
	}
}
func interactive(c *client) {
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
		if c.command(line) {
			continue
		}
		r, err := c.handleChatRequest(line)

		if err != nil {
			log.Println(err.Error())

		}
		fmt.Println(r)
	}
}
