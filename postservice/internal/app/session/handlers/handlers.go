package handlers

import (
	"context"
	"net/http"
	"post/internal/app/session/usecase"
	"post/pkg/httputils"
)

const (
	ctxUserInfo uint8 = 2
	ctxKeyReqID uint8 = 1
)

type MiddleWare struct {
	sessionUseCase usecase.UseCase
}

func New(sessionUseCase usecase.UseCase) *MiddleWare {
	return &MiddleWare{
		sessionUseCase: sessionUseCase,
	}
}

func (m *MiddleWare) CheckSession(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		reqID := r.Context().Value(ctxKeyReqID).(uint64)

		sessionID, err := r.Cookie("sessionID")
		if err != nil {
			httputils.RespondError(w, reqID, err,)
			return
		}

		u, err := m.sessionUseCase.Check(sessionID.Value, context.Background())
		if err != nil {
			httputils.RespondError(w, reqID, err,)
			return
		}
		ctx := context.WithValue(r.Context(), ctxUserInfo, u)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
