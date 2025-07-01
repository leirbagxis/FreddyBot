package channelpost

import (
	"fmt"
	"sync"
	"time"

	"github.com/leirbagxis/FreddyBot/internal/database/models"
)

// ✅ CORRIGIDO: Cache thread-safe com sync.Map
type PermissionManager struct {
	cache    sync.Map // string -> CacheEntry
	cacheTTL time.Duration
}

type CacheEntry struct {
	Value     bool
	Timestamp time.Time
}

func NewPermissionManager() *PermissionManager {
	pm := &PermissionManager{
		cacheTTL: CacheTTL,
	}

	// Cleanup automático do cache
	go pm.cleanupRoutine()

	return pm
}

func (pm *PermissionManager) cleanupRoutine() {
	ticker := time.NewTicker(pm.cacheTTL)
	defer ticker.Stop()

	for range ticker.C {
		now := time.Now()
		pm.cache.Range(func(key, value interface{}) bool {
			if entry, ok := value.(CacheEntry); ok {
				if now.Sub(entry.Timestamp) > pm.cacheTTL {
					pm.cache.Delete(key)
				}
			}
			return true
		})
	}
}

func (pm *PermissionManager) IsMessageEditAllowed(channel *models.Channel, messageType MessageType) bool {
	key := fmt.Sprintf("%d_%s_message", channel.ID, messageType)

	if cached := pm.getCached(key); cached != nil {
		return *cached
	}

	if channel.DefaultCaption == nil || channel.DefaultCaption.MessagePermission == nil {
		return pm.setCached(key, true)
	}

	messagePermission := channel.DefaultCaption.MessagePermission
	permissionKey := PermissionMap[messageType]

	if permissionKey == "" {
		return pm.setCached(key, true)
	}

	var isAllowed bool
	switch permissionKey {
	case "message":
		isAllowed = messagePermission.Message
	case "audio":
		isAllowed = messagePermission.Audio
	case "video":
		isAllowed = messagePermission.Video
	case "photo":
		isAllowed = messagePermission.Photo
	case "sticker":
		isAllowed = messagePermission.Sticker
	case "gif":
		isAllowed = messagePermission.GIF
	default:
		isAllowed = true
	}

	return pm.setCached(key, isAllowed)
}

func (pm *PermissionManager) IsButtonsAllowed(channel *models.Channel, messageType MessageType) bool {
	key := fmt.Sprintf("%d_%s_buttons", channel.ID, messageType)

	if cached := pm.getCached(key); cached != nil {
		return *cached
	}

	if channel.DefaultCaption == nil || channel.DefaultCaption.ButtonsPermission == nil {
		return pm.setCached(key, true)
	}

	buttonsPermission := channel.DefaultCaption.ButtonsPermission
	permissionKey := PermissionMap[messageType]

	if permissionKey == "" {
		return pm.setCached(key, true)
	}

	var isAllowed bool
	switch permissionKey {
	case "message":
		isAllowed = buttonsPermission.Message
	case "audio":
		isAllowed = buttonsPermission.Audio
	case "video":
		isAllowed = buttonsPermission.Video
	case "photo":
		isAllowed = buttonsPermission.Photo
	case "sticker":
		isAllowed = buttonsPermission.Sticker
	case "gif":
		isAllowed = buttonsPermission.GIF
	default:
		isAllowed = true
	}

	return pm.setCached(key, isAllowed)
}

func (pm *PermissionManager) getCached(key string) *bool {
	if value, ok := pm.cache.Load(key); ok {
		entry := value.(CacheEntry)
		if time.Since(entry.Timestamp) <= pm.cacheTTL {
			return &entry.Value
		}
		pm.cache.Delete(key)
	}
	return nil
}

func (pm *PermissionManager) setCached(key string, value bool) bool {
	pm.cache.Store(key, CacheEntry{
		Value:     value,
		Timestamp: time.Now(),
	})
	return value
}

func (pm *PermissionManager) ClearCache() {
	pm.cache.Range(func(key, value interface{}) bool {
		pm.cache.Delete(key)
		return true
	})
}

func (pm *PermissionManager) GetCacheStats() map[string]interface{} {
	count := 0
	pm.cache.Range(func(key, value interface{}) bool {
		count++
		return true
	})

	return map[string]interface{}{
		"size": count,
		"ttl":  pm.cacheTTL,
	}
}
