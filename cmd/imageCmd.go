package main

import "github.com/spf13/cobra"

var (
	output string
	prompt string
)

var imageCmd = &cobra.Command{
	Use:   "image",
	Short: "Generates an image based on a prompt",
	Long:  `This command generates an image based on the provided prompt and saves it to the specified output file.`,
	Run: func(cmd *cobra.Command, args []string) {
		imageRequest(prompt, output)
	},
}

func init() {
	conf.initConfigs()
	imageCmd.Flags().StringVarP(&output, "output", "o", "", "Output file for the generated image (required)")
	imageCmd.Flags().StringVarP(&prompt, "prompt", "p", "", "Prompt for generating the image (required)")
	imageCmd.MarkFlagRequired("output")
	imageCmd.MarkFlagRequired("prompt")
	rootCmd.AddCommand(imageCmd)
}
