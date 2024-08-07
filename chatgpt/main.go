package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

const apiURL = "https://api.openai.com/v1/chat/completions"

type GPTRequest struct {
	Messages []GPTMessage `json:"messages"`
	Model    string       `json:"model"`
}

type GPTMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type GPTResponse struct {
	Choices []struct {
		Message GPTMessage `json:"message"`
	} `json:"choices"`
}

func main() {
	apiKey := os.Getenv("OPENAI_API_KEY")
	if apiKey == "" {
		fmt.Println("Por favor, defina a variável de ambiente OPENAI_API_KEY com sua chave da API OpenAI.")
		return
	}

	var valor string
	fmt.Print("Insira um valor monetário: ")
	fmt.Scanln(&valor)

	requestBody, err := json.Marshal(GPTRequest{
		Messages: []GPTMessage{
			{Role: "user", Content: fmt.Sprintf("Escreva o valor %s por extenso.", valor)},
		},
		Model: "gpt-3.5-turbo",
	})
	if err != nil {
		fmt.Printf("Erro ao criar a requisição: %v\n", err)
		return
	}

	req, err := http.NewRequest("POST", apiURL, bytes.NewBuffer(requestBody))
	if err != nil {
		fmt.Printf("Erro ao criar a requisição HTTP: %v\n", err)
		return
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", apiKey))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("Erro ao enviar a requisição: %v\n", err)
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Erro ao ler a resposta: %v\n", err)
		return
	}

	var gptResponse GPTResponse
	err = json.Unmarshal(body, &gptResponse)
	if err != nil {
		fmt.Printf("Erro ao desserializar a resposta: %v\n", err)
		return
	}

	if len(gptResponse.Choices) > 0 {
		fmt.Printf("Valor por extenso: %s\n", gptResponse.Choices[0].Message.Content)
	} else {
		fmt.Println("Nenhuma resposta encontrada.")
	}
}
