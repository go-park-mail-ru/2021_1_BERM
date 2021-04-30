package handlers

import (
	"authorizationservice/internal/models"
	"authorizationservice/internal/profile/usecase"
	session "authorizationservice/internal/session/usecase"
	"authorizationservice/pkg/utils"
	"context"
	"encoding/json"
	"math/rand"
	"net/http"
	"time"
)

type Handler struct {
	sessionUseCase session.UseCase
	profileUseCase usecase.UseCase
}

func New(sessionUseCase session.UseCase, profileUseCase usecase.UseCase) *Handler{
	return &Handler{
		sessionUseCase: sessionUseCase,
		profileUseCase: profileUseCase,
	}
}
func (h* Handler) RegistrationProfile(w http.ResponseWriter, r *http.Request){
	rand.Seed(time.Now().UnixNano())
	reqId := rand.Uint64()
	ctx := context.WithValue(context.Background(), "ReqID", reqId)
	u := &models.NewUser{}
	if err := json.NewDecoder(r.Body).Decode(u); err != nil {
		utils.RespondError(w, reqId, err, http.StatusBadRequest)
		return
	}
	resp, err := h.profileUseCase.Create(*u, ctx)
	if err != nil{
		utils.RespondError(w, reqId, err, http.StatusBadRequest)
		return
	}
	sess, err := h.sessionUseCase.Create(resp.ID, resp.Executor, ctx);
	if err != nil{
		utils.RespondError(w, reqId, err, http.StatusBadRequest)
		return
	}
	utils.CreateCookie(sess)
	utils.Respond(w, reqId, http.StatusAccepted, resp)
}


func (h* Handler) AuthorisationProfile(w http.ResponseWriter, r *http.Request){
	rand.Seed(time.Now().UnixNano())
	reqId := rand.Uint64()
	ctx := context.WithValue(context.Background(), "ReqID", reqId)
	u := &models.LoginUser{}
	if err := json.NewDecoder(r.Body).Decode(u); err != nil {
		utils.RespondError(w, reqId, err, http.StatusBadRequest)
		return
	}
	resp, err := h.profileUseCase.Authentication(u.Email, u.Password, ctx)
	if err != nil{
		utils.RespondError(w, reqId, err, http.StatusBadRequest)
		return
	}
	sess, err := h.sessionUseCase.Create(resp.ID, resp.Executor, ctx);
	if err != nil{
		utils.RespondError(w, reqId, err, http.StatusBadRequest)
		return
	}
	utils.CreateCookie(sess)
	utils.Respond(w, reqId, http.StatusAccepted, resp)
}