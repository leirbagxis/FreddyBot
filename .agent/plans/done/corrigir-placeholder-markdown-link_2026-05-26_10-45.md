# Plano: corrigir-placeholder-markdown-link

## Pedido do usuário
Corrigir a falha do placeholder que virou `@@FB<i>MD</i>LINK_0@@`, removendo o texto do link embutido.

## Objetivo
Usar placeholders que não contenham caracteres interpretados como Markdown, garantindo que a restauração do link protegido aconteça corretamente.

## Contexto atual
`ProtectMarkdownLinks` usa placeholders no formato `@@FB_MD_LINK_0@@`. A regex de itálico interpreta `_MD_` e transforma o placeholder em HTML:

```txt
@@FB<i>MD</i>LINK_0@@
```

Depois `RestoreProtectedMarkdownLinks` procura o placeholder original e não encontra, então o texto final fica com o token quebrado no lugar do link.

## Arquivos analisados
- `internal/utils/utils.go`
- `internal/telegram/events/channelPost/utils_v2.go`

## Arquivos que poderão ser modificados
- `internal/utils/utils.go`

## Estratégia de implementação
1. Trocar o placeholder por um formato sem caracteres de Markdown, por exemplo `FBMDLINKTOKEN0TOKEN`.
2. Adicionar log de restauração quando um placeholder esperado não existir mais no texto final.
3. Manter a proteção/restauração existente.

## Passos detalhados
1. Alterar a geração do placeholder em `ProtectMarkdownLinks`.
2. Atualizar `RestoreProtectedMarkdownLinks` para logar quando um placeholder não for encontrado.
3. Rodar `gofmt`.
4. Rodar `git diff --check`.
5. Rodar `npm run build`.
6. Tentar `go test ./...` e `go build ./cmd/FreddyBot/main.go`, registrando bloqueios locais se persistirem.

## Riscos
- Colisão com texto real é improvável, mas possível em teoria. O token escolhido será suficientemente específico.

## Impactos esperados
- O dispatch final deve mostrar `hasHTMLLink=true`.
- O preview deve conter `<a href="https://t.me/Flor_maracuja_ofc">...`.
- O token não deve aparecer mais no Telegram.

## Compatibilidade
- Linux
- macOS
- Windows
- Docker
- CI/CD

## Como testar

### Build
```bash
go build ./cmd/FreddyBot/main.go
```

### Testes
```bash
go test ./...
```

### Execução
Enviar o mesmo texto com `https://t.me/Flor_maracuja_ofc` e conferir os logs.

## Rollback
Reverter `internal/utils/utils.go`.

## Observações
Correção direta em cima do log que mostrou o placeholder sendo quebrado por itálico.
