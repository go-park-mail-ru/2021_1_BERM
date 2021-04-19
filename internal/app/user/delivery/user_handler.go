package delivery

import (
	"encoding/json"
	"ff/internal/app/models"
	"ff/internal/app/user"
	"github.com/gorilla/csrf"
	"github.com/gorilla/mux"
	"net/http"
	"ff/internal/app/server/errors"
	"strconv"
)

type ctxKey uint8

const (
	ctxKeySession ctxKey = iota
	ctxKeyReqID   ctxKey = 1
)


type UserHandler struct {
	userUseCase user.UserUseCase
}

func (u *UserHandler) handleProfile(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	userReq := &models.User{}
	reqId := r.Context().Value(ctxKeyReqID).(uint64)
	if err := json.NewDecoder(r.Body).Decode(u); err != nil {
		u.error(w, reqId, InvalidJSON) //Bad json
		return
	}
	err := u.userUseCase.Create(userReq)
	if err != nil {
		u.error(w, reqId, New(err))
		return
	}
	cookies, err := u.createCookies(userReq)
	if err != nil {
		u.error(w, reqId, New(err)) //ошибка создания сессии
		return
	}
	for _, cookie := range cookies {
		http.SetCookie(w, &cookie)
	}
	u.respond(w, reqId, http.StatusCreated, userReq)
}

func (u *UserHandler) handleLogin(w http.ResponseWriter, r *http.Request) {
	userReq := &models.User{}
	reqId := r.Context().Value(ctxKeyReqID).(uint64)
	if err := json.NewDecoder(r.Body).Decode(u); err != nil {
		u.error(w, reqId, InvalidJSON) //Bad json
		return
	}
	userReq, err := u.userUseCase.UserVerification(userReq.Email, userReq.Password)
	if err != nil {
		u.error(w, reqId, New(err)) //Unauthorized
		return
	}
	cookies, err := u.createCookies(userReq)
	if err != nil {
		u.error(w, reqId, New(err)) // ошибка создания сессии
		return
	}
	for _, cookie := range cookies {
		http.SetCookie(w, &cookie)
	}
	u.respond(w, reqId, http.StatusOK, userReq)
}

func (u *UserHandler) handleChangeProfile(w http.ResponseWriter, r *http.Request) {
	reqId := r.Context().Value(ctxKeyReqID).(uint64)
	params := mux.Vars(r)
	id, err := strconv.ParseUint(params["id"], 10, 64)
	if err != nil {
		u.error(w, reqId, New(err)) //Bad json
		return
	}
	userReq := &models.User{}
	if err := json.NewDecoder(r.Body).Decode(userReq); err != nil {
		u.error(w, reqId, InvalidJSON) //Bad json
		return
	}
	userCookieID := r.Context().Value(ctxKeySession).(*models.Session).UserID
	if userCookieID != id {
		u.error(w, reqId, New(tarantoolcache.NotAuthorized))
		return
	}
	userReq.ID = id
	userReq, err = u.userUseCase.ChangeUser(*userReq)
	if err != nil {
		u.error(w, reqId, New(err))
		return
	}
	u.respond(w, reqId, http.StatusOK, userReq)
}

func (u *UserHandler) handleGetProfile(w http.ResponseWriter, r *http.Request) {
	reqId := r.Context().Value(ctxKeyReqID).(uint64)
	params := mux.Vars(r)
	id, err := strconv.ParseUint(params["id"], 10, 64)
	if err != nil {
		u.error(w, reqId, New(err))
		return
	}
	userReq, err := u.userUseCase.FindByID(id)
	if err != nil {
		u.error(w, reqId, New(err))
		return
	}
	u.respond(w, reqId, http.StatusOK, userReq)
}

func (u *UserHandler) handleCheckAuthorized(w http.ResponseWriter, r *http.Request) {
	reqId := r.Context().Value(ctxKeyReqID).(uint64)
	w.Header().Set("X-CSRF-Token", csrf.Token(r))
	session := r.Context().Value(ctxKeySession).(*models.Session)
	u.respond(w, reqId, http.StatusOK, session)
}

func (u *UserHandler) handleAddSpecialize(w http.ResponseWriter, r *http.Request) {
	reqId := r.Context().Value(ctxKeyReqID).(uint64)
	id := r.Context().Value(ctxKeySession).(*models.Session).UserID
	specialize := &models.Specialize{}
	if err := json.NewDecoder(r.Body).Decode(specialize); err != nil {
		u.error(w, reqId, InvalidJSON)
		return
	}

	if err := u.userUseCase.AddSpecialize(specialize.Name, id); err != nil {
		u.error(w, reqId, New(err))
		return
	}
	var emptyInterface interface{}
	u.respond(w, reqId, http.StatusCreated, emptyInterface)
}

func (u *UserHandler) handleDelSpecialize(w http.ResponseWriter, r *http.Request) {
	reqId := r.Context().Value(ctxKeyReqID).(uint64)
	userID := r.Context().Value(ctxKeySession).(*models.Session).UserID
	specialize := &models.Specialize{}
	if err := json.NewDecoder(r.Body).Decode(specialize); err != nil {
		u.error(w, reqId, InvalidJSON)
		return
	}
	if err := u.userUseCase.DelSpecialize(specialize.Name, userID); err != nil {
		u.error(w, reqId, New(err))
		return
	}
	var emptyInterface interface{}
	u.respond(w, reqId, http.StatusCreated, emptyInterface)
}

func (u *UserHandler) handlePutAvatar(w http.ResponseWriter, r *http.Request) {
	reqId := r.Context().Value(ctxKeyReqID).(uint64)
	defer r.Body.Close()
	userReq := &models.User{}
	err := json.NewDecoder(r.Body).Decode(userReq)
	userReq.ID = r.Context().Value(ctxKeySession).(*models.Session).UserID
	if err != nil {
		u.error(w, reqId, InvalidJSON)
		return
	}

	u, err = u.userUseCase.SetImage(u, []byte(userReq.Img))
	if err != nil {
		u.error(w, reqId, New(err))
		return
	}
	u.respond(w, reqId, http.StatusOK, userReq)

}
