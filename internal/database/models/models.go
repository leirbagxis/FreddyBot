package models

import "time"

type User struct {
	UserId    int64     `gorm:"primaryKey" json:"user_id"`
	FirstName string    `json:"first_name"`
	Channels  []Channel `gorm:"foreignKey:OwnerID" json:"channels"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoCreateTime" json:"updated_at"`
}

type Channel struct {
	TelegramChannelID int64            `gorm:"primaryKey" json:"telegramChannelId"`
	Title             string           `json:"title"`
	InviteURL         string           `json:"inviteUrl"`
	OwnerID           int64            `json:"ownerId"`
	DefaultCaptions   []DefaultCaption `gorm:"foreignKey:ChannelID" json:"defaultCaptions"`
	CreatedAt         time.Time        `gorm:"autoCreateTime" json:"createdAt"`
}

type DefaultCaption struct {
	DefaultCaptionID string             `gorm:"primaryKey" json:"defaultCaptionId"`
	ChannelID        int64              `json:"channelId"`
	Caption          string             `json:"caption"`
	MessagePerm      *MessagePermission `gorm:"foreignKey:DefaultCaptionID" json:"messagePermission,omitempty"`
	ButtonPerm       *ButtonPermission  `gorm:"foreignKey:DefaultCaptionID" json:"buttonPermission,omitempty"`
	Buttons          []Button           `gorm:"foreignKey:DefaultCaptionID" json:"buttons"`
}

type MessagePermission struct {
	MessagePermissionID string `gorm:"primaryKey" json:"messagePermissionId"`
	DefaultCaptionID    string `json:"defaultCaptionId"`
	LinkPreview         bool   `gorm:"default:true" json:"linkPreview"`
	Message             bool   `gorm:"default:true" json:"message"`
	Sticker             bool   `gorm:"default:true" json:"sticker"`
}

type ButtonPermission struct {
	ButtonsPermissionID string `gorm:"primaryKey" json:"buttonsPermissionId"`
	DefaultCaptionID    string `json:"defaultCaptionId"`
	Message             bool   `gorm:"default:true" json:"message"`
	Sticker             bool   `gorm:"default:true" json:"sticker"`
}

type Button struct {
	ButtonID         string `gorm:"primaryKey" json:"buttonId"`
	DefaultCaptionID string `json:"defaultCaptionId"`
	ButtonName       string `json:"buttonName"`
	ButtonURL        string `json:"buttonUrl"`
}
