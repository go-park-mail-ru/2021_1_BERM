package model

type Session struct {
	SessionID string `json:"sesion_id,omitempty"`
	UserID    uint64 `json:"id,omitempty"`
	Executor  bool   `json:"executor,omitempty"`
}
