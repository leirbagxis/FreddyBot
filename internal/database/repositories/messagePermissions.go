package repositories

import "gorm.io/gorm"

type MessagePermissionRepository struct {
	db *gorm.DB
}

func NewMessagePermissionRepository(db *gorm.DB) *MessagePermissionRepository {
	return &MessagePermissionRepository{db: db}
}
