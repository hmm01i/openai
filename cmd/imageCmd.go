package main

import "github.com/spf13/cobra"

var (
	output string
	prompt string
	size   string
	number int
)

var imageCmd = &cobra.Command{
	Use:   "image",
	Short: "Generates an image based on a prompt",
	Long:  `This command generates an image based on the provided prompt and saves it to the specified output file.`,
	Run: func(cmd *cobra.Command, args []string) {
		imageRequest(prompt, size, number, output)
	},
}

func init() {
	conf.initConfigs()
	imageCmd.Flags().StringVarP(&output, "output", "o", "", "Output file for the generated image (required)")
	imageCmd.Flags().StringVarP(&prompt, "prompt", "p", "", "Prompt for generating the image (required)")
	imageCmd.Flags().StringVarP(&size, "size", "s", "256", "size of image to generate defaults to options: 256 (default), 512, 1024")
	imageCmd.Flags().IntVarP(&number, "num", "n", 1, "number of images to generate. Default: 1")
	imageCmd.MarkFlagRequired("output")
	imageCmd.MarkFlagRequired("prompt")
	rootCmd.AddCommand(imageCmd)
}
