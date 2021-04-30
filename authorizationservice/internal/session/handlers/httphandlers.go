package handlers

import (
	"authorizationservice/internal/session/usecase"
	"authorizationservice/pkg/utils"
	"context"
	"github.com/gorilla/csrf"
	"math/rand"
	"net/http"
	"time"
)

type Handler struct {
	sessionUseCase usecase.UseCase
}

func New(sessionUseCase usecase.UseCase) *Handler {
	return &Handler{
		sessionUseCase: sessionUseCase,
	}
}
func (h *Handler) CheckLogin(w http.ResponseWriter, r *http.Request) {
	rand.Seed(time.Now().UnixNano())
	reqId := rand.Uint64()
	ctx := context.WithValue(context.Background(), "ReqID", reqId)
	cookie, err := r.Cookie("SessionID")
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	w.Header().Set("X-CSRF-Token", csrf.Token(r))
	sessionId := cookie.Value
	session, err := h.sessionUseCase.Get(sessionId, ctx)
	utils.Respond(w, reqId, http.StatusAccepted, session)
}

func (h *Handler) LogOut(w http.ResponseWriter, r *http.Request) {
	rand.Seed(time.Now().UnixNano())
	reqId := rand.Uint64()
	cookies := r.Cookies()
	utils.RemoveCookies(cookies)
	utils.Respond(w, reqId, http.StatusAccepted, nil)
}
