package repositories

import (
	"context"

	"github.com/leirbagxis/FreddyBot/internal/database/models"
	"gorm.io/gorm"
)

type VoteRepository struct {
	db *gorm.DB
}

func NewVoteRepository(db *gorm.DB) *VoteRepository {
	return &VoteRepository{db: db}
}

// ToggleVote adds a vote if it doesn't exist, or removes it if it does.
// Only one emoji can be voted per user per message.
// Returns true if added, false if removed, and error if any.
func (r *VoteRepository) ToggleVote(ctx context.Context, chatID int64, messageID int, userID int64, emoji string) (bool, int64, error) {
	var existingVote models.Vote
	// Check if the user has ALREADY voted for THIS SPECIFIC emoji
	err := r.db.WithContext(ctx).Where("chat_id = ? AND message_id = ? AND user_id = ? AND emoji = ?", chatID, messageID, userID, emoji).First(&existingVote).Error

	if err == nil {
		// User clicked the same emoji again -> Remove it (Toggle OFF)
		if err := r.db.WithContext(ctx).Delete(&existingVote).Error; err != nil {
			return false, 0, err
		}
		
		var count int64
		r.db.WithContext(ctx).Model(&models.Vote{}).Where("chat_id = ? AND message_id = ? AND emoji = ?", chatID, messageID, emoji).Count(&count)
		return false, count, nil
	}

	if err != nil && err != gorm.ErrRecordNotFound {
		return false, 0, err
	}

	// If we reach here, the user is either voting for the first time OR changing their vote.
	// 1. Remove ANY previous vote by this user on this message (Ensures "transfer")
	if err := r.db.WithContext(ctx).Where("chat_id = ? AND message_id = ? AND user_id = ?", chatID, messageID, userID).Delete(&models.Vote{}).Error; err != nil {
		return false, 0, err
	}

	// 2. Add the new vote
	newVote := models.Vote{
		ChatID:    chatID,
		MessageID: messageID,
		UserID:    userID,
		Emoji:     emoji,
	}
	if err := r.db.WithContext(ctx).Create(&newVote).Error; err != nil {
		return false, 0, err
	}

	// Get total count for the new emoji
	var count int64
	err = r.db.WithContext(ctx).Model(&models.Vote{}).Where("chat_id = ? AND message_id = ? AND emoji = ?", chatID, messageID, emoji).Count(&count).Error

	return true, count, err
}

func (r *VoteRepository) GetVoteCounts(ctx context.Context, chatID int64, messageID int) (map[string]int64, error) {
	type Result struct {
		Emoji string
		Count int64
	}
	var results []Result
	err := r.db.WithContext(ctx).Model(&models.Vote{}).
		Select("emoji, count(*) as count").
		Where("chat_id = ? AND message_id = ?", chatID, messageID).
		Group("emoji").
		Scan(&results).Error

	if err != nil {
		return nil, err
	}

	counts := make(map[string]int64)
	for _, res := range results {
		counts[res.Emoji] = res.Count
	}
	return counts, nil
}
