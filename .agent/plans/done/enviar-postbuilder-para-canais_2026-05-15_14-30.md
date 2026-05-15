# Plano: enviar-postbuilder-para-canais

## Pedido do usuário
Após salvar uma postagem no PostBuilder, adicionar um botão "Enviar para Canais" que permite selecionar um dos canais cadastrados do usuário e enviar a postagem diretamente para ele.

## Objetivo
Implementar um fluxo de envio direto para canais no PostBuilder após a postagem ser salva no Redis.

## Contexto atual
- O `CallbackHandler` em `internal/telegram/events/postBuilder/postBuilder.go` trata o evento `pb-save`, gerando um ID de sessão e oferecendo apenas o botão de compartilhar via modo inline.
- A função `sendFinalPost` já está preparada para enviar a postagem para um `chatID` genérico, suportando diversos tipos de mídia, incluindo stickers.
- `ChannelService.GetUserChannels` retorna os canais que o usuário configurou no bot.

## Arquivos analisados
- `internal/telegram/events/postBuilder/postBuilder.go`
- `internal/database/models/models.go`
- `internal/core/services/channels.go`

## Arquivos que poderão ser modificados
- `internal/telegram/events/postBuilder/postBuilder.go`

## Estratégia de implementação
1. **Modificar `pb-save`**: Adicionar o botão "📢 Enviar para Canais" com `CallbackData: pb-send-to-channels:<session_id>`.
2. **Adicionar Handler `pb-send-to-channels:`**:
    - Extrair o `session_id`.
    - Buscar canais do usuário via `c.ChannelService.GetUserChannels`.
    - Exibir uma lista de botões com o nome de cada canal.
    - Cada botão terá `CallbackData: pb-send-apply:<channel_id>:<session_id>`.
3. **Adicionar Handler `pb-send-apply:`**:
    - Extrair `channel_id` e `session_id`.
    - Recuperar a sessão salva via `c.CacheService.GetPostBuilderSession`.
    - Chamar `sendFinalPost` passando o `channel_id` como destino.
    - Notificar o usuário sobre o sucesso ou falha do envio.

## Passos detalhados

1. Editar `internal/telegram/events/postBuilder/postBuilder.go`:
    - No `CallbackHandler`, caso `pb-save`, atualizar o teclado `kb` para incluir o novo botão.
    - Adicionar tratamento para o prefixo `pb-send-to-channels:`:
        - Obter canais via `c.ChannelService.GetUserChannels(ctx, userID)`.
        - Se não houver canais, avisar o usuário.
        - Se houver, editar a mensagem atual (ou enviar nova) com a lista de canais.
    - Adicionar tratamento para o prefixo `pb-send-apply:`:
        - Fazer o parse de `channelID` e `sessionID` da `CallbackData`.
        - Recuperar `state` via `c.CacheService.GetPostBuilderSession`.
        - Executar `sendFinalPost(ctx, b, channelID, userID, c, state, false)`.
        - Enviar alerta/mensagem: "✅ Postagem enviada para o canal!".

## Riscos
- **Permissões**: O bot deve ser admin no canal para conseguir enviar a postagem. Se não for, o envio falhará (tratado pelo erro em `sendFinalPost`).
- **Sessão Expirada**: Se o usuário demorar muito para clicar, a sessão no Redis pode expirar (24h).

## Impactos esperados
- Facilidade para postar em canais próprios sem depender do modo inline.
- Melhor fluxo de trabalho para administradores de múltiplos canais.

## Como testar

### Execução
1. Criar um post no PostBuilder.
2. Clicar em "✅ Salvar".
3. Clicar no novo botão "📢 Enviar para Canais".
4. Selecionar um canal da lista (deve mostrar todos os canais que o usuário configurou no bot).
5. Verificar se o post chegou no canal com todos os elementos (mídia, legenda, botões, reações).
