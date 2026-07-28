package types

type CaptionButtonSnapshot struct {
	NameButton string `json:"nameButton"`
	ButtonURL  string `json:"buttonUrl"`
}

type CustomCaptionSnapshot struct {
	Code        string                 `json:"code"`
	Caption     string                 `json:"caption"`
	LinkPreview bool                   `json:"linkPreview"`
	Buttons     []CaptionButtonSnapshot `json:"buttons,omitempty"`
}

type CaptionTemplateData struct {
	DefaultCaption string                  `json:"defaultCaption"`
	CustomCaptions []CustomCaptionSnapshot `json:"customCaptions,omitempty"`
}
