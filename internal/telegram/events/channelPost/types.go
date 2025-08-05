package channelpost

import (
	"sync"
	"time"
)

type MessageType string

const (
	MessageTypeText      MessageType = "text"
	MessageTypeAudio     MessageType = "audio"
	MessageTypeSticker   MessageType = "sticker"
	MessageTypePhoto     MessageType = "photo"
	MessageTypeVideo     MessageType = "video"
	MessageTypeAnimation MessageType = "animation"
)

var PermissionMap = map[MessageType]string{
	MessageTypeText:      "message",
	MessageTypeAudio:     "audio",
	MessageTypeVideo:     "video",
	MessageTypePhoto:     "photo",
	MessageTypeSticker:   "sticker",
	MessageTypeAnimation: "gif",
}

// Constants for timeouts, retries, cache TTLs
const (
	MediaGroupTimeout = 1000 * time.Millisecond
	CleanupTimeout    = 60000 * time.Millisecond
	MaxRetryAttempts  = 3
	RetryDelay        = 1000 * time.Millisecond
	CacheTTL          = 5 * time.Minute
)

// Thread-safe struct for media group data
type MediaGroupInfo struct {
	Messages           []MediaMessage
	Processed          bool
	MessageEditAllowed bool
	FirstMessageID     int
	Timer              *time.Timer
	mu                 sync.Mutex
}

type MediaMessage struct {
	MessageID       int
	HasCaption      bool
	Caption         string
	CaptionEntities []interface{}
}

type ProcessedGroup struct {
	Timestamp time.Time
}

// MediaGroupManager controls media groups with concurrency support
type MediaGroupManager struct {
	groups          sync.Map // map[string]*MediaGroupInfo
	processedGroups sync.Map // map[string]ProcessedGroup
	newPackChannels sync.Map // map[int64]bool
}

func NewMediaGroupManager() *MediaGroupManager {
	mgm := &MediaGroupManager{}
	go mgm.cleanupRoutine()
	return mgm
}

func (mgm *MediaGroupManager) cleanupRoutine() {
	ticker := time.NewTicker(CleanupTimeout)
	defer ticker.Stop()
	for range ticker.C {
		now := time.Now()
		mgm.processedGroups.Range(func(key, value interface{}) bool {
			if group, ok := value.(ProcessedGroup); ok {
				if now.Sub(group.Timestamp) > CleanupTimeout {
					mgm.processedGroups.Delete(key)
				}
			}
			return true
		})
	}
}

func (mgm *MediaGroupManager) GetMediaGroup(groupID string) (*MediaGroupInfo, bool) {
	if val, ok := mgm.groups.Load(groupID); ok {
		return val.(*MediaGroupInfo), true
	}
	return nil, false
}

func (mgm *MediaGroupManager) SetMediaGroup(groupID string, group *MediaGroupInfo) {
	mgm.groups.Store(groupID, group)
}

func (mgm *MediaGroupManager) DeleteMediaGroup(groupID string) {
	mgm.groups.Delete(groupID)
}

func (mgm *MediaGroupManager) IsProcessed(groupID string) bool {
	_, exists := mgm.processedGroups.Load(groupID)
	return exists
}

func (mgm *MediaGroupManager) MarkProcessed(groupID string) {
	mgm.processedGroups.Store(groupID, ProcessedGroup{Timestamp: time.Now()})
}

func (mgm *MediaGroupManager) IsNewPackActive(channelID int64) bool {
	val, exists := mgm.newPackChannels.Load(channelID)
	return exists && val.(bool)
}

func (mgm *MediaGroupManager) SetNewPackActive(channelID int64, active bool) {
	if active {
		mgm.newPackChannels.Store(channelID, true)
	} else {
		mgm.newPackChannels.Delete(channelID)
	}
}
