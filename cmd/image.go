package main

import (
	"bytes"
	"context"
	"encoding/base64"
	"fmt"
	"image/png"
	"os"

	openai "github.com/sashabaranov/go-openai"
	"github.com/spf13/cobra"
)

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

func imageRequest(prompt string, outputFile string) error {
	ic := openai.NewClient(getAPIToken())
	ctx := context.Background()
	// Example image as base64
	reqBase64 := openai.ImageRequest{
		Prompt:         prompt,
		Size:           openai.CreateImageSize256x256,
		ResponseFormat: openai.CreateImageResponseFormatB64JSON,
		N:              1,
	}

	respBase64, err := ic.CreateImage(ctx, reqBase64)
	if err != nil {
		fmt.Printf("Image creation error: %v\n", err)
		return err
	}

	imgBytes, err := base64.StdEncoding.DecodeString(respBase64.Data[0].B64JSON)
	if err != nil {
		fmt.Printf("Base64 decode error: %v\n", err)
		return err
	}

	r := bytes.NewReader(imgBytes)
	imgData, err := png.Decode(r)
	if err != nil {
		fmt.Printf("PNG decode error: %v\n", err)
		return err
	}

	file, err := os.Create(outputFile)
	if err != nil {
		fmt.Printf("File creation error: %v\n", err)
		return err
	}
	defer file.Close()

	if err := png.Encode(file, imgData); err != nil {
		fmt.Printf("PNG encode error: %v\n", err)
		return err
	}

	fmt.Printf("The image was saved as %s", outputFile)

	return nil
}
