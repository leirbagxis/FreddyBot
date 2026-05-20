# Plano: add-custom-emoji-buttons-postbuilder

## Pedido do usuário
O usuário solicitou que o suporte a emojis customizados do Telegram seja estendido também para os botões do PostBuilder (botões de URL e botões de reação/voto).

## Objetivo
Permitir que usuários utilizem emojis customizados do Telegram nos botões criados via PostBuilder, utilizando o campo `IconCustomEmojiID` introduzido recentemente na API do Telegram.

## Contexto atual
- O PostBuilder suporta botões de URL e botões de reação.
- Botões de URL são armazenados na struct `PostBuilderButton` com `Text` e `URL`.
- Reações são armazenadas como uma string separada por vírgulas (`state.Reactions`).
- Atualmente, apenas emojis padrão (Unicode) funcionam corretamente nos botões.

## Arquivos analisados
- `internal/cache/types.go`
- `internal/telegram/handlers/events/postBuilder/postBuilder.go`

## Arquivos que poderão ser modificados
- `internal/cache/types.go`
- `internal/telegram/handlers/events/postBuilder/postBuilder.go`

## Estratégia de implementação
1. **Atualizar Modelo de Dados:** Adicionar o campo `CustomEmojiID` na struct `PostBuilderButton` em `internal/cache/types.go`.
2. **Processamento de Botões de URL:**
   - No handler `awaiting_button`, analisar as `Entities` da mensagem enviada pelo usuário.
   - Se houver uma entidade do tipo `custom_emoji` no texto do nome do botão (primeira linha), extrair o `CustomEmojiID` e salvá-lo.
3. **Processamento de Reações:**
   - No handler `awaiting_reactions`, permitir a entrada de emojis customizados.
   - Para armazenar os IDs sem quebrar a estrutura de string atual (evitando migrações complexas), utilizaremos um prefixo `eid:` para identificar IDs de emojis customizados na string de reações (ex: `👍,eid:1234567,❤️`).
4. **Renderização de Botões:**
   - Nas funções `sendFinalPostTelego` e `InlineHandlerTelego`, ao iterar pelos botões e reações, verificar se existe um `CustomEmojiID` ou um prefixo `eid:`.
   - Se existir, configurar o campo `IconCustomEmojiID` do objeto `telego.InlineKeyboardButton`.
5. **Ajuste no Menu:** Atualizar `showMenuTelego` para exibir uma representação amigável (ou apenas omitir o prefixo) ao listar as reações.

## Passos detalhados

1. **Em `internal/cache/types.go`:**
   - Adicionar `CustomEmojiID string `json:"custom_emoji_id,omitempty"`` na struct `PostBuilderButton`.

2. **Em `internal/telegram/handlers/events/postBuilder/postBuilder.go`:**
   - Modificar `handleTextInputTelego` -> `case "awaiting_button"`:
     - Iterar pelas entidades da mensagem. Se uma entidade `custom_emoji` estiver dentro do range da primeira linha (nome), capturar seu ID.
   - Modificar `handleTextInputTelego` -> `case "awaiting_reactions"`:
     - Refatorar a lógica para identificar se cada parte da string (separada por vírgula) contém um emoji Unicode ou um emoji customizado (via entidades).
     - Se for customizado, armazenar como `eid:<id>`.
   - Modificar `sendFinalPostTelego` e `InlineHandlerTelego`:
     - Aplicar `IconCustomEmojiID` ao criar os botões de URL.
     - Aplicar `IconCustomEmojiID` ao criar os botões de reação, detectando o prefixo `eid:`.
   - Modificar `isEmoji`: Ajustar para aceitar o formato `eid:` ou caracteres especiais de emojis customizados.

## Riscos
- **Compatibilidade:** O campo `IconCustomEmojiID` requer uma versão recente do Bot API. Como o `go build` confirmou a presença do campo no SDK, o risco é de o bot (se não for Premium/Collectible) não conseguir enviar a mensagem. O Telegram costuma apenas ignorar o campo ou retornar erro se o bot não tiver permissão.

## Impactos esperados
- Usuários poderão criar botões muito mais visualmente atraentes usando o vasto catálogo de emojis customizados do Telegram.

## Como testar
1. Criar um PostBuilder.
2. Adicionar um botão de URL enviando: `Meu Botão [emoji_custom] \n https://google.com`.
3. Adicionar reações enviando: `👍, [emoji_custom], 🔥`.
4. Verificar se os emojis aparecem corretamente no Preview e no envio Inline.
