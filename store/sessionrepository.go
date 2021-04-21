package store

import "FL_2/model"

//go:generate mockgen -destination=mock/mock_session_repo.go -package=mock FL_2/store SessionRepository
type SessionRepository interface {
	Create(session *model.Session) error
	Find(session *model.Session) error
}
