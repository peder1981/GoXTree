package utils

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
)

// OpenFileWithDefaultApp abre um arquivo com o aplicativo padrão do sistema
func OpenFileWithDefaultApp(path string) error {
	// Verificar se o arquivo existe
	_, err := os.Stat(path)
	if err != nil {
		return err
	}

	// Abrir arquivo com o aplicativo padrão
	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "windows":
		cmd = exec.Command("cmd", "/c", "start", "", path)
	case "darwin":
		cmd = exec.Command("open", path)
	default: // Linux e outros
		cmd = exec.Command("xdg-open", path)
	}

	return cmd.Start()
}

// CalculateDirectorySize calcula o tamanho total de um diretório
func CalculateDirectorySize(path string) (int64, error) {
	var size int64
	err := filepath.Walk(path, func(filePath string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			size += info.Size()
		}
		return nil
	})
	return size, err
}

// CountFilesAndDirs conta o número de arquivos e diretórios em um diretório
func CountFilesAndDirs(path string) (int, int, error) {
	var numFiles, numDirs int
	err := filepath.Walk(path, func(filePath string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if filePath == path {
			return nil
		}
		if info.IsDir() {
			numDirs++
		} else {
			numFiles++
		}
		return nil
	})
	return numFiles, numDirs, err
}

// GetFilePermissionsString retorna uma string com as permissões do arquivo
func GetFilePermissionsString(fileInfo os.FileInfo) string {
	mode := fileInfo.Mode()
	perms := fmt.Sprintf("%s", mode.Perm())
	return perms
}

// compareBuffers compara dois buffers e retorna true se forem iguais
func compareBuffers(buf1, buf2 []byte) bool {
	if len(buf1) != len(buf2) {
		return false
	}
	for i := 0; i < len(buf1); i++ {
		if buf1[i] != buf2[i] {
			return false
		}
	}
	return true
}

// CopyDir copia um diretório recursivamente
func CopyDir(src, dst string) error {
	// Obter informações do diretório de origem
	srcInfo, err := os.Stat(src)
	if err != nil {
		return err
	}
	
	// Criar diretório de destino com as mesmas permissões
	err = os.MkdirAll(dst, srcInfo.Mode())
	if err != nil {
		return err
	}
	
	// Ler o conteúdo do diretório
	entries, err := os.ReadDir(src)
	if err != nil {
		return err
	}
	
	// Copiar cada item do diretório
	for _, entry := range entries {
		srcPath := filepath.Join(src, entry.Name())
		dstPath := filepath.Join(dst, entry.Name())
		
		if entry.IsDir() {
			// Copiar subdiretório recursivamente
			err = CopyDir(srcPath, dstPath)
			if err != nil {
				return err
			}
		} else {
			// Copiar arquivo
			err = CopyFile(srcPath, dstPath)
			if err != nil {
				return err
			}
		}
	}
	
	return nil
}
