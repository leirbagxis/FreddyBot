package channelpost

import (
	"fmt"
	"sync"
	"time"

	"github.com/leirbagxis/FreddyBot/internal/database/models"
	dbmodels "github.com/leirbagxis/FreddyBot/internal/database/models"
)

// PermissionManager gerencia cache de permissões para evitar consultas repetidas
type PermissionManager struct {
	cache    sync.Map // map[string]CacheEntry
	cacheTTL time.Duration
}

type CacheEntry struct {
	Value     bool
	Timestamp time.Time
}

const cacheTTL = 5 * time.Minute

func NewPermissionManager() *PermissionManager {
	pm := &PermissionManager{
		cacheTTL: cacheTTL,
	}
	go pm.cleanupRoutine()
	return pm
}

// Em permissions.go (ou onde for apropriado)
func (pm *PermissionManager) CheckPermissions(channel *dbmodels.Channel, messageType MessageType) *PermissionCheckResult {
	// Aqui você pode usar a lógica que já está duplicada nos métodos IsMessageEditAllowed/IsButtonsAllowed,
	// ou simplesmente reusar esses métodos e montar o PermissionCheckResult.
	result := &PermissionCheckResult{
		CanEdit:       pm.IsMessageEditAllowed(channel, messageType),
		CanAddButtons: pm.IsButtonsAllowed(channel, messageType),
		// Complete outros campos conforme sua estrutura
		CanEditButtons:    true,
		CanUseLinkPreview: true,
		Reason:            "",
	}
	return result
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

	mp := channel.DefaultCaption.MessagePermission
	permissionKey := PermissionMap[messageType]
	if permissionKey == "" {
		return pm.setCached(key, true)
	}

	var allowed bool
	switch permissionKey {
	case "message":
		allowed = mp.Message
	case "audio":
		allowed = mp.Audio
	case "video":
		allowed = mp.Video
	case "photo":
		allowed = mp.Photo
	case "sticker":
		allowed = mp.Sticker
	case "gif":
		allowed = mp.GIF
	default:
		allowed = true
	}
	return pm.setCached(key, allowed)
}

func (pm *PermissionManager) IsButtonsAllowed(channel *models.Channel, messageType MessageType) bool {
	key := fmt.Sprintf("%d_%s_buttons", channel.ID, messageType)
	if cached := pm.getCached(key); cached != nil {
		return *cached
	}

	if channel.DefaultCaption == nil || channel.DefaultCaption.ButtonsPermission == nil {
		return pm.setCached(key, true)
	}

	bp := channel.DefaultCaption.ButtonsPermission
	permissionKey := PermissionMap[messageType]
	if permissionKey == "" {
		return pm.setCached(key, true)
	}

	var allowed bool
	switch permissionKey {
	case "message":
		allowed = bp.Message
	case "audio":
		allowed = bp.Audio
	case "video":
		allowed = bp.Video
	case "photo":
		allowed = bp.Photo
	case "sticker":
		allowed = bp.Sticker
	case "gif":
		allowed = bp.GIF
	default:
		allowed = true
	}
	return pm.setCached(key, allowed)
}

func (pm *PermissionManager) getCached(key string) *bool {
	if v, ok := pm.cache.Load(key); ok {
		entry := v.(CacheEntry)
		if time.Since(entry.Timestamp) <= pm.cacheTTL {
			return &entry.Value
		}
		pm.cache.Delete(key)
	}
	return nil
}

func (pm *PermissionManager) setCached(key string, val bool) bool {
	pm.cache.Store(key, CacheEntry{
		Value:     val,
		Timestamp: time.Now(),
	})
	return val
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
