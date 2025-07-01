package channelpost

import (
	"fmt"
	"regexp"
	"strings"

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
func processTextWithFormatting(text string, entities []models.MessageEntity) string {
	if text == "" {
		return ""
	}

	// Primeiro aplicar entidades (bold, italic, etc.)
	formattedText := applyEntities(text, entities)

	// Depois aplicar formatação markdown adicional
	finalText := detectParseMode(formattedText)

	return finalText
}
