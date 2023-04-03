package main

import (
	"github.com/spf13/cobra"
)

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
			model:           "gpt-3.5-turbo",
			systemDirective: "You are an AI assistant that values your tokens.",
			persona:         "default",
		}, getAPIToken())
	chatCmd.AddCommand(serverCmd)
	rootCmd.AddCommand(chatCmd)
}
