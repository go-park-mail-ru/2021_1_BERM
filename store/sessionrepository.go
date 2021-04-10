package store

import "FL_2/model"

type SessionRepository interface {
	Create(session *model.Session) error
	Find(session *model.Session) error
}
