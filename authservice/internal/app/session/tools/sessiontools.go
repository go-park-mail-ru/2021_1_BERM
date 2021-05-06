package tools

import models "authorizationservice/internal/app/models"

//go:generate mockgen -destination=../mock/mock_tools.go -package=mock authorizationservice/internal/app/session/tools SessionTools
type SessionTools interface {
	BeforeCreate(session models.Session) (models.Session, error)
	EncodingSessionToTarantool(sess *models.Session) []interface{}
	DecodingTarantoolToSession(data []interface{}) *models.Session
}
