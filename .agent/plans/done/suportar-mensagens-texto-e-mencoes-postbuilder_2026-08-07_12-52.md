# Plano: Suportar Mensagens Encaminhadas/Citações (Mídias ou Texto) no PostBuilder

## Pedido do usuário
Garantir que, quando uma mensagem é recebida ou encaminhada (citando um canal):
1. O bot verifique se a mensagem contém qualquer tipo de mídia (foto, vídeo, animação, áudio, documento ou sticker), capturando seu `FileID` e `mediaType`.
2. Se não houver anexo de mídia em arquivo, reconhecer como postagem de texto (`mediaType = "text"`).
3. Em ambos os casos, capturar o texto/legenda (com menções `@canal`, links e formatação HTML) e botões inline pré-existentes (`ReplyMarkup`), oferecendo o botão `🛠️ Post Builder`.

## Causa Raiz Identificada
- Em `loader_telego.go`, a regra de correspondência `matchPostBuilderTelego` descartava mensagens sem arquivo de mídia antes de chegarem ao handler do PostBuilder.
- Em `postBuilder.go`, a checagem `if mediaID == ""` descartava posts de texto ou citações quando o usuário não estava em uma etapa de entrada ativa.

## Objetivo
- Fazer a checagem em cascata de todas as mídias (foto, vídeo, animação, áudio, documento, sticker) e, na ausência delas, aceitar o texto como `mediaType = "text"`.
- Extrair legenda/texto (com menções a canais), emojis customizados e botões da mensagem recebida/encaminhada.

## Arquivos analisados
- `internal/telegram/loader_telego.go`
- `internal/telegram/handlers/events/postBuilder/postBuilder.go`

## Arquivos que poderão ser modificados
- `internal/telegram/loader_telego.go`
- `internal/telegram/handlers/events/postBuilder/postBuilder.go`

## Estratégia de implementação

1. **Atualizar `matchPostBuilderTelego` em `loader_telego.go`**:
   - Aceitar mensagens com mídias (foto, vídeo, animação, áudio, documento, sticker) OU texto/legenda que não sejam comandos `/`.

2. **Atualizar `HandlerTelego` em `postBuilder.go`**:
   - Verificar na sequência: Foto ➔ Vídeo ➔ Animação ➔ Áudio ➔ Documento ➔ Sticker ➔ Texto (`mediaType = "text"`).
   - Se houver sessão com etapa ativa (`state.Step != ""`), encaminhar para entrada de texto.
   - Caso contrário (nova mensagem/encaminhamento):
     - Extrair legenda (`update.Message.Caption`) ou texto (`update.Message.Text`) preservando menções a canais e formatações com `ProcessTextWithFormattingTelego`.
     - Extrair botões inline de `update.Message.ReplyMarkup`.
     - Salvar o novo `PostBuilderState` e responder com o prompt do `🛠️ Post Builder`.

---

## Passos detalhados

### Passo 1: Atualizar `matchPostBuilderTelego` em `loader_telego.go`
- Permitir a passagem de mensagens de mídia ou texto em chats privados.

### Passo 2: Atualizar `HandlerTelego` em `postBuilder.go`
- Implementar a detecção em cascata (Foto, Vídeo, Animação, Áudio, Documento, Sticker, Texto).
- Tratar entradas em etapas ativas vs novo rascunho de post.
- Capturar texto/legenda e botões.

### Passo 3: Compilação e Testes
- Executar `go build ./cmd/FreddyBot` e `go test ./...`.

---

## Riscos
- Nenhum. Preserva todas as rotas e etapas ativas.

## Impactos esperados
- Qualquer post encaminhado ou enviado citando um canal (com mídia ou texto puro) ativa o PostBuilder.
- Mídias (fotos, vídeos, áudios) são priorizadas e capturadas com sucesso.

## Compatibilidade
- Linux / macOS / Windows / Docker / Telego

## Como testar
1. Encaminhar ou enviar post com foto/vídeo citando um canal -> Bot identifica `mediaType = "photo"` (ou `"video"`), extrai a legenda e ativa o PostBuilder.
2. Encaminhar ou enviar post contendo apenas texto citando um canal -> Bot identifica `mediaType = "text"`, extrai o texto e ativa o PostBuilder.
