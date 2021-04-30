package handlers

import (
	"context"
	"net/http"
	"user/internal/session/usecase"
	"user/pkg/httputils"
)

type MidleWhare struct {
	sessionUseCase usecase.UseCase
}

func New(sessionUseCase usecase.UseCase) *MidleWhare {
	return &MidleWhare{
		sessionUseCase: sessionUseCase,
	}
}

func (m *MidleWhare) CheckSession(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		sessionID, err := r.Cookie("sessionID")
		if err != nil {
			//FIXME поправить ошибки
			httputils.Respond(w, 0, http.StatusUnauthorized, map[string]string{
				"message": "Bad cookies",
			})
		}

		u, err := m.sessionUseCase.Check(sessionID.Value, nil)
		if err != nil {
			//FIXME поправить ошибки
			httputils.Respond(w, 0, http.StatusUnauthorized, map[string]string{
				"message": "Bad cookies",
			})
		}
		ctx := context.WithValue(r.Context(), "UserInfo", u)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
