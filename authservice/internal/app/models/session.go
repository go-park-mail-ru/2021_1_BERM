package models

type Session struct {
	SessionID string `json:"-"`
	UserId    uint64 `json:"id"`
	Executor  bool   `json:"executor"`
}
