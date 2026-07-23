package types

type CreateScheduleRequest struct {
	ChannelID    int64  `json:"channelId" binding:"required"`
	ScheduleType string `json:"scheduleType" binding:"required,oneof=once daily weekly queue"`
	ScheduleTime string `json:"scheduleTime"`
	ScheduledAt  string `json:"scheduledAt"`
	ScheduleDays []int  `json:"scheduleDays"`
	RepeatUntil  string `json:"repeatUntil"`
	LoopQueue    bool   `json:"loopQueue"`
}

type UpdateScheduleStatusRequest struct {
	Status string `json:"status" binding:"required,oneof=paused pending cancelled"`
}
