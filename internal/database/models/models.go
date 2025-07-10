package models

import "time"

type User struct {
	UserId       int64     `gorm:"primaryKey" json:"id"` // ID do Telegram
	FirstName    string    `json:"firstName"`
	IsContribute bool      `gorm:"default:false" json:"isContribute"`
	Channels     []Channel `gorm:"foreignKey:OwnerID" json:"channels"`
	CreatedAt    time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt    time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

type Channel struct {
	ID             int64           `gorm:"primaryKey" json:"id"` // ID do Telegram
	Title          string          `json:"title"`
	NewPackCaption string          `json:"newPackCaption"`
	InviteURL      string          `json:"inviteUrl"`
	OwnerID        int64           `json:"ownerId"`
	DefaultCaption *DefaultCaption `gorm:"foreignKey:OwnerChannelID" json:"defaultCaption,omitempty"`
	Buttons        []Button        `gorm:"foreignKey:OwnerChannelID" json:"buttons"`
	Separator      *Separator      `gorm:"foreignKey:OwnerChannelID" json:"separator,omitempty"`
	CustomCaptions []CustomCaption `gorm:"foreignKey:OwnerChannelID" json:"customCaptions"`
	CreatedAt      time.Time       `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt      time.Time       `gorm:"autoUpdateTime" json:"updated_at"`
}

type DefaultCaption struct {
	CaptionID         string             `gorm:"type:text;primaryKey" json:"captionId"`
	Caption           string             `json:"caption"`
	MessagePermission *MessagePermission `gorm:"foreignKey:OwnerCaptionID" json:"messagePermission,omitempty"`
	ButtonsPermission *ButtonsPermission `gorm:"foreignKey:OwnerCaptionID" json:"buttonsPermission,omitempty"`
	OwnerChannelID    int64              `gorm:"unique" json:"ownerChannelId"`
	CreatedAt         time.Time          `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt         time.Time          `gorm:"autoUpdateTime" json:"updated_at"`
}

type MessagePermission struct {
	MessagePermissionID string    `gorm:"type:text;primaryKey" json:"messagePermissionId"`
	LinkPreview         bool      `gorm:"default:true" json:"linkPreview"`
	Message             bool      `gorm:"default:true" json:"message"`
	Audio               bool      `gorm:"default:true" json:"audio"`
	Video               bool      `gorm:"default:true" json:"video"`
	Photo               bool      `gorm:"default:true" json:"photo"`
	Sticker             bool      `gorm:"default:true" json:"sticker"`
	GIF                 bool      `gorm:"default:true" json:"gif"`
	OwnerCaptionID      string    `gorm:"unique" json:"ownerCaptionId"`
	CreatedAt           time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt           time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

type ButtonsPermission struct {
	ButtonsPermissionID string    `gorm:"type:text;primaryKey" json:"buttonsPermissionId"`
	Message             bool      `gorm:"default:true" json:"message"`
	Audio               bool      `gorm:"default:true" json:"audio"`
	Video               bool      `gorm:"default:true" json:"video"`
	Photo               bool      `gorm:"default:true" json:"photo"`
	Sticker             bool      `gorm:"default:true" json:"sticker"`
	GIF                 bool      `gorm:"default:true" json:"gif"`
	OwnerCaptionID      string    `gorm:"unique" json:"ownerCaptionId"`
	CreatedAt           time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt           time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

type Button struct {
	ButtonID       string    `gorm:"type:text;primaryKey" json:"buttonId"`
	NameButton     string    `json:"nameButton"`
	ButtonURL      string    `json:"buttonUrl"`
	PositionX      int       `gorm:"default:0" json:"positionX"`
	PositionY      int       `gorm:"default:0" json:"positionY"`
	OwnerChannelID int64     `json:"ownerChannelId"`
	CreatedAt      time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt      time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

type Separator struct {
	ID             string    `gorm:"type:text;primaryKey" json:"id"`
	SeparatorID    string    `json:"separatorId"`
	SeparatorURL   string    `json:"separatorUrl"`
	OwnerChannelID int64     `gorm:"unique" json:"ownerChannelId"`
	CreatedAt      time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt      time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

type CustomCaption struct {
	CaptionID      string                `gorm:"type:text;primaryKey" json:"captionId"`
	Code           string                `json:"code"`
	Caption        string                `json:"caption"`
	LinkPreview    bool                  `json:"linkPreview"`
	Buttons        []CustomCaptionButton `gorm:"foreignKey:OwnerCaptionID" json:"buttons"`
	OwnerChannelID int64                 `json:"ownerChannelId"`
	CreatedAt      time.Time             `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt      time.Time             `gorm:"autoUpdateTime" json:"updated_at"`
}

type CustomCaptionButton struct {
	ButtonID       string    `gorm:"type:text;primaryKey" json:"buttonId"`
	NameButton     string    `json:"nameButton"`
	ButtonURL      string    `json:"buttonUrl"`
	PositionX      int       `gorm:"default:0" json:"positionX"`
	PositionY      int       `gorm:"default:0" json:"positionY"`
	OwnerCaptionID string    `json:"ownerCaptionId"`
	CreatedAt      time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt      time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

func (User) TableName() string {
	return "users"
}
