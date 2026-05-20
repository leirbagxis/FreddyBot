package utils

import (
	"crypto/rand"
	"crypto/rsa"
	"regexp"
	"strings"
)

func GenerateRSAKey() (*rsa.PrivateKey, error) {
	return rsa.GenerateKey(rand.Reader, 2048)
}

var (
	htmlTagRegex = regexp.MustCompile(`<[^>]*>`)
)

func RemoveHTMLTags(input string) string {
	return htmlTagRegex.ReplaceAllString(input, "")
}

func NormalizePort(p string) string {
	if p == "" {
		return ":7000"
	}
	if !strings.HasPrefix(p, ":") {
		return ":" + p
	}
	return p
}

func MinInt(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func MarkdownToTelegramHTML(text string) string {
	// Bloco de código ``` \n texto \n ``` (processado primeiro para evitar conflitos)
	blockRegex := regexp.MustCompile("```([\\s\\S]*?)```")
	text = blockRegex.ReplaceAllString(text, "<pre><code>$1</code></pre>")

	// Código monoespaçado `texto`
	codeRegex := regexp.MustCompile("`([^`]+)`")
	text = codeRegex.ReplaceAllString(text, "<code>$1</code>")

	// Negrito **texto**
	boldRegex := regexp.MustCompile(`\*\*(.*?)\*\*`)
	text = boldRegex.ReplaceAllString(text, "<b>$1</b>")

	// Itálico __texto__
	italicRegex := regexp.MustCompile(`__(.*?)__`)
	text = italicRegex.ReplaceAllString(text, "<i>$1</i>")

	// Tachado ~~texto~~
	strikeRegex := regexp.MustCompile(`~~(.*?)~~`)
	text = strikeRegex.ReplaceAllString(text, "<s>$1</s>")

	// Spoiler ||texto||
	spoilerRegex := regexp.MustCompile(`\|\|(.*?)\|\|`)
	text = spoilerRegex.ReplaceAllString(text, "<tg-spoiler>$1</tg-spoiler>")

	// Sublinhado <u>texto</u> já é HTML, mantemos

	// Links [texto](url)
	linkRegex := regexp.MustCompile(`\[([^\]]+)\]\(([^)]+)\)`)
	text = linkRegex.ReplaceAllString(text, `<a href="$2">$1</a>`)

	return text
}
