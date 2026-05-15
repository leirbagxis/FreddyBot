# Plano: depurar-mapeamento-inline_2026-05-13_20-45.md

## Pedido do usuário
O erro de mapeamento inline persiste: `Mapeamento inline não encontrado`.

## Objetivo
Identificar por que o mapeamento entre mensagens inline e sessões do PostBuilder não está sendo registrado ou recuperado.

## Contexto atual
- O bot utiliza `ChosenInlineResultHandler` para salvar o mapeamento no Redis.
- O `vote.Handler` tenta recuperar esse mapeamento para reconstruir o teclado.
- Se o mapeamento falha, o voto é computado no banco mas o teclado não atualiza visualmente.

## Arquivos analisados
- `internal/telegram/events/postBuilder/postBuilder.go`
- `internal/telegram/callbacks/vote/vote.go`
- `internal/cache/redis.go`

## Arquivos que poderão ser modificados
- `internal/telegram/events/postBuilder/postBuilder.go`
- `internal/telegram/callbacks/vote/vote.go`

## Estratégia de implementação
1.  **Aumentar Verbosidade dos Logs**: Adicionar logs explícitos no momento em que o mapeamento é *salvo* e no momento em que ele é *buscado*.
2.  **Verificar Persistência**: Garantir que o `sessionID` (que vem do `ResultID` do Telegram) está correto.
3.  **Orientação do BotFather**: Informar ao usuário sobre a configuração necessária no Telegram.

## Passos detalhados

1.  **Modificar `internal/telegram/events/postBuilder/postBuilder.go`**
    - Adicionar um log `logger.Bot("📥 ChosenInlineResult recebido: ID=%s, MessageID=%s", sessionID, inlineMessageID)` no início do handler.

2.  **Modificar `internal/telegram/callbacks/vote/vote.go`**
    - Melhorar o log de erro para incluir o resultado da busca no Redis.

## Riscos
- **Nulo**: Apenas logs adicionais.

## Impactos esperados
- Teremos certeza se o Telegram está enviando o evento necessário para o bot.

## Como testar
1. Reiniciar o bot.
2. Usar o PostBuilder para gerar um post.
3. Compartilhar via modo inline.
4. Observar os logs: Deve aparecer "Mapeado inline_message_id ...".
5. Se NÃO aparecer, o problema é a configuração no BotFather.

## Rollback
`git checkout internal/telegram/events/postBuilder/postBuilder.go`
