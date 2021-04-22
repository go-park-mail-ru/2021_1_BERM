package handlers

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
	"post/internal/app/httputils"
	"post/internal/app/models"
	vacancyUseCase "post/internal/app/vacancy/usecase"
	"strconv"
)

var (
	InvalidJSON = &Error{
		Err:  errors.New("Invalid json. "),
		Code: http.StatusBadRequest,
		Type: TypeExternal,
		Field: map[string]interface{}{
			"error": "Invalid json",
		},
	}

	InvalidCookies = &Error{
		Err:  errors.New("Invalid cookie.\n"),
		Code: http.StatusBadRequest,
		Type: TypeExternal,
	}
)

const (
	ctxKeySession uint8 = iota
	ctxKeyReqID   uint8 = 1
)

type Handlers struct {
	useCase vacancyUseCase.UseCase
}

func NewHandler(useCase vacancyUseCase.UseCase) *Handlers {
	return &Handlers{
		useCase: useCase,
	}
}

func(h *Handlers) CreateVacancy(w http.ResponseWriter, r *http.Request) {
	reqID := r.Context().Value(ctxKeyReqID).(uint64)
	//TODO: откуда брать этот айдишник?
	id := r.Context().Value(ctxKeySession).(*models.Session).UserID
	v := &models.Vacancy{
		UserID: id,
	}
	if err := json.NewDecoder(r.Body).Decode(v); err != nil {
		httputils.RespondError(w, reqID, InvalidJSON)
		return
	}
	var err error
	if v, err = h.useCase.Create(*v); err != nil {
		//TODO: ошибка
		httputils.RespondError(w, reqID, New(err))
	}
	httputils.Respond(w, reqID, http.StatusCreated, v)
}

func(h *Handlers) GetVacancy(w http.ResponseWriter, r *http.Request) {
	reqId := r.Context().Value(ctxKeyReqID).(uint64)
	params := mux.Vars(r)
	id, err := strconv.ParseUint(params["id"], 10, 64)
	if err != nil {
		httputils.RespondError(w, reqId, InvalidJSON)
		return
	}
	v, err := h.useCase.FindByID(id)
	if err != nil {
		//TODO: ошибка
		httputils.RespondError(w, reqId, New(err))
		return
	}
	httputils.Respond(w, reqId, http.StatusOK, v)
}
