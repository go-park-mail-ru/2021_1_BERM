package handlers

import (
	session2 "authorizationservice/internal/app/session"
	"authorizationservice/pkg/utils"
	"context"
	"github.com/gorilla/csrf"
	"math/rand"
	"net/http"
	"time"
)

type Handler struct {
	sessionUseCase session2.UseCase
}

func New(sessionUseCase session2.UseCase) *Handler {
	return &Handler{
		sessionUseCase: sessionUseCase,
	}
}
func (h *Handler) CheckLogin(w http.ResponseWriter, r *http.Request) {
	rand.Seed(time.Now().UnixNano())
	reqId := rand.Uint64()
	cookie, err := r.Cookie("sessionID")
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	w.Header().Set("X-CSRF-Token", csrf.Token(r))
	sessionId := cookie.Value
	session, err := h.sessionUseCase.Get(sessionId, context.Background())
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	result, err := session.MarshalJSON()
	if err != nil {
		utils.RespondError(w, r, reqId, err)
		return
	}
	utils.Respond(w, r, reqId, http.StatusAccepted, result)
}

func (h *Handler) LogOut(w http.ResponseWriter, r *http.Request) {
	rand.Seed(time.Now().UnixNano())
	reqId := rand.Uint64()
	cookies := r.Cookies()
	utils.RemoveCookies(cookies, w)
	utils.Respond(w, r, reqId, http.StatusAccepted, nil)
}
