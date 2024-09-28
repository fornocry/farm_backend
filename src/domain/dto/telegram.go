package dto

type InitData struct {
	AuthDateRaw     int          `json:"auth_date"`
	CanSendAfterRaw int          `json:"can_send_after"`
	Chat            Chat         `json:"chat"`
	ChatType        ChatType     `json:"chat_type"`
	ChatInstance    int64        `json:"chat_instance"`
	Hash            string       `json:"hash"`
	QueryID         string       `json:"query_id"`
	Receiver        TelegramUser `json:"receiver"`
	StartParam      string       `json:"start_param"`
	TelegramUser    TelegramUser `json:"user"`
}
type TelegramUser struct {
	AddedToAttachmentMenu bool   `json:"added_to_attachment_menu"`
	AllowsWriteToPm       bool   `json:"allows_write_to_pm"`
	FirstName             string `json:"first_name"`
	ID                    int64  `json:"id"`
	IsBot                 bool   `json:"is_bot"`
	IsPremium             bool   `json:"is_premium"`
	LastName              string `json:"last_name"`
	Username              string `json:"username"`
	LanguageCode          string `json:"language_code"`
	PhotoURL              string `json:"photo_url"`
}

const (
	ChatTypeSender     ChatType = "sender"
	ChatTypePrivate    ChatType = "private"
	ChatTypeGroup      ChatType = "group"
	ChatTypeSupergroup ChatType = "supergroup"
	ChatTypeChannel    ChatType = "channel"
)

type ChatType string

func (c ChatType) Known() bool {
	switch c {
	case ChatTypeSender,
		ChatTypePrivate,
		ChatTypeGroup,
		ChatTypeSupergroup,
		ChatTypeChannel:
		return true
	default:
		return false
	}
}

type Chat struct {
	ID       int64    `json:"id"`
	Type     ChatType `json:"type"`
	Title    string   `json:"title"`
	PhotoURL string   `json:"photo_url"`
	Username string   `json:"username"`
}
