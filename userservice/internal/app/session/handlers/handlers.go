package handlers

import (
	"context"
	"net/http"
	usecase2 "user/internal/app/session/usecase"
	"user/pkg/httputils"
)
const (
	ctxUserInfo uint8 = 2
	ctxKeyReqID uint8 = 1
)

type MidleWhare struct {
	sessionUseCase usecase2.UseCase
}

func New(sessionUseCase usecase2.UseCase) *MidleWhare {
	return &MidleWhare{
		sessionUseCase: sessionUseCase,
	}
}

func (m *MidleWhare) CheckSession(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		reqID := r.Context().Value(ctxKeyReqID).(uint64)

		sessionID, err := r.Cookie("sessionID")
		if err != nil {
			httputils.RespondError(w, r, reqID, err,)
			return
		}

		u, err := m.sessionUseCase.Check(sessionID.Value, context.Background())
		if err != nil {
			httputils.RespondError(w, r, reqID, err,)
			return
		}
		ctx := context.WithValue(r.Context(), ctxUserInfo, u)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
