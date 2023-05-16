package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/hmm01i/openai/pkg/version"
)

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
