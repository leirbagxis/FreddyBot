package parser

import (
	"fmt"
	"os"
	"regexp"
	"strings"
	"sync"

	"github.com/go-telegram/bot/models"
	"gopkg.in/yaml.v3"
)

type Button struct {
	Text                         string `yaml:"text"`
	CallbackData                 string `yaml:"callback_data,omitempty"`
	SwitchInlineQuery            string `yaml:"switch_inline_query,omitempty"`
	URL                          string `yaml:"url,omitempty"`
	WebApp                       string `yaml:"web_app,omitempty"`
	SwitchInlineQueryCurrentChat string `yaml:"switch_inline_query_current_chat,omitempty"`
}

type Message struct {
	Name        string
	Text        string
	Buttons     [][]Button
	VarKeys     []string
	HasVarsText bool
}

var (
	messagesMap      map[string]*Message
	loadOnce         sync.Once
	placeholderRegex = regexp.MustCompile(`\{(\w+)\}`)
)

func loadMessages() {
	file, err := os.ReadFile("./config/messages.yml")
	if err != nil {
		panic(fmt.Sprintf("Erro ao carregar messages.yml: %v", err))
	}

	var raw []struct {
		Name    string     `yaml:"name"`
		Text    string     `yaml:"text"`
		Buttons [][]Button `yaml:"buttons,omitempty"`
	}

	if err := yaml.Unmarshal(file, &raw); err != nil {
		panic(fmt.Sprintf("Erro ao parsear YAML: %v", err))
	}

	messagesMap = make(map[string]*Message, len(raw))

	for _, item := range raw {
		varKeys := detectPlaceholders(item.Text)
		messagesMap[item.Name] = &Message{
			Name:        item.Name,
			Text:        item.Text,
			Buttons:     item.Buttons,
			VarKeys:     varKeys,
			HasVarsText: len(varKeys) > 0,
		}
	}

	fmt.Println("✅ Messages carregadas:", len(messagesMap))
}

func detectPlaceholders(text string) []string {
	matches := placeholderRegex.FindAllStringSubmatch(text, -1)
	var keys []string
	seen := make(map[string]struct{}, len(matches))

	for _, match := range matches {
		key := match[1]
		if _, exists := seen[key]; !exists {
			keys = append(keys, key)
			seen[key] = struct{}{}
		}
	}
	return keys
}

func ParseText(text string, vars map[string]string, keys []string) string {
	for _, key := range keys {
		val, ok := vars[key]
		if !ok {
			val = fmt.Sprintf("{%s}", key)
		}
		text = strings.ReplaceAll(text, "{"+key+"}", val)
	}
	return text
}

func parseButtons(buttons [][]Button, vars map[string]string) [][]Button {
	if len(buttons) == 0 || len(vars) == 0 {
		return buttons
	}

	parsed := make([][]Button, len(buttons))
	for i, row := range buttons {
		newRow := make([]Button, len(row))
		for j, btn := range row {
			newRow[j] = Button{
				Text:                         ParseText(btn.Text, vars, detectPlaceholders(btn.Text)),
				CallbackData:                 ParseText(btn.CallbackData, vars, detectPlaceholders(btn.CallbackData)),
				URL:                          ParseText(btn.URL, vars, detectPlaceholders(btn.URL)),
				WebApp:                       ParseText(btn.WebApp, vars, detectPlaceholders(btn.WebApp)),
				SwitchInlineQuery:            ParseText(btn.SwitchInlineQuery, vars, detectPlaceholders(btn.SwitchInlineQuery)),
				SwitchInlineQueryCurrentChat: ParseText(btn.SwitchInlineQueryCurrentChat, vars, detectPlaceholders(btn.SwitchInlineQueryCurrentChat)),
			}
		}
		parsed[i] = newRow
	}
	return parsed
}

func BuildInlineKeyboard(buttons [][]Button) *models.InlineKeyboardMarkup {
	if len(buttons) == 0 {
		return nil
	}

	inlineKeyboard := make([][]models.InlineKeyboardButton, len(buttons))
	for i, row := range buttons {
		btnRow := make([]models.InlineKeyboardButton, len(row))
		for j, btn := range row {
			btnRow[j] = models.InlineKeyboardButton{
				Text:                         btn.Text,
				CallbackData:                 btn.CallbackData,
				URL:                          btn.URL,
				SwitchInlineQuery:            btn.SwitchInlineQuery,
				SwitchInlineQueryCurrentChat: btn.SwitchInlineQueryCurrentChat,
				WebApp: &models.WebAppInfo{
					URL: btn.WebApp,
				},
			}
		}
		inlineKeyboard[i] = btnRow
	}

	return &models.InlineKeyboardMarkup{InlineKeyboard: inlineKeyboard}
}

func GetMessage(name string, vars map[string]string) (string, *models.InlineKeyboardMarkup) {
	loadOnce.Do(loadMessages)

	msg, ok := messagesMap[name]
	if !ok {
		return fmt.Sprintf("Mensagem '%s' não encontrada!", name), nil
	}

	text := msg.Text
	if msg.HasVarsText && len(vars) > 0 {
		text = ParseText(text, vars, msg.VarKeys)
	}

	buttons := parseButtons(msg.Buttons, vars)
	keyboard := BuildInlineKeyboard(buttons)

	return text, keyboard
}
