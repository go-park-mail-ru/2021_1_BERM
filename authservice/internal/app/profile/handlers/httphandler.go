package handlers

import (
	models2 "authorizationservice/internal/app/models"
	profile2 "authorizationservice/internal/app/profile"
	"authorizationservice/internal/app/session"
	"authorizationservice/pkg/utils"
	"context"
	"encoding/json"
	"math/rand"
	"net/http"
	"time"
)

type Handler struct {
	sessionUseCase session.UseCase
	profileUseCase profile2.UseCase
}

func New(sessionUseCase session.UseCase, profileUseCase profile2.UseCase) *Handler {
	return &Handler{
		sessionUseCase: sessionUseCase,
		profileUseCase: profileUseCase,
	}
}

func (h *Handler) RegistrationProfile(w http.ResponseWriter, r *http.Request) {
	rand.Seed(time.Now().UnixNano())
	reqId := rand.Uint64()
	u := &models2.NewUser{}
	if err := json.NewDecoder(r.Body).Decode(u); err != nil {
		utils.RespondError(w, r, reqId, err)
		return
	}
	resp, err := h.profileUseCase.Create(*u, context.Background())
	if err != nil {
		utils.RespondError(w,r, reqId, err)
		return
	}
	sess, err := h.sessionUseCase.Create(resp.ID, resp.Executor, context.Background())
	if err != nil {
		utils.RespondError(w,r, reqId, err)
		return
	}
	utils.CreateCookie(sess, w)
	utils.Respond(w,r, reqId, http.StatusAccepted, resp)
}

func (h *Handler) AuthorisationProfile(w http.ResponseWriter, r *http.Request) {
	rand.Seed(time.Now().UnixNano())
	reqId := rand.Uint64()
	ctx := context.WithValue(context.Background(), "ReqID", reqId)
	u := &models2.LoginUser{}
	if err := json.NewDecoder(r.Body).Decode(u); err != nil {
		utils.RespondError(w, r,reqId, err)
		return
	}
	resp, err := h.profileUseCase.Authentication(u.Email, u.Password, ctx)
	if err != nil {
		utils.RespondError(w, r,reqId, err)
		return
	}
	sess, err := h.sessionUseCase.Create(resp.ID, resp.Executor, ctx)
	if err != nil {
		utils.RespondError(w, r,reqId, err)
		return
	}
	utils.CreateCookie(sess, w)
	utils.Respond(w, r,reqId, http.StatusAccepted, resp)
}
