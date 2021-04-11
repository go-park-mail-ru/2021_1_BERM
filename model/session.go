package model

type Session struct {
	SessionId string
	UserId    uint64 `json:"id,omitempty"`
}
