package channelpost

import (
	"fmt"
	"html"
	"regexp"
	"sort"
	"strings"
	"unicode/utf16"

	"github.com/go-telegram/bot/models"
)

// ✅ CORRIGIDO: processTextWithFormatting - priorizar entities do Telegram
func processTextWithFormatting(text string, entities []models.MessageEntity) string {
	if text == "" {
		return ""
	}

	// ✅ SE TEM ENTITIES, PROCESSAR APENAS ELAS (não aplicar markdown)
	if len(entities) > 0 {
		return processEntitiesOnly(text, entities)
	}

	// ✅ SE NÃO TEM ENTITIES, APLICAR FORMATAÇÃO MARKDOWN
	return detectParseMode(html.EscapeString(text))
}

// ✅ NOVA FUNÇÃO: Processar apenas entities do Telegram
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

		// ✅ APLICAR FORMATAÇÃO BASEADA NO TIPO DE ENTIDADE
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
			// ✅ BLOCKQUOTE - como no seu exemplo
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
			// ✅ CUSTOM EMOJI
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
			// Para tipos desconhecidos, apenas escapar o texto
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

// ✅ MANTIDO: detectParseMode para quando NÃO há entities
func detectParseMode(text string) string {
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
			codeText := strings.Join(codeContent, "\n")
			if codeLanguage != "" {
				result = append(result, fmt.Sprintf(`<pre><code class="%s">%s</code></pre>`, html.EscapeString(codeLanguage), codeText))
			} else {
				result = append(result, fmt.Sprintf(`<pre><code>%s</code></pre>`, codeText))
			}
			i++
			continue
		} else if inCodeBlock {
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

// ✅ MANTIDO: processInlineFormatting para markdown manual
func processInlineFormatting(line string) string {
	if line == "" {
		return ""
	}

	// 1. Code inline (maior prioridade)
	codeRegex := regexp.MustCompile("`([^`]+)`")
	line = codeRegex.ReplaceAllStringFunc(line, func(match string) string {
		content := match[1 : len(match)-1]
		return fmt.Sprintf("<code>%s</code>", content)
	})

	// 2. Spoiler
	spoilerRegex := regexp.MustCompile(`\|\|([^|]+)\|\|`)
	line = spoilerRegex.ReplaceAllString(line, `<tg-spoiler>$1</tg-spoiler>`)

	// 3. Bold
	boldRegex := regexp.MustCompile(`\*\*([^*]+)\*\*`)
	line = boldRegex.ReplaceAllString(line, "<b>$1</b>")

	// 4. Underline
	underlineRegex := regexp.MustCompile(`__([^_]+)__`)
	line = underlineRegex.ReplaceAllString(line, "<u>$1</u>")

	// 5. Strikethrough
	strikeRegex := regexp.MustCompile(`~~([^~]+)~~`)
	line = strikeRegex.ReplaceAllString(line, "<s>$1</s>")

	// 6. Italic (deve vir por último)
	italicRegex := regexp.MustCompile(`\*([^*]+)\*`)
	line = italicRegex.ReplaceAllString(line, "<i>$1</i>")

	return line
}

// ✅ FUNÇÃO AUXILIAR: Debug para verificar entities
func debugEntities(text string, entities []models.MessageEntity) {
	fmt.Printf("=== DEBUG ENTITIES ===\n")
	fmt.Printf("Text: %q\n", text)
	fmt.Printf("Entities count: %d\n", len(entities))
	for i, entity := range entities {
		fmt.Printf("Entity %d: Type=%s, Offset=%d, Length=%d, URL=%s, Language=%s, CustomEmojiID=%s\n",
			i, entity.Type, entity.Offset, entity.Length, entity.URL, entity.Language, entity.CustomEmojiID)
	}
	formatted := processTextWithFormatting(text, entities)
	fmt.Printf("Formatted: %q\n", formatted)
	fmt.Printf("=== END DEBUG ===\n")
}
