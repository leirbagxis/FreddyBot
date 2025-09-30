package channelpost

import (
	"fmt"
	"html"
	"regexp"
	"sort"
	"strings"
	"unicode/utf16"

	"github.com/go-telegram/bot/models"
	dbmodels "github.com/leirbagxis/FreddyBot/internal/database/models"
)

// ✅ REGEX para conversão de Markdown para HTML
var (
	boldRegex          = regexp.MustCompile(`\*\*([^*]+)\*\*`)
	italicRegex        = regexp.MustCompile(`\*([^*]+)\*`)
	underlineRegex     = regexp.MustCompile(`__([^_]+)__`)
	strikethroughRegex = regexp.MustCompile(`~~([^~]+)~~`)
	codeRegex          = regexp.MustCompile("`([^`]+)`")
	linkRegex          = regexp.MustCompile(`\[([^\]]+)\]\(([^)]+)\)`)
	spoilerRegex       = regexp.MustCompile(`\|\|([^|]+)\|\|`)
	blockquoteRegex    = regexp.MustCompile(`(?m)^>\s*(.+)$`)
)

// ✅ CORRIGIDO: processTextWithFormatting
func processTextWithFormatting(text string, entities []models.MessageEntity) string {
	if text == "" {
		return ""
	}

	// ✅ SE TEM ENTITIES, PROCESSAR APENAS ELAS (não aplicar markdown)
	if len(entities) > 0 {
		return processEntitiesOnly(text, entities)
	}

	// ✅ SE NÃO TEM ENTITIES, APLICAR FORMATAÇÃO MARKDOWN
	return detectParseMode(text)
}

// ✅ CORRIGIDO: Processar entities para HTML
// func processEntitiesOnly(text string, entities []models.MessageEntity) string {
// 	if len(entities) == 0 {
// 		return html.EscapeString(text)
// 	}

// 	// Converter string para UTF-16 para cálculos corretos de offset
// 	textRunes := []rune(text)
// 	textUTF16 := utf16.Encode(textRunes)

// 	// Ordenar entidades por offset
// 	sort.Slice(entities, func(i, j int) bool {
// 		return entities[i].Offset < entities[j].Offset
// 	})

// 	var result strings.Builder
// 	lastOffset := 0

// 	for _, entity := range entities {
// 		// Adicionar texto antes da entidade (escapado)
// 		if entity.Offset > lastOffset {
// 			beforeText := string(utf16.Decode(textUTF16[lastOffset:entity.Offset]))
// 			result.WriteString(html.EscapeString(beforeText))
// 		}

// 		// Extrair o texto da entidade
// 		entityEnd := entity.Offset + entity.Length
// 		if entityEnd > len(textUTF16) {
// 			entityEnd = len(textUTF16)
// 		}

// 		entityText := string(utf16.Decode(textUTF16[entity.Offset:entityEnd]))

// 		// ✅ CORRIGIDO: GERAR HTML AO INVÉS DE MARKDOWN
// 		switch entity.Type {
// 		case "bold":
// 			result.WriteString("<b>")
// 			result.WriteString(html.EscapeString(entityText))
// 			result.WriteString("</b>")
// 		case "italic":
// 			result.WriteString("<i>")
// 			result.WriteString(html.EscapeString(entityText))
// 			result.WriteString("</i>")
// 		case "underline":
// 			result.WriteString("<u>")
// 			result.WriteString(html.EscapeString(entityText))
// 			result.WriteString("</u>")
// 		case "strikethrough":
// 			result.WriteString("<s>")
// 			result.WriteString(html.EscapeString(entityText))
// 			result.WriteString("</s>")
// 		case "spoiler":
// 			result.WriteString("<tg-spoiler>")
// 			result.WriteString(html.EscapeString(entityText))
// 			result.WriteString("</tg-spoiler>")
// 		case "code":
// 			result.WriteString("<code>")
// 			result.WriteString(html.EscapeString(entityText))
// 			result.WriteString("</code>")
// 		case "pre":
// 			if entity.Language != "" {
// 				result.WriteString("<pre><code class=\"")
// 				result.WriteString(html.EscapeString(entity.Language))
// 				result.WriteString("\">")
// 			} else {
// 				result.WriteString("<pre><code>")
// 			}
// 			result.WriteString(html.EscapeString(entityText))
// 			result.WriteString("</code></pre>")
// 		case "blockquote":
// 			result.WriteString("<blockquote>")
// 			result.WriteString(html.EscapeString(entityText))
// 			result.WriteString("</blockquote>")
// 		case "expandable_blockquote":
// 			result.WriteString("<blockquote expandable>")
// 			result.WriteString(html.EscapeString(entityText))
// 			result.WriteString("</blockquote>")
// 		case "text_link":
// 			result.WriteString("<a href=\"")
// 			result.WriteString(html.EscapeString(entity.URL))
// 			result.WriteString("\">")
// 			result.WriteString(html.EscapeString(entityText))
// 			result.WriteString("</a>")
// 		case "url":
// 			result.WriteString("<a href=\"")
// 			result.WriteString(html.EscapeString(entityText))
// 			result.WriteString("\">")
// 			result.WriteString(html.EscapeString(entityText))
// 			result.WriteString("</a>")
// 		case "email":
// 			result.WriteString("<a href=\"mailto:")
// 			result.WriteString(html.EscapeString(entityText))
// 			result.WriteString("\">")
// 			result.WriteString(html.EscapeString(entityText))
// 			result.WriteString("</a>")
// 		case "phone_number":
// 			result.WriteString("<a href=\"tel:")
// 			result.WriteString(html.EscapeString(entityText))
// 			result.WriteString("\">")
// 			result.WriteString(html.EscapeString(entityText))
// 			result.WriteString("</a>")
// 		case "mention":
// 			result.WriteString("<a href=\"https://t.me/")
// 			result.WriteString(html.EscapeString(strings.TrimPrefix(entityText, "@")))
// 			result.WriteString("\">")
// 			result.WriteString(html.EscapeString(entityText))
// 			result.WriteString("</a>")
// 		case "hashtag":
// 			result.WriteString("<a href=\"https://t.me/hashtag/")
// 			result.WriteString(html.EscapeString(strings.TrimPrefix(entityText, "#")))
// 			result.WriteString("\">")
// 			result.WriteString(html.EscapeString(entityText))
// 			result.WriteString("</a>")
// 		case "cashtag":
// 			result.WriteString("<a href=\"https://t.me/cashtag/")
// 			result.WriteString(html.EscapeString(strings.TrimPrefix(entityText, "$")))
// 			result.WriteString("\">")
// 			result.WriteString(html.EscapeString(entityText))
// 			result.WriteString("</a>")
// 		case "bot_command":
// 			result.WriteString("<code>")
// 			result.WriteString(html.EscapeString(entityText))
// 			result.WriteString("</code>")
// 		case "custom_emoji":
// 			if entity.CustomEmojiID != "" {
// 				result.WriteString("<tg-emoji emoji-id=\"")
// 				result.WriteString(html.EscapeString(entity.CustomEmojiID))
// 				result.WriteString("\">")
// 				result.WriteString(html.EscapeString(entityText))
// 				result.WriteString("</tg-emoji>")
// 			} else {
// 				result.WriteString(html.EscapeString(entityText))
// 			}
// 		default:
// 			result.WriteString(html.EscapeString(entityText))
// 		}

// 		lastOffset = entityEnd
// 	}

// 	// Adicionar texto restante (escapado)
// 	if lastOffset < len(textUTF16) {
// 		remainingText := string(utf16.Decode(textUTF16[lastOffset:]))
// 		result.WriteString(html.EscapeString(remainingText))
// 	}

// 	return result.String()
// }

func processEntitiesOnly(text string, entities []models.MessageEntity) string {
	if text == "" {
		return ""
	}
	if len(entities) == 0 {
		return html.EscapeString(text)
	}

	// Telegram usa offsets/lengths baseados em UTF-16
	r := []rune(text)
	u := utf16.Encode(r)

	type ent struct {
		t     string
		o     int
		end   int
		url   string
		lang  string
		emoji string
		uid   int64 // para text_mention
	}

	// Normaliza entities e delimita dentro do range
	es := make([]ent, 0, len(entities))
	for _, e := range entities {
		o := e.Offset
		end := e.Offset + e.Length
		if o < 0 {
			o = 0
		}
		if end > len(u) {
			end = len(u)
		}
		if o >= end {
			continue
		}
		var uid int64
		if e.User != nil {
			uid = e.User.ID
		}
		es = append(es, ent{
			t:     string(e.Type),
			o:     o,
			end:   end,
			url:   e.URL,
			lang:  e.Language,
			emoji: e.CustomEmojiID,
			uid:   uid,
		})
	}
	if len(es) == 0 {
		return html.EscapeString(text)
	}

	// Calcula boundaries para spans mínimos não sobrepostos
	bset := map[int]struct{}{0: {}, len(u): {}}
	for _, e := range es {
		bset[e.o] = struct{}{}
		bset[e.end] = struct{}{}
	}
	bounds := make([]int, 0, len(bset))
	for b := range bset {
		bounds = append(bounds, b)
	}
	sort.Ints(bounds)

	// Prioridade estável das entities
	// 0: code/pre (exclusivos) → 1: blockquote (outer) → 2: links (text_link, text_mention, url, email, phone)
	// 3..N: estilos inline (bold, italic, underline, strikethrough, spoiler) → custom_emoji dentro do conteúdo
	pri := func(t string) int {
		switch t {
		case "code", "pre":
			return 0
		case "blockquote":
			return 1
		case "text_link", "text_mention":
			return 2
		case "url", "email", "phone_number":
			return 2
		case "bold":
			return 3
		case "italic":
			return 4
		case "underline":
			return 5
		case "strikethrough":
			return 6
		case "spoiler":
			return 7
		case "custom_emoji":
			return 8
		default:
			// mention/hashtag/cashtag/bot_command/bank_card etc. são mantidos como texto
			return 9
		}
	}

	// Renderiza um span dado o conjunto de entities ativas
	spanHTML := func(seg string, active []ent) string {
		// code/pre são exclusivos e retornam imediatamente
		for _, a := range active {
			if a.t == "code" {
				return "<code>" + html.EscapeString(seg) + "</code>"
			}
			if a.t == "pre" {
				langClass := ""
				if a.lang != "" {
					// Telegram aceita class no <code> para linguagem
					langClass = ` class="language-` + html.EscapeString(a.lang) + `"`
				}
				return "<pre><code" + langClass + ">" + html.EscapeString(seg) + "</code></pre>"
			}
		}

		// Base escapada
		inner := html.EscapeString(seg)

		// custom_emoji envolve o conteúdo
		for _, a := range active {
			if a.t == "custom_emoji" && a.emoji != "" {
				inner = `<tg-emoji emoji-id="` + html.EscapeString(a.emoji) + `">` + inner + "</tg-emoji>"
				break
			}
		}

		// Estilos inline (ordem determinística, aninhados)
		wrapIf := func(t, open, close string) {
			for _, a := range active {
				if a.t == t {
					inner = open + inner + close
					break
				}
			}
		}
		wrapIf("bold", "<b>", "</b>")
		wrapIf("italic", "<i>", "</i>")
		wrapIf("underline", "<u>", "</u>")
		wrapIf("strikethrough", "<s>", "</s>")
		// spoiler: Telegram suporta <tg-spoiler>
		for _, a := range active {
			if a.t == "spoiler" {
				inner = "<tg-spoiler>" + inner + "</tg-spoiler>"
				break
			}
		}

		// Links (externos ao conteúdo estilizado)
		// text_link
		for _, a := range active {
			if a.t == "text_link" && a.url != "" {
				inner = `<a href="` + html.EscapeString(a.url) + `">` + inner + "</a>"
				// apenas um text_link deve envolver
				break
			}
		}
		// text_mention → tg://user?id=...
		for _, a := range active {
			if a.t == "text_mention" && a.uid != 0 {
				inner = `<a href="tg://user?id=` + html.EscapeString(int64ToStr(a.uid)) + `">` + inner + "</a>"
				break
			}
		}
		// url/email/phone_number
		for _, a := range active {
			switch a.t {
			case "url":
				inner = `<a href="` + inner + `">` + inner + "</a>"
			case "email":
				inner = `<a href="mailto:` + inner + `">` + inner + "</a>"
			case "phone_number":
				inner = `<a href="tel:` + inner + `">` + inner + "</a>"
			}
		}

		// blockquote como wrapper mais externo (após links)
		for _, a := range active {
			if a.t == "blockquote" {
				inner = "<blockquote>" + inner + "</blockquote>"
				break
			}
		}

		return inner
	}

	var b strings.Builder
	for i := 0; i < len(bounds)-1; i++ {
		s, e := bounds[i], bounds[i+1]
		if e <= s {
			continue
		}
		seg := string(utf16.Decode(u[s:e]))

		// Entities que cobrem integralmente este span
		active := make([]ent, 0, 6)
		for _, en := range es {
			if en.o <= s && en.end >= e {
				active = append(active, en)
			}
		}
		if len(active) == 0 {
			b.WriteString(html.EscapeString(seg))
			continue
		}

		// Ordena por prioridade estável (e depois por abrangência/tie-break)
		sort.SliceStable(active, func(i, j int) bool {
			pi, pj := pri(active[i].t), pri(active[j].t)
			if pi != pj {
				return pi < pj
			}
			di, dj := (active[i].end - active[i].o), (active[j].end - active[j].o)
			if di != dj {
				return di > dj
			}
			if active[i].o != active[j].o {
				return active[i].o < active[j].o
			}
			return active[i].t < active[j].t
		})

		b.WriteString(spanHTML(seg, active))
	}
	return b.String()
}

// Pequena utilidade para evitar fmt import só por isso
func int64ToStr(v int64) string {
	// conversão rápida sem fmt
	// para simplificar e manter dependências baixas
	if v == 0 {
		return "0"
	}
	neg := v < 0
	if neg {
		v = -v
	}
	var buf [20]byte
	i := len(buf)
	for v > 0 {
		i--
		buf[i] = byte('0' + v%10)
		v /= 10
	}
	if neg {
		i--
		buf[i] = '-'
	}
	return string(buf[i:])
}

// ✅ CORRIGIDO: processInlineFormatting para gerar HTML COM ESCAPE ADEQUADO
func processInlineFormatting(line string) string {
	if line == "" {
		return ""
	}

	// ✅ PROCESSAR EM ORDEM DE PRIORIDADE, ESCAPANDO ADEQUADAMENTE

	// 1. Code inline (maior prioridade) - processa primeiro para evitar conflitos
	codeRegex := regexp.MustCompile("`([^`]+)`")
	line = codeRegex.ReplaceAllStringFunc(line, func(match string) string {
		content := match[1 : len(match)-1]
		// ✅ ESCAPAR CONTEÚDO DO CODE
		return fmt.Sprintf("<code>%s</code>", html.EscapeString(content))
	})

	// 2. Spoiler
	spoilerRegex := regexp.MustCompile(`\|\|([^|]+)\|\|`)
	line = spoilerRegex.ReplaceAllStringFunc(line, func(match string) string {
		content := spoilerRegex.FindStringSubmatch(match)[1]
		// ✅ ESCAPAR CONTEÚDO DO SPOILER
		return fmt.Sprintf("<tg-spoiler>%s</tg-spoiler>", html.EscapeString(content))
	})

	// 3. Bold
	boldRegex := regexp.MustCompile(`\*\*([^*]+)\*\*`)
	line = boldRegex.ReplaceAllStringFunc(line, func(match string) string {
		content := boldRegex.FindStringSubmatch(match)[1]
		// ✅ ESCAPAR CONTEÚDO DO BOLD
		return fmt.Sprintf("<b>%s</b>", html.EscapeString(content))
	})

	// 4. Underline
	underlineRegex := regexp.MustCompile(`__([^_]+)__`)
	line = underlineRegex.ReplaceAllStringFunc(line, func(match string) string {
		content := underlineRegex.FindStringSubmatch(match)[1]
		// ✅ ESCAPAR CONTEÚDO DO UNDERLINE
		return fmt.Sprintf("<u>%s</u>", html.EscapeString(content))
	})

	// 5. Strikethrough
	strikeRegex := regexp.MustCompile(`~~([^~]+)~~`)
	line = strikeRegex.ReplaceAllStringFunc(line, func(match string) string {
		content := strikeRegex.FindStringSubmatch(match)[1]
		// ✅ ESCAPAR CONTEÚDO DO STRIKETHROUGH
		return fmt.Sprintf("<s>%s</s>", html.EscapeString(content))
	})

	// 6. Italic (deve vir por último para evitar conflitos com **)
	italicRegex := regexp.MustCompile(`\*([^*]+)\*`)
	line = italicRegex.ReplaceAllStringFunc(line, func(match string) string {
		content := italicRegex.FindStringSubmatch(match)[1]
		// ✅ ESCAPAR CONTEÚDO DO ITALIC
		return fmt.Sprintf("<i>%s</i>", html.EscapeString(content))
	})

	// ✅ ESCAPAR QUALQUER TEXTO RESTANTE QUE NÃO ESTEJA DENTRO DE TAGS HTML
	line = escapeRemainingText(line)

	return line
}

// ✅ NOVA FUNÇÃO: Escapar texto que não está dentro de tags HTML
func escapeRemainingText(text string) string {
	// Regex para encontrar texto fora de tags HTML
	htmlTagRegex := regexp.MustCompile(`(<[^>]+>)`)

	// Dividir o texto em partes: tags HTML e texto normal
	parts := htmlTagRegex.Split(text, -1)
	tags := htmlTagRegex.FindAllString(text, -1)

	var result strings.Builder

	for i, part := range parts {
		// Escapar apenas as partes que não são tags HTML
		if part != "" {
			result.WriteString(html.EscapeString(part))
		}

		// Adicionar a tag HTML (se existir) sem escapar
		if i < len(tags) {
			result.WriteString(tags[i])
		}
	}

	return result.String()
}

// ✅ FUNÇÃO PRINCIPAL: Converter Markdown para HTML
func convertMarkdownToHTML(text string) string {
	if text == "" {
		return text
	}

	// Aplicar conversões na ordem correta
	result := text

	// 1. Blockquotes (deve ser primeiro para processar linhas inteiras)
	result = blockquoteRegex.ReplaceAllString(result, "<blockquote>$1</blockquote>")

	// 2. Links
	result = linkRegex.ReplaceAllString(result, `<a href="$2">$1</a>`)

	// 3. Formatação de texto
	result = boldRegex.ReplaceAllString(result, "<b>$1</b>")
	result = italicRegex.ReplaceAllString(result, "<i>$1</i>")
	result = underlineRegex.ReplaceAllString(result, "<u>$1</u>")
	result = strikethroughRegex.ReplaceAllString(result, "<s>$1</s>")
	result = spoilerRegex.ReplaceAllString(result, `<span class="tg-spoiler">$1</span>`)
	result = codeRegex.ReplaceAllString(result, "<code>$1</code>")

	return result
}

// ✅ FUNÇÃO MELHORADA: Detectar formato e converter para HTML
func detectParseMode(text string) string {
	if text == "" {
		return text
	}

	// Detectar se é Markdown
	if isMarkdown(text) {
		converted := convertMarkdownToHTML(text)
		return converted
	}

	if isHTML(text) {

		return text
	}

	return text
}

// ✅ FUNÇÃO: Detectar se texto é Markdown
func isMarkdown(text string) bool {
	if text == "" {
		return false
	}

	markdownPatterns := []string{
		`\*\*[^*]+\*\*`, // **bold**
		`\*[^*]+\*`,     // *italic*
		`__[^_]+__`,     // __underline__
		`~~[^~]+~~`,     // ~~strikethrough~~
		"`[^`]+`",       // `code`
		`\[.+\]\(.+\)`,  // [link](url)
		`\|\|[^|]+\|\|`, // ||spoiler||
		`(?m)^>\s*.+$`,  // > blockquote
	}

	for _, pattern := range markdownPatterns {
		if matched, _ := regexp.MatchString(pattern, text); matched {
			return true
		}
	}

	return false
}

// ✅ FUNÇÃO: Detectar se texto é HTML
func isHTML(text string) bool {
	if text == "" {
		return false
	}

	htmlPatterns := []string{
		`<b>.*</b>`,
		`<i>.*</i>`,
		`<u>.*</u>`,
		`<s>.*</s>`,
		`<code>.*</code>`,
		`<pre>.*</pre>`,
		`<a href=.*>.*</a>`,
		`<blockquote>.*</blockquote>`,
		`<span class="tg-spoiler">.*</span>`,
	}

	for _, pattern := range htmlPatterns {
		if matched, _ := regexp.MatchString(pattern, text); matched {
			return true
		}
	}

	return false
}

func removeHashtag(text, hashtag string) string {
	if text == "" || hashtag == "" {
		return text
	}
	var re *regexp.Regexp
	if value, ok := removeHashRegexCache.Load(hashtag); ok {
		re = value.(*regexp.Regexp)
	} else {
		re = regexp.MustCompile(`#` + regexp.QuoteMeta(hashtag) + `\s*`)
		removeHashRegexCache.Store(hashtag, re)
	}
	return strings.TrimSpace(re.ReplaceAllString(text, ""))
}

// ✅ FUNÇÃO ATUALIZADA: processMessageWithHashtag com conversão para HTML
func (mp *MessageProcessor) processMessageWithHashtag(text string, channel *dbmodels.Channel) (string, *dbmodels.CustomCaption) {
	hashtag := extractHashtag(text)

	if hashtag == "" {
		defaultCaption := ""
		if channel.DefaultCaption != nil {
			// ✅ CONVERTER CAPTION PADRÃO PARA HTML
			defaultCaption = detectParseMode(channel.DefaultCaption.Caption)
		}
		return fmt.Sprintf("%s\n\n%s", text, defaultCaption), nil
	}

	customCaption := findCustomCaption(channel, hashtag)
	if customCaption == nil {
		defaultCaption := ""
		if channel.DefaultCaption != nil {
			// ✅ CONVERTER CAPTION PADRÃO PARA HTML
			defaultCaption = detectParseMode(channel.DefaultCaption.Caption)
		}
		return fmt.Sprintf("%s\n\n%s", text, defaultCaption), nil
	}

	cleanText := removeHashtag(text, hashtag)

	// ✅ CONVERTER CUSTOM CAPTION PARA HTML
	formattedCustomCaption := detectParseMode(customCaption.Caption)

	return fmt.Sprintf("%s\n\n%s", cleanText, formattedCustomCaption), customCaption
}
