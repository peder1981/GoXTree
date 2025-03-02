Este diretório contém recursos gráficos e de mídia para o projeto GoXTree.

- banner.txt: Arte ASCII para o banner do projeto
- screenshot.png: Captura de tela do GoXTree em execução
- logo.png: Logo do projeto para uso em documentação e site

Para capturar uma nova screenshot do GoXTree em execução, use:
1. Execute o GoXTree em um terminal
2. Use a ferramenta 'asciinema' para gravar uma sessão
3. Ou use a ferramenta 'terminalizer' para criar GIFs animados

Exemplo com asciinema:
```
asciinema rec -t "GoXTree Demo" goxTree_demo.cast
```

Exemplo com terminalizer:
```
terminalizer record goxTree_demo
terminalizer render goxTree_demo
```
