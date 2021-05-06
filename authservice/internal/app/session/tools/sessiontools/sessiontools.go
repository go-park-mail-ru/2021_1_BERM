package sessiontools

import (
	models2 "authorizationservice/internal/app/models"
	"golang.org/x/crypto/bcrypt"
	"math/rand"
	"strconv"
	"time"
)

const (
	saltLength = 8
)


type SessionTools struct {

}


func (s SessionTools)BeforeCreate(session models2.Session) (models2.Session, error) {
	session.SessionID = strconv.FormatUint(session.UserId, 10) + time.Now().String()
	salt := make([]byte, saltLength)
	_, err := rand.Read(salt)
	if err != nil {
		return models2.Session{},err
	}

	session.SessionID = string(hashSessionId(salt, session.SessionID))
	return session, nil
}

func hashSessionId(salt []byte, plainPassword string) []byte {
	//TODO: обрабатывать ошибку
	hashedPass, _ := bcrypt.GenerateFromPassword([]byte(plainPassword+string(salt)), bcrypt.MinCost)
	return hashedPass
}
func (s SessionTools)EncodingSessionToTarantool(sess *models2.Session) []interface{} {
	return []interface{}{sess.SessionID, sess.UserId, sess.Executor}
}

func (s SessionTools) DecodingTarantoolToSession(data []interface{}) *models2.Session {
	sess := &models2.Session{}
	sess.SessionID, _ = data[0].(string)
	sess.UserId, _ = data[1].(uint64)
	sess.Executor, _ = data[2].(bool)
	return sess
}
