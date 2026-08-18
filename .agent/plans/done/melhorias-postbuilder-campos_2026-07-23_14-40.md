# Plano: melhorias-postbuilder-campos

## Pedido do usuário
Três melhorias no Post Builder:

1. **Capturar caption/entidades da mídia original** — quando o usuário envia uma mídia que já tem legenda com formatação (entities), o bot deve identificar, preservar e pré-popular os campos do Post Builder para o usuário poder modificar depois.

2. **Botão "Voltar" nos prompts de edição** — quando o usuário clica acidentalmente em um campo (ex: Título), ele entra em estado `awaiting_*` e precisa enviar texto. Quer um botão no prompt para voltar ao menu sem precisar enviar texto.

3. **Ajustar campos conforme o tipo de mídia** — sticker não suporta legenda no Telegram. Os campos de texto (Título, Corpo, Rodapé) devem ser ocultados/desabilitados para stickers.

## Objetivo
Melhorar a UX do Post Builder: capturar dados existentes da mídia, oferecer escape de estados de espera, e adaptar a UI ao tipo de mídia.

## Contexto atual
- `HandlerTelego()` captura `MediaType` + `MediaFileID` mas ignora `update.Message.Caption` e `update.Message.CaptionEntities`
- `showMenuTelego()` sempre exibe Título/Corpo/Rodapé independente do `MediaType`
- Prompts de edição (`pb-edit-title`, `pb-edit-body`, etc.) enviam mensagem sem `ReplyMarkup`, sem botão de cancelamento
- `sendFinalPostTelego()` já trata sticker corretamente (sem caption), mas o menu engana o usuário

## Arquivos analisados
- `internal/telegram/handlers/events/postBuilder/postBuilder.go`
- `internal/cache/types.go`
- `internal/telegram/loader_telego.go`
- `internal/telegram/events/channelPost/formatting_telego.go`
- `internal/telegram/events/channelPost/entities_telego.go`

## Arquivos que poderão ser modificados
- `internal/telegram/handlers/events/postBuilder/postBuilder.go`

## Estratégia de implementação

### 1. Capturar caption/entidades (HandlerTelego)
- Após detectar o tipo de mídia, verificar `update.Message.Caption`
- Se houver caption, processar com `ProcessTextWithFormattingTelego(caption, captionEntities)`
- Armazenar no `state.Body` (campo principal de conteúdo)
- Exibir mensagem informando que a legenda foi capturada

### 2. Botão Voltar nos prompts (CallbackHandlerTelego)
- Em cada caso de edição (`pb-edit-title`, `pb-edit-body`, `pb-edit-footer`, `pb-edit-reactions`, `pb-add-button`), adicionar `ReplyMarkup` com botão "🔙 Voltar ao Menu" com `CallbackData: "pb-start"`
- `pb-start` já funciona corretamente: mantém o state e chama `showMenuTelego`

### 3. Ajustar campos por tipo de mídia (showMenuTelego + CallbackHandlerTelego)
- Em `showMenuTelego()`, se `state.MediaType == "sticker"`:
  - Mostrar "📝 Texto: não suportado para stickers" no lugar dos campos
  - Remover botões de edição de Título/Corpo/Rodapé do teclado
- Em `CallbackHandlerTelego()`, ao entrar em prompt de edição, verificar se o MediaType permite o campo editado

## Passos detalhados

1. **Editar `HandlerTelego` (linha ~168)**
   - Extrair `update.Message.Caption` com `ProcessTextWithFormattingTelego`
   - Incluir no `PostBuilderState` inicial como `Body`
   - Ajustar mensagem de confirmação

2. **Editar `CallbackHandlerTelego` — prompts de edição (linhas ~706-760)**
   - Adicionar `ReplyMarkup` com `InlineKeyboardMarkup` e botão "🔙 Voltar ao Menu" (`pb-start`)
   - Em prompts de sticker, verificar se o campo editado é suportado

3. **Editar `showMenuTelego` (linha ~341)**
   - Adicionar verificação de `state.MediaType`
   - Para sticker: modificar exibição e teclado

## Riscos
- Nenhum risco de regressão: mudanças são aditivas (adicionam comportamento sem remover existente)
- `pb-start` já é testado em produção; adicionar ao prompt não quebra fluxo
- Sticker sem caption já funciona; esconder campos só melhora UX

## Impactos esperados
- Usuário não precisa redigitar legenda que já veio com a mídia
- Usuário pode sair de estados `awaiting_*` sem enviar texto
- Usuário não tenta adicionar texto a sticker (que seria ignorado)

## Compatibilidade
- Linux
- macOS
- Windows
- Docker
- CI/CD

## Como testar

### Build
```bash
go build ./...
```

### Testes
```bash
go test ./...
```

### Execução
```bash
go run cmd/FreddyBot/main.go
```

## Rollback
Reverter commits ou restaurar `postBuilder.go` do git.

## Observações
- O botão usa `CallbackData: "pb-start"` em vez de `"pb-cancel"` para preservar o state atual (volta ao menu sem perder dados)
- Para stickers, `sendFinalPostTelego` já trata corretamente (sem caption). Apenas o menu será ajustado.
