# Plano: corrigir-erro-sessao-expirada-salvar-rascunho_2026-08-16_12-05

## Pedido do usuário
Corrigir a mensagem de erro "Sessão expirada" que acontece ao clicar no botão `pb-save-template:current` (`💾 Salvar Rascunho`) durante a criação de um post.

## Causa Raiz
Quando o botão `💾 Salvar Rascunho` envia a callback data `pb-save-template:current`, a variável `sessionID` assume o valor `"current"`. O bot tentava buscar no Redis a sessão com a chave literal `"current"`, o que falhava por não ser um ID de sessão numérico, recaindo em um objeto de estado nulo ou incompleto e acusando erro.

## Objetivo
Garantir que quando a `sessionID` for `"current"` ou quando não houver ID explícito na URL, o bot busque a sessão ativa diretamente através da chave do usuário (`c.CacheService.GetPostBuilderState(ctx, userID)`), salvando o rascunho instantaneamente.

## Arquivos analisados
- `internal/telegram/handlers/events/postBuilder/postBuilder.go`

## Arquivos que poderão ser modificados
- `internal/telegram/handlers/events/postBuilder/postBuilder.go`

## Estratégia de implementação

1. **Ajuste em `postBuilder.go` no handler `pb-save-template:`**:
   - Se `sessionID == "current"` ou `sessionID == ""`, consultar diretamente o estado ativo no Redis via `c.CacheService.GetPostBuilderState(context.Background(), userID)`.
   - Se a sessão for válida, extrair o título/corpo/mídia e criar o rascunho com o nome automático (`Rascunho - DD/MM às HH:MM`) ou Título do post.
   - Retornar o `AnswerCallbackQuery` com o alerta `"💾 Rascunho '[Nome]' salvo com sucesso na sua biblioteca!"`.

2. **Validação & Testes**:
   - Executar a suíte de testes do Go (`go test ./cmd/... ./internal/... ./pkg/...`).
   - Compilar o binário Go (`go build ./cmd/FreddyBot`).

## Passos detalhados
1. Salvar o plano em `.agent/plans/pending/corrigir-erro-sessao-expirada-salvar-rascunho_2026-08-16_12-05.md`.
2. Apresentar o resumo ao usuário e solicitar aprovação explícita.
3. Atualizar `postBuilder.go`.
4. Rodar testes e compilação (`go test ./...`, `go build ./cmd/FreddyBot`).

## Riscos
- Mínimo. Correção cirúrgica na resolução do estado da sessão do PostBuilder no Telegram.

## Compatibilidade
- Linux, macOS, Windows, Docker, CI/CD.

## Como testar

### Testes Go
```bash
go test ./cmd/... ./internal/... ./pkg/...
go build ./cmd/FreddyBot
```

## Rollback
`git checkout internal/telegram/handlers/events/postBuilder/postBuilder.go`
