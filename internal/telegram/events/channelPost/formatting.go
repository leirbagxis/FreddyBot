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

// ✅ PORTADO DO NODE.JS: detectParseMode
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

		// Code blocks
		if strings.HasPrefix(line, "```") && !inCodeBlock {
			inCodeBlock = true
			if len(line) > 3 {
				codeLanguage = line[3:]
				codeContent = []string{}
			} else {
				codeLanguage = ""
				codeContent = []string{}
			}
			i++
			continue
		} else if strings.HasPrefix(line, "```") && inCodeBlock {
			inCodeBlock = false
			if codeLanguage != "" {
				result = append(result, fmt.Sprintf(`<pre><code class="language-%s">%s</code></pre>`, codeLanguage, strings.Join(codeContent, "\n")))
			} else {
				result = append(result, fmt.Sprintf(`<pre>%s</pre>`, strings.Join(codeContent, "\n")))
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
			expandableBlockQuoteLines = []string{line[3 : len(line)-2]}
			i++
			for i < len(lines) && strings.HasPrefix(lines[i], ">") {
				expandableBlockQuoteLines = append(expandableBlockQuoteLines, strings.TrimSpace(lines[i][1:]))
				i++
			}
			result = append(result, fmt.Sprintf(`<blockquote expandable>%s</blockquote>`, strings.Join(expandableBlockQuoteLines, "\n")))
			continue
		}

		// Regular blockquotes
		if strings.HasPrefix(line, ">") {
			blockQuoteLines = []string{strings.TrimSpace(line[1:])}
			i++
			for i < len(lines) && strings.HasPrefix(lines[i], ">") {
				blockQuoteLines = append(blockQuoteLines, strings.TrimSpace(lines[i][1:]))
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

// ✅ PORTADO DO NODE.JS: processInlineFormatting
func processInlineFormatting(line string) string {
	// Bold
	boldRegex := regexp.MustCompile(`\*\*(.*?)\*\*`)
	line = boldRegex.ReplaceAllString(line, "<b>$1</b>")

	// Italic
	italicRegex := regexp.MustCompile(`\*(.*?)\*`)
	line = italicRegex.ReplaceAllString(line, "<i>$1</i>")

	// Underline
	underlineRegex := regexp.MustCompile(`__(.*?)__`)
	line = underlineRegex.ReplaceAllString(line, "<u>$1</u>")

	// Strikethrough
	strikeRegex := regexp.MustCompile(`~~(.*?)~~`)
	line = strikeRegex.ReplaceAllString(line, "<s>$1</s>")

	// Spoiler
	spoilerRegex := regexp.MustCompile(`\|\|(.*?)\|\|`)
	line = spoilerRegex.ReplaceAllString(line, `<span class="tg-spoiler">$1</span>`)

	// Code
	codeRegex := regexp.MustCompile("`(.*?)`")
	line = codeRegex.ReplaceAllString(line, "<code>$1</code>")

	return line
}

// ✅ PORTADO DO NODE.JS: applyEntities
func applyEntities(text string, entities []models.MessageEntity) string {
	if len(entities) == 0 {
		return text
	}

	openTags := make(map[int][]string)
	closeTags := make(map[int][]string)

	tagMap := map[string]string{
		"bold":          "b",
		"blockquote":    "blockquote",
		"italic":        "i",
		"underline":     "u",
		"strikethrough": "s",
		"code":          "code",
		"spoiler":       "tg-spoiler",
	}

	for _, entity := range entities {
		offset := entity.Offset
		length := entity.Length
		entityType := entity.Type
		url := entity.URL

		if entityType == "text_link" {
			if openTags[offset] == nil {
				openTags[offset] = []string{}
			}
			if closeTags[offset+length] == nil {
				closeTags[offset+length] = []string{}
			}

			openTags[offset] = append(openTags[offset], fmt.Sprintf(`<a href='%s'>`, url))
			closeTags[offset+length] = append([]string{"</a>"}, closeTags[offset+length]...)
			continue
		}

		tag, exists := tagMap[string(entityType)]
		if !exists {
			continue
		}

		if openTags[offset] == nil {
			openTags[offset] = []string{}
		}
		if closeTags[offset+length] == nil {
			closeTags[offset+length] = []string{}
		}

		openTags[offset] = append(openTags[offset], fmt.Sprintf("<%s>", tag))
		closeTags[offset+length] = append([]string{fmt.Sprintf("</%s>", tag)}, closeTags[offset+length]...)
	}

	var result strings.Builder
	textRunes := []rune(text)

	for i := 0; i < len(textRunes); i++ {
		if tags, exists := openTags[i]; exists {
			result.WriteString(strings.Join(tags, ""))
		}
		result.WriteRune(textRunes[i])
		if tags, exists := closeTags[i+1]; exists {
			result.WriteString(strings.Join(tags, ""))
		}
	}

	return result.String()
}

// ✅ FUNÇÃO PARA PROCESSAR TEXTO COM FORMATAÇÃO COMPLETA
func aprocessTextWithFormatting(text string, entities []models.MessageEntity) string {
	if text == "" {
		return ""
	}

	// Primeiro aplicar entidades (bold, italic, etc.)
	formattedText := applyEntities(text, entities)

	// Depois aplicar formatação markdown adicional
	finalText := detectParseMode(formattedText)

	return finalText
}

// processTextWithFormatting converte MessageEntities do Telegram para HTML
func processTextWithFormatting(text string, entities []models.MessageEntity) string {
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

		// Aplicar formatação baseada no tipo de entidade
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
			result.WriteString("<span class=\"tg-spoiler\">")
			result.WriteString(html.EscapeString(entityText))
			result.WriteString("</span>")
		case "code":
			result.WriteString("<code>")
			result.WriteString(html.EscapeString(entityText))
			result.WriteString("</code>")
		case "pre":
			if entity.Language != "" {
				result.WriteString("<pre><code class=\"language-")
				result.WriteString(html.EscapeString(entity.Language))
				result.WriteString("\">")
				result.WriteString(html.EscapeString(entityText))
				result.WriteString("</code></pre>")
			} else {
				result.WriteString("<pre>")
				result.WriteString(html.EscapeString(entityText))
				result.WriteString("</pre>")
			}
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

// Função auxiliar para processar entidades sobrepostas (se necessário)
func mergeOverlappingEntities(entities []models.MessageEntity) []models.MessageEntity {
	if len(entities) <= 1 {
		return entities
	}

	// Ordenar por offset
	sort.Slice(entities, func(i, j int) bool {
		if entities[i].Offset == entities[j].Offset {
			return entities[i].Length > entities[j].Length // Entidades maiores primeiro
		}
		return entities[i].Offset < entities[j].Offset
	})

	var merged []models.MessageEntity
	for _, entity := range entities {
		// Verificar se a entidade atual se sobrepõe com alguma já processada
		overlaps := false
		for i := range merged {
			existingEnd := merged[i].Offset + merged[i].Length
			entityEnd := entity.Offset + entity.Length

			// Se há sobreposição, mesclar ou pular
			if entity.Offset < existingEnd && merged[i].Offset < entityEnd {
				overlaps = true
				break
			}
		}

		if !overlaps {
			merged = append(merged, entity)
		}
	}

	return merged
}
