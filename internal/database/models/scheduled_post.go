package models

import "time"

type ScheduledPost struct {
	ID             string     `gorm:"type:uuid;primaryKey" json:"id"`
	OwnerID        int64      `gorm:"index:idx_schedule_owner" json:"ownerId"`
	ChannelID      int64      `gorm:"index" json:"channelId"`
	ChannelTitle   string     `json:"channelTitle"`
	PostData       string     `gorm:"type:text" json:"postData"` // JSON do PostBuilderState

	ScheduleType  string  `json:"scheduleType"`  // "once" | "daily" | "weekly" | "queue"
	ScheduleTime  string  `json:"scheduleTime"`  // "HH:MM" UTC para recorrente
	ScheduledAt   *time.Time `json:"scheduledAt"` // datetime exato para "once"
	ScheduleDays  string  `json:"scheduleDays"`  // JSON "[1,3,5]" para weekly
	NextRunAt     time.Time `gorm:"index:idx_schedule_next_run" json:"nextRunAt"`
	RepeatUntil   *time.Time `json:"repeatUntil"`

	QueueGroupID  string `gorm:"index:idx_schedule_queue" json:"queueGroupId"`
	QueuePosition int    `json:"queuePosition"`
	LoopQueue     bool   `json:"loopQueue"`

	Status    string     `gorm:"index:idx_schedule_next_run" json:"status"` // "pending"|"sent"|"cancelled"|"paused"|"failed"
	SentAt    *time.Time `json:"sentAt"`
	SentCount int        `json:"sentCount"`
	LastError string     `json:"lastError"`
	RetryCount int       `json:"retryCount"`
	CreatedAt time.Time  `gorm:"autoCreateTime" json:"createdAt"`
	UpdatedAt time.Time  `gorm:"autoUpdateTime" json:"updatedAt"`
}
