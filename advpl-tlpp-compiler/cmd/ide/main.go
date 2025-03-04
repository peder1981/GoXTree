package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	
	"advpl-tlpp-compiler/pkg/compiler"
)

var (
	projectDir  string
	configFile  string
	theme       string
	showVersion bool
)

const version = "0.1.0"

func init() {
	flag.StringVar(&configFile, "config", "", "Arquivo de configuração (padrão: .advplide.json no diretório do projeto)")
	flag.StringVar(&theme, "theme", "default", "Tema do IDE (default, dark, light)")
	flag.BoolVar(&showVersion, "version", false, "Mostrar versão do IDE")

	// Uso personalizado
	flag.Usage = func() {
		fmt.Println("AdvPL/TLPP IDE - Um IDE baseado em terminal para AdvPL e TLPP")
		fmt.Println("")
		fmt.Println("Uso:")
		fmt.Println("  advpl-ide [opções] [diretório_do_projeto]")
		fmt.Println("")
		fmt.Println("Opções:")
		flag.PrintDefaults()
	}
}

func main() {
	flag.Parse()

	if showVersion {
		fmt.Printf("AdvPL/TLPP IDE versão %s\n", version)
		return
	}

	// Verificar se foi fornecido um diretório de projeto
	args := flag.Args()
	if len(args) > 0 {
		projectDir = args[0]
	} else {
		// Usar o diretório atual como diretório do projeto
		var err error
		projectDir, err = os.Getwd()
		if err != nil {
			fmt.Printf("Erro ao obter diretório atual: %v\n", err)
			os.Exit(1)
		}
	}

	// Verificar se o diretório existe
	if _, err := os.Stat(projectDir); os.IsNotExist(err) {
		fmt.Printf("Erro: Diretório não encontrado: %s\n", projectDir)
		os.Exit(1)
	}

	// Determinar o arquivo de configuração se não foi especificado
	if configFile == "" {
		configFile = filepath.Join(projectDir, ".advplide.json")
	}

	// Carregar configuração
	config, err := utils.LoadConfig(configFile)
	if err != nil && !os.IsNotExist(err) {
		fmt.Printf("Aviso: Erro ao carregar configuração: %v\n", err)
		fmt.Println("Usando configuração padrão.")
	}

	// Sobrescrever tema se especificado na linha de comando
	if theme != "default" {
		config.Theme = theme
	}

	// Inicializar o IDE
	app := ide.NewApp(projectDir, config)

	// Configurar manipuladores de eventos
	app.SetEventHandler(func(event tcell.Event) bool {
		switch event := event.(type) {
		case *tcell.EventKey:
			// Verificar atalhos globais
			if event.Key() == tcell.KeyCtrlQ {
				app.Quit()
				return true
			}
		}
		return false
	})

	// Iniciar a interface do usuário
	if err := app.Run(); err != nil {
		fmt.Printf("Erro ao executar o IDE: %v\n", err)
		os.Exit(1)
	}
}
