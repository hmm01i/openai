package main

import (
	"fmt"
	"log"
	"os"
	"path"
	"strings"

	openai "github.com/sashabaranov/go-openai"
)

type appFiles struct {
	configDir       string
	personasDir     string
	conversationDir string
	apiTokenFile    string
	imageSaveDir    string
}

var conf appFiles

func main() {
	Execute()
}

func (c *appFiles) initConfigs() {

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

type chatClient struct {
	model           string
	persona         string
	client          *openai.Client
	systemDirective string
	history         []openai.ChatCompletionMessage
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
