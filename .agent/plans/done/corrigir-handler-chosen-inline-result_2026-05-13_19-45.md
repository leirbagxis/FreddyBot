# Plano: corrigir-handler-chosen-inline-result_2026-05-13_19-45.md

## Pedido do usuário
O bot não estava compilando devido ao erro: `internal/telegram/events/loader.go:27:24: undefined: bot.HandlerTypeChosenInlineResult`.

## Objetivo
Corrigir o erro de compilação substituindo a constante inexistente por uma função de match customizada.

## Contexto atual
- A biblioteca `github.com/go-telegram/bot` v1.19.0 não possui a constante `HandlerTypeChosenInlineResult`.
- O bot utiliza `ChosenInlineResult` para mapear mensagens inline a sessões do Post Builder.

## Arquivos analisados
- `internal/telegram/events/loader.go`
- `go.mod`

## Arquivos modificados
- `internal/telegram/events/loader.go`

## Estratégia de implementação
1.  Remover o uso de `bot.HandlerTypeChosenInlineResult`.
2.  Implementar a função `matchChosenInlineResult` que verifica se `update.ChosenInlineResult != nil`.
3.  Registrar o handler usando `b.RegisterHandlerMatchFunc`.

## Passos detalhados
1.  **Modificar `internal/telegram/events/loader.go`**:
    - Criar função `matchChosenInlineResult`.
    - Atualizar `LoadEvents` para usar `RegisterHandlerMatchFunc` com a nova função de match.

## Riscos
- **Nulo:** A funcionalidade permanece idêntica, apenas altera a forma de registro do handler.

## Impactos esperados
- O projeto volta a compilar normalmente.

## Compatibilidade
- Go 1.24+
- go-telegram/bot v1.19.0

## Como testar
### Build
```bash
go build ./cmd/FreddyBot/main.go
```

## Rollback
`git checkout internal/telegram/events/loader.go`
