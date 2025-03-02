package utils

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

// CompressFile compacta um arquivo para um arquivo zip
func CompressFile(src, dst string) error {
	// Abrir o arquivo de origem
	srcFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	// Obter informações do arquivo
	srcInfo, err := srcFile.Stat()
	if err != nil {
		return err
	}

	// Criar o arquivo zip de destino
	zipFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer zipFile.Close()

	// Criar um novo escritor zip
	zipWriter := zip.NewWriter(zipFile)
	defer zipWriter.Close()

	// Criar um cabeçalho para o arquivo dentro do zip
	header, err := zip.FileInfoHeader(srcInfo)
	if err != nil {
		return err
	}

	// Usar apenas o nome do arquivo, não o caminho completo
	header.Name = filepath.Base(src)

	// Usar compressão máxima
	header.Method = zip.Deflate

	// Criar um escritor para o arquivo dentro do zip
	writer, err := zipWriter.CreateHeader(header)
	if err != nil {
		return err
	}

	// Copiar o conteúdo do arquivo para o zip
	_, err = io.Copy(writer, srcFile)
	if err != nil {
		return err
	}

	return nil
}

// CompressDirectory compacta um diretório para um arquivo zip
func CompressDirectory(src, dst string) error {
	// Criar o arquivo zip de destino
	zipFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer zipFile.Close()

	// Criar um novo escritor zip
	zipWriter := zip.NewWriter(zipFile)
	defer zipWriter.Close()

	// Percorrer o diretório recursivamente
	err = filepath.Walk(src, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Criar um cabeçalho para o arquivo/diretório
		header, err := zip.FileInfoHeader(info)
		if err != nil {
			return err
		}

		// Ajustar o nome para ser relativo ao diretório de origem
		relPath, err := filepath.Rel(src, path)
		if err != nil {
			return err
		}
		if relPath == "." {
			return nil // Pular o diretório raiz
		}
		header.Name = relPath

		// Ajustar o separador de caminho para o padrão do zip (sempre /)
		header.Name = strings.ReplaceAll(header.Name, string(os.PathSeparator), "/")

		// Adicionar barra no final para diretórios
		if info.IsDir() {
			header.Name += "/"
			// Criar uma entrada para o diretório
			_, err = zipWriter.CreateHeader(header)
			return err
		}

		// Usar compressão máxima para arquivos
		header.Method = zip.Deflate

		// Criar um escritor para o arquivo dentro do zip
		writer, err := zipWriter.CreateHeader(header)
		if err != nil {
			return err
		}

		// Se for um arquivo, copiar o conteúdo
		file, err := os.Open(path)
		if err != nil {
			return err
		}
		defer file.Close()

		_, err = io.Copy(writer, file)
		return err
	})

	return err
}

// FormatHexDump formata um array de bytes em um dump hexadecimal
func FormatHexDump(data []byte) string {
	var result strings.Builder
	var asciiLine strings.Builder

	for i, b := range data {
		// Adicionar endereço no início de cada linha
		if i%16 == 0 {
			if i > 0 {
				result.WriteString("  " + asciiLine.String() + "\n")
				asciiLine.Reset()
			}
			result.WriteString(fmt.Sprintf("%08x: ", i))
		}

		// Adicionar byte em hexadecimal
		result.WriteString(fmt.Sprintf("%02x ", b))

		// Adicionar caractere ASCII se for imprimível, ou ponto se não for
		if b >= 32 && b <= 126 {
			asciiLine.WriteByte(b)
		} else {
			asciiLine.WriteString(".")
		}

		// Adicionar espaço extra no meio da linha
		if i%16 == 7 {
			result.WriteString(" ")
		}
	}

	// Preencher a última linha com espaços se necessário
	lastLineLen := len(data) % 16
	if lastLineLen > 0 {
		// Calcular quantos espaços precisamos adicionar
		spaces := (16 - lastLineLen) * 3
		if lastLineLen <= 8 {
			spaces++ // Adicionar o espaço extra do meio da linha
		}
		result.WriteString(strings.Repeat(" ", spaces))
	}

	// Adicionar a última linha ASCII
	if asciiLine.Len() > 0 {
		result.WriteString("  " + asciiLine.String())
	}

	return result.String()
}
