package ui_test

import (
	"testing"
	
	"github.com/peder1981/GoXTree/pkg/ui"
)

func TestNewApp(t *testing.T) {
	app := ui.NewApp()
	if app == nil {
		t.Error("NewApp() returned nil")
	}
}

func TestAppInitialization(t *testing.T) {
	app := ui.NewApp()
	if app == nil {
		t.Skip("NewApp() returned nil, skipping initialization test")
	}
	
	// Verificar se os componentes básicos foram inicializados
	// Nota: Este teste pode precisar ser adaptado com base na implementação real
}
