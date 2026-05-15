# Plano: verificar-recebimento-updates_2026-05-13_20-55.md

## Pedido do usuário
O mapeamento inline continua falhando (cache miss).

## Objetivo
Confirmar se o Telegram está enviando QUALQUER dado de `ChosenInlineResult` para o bot.

## Contexto atual
- Logs de depuração foram adicionados, mas o log de "recebido" não aparece.
- Suspeita de configuração no BotFather ou erro no MatchFunc do handler.

## Arquivos analisados
- `internal/telegram/events/loader.go`
- `internal/telegram/client.go`

## Arquivos que poderão ser modificados
- `internal/telegram/client.go`
- `internal/telegram/events/loader.go`

## Estratégia de implementação
1.  **Log Global de Updates**: Adicionar um log no `StartBot` que imprime o tipo de cada update recebido antes de processar pelos handlers.
2.  **Revisar MatchFunc**: Garantir que `matchChosenInlineResult` em `loader.go` está correto.

## Passos detalhados

1.  **Modificar `internal/telegram/client.go`**
    - Adicionar um log no middleware ou no início do processamento para mostrar o conteúdo bruto de `update.ChosenInlineResult`.

2.  **Modificar `internal/telegram/events/loader.go`**
    - Verificar a ordem dos handlers e se nada está "engolindo" o update antes.

## Riscos
- **Nulo**: Apenas logs.

## Como testar
1. Reiniciar o bot.
2. Compartilhar via modo inline.
3. Verificar se aparece no log algo como: "Update recebido: ChosenInlineResult={...}".

## Rollback
`git checkout internal/telegram/client.go`
