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
	if err := Execute(); err != nil {
		log.Fatalf("Error executing command: %v", err)
	}
}

func (c *appFiles) initConfigs() error {
	home, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("failed to get user home directory: %w", err)
	}

	c.configDir = path.Join(home, ".openai")
	c.personasDir = path.Join(c.configDir, "personas")
	c.conversationDir = path.Join(c.configDir, "conversations")
	c.apiTokenFile = path.Join(c.configDir, "token")

	// Create directories with more restrictive permissions
	for _, dir := range []string{c.configDir, c.personasDir, c.conversationDir} {
		if err := os.MkdirAll(dir, 0700); err != nil {
			return fmt.Errorf("failed to create directory %s: %w", dir, err)
		}
	}
	return nil
}

func getAPIToken() string {
	// Try environment variable first
	if token := os.Getenv("OPENAI_API_TOKEN"); token != "" {
		log.Println("Using API token from OPENAI_API_TOKEN environment variable")
		return token
	}

	// Try token file
	b, err := os.ReadFile(conf.apiTokenFile)
	if err != nil {
		log.Fatalln("Unable to load token")
	}

	token := strings.TrimSpace(string(b))
	if token == "" {
		log.Fatalln("Token file is empty")
	}

	return token
}

var rootCmd = &cobra.Command{
	Use:   "oai",
	Short: "OpenAI CLI client",
	Long: `A command-line interface for interacting with OpenAI services.
Supports managing conversations, personas, and various AI interactions.`,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		if err := conf.initConfigs(); err != nil {
			return err
		}
		return nil
	},
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("OAI v%s\n", version.Current)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}

func Execute() error {
	return rootCmd.Execute()
}
