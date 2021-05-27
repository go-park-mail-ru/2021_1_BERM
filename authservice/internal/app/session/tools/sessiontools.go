package sessiontools

import models "authorizationservice/internal/app/models"

//nolint:lll    //go:generate mockgen -destination=../mock/mock_tools.go -package=mock authorizationservice/internal/app/session/sessiontools SessionTools
type SessionTools interface {
	BeforeCreate(session models.Session) (models.Session, error)
	EncodingSessionToTarantool(sess *models.Session) []interface{}
	DecodingTarantoolToSession(data []interface{}) *models.Session
}
