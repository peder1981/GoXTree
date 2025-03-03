package utils

import (
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// SyncOptions define as opções para sincronização de diretórios
type SyncOptions struct {
	SourceDir      string
	DestDir        string
	DeleteOrphaned bool
	PreviewOnly    bool
	SkipNewer      bool
	SkipExisting   bool
	IncludeHidden  bool
}

// SyncAction representa uma ação de sincronização
type SyncAction struct {
	Action     string
	SourcePath string
	DestPath   string
	IsDir      bool
	Size       int64
	ModTime    time.Time
	Reason     string
}

// SyncDirectories sincroniza dois diretórios
func SyncDirectories(options SyncOptions) ([]SyncAction, error) {
	var actions []SyncAction

	// Verificar se os diretórios existem
	srcInfo, err := os.Stat(options.SourceDir)
	if err != nil {
		return nil, fmt.Errorf("erro ao acessar diretório de origem: %v", err)
	}
	if !srcInfo.IsDir() {
		return nil, fmt.Errorf("origem não é um diretório: %s", options.SourceDir)
	}

	// Verificar se o diretório de destino existe, se não, criar
	destInfo, err := os.Stat(options.DestDir)
	if err != nil {
		if os.IsNotExist(err) {
			if !options.PreviewOnly {
				if err := os.MkdirAll(options.DestDir, srcInfo.Mode()); err != nil {
					return nil, fmt.Errorf("erro ao criar diretório de destino: %v", err)
				}
			}
			actions = append(actions, SyncAction{
				Action:     "create_dir",
				DestPath:   options.DestDir,
				IsDir:      true,
				ModTime:    srcInfo.ModTime(),
				Reason:     "Diretório de destino não existe",
			})
		} else {
			return nil, fmt.Errorf("erro ao acessar diretório de destino: %v", err)
		}
	} else if !destInfo.IsDir() {
		return nil, fmt.Errorf("destino não é um diretório: %s", options.DestDir)
	}

	// Mapear arquivos no diretório de origem
	srcFiles, err := mapDirectoryFiles(options.SourceDir, options.IncludeHidden)
	if err != nil {
		return nil, fmt.Errorf("erro ao mapear diretório de origem: %v", err)
	}

	// Mapear arquivos no diretório de destino
	destFiles, err := mapDirectoryFiles(options.DestDir, options.IncludeHidden)
	if err != nil {
		return nil, fmt.Errorf("erro ao mapear diretório de destino: %v", err)
	}

	// Comparar e sincronizar
	for relPath, srcFile := range srcFiles {
		srcFullPath := filepath.Join(options.SourceDir, relPath)
		destFullPath := filepath.Join(options.DestDir, relPath)

		// Verificar se o arquivo existe no destino
		destFile, exists := destFiles[relPath]

		if !exists {
			// Arquivo não existe no destino, criar
			if srcFile.IsDir {
				actions = append(actions, SyncAction{
					Action:     "create_dir",
					SourcePath: srcFullPath,
					DestPath:   destFullPath,
					IsDir:      true,
					ModTime:    srcFile.ModTime,
					Reason:     "Diretório não existe no destino",
				})

				if !options.PreviewOnly {
					if err := os.MkdirAll(destFullPath, os.ModePerm); err != nil {
						return nil, fmt.Errorf("erro ao criar diretório: %v", err)
					}
				}
			} else {
				actions = append(actions, SyncAction{
					Action:     "copy",
					SourcePath: srcFullPath,
					DestPath:   destFullPath,
					IsDir:      false,
					Size:       srcFile.Size,
					ModTime:    srcFile.ModTime,
					Reason:     "Arquivo não existe no destino",
				})

				if !options.PreviewOnly {
					// Criar diretório pai se necessário
					destDir := filepath.Dir(destFullPath)
					if err := os.MkdirAll(destDir, os.ModePerm); err != nil {
						return nil, fmt.Errorf("erro ao criar diretório pai: %v", err)
					}

					if err := CopyFile(srcFullPath, destFullPath); err != nil {
						return nil, fmt.Errorf("erro ao copiar arquivo: %v", err)
					}
				}
			}
		} else {
			// Arquivo existe no destino, verificar se precisa atualizar
			if !srcFile.IsDir && !destFile.IsDir {
				needsUpdate := false
				reason := ""

				// Verificar se o arquivo de origem é mais novo
				if srcFile.ModTime.After(destFile.ModTime) {
					if options.SkipNewer {
						// Pular arquivos mais novos no destino
						continue
					}
					needsUpdate = true
					reason = "Arquivo de origem é mais recente"
				} else if srcFile.ModTime.Equal(destFile.ModTime) && srcFile.Size != destFile.Size {
					needsUpdate = true
					reason = "Tamanhos diferentes com mesma data"
				}

				// Pular arquivos existentes se a opção estiver ativada
				if options.SkipExisting {
					continue
				}

				if needsUpdate {
					actions = append(actions, SyncAction{
						Action:     "update",
						SourcePath: srcFullPath,
						DestPath:   destFullPath,
						IsDir:      false,
						Size:       srcFile.Size,
						ModTime:    srcFile.ModTime,
						Reason:     reason,
					})

					if !options.PreviewOnly {
						if err := CopyFile(srcFullPath, destFullPath); err != nil {
							return nil, fmt.Errorf("erro ao atualizar arquivo: %v", err)
						}
					}
				}
			} else if srcFile.IsDir && !destFile.IsDir {
				// Diretório na origem, arquivo no destino
				actions = append(actions, SyncAction{
					Action:     "replace",
					SourcePath: srcFullPath,
					DestPath:   destFullPath,
					IsDir:      true,
					ModTime:    srcFile.ModTime,
					Reason:     "Diretório na origem, arquivo no destino",
				})

				if !options.PreviewOnly {
					if err := os.Remove(destFullPath); err != nil {
						return nil, fmt.Errorf("erro ao remover arquivo: %v", err)
					}
					if err := os.MkdirAll(destFullPath, os.ModePerm); err != nil {
						return nil, fmt.Errorf("erro ao criar diretório: %v", err)
					}
				}
			} else if !srcFile.IsDir && destFile.IsDir {
				// Arquivo na origem, diretório no destino
				actions = append(actions, SyncAction{
					Action:     "replace",
					SourcePath: srcFullPath,
					DestPath:   destFullPath,
					IsDir:      false,
					Size:       srcFile.Size,
					ModTime:    srcFile.ModTime,
					Reason:     "Arquivo na origem, diretório no destino",
				})

				if !options.PreviewOnly {
					if err := os.RemoveAll(destFullPath); err != nil {
						return nil, fmt.Errorf("erro ao remover diretório: %v", err)
					}
					if err := CopyFile(srcFullPath, destFullPath); err != nil {
						return nil, fmt.Errorf("erro ao copiar arquivo: %v", err)
					}
				}
			}
		}

		// Marcar como processado
		delete(destFiles, relPath)
	}

	// Processar arquivos órfãos no destino
	if options.DeleteOrphaned {
		// Ordenar os caminhos para garantir que os arquivos sejam excluídos antes dos diretórios
		var orphanedPaths []string
		for relPath := range destFiles {
			orphanedPaths = append(orphanedPaths, relPath)
		}

		// Processar arquivos órfãos
		for _, relPath := range orphanedPaths {
			destFile := destFiles[relPath]
			destFullPath := filepath.Join(options.DestDir, relPath)

			actions = append(actions, SyncAction{
				Action:     "delete",
				DestPath:   destFullPath,
				IsDir:      destFile.IsDir,
				Size:       destFile.Size,
				ModTime:    destFile.ModTime,
				Reason:     "Arquivo órfão no destino",
			})

			if !options.PreviewOnly {
				if destFile.IsDir {
					if err := os.RemoveAll(destFullPath); err != nil {
						return nil, fmt.Errorf("erro ao remover diretório órfão: %v", err)
					}
				} else {
					if err := os.Remove(destFullPath); err != nil {
						return nil, fmt.Errorf("erro ao remover arquivo órfão: %v", err)
					}
				}
			}
		}
	}

	return actions, nil
}

// Estrutura auxiliar para mapear arquivos
type fileMapEntry struct {
	IsDir   bool
	Size    int64
	ModTime time.Time
}

// Mapeia todos os arquivos em um diretório recursivamente
func mapDirectoryFiles(rootDir string, includeHidden bool) (map[string]fileMapEntry, error) {
	files := make(map[string]fileMapEntry)

	err := filepath.Walk(rootDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Pular o diretório raiz
		if path == rootDir {
			return nil
		}

		// Verificar se é um arquivo oculto
		if !includeHidden && isHidden(path) {
			if info.IsDir() {
				return filepath.SkipDir
			}
			return nil
		}

		// Obter caminho relativo
		relPath, err := filepath.Rel(rootDir, path)
		if err != nil {
			return err
		}

		files[relPath] = fileMapEntry{
			IsDir:   info.IsDir(),
			Size:    info.Size(),
			ModTime: info.ModTime(),
		}

		return nil
	})

	return files, err
}
