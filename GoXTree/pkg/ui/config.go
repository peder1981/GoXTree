package ui

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// Config representa as configurações do aplicativo
type Config struct {
	Theme         string            `json:"theme"`
	ColorScheme   map[string]string `json:"colorScheme"`
	ShowHidden    bool              `json:"showHidden"`
	CustomHotkeys map[string]string `json:"customHotkeys"`
}

// DefaultConfig retorna as configurações padrão
func DefaultConfig() Config {
	return Config{
		Theme:       "retro",
		ShowHidden:  false,
		ColorScheme: make(map[string]string),
		CustomHotkeys: map[string]string{
			"help":         "F1",
			"rename":       "F2",
			"search":       "F3",
			"advSearch":    "F4",
			"createDir":    "F7",
			"delete":       "F8",
			"sync":         "F9",
			"exit":         "F10",
			"selectAll":    "Ctrl+A",
			"deselectAll":  "Ctrl+D",
			"toggleHidden": "Ctrl+H",
		},
	}
}

// LoadConfig carrega as configurações do arquivo
func LoadConfig() (Config, error) {
	// Obter diretório de configuração
	configDir, err := getConfigDir()
	if err != nil {
		return DefaultConfig(), err
	}

	// Verificar se o arquivo de configuração existe
	configFile := filepath.Join(configDir, "config.json")
	if _, err := os.Stat(configFile); os.IsNotExist(err) {
		// Se não existir, criar com as configurações padrão
		config := DefaultConfig()
		if err := SaveConfig(config); err != nil {
			return config, err
		}
		return config, nil
	}

	// Ler arquivo de configuração
	data, err := os.ReadFile(configFile)
	if err != nil {
		return DefaultConfig(), err
	}

	// Decodificar JSON
	var config Config
	if err := json.Unmarshal(data, &config); err != nil {
		return DefaultConfig(), err
	}

	return config, nil
}

// SaveConfig salva as configurações no arquivo
func SaveConfig(config Config) error {
	// Obter diretório de configuração
	configDir, err := getConfigDir()
	if err != nil {
		return err
	}

	// Criar diretório se não existir
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return err
	}

	// Codificar configurações em JSON
	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return err
	}

	// Salvar no arquivo
	configFile := filepath.Join(configDir, "config.json")
	return os.WriteFile(configFile, data, 0644)
}

// getConfigDir retorna o diretório de configuração
func getConfigDir() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(homeDir, ".gxtree"), nil
}

// ApplyTheme aplica o tema especificado
func ApplyTheme(app *App, themeName string) error {
	switch themeName {
	case "retro":
		ApplyRetroThemeToApp(app)
	case "modern":
		ApplyModernThemeToApp(app)
	case "dark":
		ApplyDarkThemeToApp(app)
	case "light":
		ApplyLightThemeToApp(app)
	default:
		return fmt.Errorf("tema desconhecido: %s", themeName)
	}
	return nil
}
