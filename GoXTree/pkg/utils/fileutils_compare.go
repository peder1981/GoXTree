package utils

import (
	"os"
	"sort"
	"strings"
)

// CompareDirectories compara dois diretórios e retorna arquivos únicos em cada um
func CompareDirectories(dir1, dir2 string) ([]FileInfo, []FileInfo, []FileInfo, error) {
	// Listar arquivos no primeiro diretório
	files1, err := ListFiles(dir1, true)
	if err != nil {
		return nil, nil, nil, err
	}

	// Listar arquivos no segundo diretório
	files2, err := ListFiles(dir2, true)
	if err != nil {
		return nil, nil, nil, err
	}

	// Mapear arquivos por nome para facilitar a comparação
	map1 := make(map[string]FileInfo)
	map2 := make(map[string]FileInfo)

	for _, file := range files1 {
		map1[file.Name] = file
	}

	for _, file := range files2 {
		map2[file.Name] = file
	}

	// Encontrar arquivos únicos em cada diretório e comuns
	var uniqueInDir1 []FileInfo
	var uniqueInDir2 []FileInfo
	var common []FileInfo

	for name, file := range map1 {
		if _, exists := map2[name]; !exists {
			uniqueInDir1 = append(uniqueInDir1, file)
		} else {
			common = append(common, file)
		}
	}

	for name, file := range map2 {
		if _, exists := map1[name]; !exists {
			uniqueInDir2 = append(uniqueInDir2, file)
		}
	}

	// Ordenar os resultados por nome
	sort.Slice(uniqueInDir1, func(i, j int) bool {
		return uniqueInDir1[i].Name < uniqueInDir1[j].Name
	})

	sort.Slice(uniqueInDir2, func(i, j int) bool {
		return uniqueInDir2[i].Name < uniqueInDir2[j].Name
	})

	sort.Slice(common, func(i, j int) bool {
		return common[i].Name < common[j].Name
	})

	return uniqueInDir1, uniqueInDir2, common, nil
}

// CompareFiles compara dois arquivos e retorna as diferenças
func CompareFiles(file1, file2 string) ([]string, error) {
	// Ler o conteúdo dos arquivos
	content1, err := os.ReadFile(file1)
	if err != nil {
		return nil, err
	}

	content2, err := os.ReadFile(file2)
	if err != nil {
		return nil, err
	}

	// Dividir o conteúdo em linhas
	lines1 := strings.Split(string(content1), "\n")
	lines2 := strings.Split(string(content2), "\n")

	// Encontrar a subsequência comum mais longa
	lcs := longestCommonSubsequence(lines1, lines2)

	// Gerar as diferenças
	var diff []string

	i, j := 0, 0
	for k := 0; k < len(lcs); k++ {
		for i < len(lines1) && lines1[i] != lcs[k] {
			diff = append(diff, "- "+lines1[i])
			i++
		}
		for j < len(lines2) && lines2[j] != lcs[k] {
			diff = append(diff, "+ "+lines2[j])
			j++
		}
		diff = append(diff, "  "+lcs[k])
		i++
		j++
	}

	// Adicionar linhas restantes
	for i < len(lines1) {
		diff = append(diff, "- "+lines1[i])
		i++
	}
	for j < len(lines2) {
		diff = append(diff, "+ "+lines2[j])
		j++
	}

	return diff, nil
}

// longestCommonSubsequence encontra a subsequência comum mais longa entre duas slices
func longestCommonSubsequence(a, b []string) []string {
	// Criar a tabela de programação dinâmica
	m, n := len(a), len(b)
	dp := make([][]int, m+1)
	for i := range dp {
		dp[i] = make([]int, n+1)
	}

	// Preencher a tabela
	for i := 1; i <= m; i++ {
		for j := 1; j <= n; j++ {
			if a[i-1] == b[j-1] {
				dp[i][j] = dp[i-1][j-1] + 1
			} else {
				dp[i][j] = max(dp[i-1][j], dp[i][j-1])
			}
		}
	}

	// Reconstruir a subsequência
	var lcs []string
	i, j := m, n
	for i > 0 && j > 0 {
		if a[i-1] == b[j-1] {
			lcs = append([]string{a[i-1]}, lcs...)
			i--
			j--
		} else if dp[i-1][j] > dp[i][j-1] {
			i--
		} else {
			j--
		}
	}

	return lcs
}

// max retorna o maior de dois inteiros
func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
