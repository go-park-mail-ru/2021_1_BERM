package handlers

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/pkg/errors"
	"net/http"
	"strconv"
	"user/Error"
	"user/internal/app/models"
	"user/internal/app/user/usecase"
	"user/pkg/httputils"
)

type Handlers struct {
	userUseCase usecase.UseCase
}

func New(userUseCase usecase.UseCase) *Handlers {
	return &Handlers{
		userUseCase: userUseCase,
	}
}

func (h *Handlers) ChangeProfile(w http.ResponseWriter, r *http.Request) {
	reqId, err := strconv.ParseUint(r.Header.Get("X_Request_Id"), 10, 64)
	if err != nil {
		httputils.RespondError(w, reqId, err, http.StatusInternalServerError)
		return
	}
	params := mux.Vars(r)
	id, err := strconv.ParseUint(params["id"], 10, 64)
	if err != nil {
		httputils.RespondError(w, reqId, err, http.StatusInternalServerError)
		return
	}
	u := &models.ChangeUser{}
	if err := json.NewDecoder(r.Body).Decode(u); err != nil {
		httputils.RespondError(w, reqId, err, http.StatusInternalServerError)
		return
	}
	u.ID = id
	response, err := h.userUseCase.Change(*u)
	if err != nil {
		httpErr := &Error.Error{}
		errors.As(err, httpErr)
		if httpErr.InternalError {
			httputils.RespondError(w, reqId, err, http.StatusInternalServerError)
		} else {
			httputils.RespondError(w, reqId, err, http.StatusBadRequest)
		}
		return
	}
	httputils.Respond(w, reqId, http.StatusOK, response)
}

func (h *Handlers) GetUserInfo(w http.ResponseWriter, r *http.Request) {
	reqID, err := strconv.ParseUint(r.Header.Get("X_Request_Id"), 10, 64)
	if err != nil {
		httputils.RespondError(w, reqID, err, http.StatusInternalServerError)
		return
	}
	params := mux.Vars(r)
	ID, err := strconv.ParseUint(params["ID"], 10, 64)
	if err != nil {
		httputils.RespondError(w, reqID, err, http.StatusInternalServerError)
		return
	}
	u, err := h.userUseCase.GetById(ID)
	if err != nil {
		httpErr := &Error.Error{}
		errors.As(err, httpErr)
		if httpErr.InternalError {
			httputils.RespondError(w, reqID, err, http.StatusInternalServerError)
		} else {
			httputils.RespondError(w, reqID, err, http.StatusBadRequest)
		}
		return
	}
	httputils.Respond(w, reqID, http.StatusOK, u)
}
