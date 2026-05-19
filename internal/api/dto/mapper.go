package dto

import "github.com/leirbagxis/FreddyBot/internal/database/models"

func ToUserDTO(u *models.User) UserDTO {
	if u == nil {
		return UserDTO{}
	}
	dto := UserDTO{
		ID:            u.UserId,
		FirstName:     u.FirstName,
		Username:      u.Username,
		IsAdmin:       u.IsAdmin,
		IsBlacklisted: u.IsBlacklisted,
		IsContribute:  u.IsContribute,
	}

	if len(u.Channels) > 0 {
		for _, c := range u.Channels {
			dto.Channels = append(dto.Channels, ToChannelDTO(&c))
		}
	}

	return dto
}

func ToChannelDTO(c *models.Channel) ChannelDTO {
	if c == nil {
		return ChannelDTO{}
	}
	dto := ChannelDTO{
		ID:               c.ID,
		Title:            c.Title,
		NewPackCaption:   c.NewPackCaption,
		InviteURL:        c.InviteURL,
		OwnerID:          c.OwnerID,
		Reactions:        c.Reactions,
		ReactionPosition: c.ReactionPosition,
		DynamicLinks:     c.DynamicLinks,
		DLBotButtons:     c.DLBotButtons,
		DLBotCaptions:    c.DLBotCaptions,
		DLBotReactions:   c.DLBotReactions,
		CreatedAt:        c.CreatedAt,
		UpdatedAt:        c.UpdatedAt,
	}

	if c.DefaultCaption != nil {
		dto.DefaultCaption = ToDefaultCaptionDTO(c.DefaultCaption)
	}

	if len(c.Buttons) > 0 {
		for _, b := range c.Buttons {
			dto.Buttons = append(dto.Buttons, ToButtonDTO(&b))
		}
	}

	if len(c.CustomCaptions) > 0 {
		for _, cc := range c.CustomCaptions {
			dto.CustomCaptions = append(dto.CustomCaptions, ToCustomCaptionDTO(&cc))
		}
	}

	return dto
}

func ToDefaultCaptionDTO(dc *models.DefaultCaption) *DefaultCaptionDTO {
	dto := &DefaultCaptionDTO{
		CaptionID: dc.CaptionID,
		Caption:   dc.Caption,
		CreatedAt: dc.CreatedAt,
	}

	if dc.MessagePermission != nil {
		dto.MessagePermission = &PermissionDTO{
			LinkPreview: dc.MessagePermission.LinkPreview,
			Message:     dc.MessagePermission.Message,
			Audio:       dc.MessagePermission.Audio,
			Video:       dc.MessagePermission.Video,
			Photo:       dc.MessagePermission.Photo,
			Document:    dc.MessagePermission.Document,
			Sticker:     dc.MessagePermission.Sticker,
			GIF:         dc.MessagePermission.GIF,
			Reactions:   dc.MessagePermission.Reactions,
		}
	}

	if dc.ButtonsPermission != nil {
		dto.ButtonsPermission = &PermissionDTO{
			Message:  dc.ButtonsPermission.Message,
			Audio:    dc.ButtonsPermission.Audio,
			Video:    dc.ButtonsPermission.Video,
			Photo:    dc.ButtonsPermission.Photo,
			Document: dc.ButtonsPermission.Document,
			Sticker:  dc.ButtonsPermission.Sticker,
			GIF:      dc.ButtonsPermission.GIF,
		}
	}

	return dto
}

func ToButtonDTO(b *models.Button) ButtonDTO {
	return ButtonDTO{
		ButtonID:  b.ButtonID,
		Name:      b.NameButton,
		URL:       b.ButtonURL,
		PositionX: b.PositionX,
		PositionY: b.PositionY,
	}
}

func ToCustomCaptionButtonDTO(b *models.CustomCaptionButton) ButtonDTO {
	return ButtonDTO{
		ButtonID:  b.ButtonID,
		Name:      b.NameButton,
		URL:       b.ButtonURL,
		PositionX: b.PositionX,
		PositionY: b.PositionY,
	}
}

func ToCustomCaptionDTO(cc *models.CustomCaption) CustomCaptionDTO {
	dto := CustomCaptionDTO{
		CaptionID:   cc.CaptionID,
		Code:        cc.Code,
		Caption:     cc.Caption,
		LinkPreview: cc.LinkPreview,
		CreatedAt:   cc.CreatedAt,
	}

	if len(cc.Buttons) > 0 {
		for _, b := range cc.Buttons {
			dto.Buttons = append(dto.Buttons, ToCustomCaptionButtonDTO(&b))
		}
	}

	return dto
}
