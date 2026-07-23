package services

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/leirbagxis/FreddyBot/internal/cache"
	"github.com/leirbagxis/FreddyBot/internal/database/models"
	"github.com/leirbagxis/FreddyBot/internal/database/repositories"
	"github.com/leirbagxis/FreddyBot/pkg/logger"
	"github.com/mymmrac/telego"
)

type ScheduleOptions struct {
	ScheduleType  string
	ScheduleTime  string
	ScheduledAt   *time.Time
	ScheduleDays  []int
	RepeatUntil   *time.Time
	QueueGroupID  string
	QueuePosition int
	LoopQueue     bool
}

type SchedulerService struct {
	repo        *repositories.ScheduledPostRepository
	cacheService *cache.Service
	bot         *telego.Bot
}

func NewSchedulerService(
	repo *repositories.ScheduledPostRepository,
	cacheService *cache.Service,
	bot *telego.Bot,
) *SchedulerService {
	return &SchedulerService{
		repo:         repo,
		cacheService: cacheService,
		bot:          bot,
	}
}

func (s *SchedulerService) Start(ctx context.Context) {
	logger.Info("SCHEDULER", "Scheduler iniciado")
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	s.processDuePosts()

	for {
		select {
		case <-ctx.Done():
			logger.Info("SCHEDULER", "Scheduler encerrado")
			return
		case <-ticker.C:
			s.processDuePosts()
		}
	}
}

func (s *SchedulerService) processDuePosts() {
	ctx := context.Background()
	posts, err := s.repo.GetDuePosts(ctx, time.Now())
	if err != nil {
		logger.Error("SCHEDULER", "Erro ao buscar posts pendentes: %v", err)
		return
	}

	for _, post := range posts {
		s.sendScheduledPost(ctx, &post)
		time.Sleep(1 * time.Second)
	}
}

func (s *SchedulerService) sendScheduledPost(ctx context.Context, post *models.ScheduledPost) {
	logger.Info("SCHEDULER", "Enviando post agendado %s no canal %d", post.ID, post.ChannelID)

	var state cache.PostBuilderState
	if err := json.Unmarshal([]byte(post.PostData), &state); err != nil {
		logger.Error("SCHEDULER", "Erro ao desserializar PostData: %v", err)
		s.repo.UpdateError(ctx, post.ID, err.Error(), post.RetryCount+1)
		return
	}

	err := s.buildAndSend(ctx, post.ChannelID, &state)
	if err != nil {
		logger.Error("SCHEDULER", "Erro ao enviar post %s: %v", post.ID, err)
		s.repo.UpdateError(ctx, post.ID, err.Error(), post.RetryCount+1)
		return
	}

	sentAt := time.Now()
	s.repo.MarkSent(ctx, post.ID, sentAt)

	switch post.ScheduleType {
	case "once":
		logger.Info("SCHEDULER", "Post one-shot %s enviado com sucesso", post.ID)
	case "daily":
		nextRun := s.calculateNextDaily(post)
		if nextRun != nil && (post.RepeatUntil == nil || nextRun.Before(*post.RepeatUntil)) {
			s.repo.UpdateStatus(ctx, post.ID, "pending")
			s.repo.UpdateNextRunAt(ctx, post.ID, *nextRun)
			logger.Info("SCHEDULER", "Post diário %s: próximo envio %s", post.ID, nextRun.Format("02/01 15:04"))
		} else {
			logger.Info("SCHEDULER", "Post diário %s finalizado (repeat_until atingido)", post.ID)
		}
	case "weekly":
		nextRun := s.calculateNextWeekly(post)
		if nextRun != nil && (post.RepeatUntil == nil || nextRun.Before(*post.RepeatUntil)) {
			s.repo.UpdateStatus(ctx, post.ID, "pending")
			s.repo.UpdateNextRunAt(ctx, post.ID, *nextRun)
			logger.Info("SCHEDULER", "Post semanal %s: próximo envio %s", post.ID, nextRun.Format("02/01 15:04"))
		} else {
			logger.Info("SCHEDULER", "Post semanal %s finalizado", post.ID)
		}
	case "queue":
		s.advanceQueue(ctx, post)
	}
}

func (s *SchedulerService) buildAndSend(ctx context.Context, chatID int64, state *cache.PostBuilderState) error {
	var sb strings.Builder
	if state.Title != "" {
		sb.WriteString(state.Title + "\n\n")
	}
	if state.Body != "" {
		sb.WriteString(state.Body + "\n\n")
	}
	if state.Footer != "" {
		sb.WriteString(state.Footer)
	}
	caption := sb.String()

	if strings.TrimSpace(caption) == "" && state.MediaType == "" {
		return fmt.Errorf("post vazio (sem mídia e sem texto)")
	}

	var kb telego.ReplyMarkup
	if len(state.Buttons) > 0 || state.Reactions != "" {
		ikb := &telego.InlineKeyboardMarkup{}
		for _, btn := range state.Buttons {
			ikb.InlineKeyboard = append(ikb.InlineKeyboard, []telego.InlineKeyboardButton{
				{Text: btn.Text, URL: btn.URL},
			})
		}
		if state.Reactions != "" {
			reactions := strings.Split(state.Reactions, ",")
			var reactionRow []telego.InlineKeyboardButton
			for _, r := range reactions {
				val := strings.TrimSpace(r)
				if val != "" {
					reactionRow = append(reactionRow, telego.InlineKeyboardButton{
						CallbackData: "vote:" + val,
						Text:         val,
					})
				}
			}
			if len(reactionRow) > 0 {
				ikb.InlineKeyboard = append(ikb.InlineKeyboard, reactionRow)
			}
		}
		kb = ikb
	}

	switch state.MediaType {
	case "photo":
		params := &telego.SendPhotoParams{
			ChatID:    telego.ChatID{ID: chatID},
			Photo:     telego.InputFile{FileID: state.MediaFileID},
			Caption:   caption,
			ParseMode: telego.ModeHTML,
		}
		if kb != nil {
			params.ReplyMarkup = kb
		}
		_, err := s.bot.SendPhoto(ctx, params)
		return err
	case "video":
		params := &telego.SendVideoParams{
			ChatID:    telego.ChatID{ID: chatID},
			Video:     telego.InputFile{FileID: state.MediaFileID},
			Caption:   caption,
			ParseMode: telego.ModeHTML,
		}
		if kb != nil {
			params.ReplyMarkup = kb
		}
		_, err := s.bot.SendVideo(ctx, params)
		return err
	case "animation":
		params := &telego.SendAnimationParams{
			ChatID:    telego.ChatID{ID: chatID},
			Animation: telego.InputFile{FileID: state.MediaFileID},
			Caption:   caption,
			ParseMode: telego.ModeHTML,
		}
		if kb != nil {
			params.ReplyMarkup = kb
		}
		_, err := s.bot.SendAnimation(ctx, params)
		return err
	case "audio":
		params := &telego.SendAudioParams{
			ChatID:    telego.ChatID{ID: chatID},
			Audio:     telego.InputFile{FileID: state.MediaFileID},
			Caption:   caption,
			ParseMode: telego.ModeHTML,
		}
		if kb != nil {
			params.ReplyMarkup = kb
		}
		_, err := s.bot.SendAudio(ctx, params)
		return err
	case "document":
		params := &telego.SendDocumentParams{
			ChatID:    telego.ChatID{ID: chatID},
			Document:  telego.InputFile{FileID: state.MediaFileID},
			Caption:   caption,
			ParseMode: telego.ModeHTML,
		}
		if kb != nil {
			params.ReplyMarkup = kb
		}
		_, err := s.bot.SendDocument(ctx, params)
		return err
	default:
		params := &telego.SendMessageParams{
			ChatID:    telego.ChatID{ID: chatID},
			Text:      caption,
			ParseMode: telego.ModeHTML,
		}
		if kb != nil {
			params.ReplyMarkup = kb
		}
		_, err := s.bot.SendMessage(ctx, params)
		return err
	}
}

func (s *SchedulerService) calculateNextDaily(post *models.ScheduledPost) *time.Time {
	now := time.Now().UTC()
	hour, min := parseHHMM(post.ScheduleTime)
	next := time.Date(now.Year(), now.Month(), now.Day()+1, hour, min, 0, 0, time.UTC)
	return &next
}

func (s *SchedulerService) calculateNextWeekly(post *models.ScheduledPost) *time.Time {
	now := time.Now().UTC()
	hour, min := parseHHMM(post.ScheduleTime)

	var days []int
	if post.ScheduleDays != "" {
		json.Unmarshal([]byte(post.ScheduleDays), &days)
	}
	if len(days) == 0 {
		days = []int{0, 1, 2, 3, 4, 5, 6}
	}

	for i := 1; i <= 7; i++ {
		nextDay := now.AddDate(0, 0, i)
		for _, d := range days {
			if int(nextDay.Weekday()) == d {
				result := time.Date(nextDay.Year(), nextDay.Month(), nextDay.Day(), hour, min, 0, 0, time.UTC)
				return &result
			}
		}
	}
	return nil
}

func (s *SchedulerService) advanceQueue(ctx context.Context, sentPost *models.ScheduledPost) {
	if sentPost.QueueGroupID == "" {
		return
	}

	posts, err := s.repo.GetQueueGroup(ctx, sentPost.QueueGroupID)
	if err != nil {
		logger.Error("SCHEDULER", "Erro ao buscar fila %s: %v", sentPost.QueueGroupID, err)
		return
	}

	var nextPost *models.ScheduledPost
	for i := range posts {
		if posts[i].QueuePosition > sentPost.QueuePosition && posts[i].Status == "pending" {
			nextPost = &posts[i]
			break
		}
	}

	if nextPost != nil {
		hour, min := parseHHMM(sentPost.ScheduleTime)
		nextRun := time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day()+1, hour, min, 0, 0, time.UTC)
		s.repo.UpdateNextRunAt(ctx, nextPost.ID, nextRun)
		logger.Info("SCHEDULER", "Fila %s: próximo post %s agendado para %s", sentPost.QueueGroupID, nextPost.ID, nextRun.Format("02/01 15:04"))
	} else if sentPost.LoopQueue {
		s.resetQueue(ctx, sentPost.QueueGroupID, posts, sentPost.ScheduleTime)
	} else {
		logger.Info("SCHEDULER", "Fila %s concluída", sentPost.QueueGroupID)
	}
}

func (s *SchedulerService) resetQueue(ctx context.Context, queueGroupID string, posts []models.ScheduledPost, scheduleTime string) {
	hour, min := parseHHMM(scheduleTime)
	nextRun := time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day()+1, hour, min, 0, 0, time.UTC)

	for i := range posts {
		if posts[i].Status == "sent" {
			s.repo.UpdateStatus(ctx, posts[i].ID, "pending")
		}
	}

	if len(posts) > 0 {
		s.repo.UpdateNextRunAt(ctx, posts[0].ID, nextRun)
	}
	logger.Info("SCHEDULER", "Fila %s reiniciada, próximo: %s", queueGroupID, nextRun.Format("02/01 15:04"))
}

func parseHHMM(s string) (int, int) {
	parts := strings.Split(s, ":")
	if len(parts) != 2 {
		return 12, 0
	}
	hour := 0
	min := 0
	fmt.Sscanf(parts[0], "%d", &hour)
	fmt.Sscanf(parts[1], "%d", &min)
	return hour, min
}

// CRUD methods

func (s *SchedulerService) CreateScheduledPost(ctx context.Context, ownerID, channelID int64, channelTitle, postData string, opts ScheduleOptions) (*models.ScheduledPost, error) {
	limit, _ := s.repo.CountByOwner(ctx, ownerID)
	if limit >= 50 {
		return nil, fmt.Errorf("limite de 50 agendamentos ativos atingido")
	}

	id := fmt.Sprintf("sch_%d_%d", time.Now().UnixNano(), ownerID)

	nextRunAt := time.Now().UTC()
	if opts.ScheduleType == "once" && opts.ScheduledAt != nil {
		nextRunAt = *opts.ScheduledAt
	} else if opts.ScheduleType == "daily" || opts.ScheduleType == "weekly" {
		hour, min := parseHHMM(opts.ScheduleTime)
		nextRun := time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), hour, min, 0, 0, time.UTC)
		if nextRun.Before(time.Now().UTC()) {
			nextRun = nextRun.AddDate(0, 0, 1)
		}
		nextRunAt = nextRun
	}

	var scheduleDaysStr string
	if len(opts.ScheduleDays) > 0 {
		b, _ := json.Marshal(opts.ScheduleDays)
		scheduleDaysStr = string(b)
	}

	post := &models.ScheduledPost{
		ID:             id,
		OwnerID:        ownerID,
		ChannelID:      channelID,
		ChannelTitle:   channelTitle,
		PostData:       postData,
		ScheduleType:   opts.ScheduleType,
		ScheduleTime:   opts.ScheduleTime,
		ScheduledAt:    opts.ScheduledAt,
		ScheduleDays:   scheduleDaysStr,
		NextRunAt:      nextRunAt,
		RepeatUntil:    opts.RepeatUntil,
		QueueGroupID:   opts.QueueGroupID,
		QueuePosition:  opts.QueuePosition,
		LoopQueue:      opts.LoopQueue,
		Status:         "pending",
		SentCount:      0,
	}

	if err := s.repo.Create(ctx, post); err != nil {
		return nil, err
	}
	return post, nil
}

func (s *SchedulerService) CancelScheduledPost(ctx context.Context, id string, ownerID int64) error {
	post, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if post.OwnerID != ownerID {
		return fmt.Errorf("não autorizado")
	}
	return s.repo.UpdateStatus(ctx, id, "cancelled")
}

func (s *SchedulerService) PauseScheduledPost(ctx context.Context, id string, ownerID int64) error {
	post, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if post.OwnerID != ownerID {
		return fmt.Errorf("não autorizado")
	}
	return s.repo.UpdateStatus(ctx, id, "paused")
}

func (s *SchedulerService) ResumeScheduledPost(ctx context.Context, id string, ownerID int64) error {
	post, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if post.OwnerID != ownerID {
		return fmt.Errorf("não autorizado")
	}
	return s.repo.UpdateStatus(ctx, id, "pending")
}

func (s *SchedulerService) GetUserSchedules(ctx context.Context, ownerID int64) ([]models.ScheduledPost, error) {
	return s.repo.GetByOwner(ctx, ownerID)
}

func (s *SchedulerService) DeleteScheduledPost(ctx context.Context, id string, ownerID int64) error {
	post, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if post.OwnerID != ownerID {
		return fmt.Errorf("não autorizado")
	}
	return s.repo.Delete(ctx, id)
}

func (s *SchedulerService) GetScheduleByID(ctx context.Context, id string) (*models.ScheduledPost, error) {
	return s.repo.GetByID(ctx, id)
}
