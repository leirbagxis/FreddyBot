package channelpost

import (
	"fmt"
	"html"
	"log"
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
func processEntitiesOnly(text string, entities []models.MessageEntity) string {
	if len(entities) == 0 {
		return html.EscapeString(text)
	}

	// Converter string para UTF-16 para cálculos corretos de offset
	textRunes := []rune(text)
	textUTF16 := utf16.Encode(textRunes)

	// Ordenar entidades por offset
	sort.Slice(entities, func(i, j int) bool {
		return entities[i].Offset < entities[j].Offset
	})

	var result strings.Builder
	lastOffset := 0

	for _, entity := range entities {
		// Adicionar texto antes da entidade (escapado)
		if entity.Offset > lastOffset {
			beforeText := string(utf16.Decode(textUTF16[lastOffset:entity.Offset]))
			result.WriteString(html.EscapeString(beforeText))
		}

		// Extrair o texto da entidade
		entityEnd := entity.Offset + entity.Length
		if entityEnd > len(textUTF16) {
			entityEnd = len(textUTF16)
		}

		entityText := string(utf16.Decode(textUTF16[entity.Offset:entityEnd]))

		// ✅ CORRIGIDO: GERAR HTML AO INVÉS DE MARKDOWN
		switch entity.Type {
		case "bold":
			result.WriteString("<b>")
			result.WriteString(html.EscapeString(entityText))
			result.WriteString("</b>")
		case "italic":
			result.WriteString("<i>")
			result.WriteString(html.EscapeString(entityText))
			result.WriteString("</i>")
		case "underline":
			result.WriteString("<u>")
			result.WriteString(html.EscapeString(entityText))
			result.WriteString("</u>")
		case "strikethrough":
			result.WriteString("<s>")
			result.WriteString(html.EscapeString(entityText))
			result.WriteString("</s>")
		case "spoiler":
			result.WriteString("<tg-spoiler>")
			result.WriteString(html.EscapeString(entityText))
			result.WriteString("</tg-spoiler>")
		case "code":
			result.WriteString("<code>")
			result.WriteString(html.EscapeString(entityText))
			result.WriteString("</code>")
		case "pre":
			if entity.Language != "" {
				result.WriteString("<pre><code class=\"")
				result.WriteString(html.EscapeString(entity.Language))
				result.WriteString("\">")
			} else {
				result.WriteString("<pre><code>")
			}
			result.WriteString(html.EscapeString(entityText))
			result.WriteString("</code></pre>")
		case "blockquote":
			result.WriteString("<blockquote>")
			result.WriteString(html.EscapeString(entityText))
			result.WriteString("</blockquote>")
		case "expandable_blockquote":
			result.WriteString("<blockquote expandable>")
			result.WriteString(html.EscapeString(entityText))
			result.WriteString("</blockquote>")
		case "text_link":
			result.WriteString("<a href=\"")
			result.WriteString(html.EscapeString(entity.URL))
			result.WriteString("\">")
			result.WriteString(html.EscapeString(entityText))
			result.WriteString("</a>")
		case "url":
			result.WriteString("<a href=\"")
			result.WriteString(html.EscapeString(entityText))
			result.WriteString("\">")
			result.WriteString(html.EscapeString(entityText))
			result.WriteString("</a>")
		case "email":
			result.WriteString("<a href=\"mailto:")
			result.WriteString(html.EscapeString(entityText))
			result.WriteString("\">")
			result.WriteString(html.EscapeString(entityText))
			result.WriteString("</a>")
		case "phone_number":
			result.WriteString("<a href=\"tel:")
			result.WriteString(html.EscapeString(entityText))
			result.WriteString("\">")
			result.WriteString(html.EscapeString(entityText))
			result.WriteString("</a>")
		case "mention":
			result.WriteString("<a href=\"https://t.me/")
			result.WriteString(html.EscapeString(strings.TrimPrefix(entityText, "@")))
			result.WriteString("\">")
			result.WriteString(html.EscapeString(entityText))
			result.WriteString("</a>")
		case "hashtag":
			result.WriteString("<a href=\"https://t.me/hashtag/")
			result.WriteString(html.EscapeString(strings.TrimPrefix(entityText, "#")))
			result.WriteString("\">")
			result.WriteString(html.EscapeString(entityText))
			result.WriteString("</a>")
		case "cashtag":
			result.WriteString("<a href=\"https://t.me/cashtag/")
			result.WriteString(html.EscapeString(strings.TrimPrefix(entityText, "$")))
			result.WriteString("\">")
			result.WriteString(html.EscapeString(entityText))
			result.WriteString("</a>")
		case "bot_command":
			result.WriteString("<code>")
			result.WriteString(html.EscapeString(entityText))
			result.WriteString("</code>")
		case "custom_emoji":
			if entity.CustomEmojiID != "" {
				result.WriteString("<tg-emoji emoji-id=\"")
				result.WriteString(html.EscapeString(entity.CustomEmojiID))
				result.WriteString("\">")
				result.WriteString(html.EscapeString(entityText))
				result.WriteString("</tg-emoji>")
			} else {
				result.WriteString(html.EscapeString(entityText))
			}
		default:
			result.WriteString(html.EscapeString(entityText))
		}

		lastOffset = entityEnd
	}

	// Adicionar texto restante (escapado)
	if lastOffset < len(textUTF16) {
		remainingText := string(utf16.Decode(textUTF16[lastOffset:]))
		result.WriteString(html.EscapeString(remainingText))
	}

	return result.String()
}

// ✅ CORRIGIDO: detectParseMode para markdown manual COM ESCAPE ADEQUADO
func detectParseModea(text string) string {
	if text == "" {
		return ""
	}

	lines := strings.Split(text, "\n")
	var result []string
	inCodeBlock := false
	var codeLanguage string
	var codeContent []string
	var blockQuoteLines []string
	var expandableBlockQuoteLines []string

	i := 0
	for i < len(lines) {
		line := lines[i]

		// Code blocks têm prioridade máxima
		if strings.HasPrefix(line, "```") && !inCodeBlock {
			inCodeBlock = true
			if len(line) > 3 {
				codeLanguage = strings.TrimSpace(line[3:])
				codeContent = []string{}
			} else {
				codeLanguage = ""
				codeContent = []string{}
			}
			i++
			continue
		} else if strings.HasPrefix(line, "```") && inCodeBlock {
			inCodeBlock = false
			// ✅ ESCAPAR CONTEÚDO DO CODE BLOCK
			codeText := html.EscapeString(strings.Join(codeContent, "\n"))
			if codeLanguage != "" {
				result = append(result, fmt.Sprintf(`<pre><code class="%s">%s</code></pre>`, html.EscapeString(codeLanguage), codeText))
			} else {
				result = append(result, fmt.Sprintf(`<pre><code>%s</code></pre>`, codeText))
			}
			i++
			continue
		} else if inCodeBlock {
			// ✅ NÃO ESCAPAR AQUI - será escapado quando o bloco terminar
			codeContent = append(codeContent, line)
			i++
			continue
		}

		// Expandable blockquotes
		if strings.HasPrefix(line, "**>") && strings.HasSuffix(line, "||") {
			expandableBlockQuoteLines = []string{processInlineFormatting(line[3 : len(line)-2])}
			i++
			for i < len(lines) && strings.HasPrefix(lines[i], ">") {
				expandableBlockQuoteLines = append(expandableBlockQuoteLines, processInlineFormatting(strings.TrimSpace(lines[i][1:])))
				i++
			}
			result = append(result, fmt.Sprintf(`<blockquote expandable>%s</blockquote>`, strings.Join(expandableBlockQuoteLines, "\n")))
			continue
		}

		// Regular blockquotes
		if strings.HasPrefix(line, ">") {
			blockQuoteLines = []string{processInlineFormatting(strings.TrimSpace(line[1:]))}
			i++
			for i < len(lines) && strings.HasPrefix(lines[i], ">") {
				blockQuoteLines = append(blockQuoteLines, processInlineFormatting(strings.TrimSpace(lines[i][1:])))
				i++
			}
			result = append(result, fmt.Sprintf(`<blockquote>%s</blockquote>`, strings.Join(blockQuoteLines, "\n")))
			continue
		}

		// Process inline formatting
		processedLine := processInlineFormatting(line)
		result = append(result, processedLine)
		i++
	}

	return strings.Join(result, "\n")
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
		log.Printf("📝 Convertido Markdown -> HTML: %q -> %q", text, converted)
		return converted
	}

	// Detectar se já é HTML
	if isHTML(text) {
		log.Printf("📝 Texto já é HTML: %q", text)
		return text
	}

	// Texto plano
	log.Printf("📝 Texto plano: %q", text)
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
