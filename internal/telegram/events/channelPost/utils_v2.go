package channelpost

import (
	"regexp"
	"strconv"
	"strings"

	"github.com/go-telegram/bot/models"
	dbmodels "github.com/leirbagxis/FreddyBot/internal/database/models"
)

func GetMessageType(post *models.Message) MessageType {
	switch {
	case post.Text != "":
		return MessageTypeText
	case post.Audio != nil:
		return MessageTypeAudio
	case post.Sticker != nil:
		return MessageTypeSticker
	case post.Photo != nil:
		return MessageTypePhoto
	case post.Video != nil:
		return MessageTypeVideo
	case post.Animation != nil:
		return MessageTypeAnimation
	case post.Document != nil:
		return MessageTypeDocument
	default:
		return ""
	}
}

var (
	retryAfterRegex = regexp.MustCompile(`retry after (\d+)`)
	hashtagRegex    = regexp.MustCompile(`#(\w+)`)
)

func extractRetryAfter(errorMsg string) int {
	matches := retryAfterRegex.FindStringSubmatch(errorMsg)
	if len(matches) > 1 {
		if retryAfter, err := strconv.Atoi(matches[1]); err == nil {
			return retryAfter
		}
	}
	return 0
}

func extractHashtag(text string) string {
	m := hashtagRegex.FindStringSubmatch(text)
	if len(m) > 1 {
		return strings.ToLower(m[1])
	}
	return ""
}

func findCustomCaption(channel *dbmodels.Channel, code string) *dbmodels.CustomCaption {
	if channel == nil || len(channel.CustomCaptions) == 0 {
		return nil
	}
	code = strings.ToLower(strings.TrimSpace(code))
	for i := range channel.CustomCaptions {
		if strings.ToLower(channel.CustomCaptions[i].Code) == code {
			return &channel.CustomCaptions[i]
		}
	}
	return nil
}

func convertMessageEntitiesToInterface(ents []models.MessageEntity) []interface{} {
	out := make([]interface{}, 0, len(ents))
	for _, e := range ents {
		out = append(out, e)
	}
	return out
}

func convertInterfaceToEntities(anys []interface{}) []models.MessageEntity {
	out := make([]models.MessageEntity, 0, len(anys))
	for _, v := range anys {
		if e, ok := v.(models.MessageEntity); ok {
			out = append(out, e)
		}
	}
	return out
}

// composeMessage combina o conteúdo original com uma legenda do banco.
// order: "append" -> original + sep + db; "prepend" -> db + sep + original
func composeMessage(original, fromDB, sep, order string) string {
	o := strings.TrimSpace(original)
	d := strings.TrimSpace(fromDB)
	if o == "" && d == "" {
		return ""
	}
	if o == "" {
		return d
	}
	if d == "" {
		return o
	}
	if sep == "" {
		sep = "\n\n"
	}
	if order == "prepend" {
		return d + sep + o
	}
	return o + sep + d
}
