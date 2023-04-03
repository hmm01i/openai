package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path"

	"github.com/chzyer/readline"
	openai "github.com/sashabaranov/go-openai"
)

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
