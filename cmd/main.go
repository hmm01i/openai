package main

import (
	"fmt"
	"log"
	"os"
	"path"
	"strings"

	"github.com/hmm01i/openai/pkg/version"
	"github.com/spf13/cobra"
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

var rootCmd = &cobra.Command{
	Use:   "oai",
	Short: "A brief description of your application",
	Long: `A longer description of your application, which can span multiple lines.
You can include more information about the app here.`,
}

var versionFlag bool

func init() {
	rootCmd.PersistentFlags().BoolVarP(&versionFlag, "version", "v", false, "display the current version")
}

func Execute() {
	cobra.OnInitialize(func() {
		if versionFlag {
			fmt.Printf("OAI v%s\n", version.Current)
			os.Exit(0)
		}
	})
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
