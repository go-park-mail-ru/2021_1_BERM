package models

type Session struct {
	SessionID string
	UserId    uint64 `json:"user_id"`
	Executor  bool   `json:"executor"`
}
