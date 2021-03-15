package teststore

import "fl_ru/model"

type SessionRepository struct {
	store    *Store
	sessions map[string]*model.Session
}

func (r *SessionRepository) Create(session *model.Session) error {
	r.sessions[session.SessionID] = session

	return nil
}

func (r *SessionRepository) Find(session *model.Session) error {
	s, ok := r.sessions[session.SessionID]
	if !ok {
		return nil
	}
	*session = *s

	return nil
}
