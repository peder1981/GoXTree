package main

import (
    "bufio"
    "fmt"
    "os"
    "regexp"
    "strings"
)

func compressText(inputFile, outputFile string) error {
    // Abrir o arquivo de entrada
    file, err := os.Open(inputFile)
    if err != nil {
        return err
    }
    defer file.Close()

    // Ler todo o conteúdo
    scanner := bufio.NewScanner(file)
    var sb strings.Builder
    for scanner.Scan() {
        sb.WriteString(scanner.Text())
        sb.WriteString(" ")
    }

    if err := scanner.Err(); err != nil {
        return err
    }

    text := sb.String()

    // Substituir \r, \n, \t por espaço
    text = strings.ReplaceAll(text, "\r", " ")
    text = strings.ReplaceAll(text, "\n", " ")
    text = strings.ReplaceAll(text, "\t", " ")

    // Remover múltiplos espaços usando regex
    re := regexp.MustCompile(` +`)
    text = re.ReplaceAllString(text, " ")

    // Remover espaços no início e no fim
    text = strings.TrimSpace(text)

    // Escrever no arquivo de saída
    outFile, err := os.Create(outputFile)
    if err != nil {
        return err
    }
    defer outFile.Close()

    _, err = outFile.WriteString(text)
    if err != nil {
        return err
    }

    return nil
}

func main() {
    if len(os.Args) < 3 {
        fmt.Println("Uso: go run compress_text.go caminho_para_entrada caminho_para_saida")
        os.Exit(1)
    }

    inputFile := os.Args[1]
    outputFile := os.Args[2]

    err := compressText(inputFile, outputFile)
    if err != nil {
        fmt.Printf("Erro: %v\n", err)
        os.Exit(1)
    }

    fmt.Printf("Arquivo comprimido salvo como: %s\n", outputFile)
}

