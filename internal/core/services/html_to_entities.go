package services

import (
	"strconv"
	"strings"
	"unicode/utf8"

	"github.com/gotd/td/tg"
)

func HTMLToEntities(html string) (plain string, entities []tg.MessageEntityClass) {
	if html == "" {
		return "", nil
	}

	var buf strings.Builder
	var runes []rune
	var stack []struct {
		name   string
		attrs  map[string]string
		runeStart int
	}

	i := 0
	for i < len(html) {
		if html[i] == '<' {
			end := strings.IndexByte(html[i:], '>')
			if end == -1 {
				r, size := utf8.DecodeRuneInString(html[i:])
				buf.WriteRune(r)
				runes = append(runes, r)
				i += size
				continue
			}
			tagStr := html[i+1 : i+end]
			if strings.HasPrefix(tagStr, "/") {
				tagName := tagStr[1:]
			if len(stack) > 0 && stack[len(stack)-1].name == tagName {
				opening := stack[len(stack)-1]
				stack = stack[:len(stack)-1]
				plainText := runes[opening.runeStart:]
				utf16Len := utf16LenRunes(plainText)
				if utf16Len > 0 {
					utf16Start := utf16LenRunes(runes[:opening.runeStart])
					entities = appendEntity(entities, opening.name, opening.attrs, utf16Start, utf16Len)
				}
				}
			} else {
				tagName, attrs := parseTag(tagStr)
				if !isVoidTag(tagName) {
					stack = append(stack, struct {
						name   string
						attrs  map[string]string
						runeStart int
					}{tagName, attrs, len(runes)})
				}
			}
			i += end + 1
		} else {
			r, size := utf8.DecodeRuneInString(html[i:])
			buf.WriteRune(r)
			runes = append(runes, r)
			i += size
		}
	}

	for len(stack) > 0 {
		opening := stack[0]
		stack = stack[1:]
		plainText := runes[opening.runeStart:]
		utf16Len := utf16LenRunes(plainText)
		if utf16Len > 0 {
			utf16Start := utf16LenRunes(runes[:opening.runeStart])
			entities = appendEntity(entities, opening.name, opening.attrs, utf16Start, utf16Len)
		}
	}

	return buf.String(), entities
}

func utf16LenRunes(runes []rune) int {
	count := 0
	for _, r := range runes {
		if r >= 0x10000 {
			count += 2
		} else {
			count++
		}
	}
	return count
}

func parseTag(tagStr string) (string, map[string]string) {
	parts := strings.Fields(tagStr)
	if len(parts) == 0 {
		return "", nil
	}
	name := strings.ToLower(parts[0])
	attrs := make(map[string]string)
	for _, p := range parts[1:] {
		eq := strings.IndexByte(p, '=')
		if eq > 0 {
			key := p[:eq]
			val := p[eq+1:]
			val = strings.Trim(val, `"'`)
			attrs[key] = val
		}
	}
	return name, attrs
}

func appendEntity(entities []tg.MessageEntityClass, tagName string, attrs map[string]string, offset, length int) []tg.MessageEntityClass {
	if length <= 0 {
		return entities
	}
	switch tagName {
	case "b", "strong":
		return append(entities, &tg.MessageEntityBold{Offset: offset, Length: length})
	case "i", "em":
		return append(entities, &tg.MessageEntityItalic{Offset: offset, Length: length})
	case "u":
		return append(entities, &tg.MessageEntityUnderline{Offset: offset, Length: length})
	case "s", "del", "strike":
		return append(entities, &tg.MessageEntityStrike{Offset: offset, Length: length})
	case "tg-spoiler":
		return append(entities, &tg.MessageEntitySpoiler{Offset: offset, Length: length})
	case "code":
		return append(entities, &tg.MessageEntityCode{Offset: offset, Length: length})
	case "pre":
		lang := strings.TrimPrefix(attrs["class"], "language-")
		return append(entities, &tg.MessageEntityPre{Offset: offset, Length: length, Language: lang})
	case "blockquote":
		return append(entities, &tg.MessageEntityBlockquote{Offset: offset, Length: length})
	case "a":
		url := attrs["href"]
		if url == "" {
			return entities
		}
		return append(entities, &tg.MessageEntityTextURL{Offset: offset, Length: length, URL: url})
	case "tg-emoji":
		emojiIDStr := attrs["emoji-id"]
		if emojiIDStr != "" {
			docID, err := strconv.ParseInt(emojiIDStr, 10, 64)
			if err == nil && docID > 0 {
				return append(entities, &tg.MessageEntityCustomEmoji{Offset: offset, Length: length, DocumentID: docID})
			}
		}
		return entities
	}
	return entities
}

func isVoidTag(name string) bool {
	switch name {
	case "br", "hr", "img", "input", "meta", "link", "area", "base", "col", "embed", "source", "track", "wbr":
		return true
	}
	return false
}

func StripHTML(s string) string {
	var buf strings.Builder
	i := 0
	for i < len(s) {
		if s[i] == '<' {
			end := strings.IndexByte(s[i:], '>')
			if end != -1 {
				i += end + 1
				continue
			}
		}
		if s[i] == '&' {
			semi := strings.IndexByte(s[i:], ';')
			if semi != -1 && semi <= 6 {
				ent := s[i : i+semi+1]
				switch ent {
				case "&amp;":
					buf.WriteByte('&')
				case "&lt;":
					buf.WriteByte('<')
				case "&gt;":
					buf.WriteByte('>')
				case "&quot;":
					buf.WriteByte('"')
				case "&#39;", "&#x27;":
					buf.WriteByte('\'')
				case "&nbsp;":
					buf.WriteByte(' ')
				default:
					i += semi + 1
					continue
				}
				i += semi + 1
				continue
			}
		}
		buf.WriteByte(s[i])
		i++
	}
	return buf.String()
}

func KeyboardRow(buttons ...tg.KeyboardButtonClass) tg.KeyboardButtonRow {
	return tg.KeyboardButtonRow{Buttons: buttons}
}

func URLButton(text, url string) tg.KeyboardButtonClass {
	return &tg.KeyboardButtonURL{Text: text, URL: url}
}
