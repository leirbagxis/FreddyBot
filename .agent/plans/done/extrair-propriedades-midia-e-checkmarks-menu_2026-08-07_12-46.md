# Plano: Extração Automática de Mídias (Mencções/Canais) e Checkmarks no Menu

## Pedido do usuário
1. **Extração Completa de Mídias (com menções/canais, entidades e botões)**:
   - Ao receber qualquer mídia (foto, vídeo, documento, animação, áudio, sticker), seja de um usuário ou encaminhada de canais (mídias citando/marcando canais), identificar o PostBuilder e capturar **todas** as propriedades: texto da legenda, menções a canais (`@canal`), links, entidades de formatação HTML, emojis premium/customizados e botões inline existentes na mídia.
   - Aplicar todas as propriedades extraídas imediatamente no rascunho do PostBuilder.
2. **Indicadores de Checkmark (`✅` / `❌`) no Menu de Edição**:
   - No resumo do menu e nos botões de edição de `showMenuTelego`, substituir os textos brutos longos e a contagem de botões por um emoji checkmark (`✅` se o item estiver preenchido/configurado, `❌` se estiver vazio).

## Objetivo
- Importar 100% da legenda, menções, links, emojis premium e botões de qualquer mídia recebida.
- Tornar o menu do PostBuilder extremamente visual e legível com status `✅` / `❌` em cada opção.

## Contexto atual
- `HandlerTelego` já usa `ProcessTextWithFormattingTelego` para converter a legenda em HTML (suportando menções a canais, links e custom emojis).
- Botões inline vindos com a mídia recebida (`update.Message.ReplyMarkup`) não eram lidos no estado inicial do PostBuilder.
- `showMenuTelego` exibia o texto bruto de Título, Corpo, Rodapé, Reações e contagem de botões, gerando poluição no menu.

## Arquivos analisados
- `internal/telegram/handlers/events/postBuilder/postBuilder.go`
- `internal/telegram/events/channelPost/formatting_telego.go`

## Arquivos que poderão ser modificados
- `internal/telegram/handlers/events/postBuilder/postBuilder.go`

## Estratégia de implementação

1. **Leitura completa da mídia em `HandlerTelego`**:
   - Capturar legenda com `ProcessTextWithFormattingTelego` (preservando menções a canais `@canal`, links `t.me/`, marcas d'água e emojis premium `<tg-emoji>`).
   - Ler `update.Message.ReplyMarkup`: se a mídia contiver botões inline com URL, extrair cada botão para `cache.PostBuilderButton{Text, URL, CustomEmojiID}`.
   - Inicializar `PostBuilderState` com `Body: captionText` e `Buttons: initialButtons`.

2. **Indicadores de Checkmark (`✅` / `❌`) em `showMenuTelego`**:
   - Helper `check(filled bool) string`: retorna `"✅"` se preenchido e `"❌"` se vazio.
   - Resumo do texto:
     - 📝 **Título:** `✅` ou `❌`
     - 📄 **Corpo:** `✅` ou `❌`
     - 👣 **Rodapé:** `✅` ou `❌`
     - 🎭 **Reações:** `✅` ou `❌`
     - 🔘 **Botões:** `✅` ou `❌`
   - Botões inline do teclado de edição:
     - `📝 Título ✅` / `📝 Título ❌`
     - `📄 Corpo ✅` / `📄 Corpo ❌`
     - `👣 Rodapé ✅` / `👣 Rodapé ❌`
     - `🎭 Reações ✅` / `🎭 Reações ❌`
     - `🔘 Botões ✅` / `🔘 Botões ❌`

---

## Passos detalhados

### Passo 1: Atualizar `HandlerTelego` em `postBuilder.go`
- Ler a legenda com `ProcessTextWithFormattingTelego`.
- Ler `update.Message.ReplyMarkup` para extrair botões inline preexistentes.
- Salvar o estado inicial no Redis.

### Passo 2: Atualizar `showMenuTelego` em `postBuilder.go`
- Substituir o texto bruto e contagens do menu pelo formato `✅` / `❌`.
- Adicionar os sufixos `✅` / `❌` nos botões do menu de edição.

### Passo 3: Compilação e Testes
- Executar `go build ./cmd/FreddyBot` e `go test ./...`.

---

## Riscos
- Nenhum. Preserva 100% da compatibilidade com a Bot API.

## Impactos esperados
- Mídias recebidas com menções a canais, emojis premium e botões têm tudo extraído automaticamente no primeiro clique.
- Menu visualmente limpo e fácil de verificar o que já foi preenchido.

## Compatibilidade
- Linux / macOS / Windows / Docker / Telego

## Como testar
1. Enviar mídia com legenda contendo menção a canal `@canal` e botões inline -> Clicar em "🛠️ Post Builder".
2. Verificar que o Corpo e Botões aparecem com o status `✅` no menu.
3. Clicar em `👁️ Preview` para ver a mídia montada com as propriedades importadas.
