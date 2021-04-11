package model

type Session struct {
	SessionId string `json:"sesion_id,omitempty"`
	UserId    uint64 `json:"id,omitempty"`
	Executor  bool   `json:"executor,omitempty"`
}
