package main

import (
	"fmt"
	"os"

	"github.com/peder1981/GoXTree/pkg/ui"
)

func main() {
	// Criar e executar a aplicação
	app := ui.NewApp()

	if err := app.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Erro ao executar o GoXTree: %v\n", err)
		os.Exit(1)
	}
}
