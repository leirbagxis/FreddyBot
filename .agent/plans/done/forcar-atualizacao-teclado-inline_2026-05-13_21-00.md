# Plano: forcar-atualizacao-teclado-inline_2026-05-13_21-00.md

## Pedido do usuário
O log diz que o teclado foi reconstruído, mas a contagem de votos não aparece na mensagem.

## Objetivo
Garantir que, após reconstruir o teclado para mensagens inline, o bot efetivamente envie o comando `EditMessageReplyMarkup` para o Telegram.

## Contexto atual
- O `ChosenInlineResult` já está funcionando (mapeamento salvo no Redis).
- O `vote.Handler` reconstrói o teclado corretamente, mas a flag `updated` permanece `false` porque o teclado reconstruído já "nasce" com as contagens certas, falhando na comparação de diferença.

## Arquivos analisados
- `internal/telegram/callbacks/vote/vote.go`

## Arquivos que poderão ser modificados
- `internal/telegram/callbacks/vote/vote.go`

## Estratégia de implementação
1.  **Forçar Flag de Atualização**: Definir `updated = true` imediatamente após a reconstrução bem-sucedida do teclado para mensagens inline.
2.  **Logs de Confirmação**: Adicionar um log para confirmar o envio da edição do teclado.

## Passos detalhados

1.  **Modificar `internal/telegram/callbacks/vote/vote.go`**
    - Dentro do bloco `if ikb == nil && inlineMessageID != ""`, após a reconstrução do objeto `ikb`, adicionar `updated = true`.
    - Inicializar a variável `updated` antes do bloco de reconstrução.

## Riscos
- **Nulo**: Apenas garante que o comando de edição seja enviado quando necessário.

## Como testar
1. Reiniciar o bot.
2. Compartilhar postagem via modo inline.
3. Votar em um emoji.
4. O log deve mostrar "Teclado inline reconstruído" e a mensagem no Telegram deve atualizar para "Emoji 1".

## Rollback
`git checkout internal/telegram/callbacks/vote/vote.go`
