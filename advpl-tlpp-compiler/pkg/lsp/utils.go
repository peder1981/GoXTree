package lsp

import (
	"encoding/json"
)

// applyChanges aplica as mudanças ao texto do documento
func applyChanges(text string, changes []TextDocumentContentChangeEvent) string {
	if len(changes) == 0 {
		return text
	}

	// Por enquanto, apenas suporta substituição completa do documento
	return changes[0].Text
}

// parseSymbols converte os símbolos do IDE Integration para o formato LSP
func parseSymbols(symbolsJSON string) []DocumentSymbol {
	var symbols []DocumentSymbol
	var ideSymbols []struct {
		Name        string   `json:"name"`
		Kind        int      `json:"kind"`
		Line        int      `json:"line"`
		Column      int      `json:"column"`
		EndLine     int      `json:"endLine"`
		EndColumn   int      `json:"endColumn"`
		Description string   `json:"description"`
		Children    []string `json:"children,omitempty"`
	}

	if err := json.Unmarshal([]byte(symbolsJSON), &ideSymbols); err != nil {
		return symbols
	}

	for _, s := range ideSymbols {
		symbol := DocumentSymbol{
			Name:   s.Name,
			Detail: s.Description,
			Kind:   mapSymbolKind(s.Kind),
			Range: Range{
				Start: Position{Line: s.Line, Character: s.Column},
				End:   Position{Line: s.EndLine, Character: s.EndColumn},
			},
			SelectionRange: Range{
				Start: Position{Line: s.Line, Character: s.Column},
				End:   Position{Line: s.EndLine, Character: s.EndColumn},
			},
		}
		symbols = append(symbols, symbol)
	}

	return symbols
}

// parseCompletions converte os itens de completação do IDE Integration para o formato LSP
func parseCompletions(completionsJSON string) []CompletionItem {
	var items []CompletionItem
	var ideCompletions []struct {
		Label         string `json:"label"`
		Kind         int    `json:"kind"`
		Detail       string `json:"detail"`
		Documentation string `json:"documentation"`
		InsertText   string `json:"insertText"`
		SortText     string `json:"sortText"`
	}

	if err := json.Unmarshal([]byte(completionsJSON), &ideCompletions); err != nil {
		return items
	}

	for _, c := range ideCompletions {
		item := CompletionItem{
			Label:         c.Label,
			Kind:         mapCompletionKind(c.Kind),
			Detail:       c.Detail,
			Documentation: c.Documentation,
			InsertText:   c.InsertText,
			SortText:     c.SortText,
		}
		items = append(items, item)
	}

	return items
}

// mapSymbolKind mapeia os tipos de símbolos do IDE Integration para os tipos LSP
func mapSymbolKind(kind int) int {
	// https://microsoft.github.io/language-server-protocol/specifications/specification-current/#symbolKind
	switch kind {
	case 0: // Function
		return 12 // Function
	case 1: // Class
		return 5 // Class
	case 2: // Method
		return 6 // Method
	case 3: // Variable
		return 13 // Variable
	case 4: // Parameter
		return 20 // TypeParameter
	case 5: // Data
		return 7 // Property
	default:
		return 1 // File
	}
}

// mapCompletionKind mapeia os tipos de completação do IDE Integration para os tipos LSP
func mapCompletionKind(kind int) int {
	// https://microsoft.github.io/language-server-protocol/specifications/specification-current/#completionItemKind
	switch kind {
	case 0: // Function
		return 3 // Function
	case 1: // Class
		return 7 // Class
	case 2: // Method
		return 2 // Method
	case 3: // Variable
		return 6 // Variable
	case 4: // Parameter
		return 6 // Variable
	case 5: // Data
		return 10 // Property
	default:
		return 1 // Text
	}
}
