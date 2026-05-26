# Plano: remover-logs-markdown-link

## Pedido do usuário
Remover os logs temporários adicionados para diagnosticar/corrigir o link Markdown embutido.

## Objetivo
Limpar os logs de diagnóstico mantendo intacta a correção do parser de Markdown e a proteção por placeholders.

## Contexto atual
A correção funcionou. Ainda existem logs temporários em:
- `internal/utils/utils.go`, para link Markdown convertido/ignorado e placeholder não restaurado.
- `internal/telegram/events/channelPost/dispatch_telego.go`, para preview do HTML final antes do envio.

## Arquivos analisados
- `internal/utils/utils.go`
- `internal/telegram/events/channelPost/dispatch_telego.go`

## Arquivos que poderão ser modificados
- `internal/utils/utils.go`
- `internal/telegram/events/channelPost/dispatch_telego.go`

## Estratégia de implementação
1. Remover logs temporários de conversão/restauração de links Markdown em `utils.go`.
2. Remover o helper `dispatchTextPreview` e o log de dispatch final em `dispatch_telego.go`.
3. Limpar imports que ficarem sem uso, especialmente `strings` em `dispatch_telego.go` se não houver outro uso.
4. Preservar `NormalizeMarkdownLinks`, `ProtectMarkdownLinks` e `RestoreProtectedMarkdownLinks`.

## Passos detalhados
1. Editar `internal/utils/utils.go` removendo chamadas a `logger.Bot` do fluxo de Markdown link.
2. Remover import `pkg/logger` de `utils.go` se não for mais usado.
3. Editar `internal/telegram/events/channelPost/dispatch_telego.go` removendo `dispatchTextPreview` e o log antes de `EditMessageText`.
4. Remover import `strings` de `dispatch_telego.go` se ficar sem uso.
5. Rodar `gofmt`.
6. Rodar `git diff --check`.
7. Rodar `npm run build`.
8. Tentar `go test ./...` e `go build ./cmd/FreddyBot/main.go`, registrando bloqueios locais se persistirem.

## Riscos
- Remover logs reduz a visibilidade do parser, mas a correção já foi validada pelo usuário.

## Impactos esperados
- Logs do bot ficam limpos novamente.
- Link Markdown embutido continua funcionando.

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
Enviar o mesmo post com `[texto](https://t.me/Flor_maracuja_ofc)` e confirmar que funciona sem logs de diagnóstico.

## Rollback
Reverter os dois arquivos modificados neste plano se for necessário recuperar logs de diagnóstico.

## Observações
Não reverter as correções anteriores nem as demais mudanças pendentes no worktree.
