package tools

import (
	"authorizationservice/internal/models"
	"authorizationservice/pkg/Error"
	"golang.org/x/crypto/argon2"
	"math/rand"
)

const(
	saltLength = 8
)

func BeforeCreate(session* models.Session) error {
	salt := make([]byte, saltLength)
	_, err := rand.Read(salt)
	if err != nil {
		return &Error.Error{
			Err: err,
			InternalError: true,
			ErrorDescription: Error.InternalServerErrorDescription,
		}
	}

	session.SessionID = string(hashSessionId(salt, session.SessionID))
	return nil
}

func hashSessionId(salt []byte, plainPassword string) []byte {
	hashedPass := argon2.IDKey([]byte(plainPassword), []byte(salt), 1, 64*1024, 4, 32)
	return append(salt, hashedPass...)
}


func EncodingSessionToTarantool(sess *models.Session) []interface{}{
	return []interface{}{sess.SessionID, sess.UserId, sess.Executor}
}


func DecodingTarantoolToSession(data []interface{}) *models.Session{
	s := &models.Session{}
	s.SessionID, _ = data[0].(string)
	s.UserId, _ = data[1].(uint64)
	s.Executor, _ = data[2].(bool)
	return s
}
