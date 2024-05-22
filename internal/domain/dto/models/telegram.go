package models

type TelegramInputDTO struct {
	Username string `json:"username"`
	UserID   uint64 `json:"user_id"`
}

type TelegramOutputDTO struct {
	ID       uint64 `json:"id"`
	UserID   uint64 `json:"user_id"`
	ChatID   uint64 `json:"chat_id"`
	Username string `json:"username"`
	DoAlerts bool   `json:"do_alerts"`
}
