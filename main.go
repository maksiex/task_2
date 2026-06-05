package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

type ContentItem struct {
	Type string `json:"type"`
	Text string `json:"text"`
}

type InputMessage struct {
	Role    string        `json:"role"`
	Content []ContentItem `json:"content"`
}

type RequestBody struct {
	Model           string         `json:"model"`
	Input           []InputMessage `json:"input"`
	MaxOutputTokens int            `json:"max_output_tokens"`
}

type Response struct {
	Output []struct {
		Type    string `json:"type"`
		Content []struct {
			Type string `json:"type"`
			Text string `json:"text"`
		} `json:"content"`
	} `json:"output"`
}

func main() {
	_ = godotenv.Load()

	apiKey := os.Getenv("OPENAI_API_KEY")

	if apiKey == "" {
		fmt.Println("OPENAI_API_KEY is not set")
		return
	}

	reader := bufio.NewReader(os.Stdin)

	fmt.Print("Enter your question: ")

	userInput, err := reader.ReadString('\n')
	if err != nil {
		panic(err)
	}

	userInput = strings.TrimSpace(userInput)

	fmt.Print("Use response constraints? (y/n): ")

	mode, err := reader.ReadString('\n')
	if err != nil {
		panic(err)
	}

	mode = strings.TrimSpace(strings.ToLower(mode))

	systemPrompt := "You are a helpful AI assistant."

	if mode == "y" {
		systemPrompt = `
Answer in the following format:
- Use bullet points.
- Maximum 80 words total.
- End the response with the word END.
`
	}

	body := RequestBody{
		Model:           "gpt-5",
		MaxOutputTokens: 1200,
		Input: []InputMessage{
			{
				Role: "system",
				Content: []ContentItem{
					{
						Type: "input_text",
						Text: systemPrompt,
					},
				},
			},
			{
				Role: "user",
				Content: []ContentItem{
					{
						Type: "input_text",
						Text: userInput,
					},
				},
			},
		},
	}

	jsonData, err := json.Marshal(body)
	if err != nil {
		panic(err)
	}

	req, err := http.NewRequest(
		http.MethodPost,
		"https://api.openai.com/v1/responses",
		bytes.NewBuffer(jsonData),
	)
	if err != nil {
		panic(err)
	}

	req.Header.Set("Authorization", "Bearer "+apiKey)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	if resp.StatusCode >= 400 {
		fmt.Printf("API error: %s\n", resp.Status)
		fmt.Println(string(responseBody))
		return
	}

	var response Response

	if err := json.Unmarshal(responseBody, &response); err != nil {
		panic(err)
	}

	for _, output := range response.Output {

		if output.Type != "message" {
			continue
		}

		for _, content := range output.Content {

			if content.Type == "output_text" {

				fmt.Println("\nAI Response:")
				fmt.Println(content.Text)

				return
			}
		}
	}

	fmt.Println("No text response found")
}
