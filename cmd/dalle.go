package main

import (
	"bytes"
	"context"
	"encoding/base64"
	"fmt"
	"image/png"
	"os"

	openai "github.com/sashabaranov/go-openai"
)

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
