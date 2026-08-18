# Plano: separador-emoji-customizado

## Pedido do usuário
Permitir que o usuário use emojis customizados (Premium) como separador entre postagens, além dos stickers normais. Só funciona quando o usuário tem conta Telegram conectada. O usuário pode enviar sticker normal OU emoji premium para usar como separador.

## Objetivo
- Estender o Separator model para suportar dois tipos: "sticker" e "custom_emoji"
- Modificar o handler de captura para detectar se o usuário enviou um sticker ou uma mensagem com custom emoji
- Modificar o envio do separador para: se for sticker → `Bot.SendSticker()`; se for custom emoji → enviar mensagem de texto com a entity `custom_emoji`
- Usar MTProto quando disponível para enviar o custom emoji (já que o usuário tem conta conectada)

## Contexto atual
- `Separator` model: `ID`, `SeparatorID` (sticker FileID), `SeparatorURL`, `OwnerChannelID`
- `SetStickerSeparatorHandlerTelego`: captura STICKER enviado pelo usuário, extrai FileID, salva
- `ProcessSeparatorTelego`: envia sticker via `Bot.SendSticker()`
- Cache do canal não era invalidado ao salvar separator (corrigido)

## Arquivos analisados
- `internal/database/models/models.go` - Separator model
- `internal/telegram/handlers/callbacks/my_channel/channel_actions.go` - SetStickerSeparatorHandlerTelego
- `internal/telegram/events/channelPost/dispatch_telego.go` - ProcessSeparatorTelego
- `config/messages.yml` - Mensagens do separador
- `internal/telegram/handlers/events/postBuilder/postBuilder.go` - Referência de custom emoji
- `internal/telegram/executor/entities.go` - Conversão de entities para gotd

## Arquivos que poderão ser modificados
- `internal/database/models/models.go` - Adicionar campos `Type`, `EmojiText`, `EmojiID`
- `internal/telegram/handlers/callbacks/my_channel/channel_actions.go` - Modificar handler para aceitar sticker OU texto com emoji
- `internal/telegram/events/channelPost/dispatch_telego.go` - Modificar `ProcessSeparatorTelego` para enviar custom emoji
- `config/messages.yml` - Atualizar mensagens (require-separator-message, success-save-separator)

## Estrategia de implementacao

### 1. Model (`Separator`)
Adicionar campos:
```go
Type      string  `gorm:"default:sticker"` // "sticker" ou "custom_emoji"
EmojiText string  // caracter do emoji (ex: "🌟")
EmojiID   string  // ID do emoji customizado
```
`SeparatorID` continua sendo usado para sticker `FileID`.

### 2. Handler de captura (`SetStickerSeparatorHandlerTelego`)
- Renomear para `SetSeparatorHandlerTelego` (mais generico)
- Se `update.Message.Sticker != nil`: fluxo atual (salva como "sticker")
- Se `update.Message.Text != ""` e tiver `custom_emoji` entity: extrair emoji text + emoji ID, salvar como "custom_emoji"
- Se `update.Message.Caption != ""` e tiver `custom_emoji` entity: extrair da caption
- Exigir conta conectada para custom emoji (emoji premium precisa de conta Telegram)
- Atualizar predicate `matchAwaitingStickerSeparatorTelego` se necessario

### 3. Envio (`ProcessSeparatorTelego`)
- Se `channel.Separator.Type == "sticker"`: `Bot.SendSticker()` (atual)
- Se `channel.Separator.Type == "custom_emoji"`: 
  - Se tiver executor (MTProto): enviar via MTProto com a entity `custom_emoji`
  - Fallback: `Bot.SendMessage()` com `Entities` contendo `custom_emoji`

### 4. Mensagens (YAML)
- Atualizar `require-separator-message` para mencionar que pode enviar sticker OU emoji
- Atualizar `success-save-separator` para mencionar o tipo salvo

## Passos detalhados

1. Alterar `Separator` model: adicionar `Type`, `EmojiText`, `EmojiID`
2. Renomear/expandir `SetStickerSeparatorHandlerTelego` para aceitar sticker OU texto com custom emoji
3. Atualizar predicate `matchAwaitingStickerSeparatorTelego` (ou criar novo)
4. Atualizar `ProcessSeparatorTelego` para enviar custom emoji via Bot API (Entities) ou MTProto
5. Atualizar `config/messages.yml` com as novas mensagens
6. Adicionar migration ou auto-migrate para as novas colunas
7. Build e teste

## Riscos
- Quebrar compatibilidade com separadores existentes (tipo "sticker") - Migrar automaticamente: `Type` default "sticker" para registros existentes
- Custom emoji entity pode nao ser suportada pela versao do telego - Verificar se telego suporta `MessageEntityCustomEmoji`
- Usuario pode enviar texto com multiplos emojis - Usar apenas o primeiro `custom_emoji` entity encontrado

## Impactos esperados
- Separadores existentes (sticker) continuam funcionando sem alteracao
- Usuarios com conta conectada podem usar emoji premium como separador
- Nenhuma regressao no fluxo de legendas ou entities

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
go vet ./...
```

### Execucao
```bash
./Release
```

1. Configurar separador com sticker (fluxo normal) → deve continuar funcionando
2. Configurar separador com emoji customizado:
   - Clicar em "Sticker Separador"
   - Enviar mensagem de TEXTO com um emoji premium
   - Verificar se salvou como "custom_emoji"
3. Postar no canal → verificar se o emoji aparece como separador

## Rollback
Reverter alteracoes nos arquivos:
- `internal/database/models/models.go` - Remover campos novos
- `internal/telegram/handlers/callbacks/my_channel/channel_actions.go` - Reverter para versao anterior
- `internal/telegram/events/channelPost/dispatch_telego.go` - Reverter `ProcessSeparatorTelego`
- `config/messages.yml` - Reverter mensagens

## Observacoes
- O custom emoji entity `custom_emoji` tem os campos `Offset`, `Length`, `CustomEmojiID`
- Para enviar via Bot API, usar `SendMessageParams.Entities` com `telego.MessageEntity{Type: "custom_emoji", Offset: 0, Length: len(emojiText), CustomEmojiID: emojiID}`
- Para enviar via MTProto, usar `MessageEntityCustomEmoji{CustomEmojiID: emojiID}` com offset/length corretos
- O texto da mensagem deve conter o caractere do emoji (ex: "��") para o offset/length funcionarem
