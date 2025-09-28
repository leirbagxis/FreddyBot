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

// Regex e helpers para conversões (mantidos)

var (
	boldRegex          = regexp.MustCompile(`\*\*(.+?)\*\*`)
	italicRegex        = regexp.MustCompile(`\*(.+?)\*`)
	underlineRegex     = regexp.MustCompile(`__(.+?)__`)
	strikethroughRegex = regexp.MustCompile(`~~(.+?)~~`)
	codeRegex          = regexp.MustCompile("`([^`]+)`")
	linkRegex          = regexp.MustCompile(`\[(.+?)\]\((https?://[^\s)]+)\)`)
	spoilerRegex       = regexp.MustCompile(`\|\|(.+?)\|\|`)
	blockquoteRegex    = regexp.MustCompile(`(?m)^> (.+)$`)
)

// Processa texto com Entities se existirem; caso contrário, aplica markdown → HTML
func processTextWithFormatting(text string, entities []models.MessageEntity) string {
	if text == "" {
		return ""
	}
	// Se tem entities, processar apenas elas (não aplicar markdown)
	if len(entities) > 0 {
		return processEntitiesOnly(text, entities)
	}
	// Senão, aplicar heurística de markdown simples
	return detectParseMode(text)
}

// Converte somente entities para HTML escapando demais
func processEntitiesOnly(text string, entities []models.MessageEntity) string {
	if len(entities) == 0 {
		return html.EscapeString(text)
	}

	// Converter em UTF-16 para offsets corretos
	runes := []rune(text)
	utf16Data := utf16.Encode(runes)

	type span struct {
		start int
		end   int
		tag   string
		attr  string
	}
	var spans []span

	// Cria spans a partir das entities
	for _, e := range entities {
		start := int(e.Offset)
		end := start + int(e.Length)
		if start < 0 || end > len(utf16Data) || start >= end {
			continue
		}
		tag := ""
		attr := ""
		switch e.Type {
		case "bold":
			tag = "b"
		case "italic":
			tag = "i"
		case "underline":
			tag = "u"
		case "strikethrough":
			tag = "s"
		case "code":
			tag = "code"
		case "blockquote":
			tag = "blockquote"
		case "text_link":
			tag = "a"
			if e.URL != "" {
				attr = fmt.Sprintf(` href="%s"`, html.EscapeString(e.URL))
			}
		case "spoiler":
			tag = "span"
			attr = ` class="tg-spoiler"`
		default:
			// outros tipos ignorados
			continue
		}
		spans = append(spans, span{start: start, end: end, tag: tag, attr: attr})
	}

	// Ordenar por começo asc e fim desc para aninhamento correto
	sort.Slice(spans, func(i, j int) bool {
		if spans[i].start == spans[j].start {
			return spans[i].end > spans[j].end
		}
		return spans[i].start < spans[j].start
	})

	// Montagem
	var b strings.Builder
	cur := 0
	for _, s := range spans {
		if s.start > cur {
			b.WriteString(html.EscapeString(decodeUTF16(utf16Data[cur:s.start])))
		}
		open := "<" + s.tag + s.attr + ">"
		close := "</" + s.tag + ">"
		b.WriteString(open)
		b.WriteString(html.EscapeString(decodeUTF16(utf16Data[s.start:s.end])))
		b.WriteString(close)
		cur = s.end
	}
	if cur < len(utf16Data) {
		b.WriteString(html.EscapeString(decodeUTF16(utf16Data[cur:])))
	}
	return b.String()
}

func detectParseMode(text string) string {
	escaped := html.EscapeString(text)
	escaped = boldRegex.ReplaceAllString(escaped, "<b>$1</b>")
	escaped = italicRegex.ReplaceAllString(escaped, "<i>$1</i>")
	escaped = underlineRegex.ReplaceAllString(escaped, "<u>$1</u>")
	escaped = strikethroughRegex.ReplaceAllString(escaped, "<s>$1</s>")
	escaped = codeRegex.ReplaceAllString(escaped, "<code>$1</code>")
	escaped = spoilerRegex.ReplaceAllString(escaped, `<span class="tg-spoiler">$1</span>`)
	escaped = blockquoteRegex.ReplaceAllString(escaped, "<blockquote>$1</blockquote>")
	escaped = linkRegex.ReplaceAllStringFunc(escaped, func(s string) string {
		m := linkRegex.FindStringSubmatch(s)
		if len(m) != 3 {
			return s
		}
		return fmt.Sprintf(`<a href="%s">%s</a>`, html.EscapeString(m[2]), html.EscapeString(m[1]))
	})
	return escaped
}

func decodeUTF16(u16 []uint16) string {
	runes := utf16.Decode(u16)
	return string(runes)
}

// Placeholders para tipos externos referenciados
// Deixe como estão caso já existam em seu projeto real
func _useDbModels(_ *dbmodels.Channel) {}
