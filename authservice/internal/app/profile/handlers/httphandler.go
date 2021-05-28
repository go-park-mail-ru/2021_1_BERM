package handlers

import (
	models2 "authorizationservice/internal/app/models"
	profile2 "authorizationservice/internal/app/profile"
	"authorizationservice/internal/app/session"
	"authorizationservice/pkg/utils"
	"context"
	"io/ioutil"
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
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		utils.RespondError(w, r, reqId, err)
		return
	}
	if err := u.UnmarshalJSON(body); err != nil {
		utils.RespondError(w, r, reqId, err)
		return
	}
	resp, err := h.profileUseCase.Create(*u, context.Background())
	if err != nil {
		utils.RespondError(w, r, reqId, err)
		return
	}
	sess, err := h.sessionUseCase.Create(resp.ID, resp.Executor, context.Background())
	if err != nil {
		utils.RespondError(w, r, reqId, err)
		return
	}
	utils.CreateCookie(sess, w)
	result, err := resp.MarshalJSON()
	if err != nil {
		utils.RespondError(w, r, reqId, err)
		return
	}
	utils.Respond(w, r, reqId, http.StatusAccepted, result)
}

func (h *Handler) AuthorisationProfile(w http.ResponseWriter, r *http.Request) {
	rand.Seed(time.Now().UnixNano())
	reqId := rand.Uint64()
	u := &models2.LoginUser{}
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		utils.RespondError(w, r, reqId, err)
		return
	}
	if err := u.UnmarshalJSON(body); err != nil {
		utils.RespondError(w, r, reqId, err)
		return
	}
	resp, err := h.profileUseCase.Authentication(u.Email, u.Password, r.Context())
	if err != nil {
		utils.RespondError(w, r, reqId, err)
		return
	}
	sess, err := h.sessionUseCase.Create(resp.ID, resp.Executor, r.Context())
	if err != nil {
		utils.RespondError(w, r, reqId, err)
		return
	}
	utils.CreateCookie(sess, w)
	result, err := resp.MarshalJSON()
	if err != nil {
		utils.RespondError(w, r, reqId, err)
		return
	}
	utils.Respond(w, r, reqId, http.StatusAccepted, result)
}
