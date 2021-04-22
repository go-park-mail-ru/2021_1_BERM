package handlers

import (
	"FL_2/model"
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
	"post/internal/app/httputils"
	"post/internal/app/models"
	responseUseCase "post/internal/app/response/usecase"
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
	useCase responseUseCase.UseCase
}

func NewHandler(useCase responseUseCase.UseCase) *Handlers {
	return &Handlers{
		useCase: useCase,
	}
}

func(h *Handlers) CreatePostResponse(w http.ResponseWriter, r *http.Request) {
	reqID := r.Context().Value(ctxKeyReqID).(uint64)
	response := &models.Response{}
	if err := json.NewDecoder(r.Body).Decode(response); err != nil {
		httputils.RespondError(w, http.StatusBadRequest, InvalidJSON) //Bad json
		return
	}
	params := mux.Vars(r)
	id, err := strconv.ParseUint(params["id"], 10, 64)
	if err != nil {
		//TODO: ошибка
		httputils.RespondError(w, http.StatusBadRequest, New(err)) //Bad json
		return
	}
	response.PostID = id
	response, err = h.useCase.Create(*response)
	if err != nil {
		//TODO: ошибка
		httputils.RespondError(w, http.StatusBadRequest, New(err))
		return
	}
	httputils.Respond(w, reqID, http.StatusCreated, response)
}

func(h *Handlers) GetAllPostResponses(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	reqId := r.Context().Value(ctxKeyReqID).(uint64)
	id, err := strconv.ParseUint(params["id"], 10, 64)
	if err != nil {
		//TODO: ошибка
		httputils.RespondError(w, reqId, New(err)) //Bad json
		return
	}
	//TODO: откуда-то брать инфу о том ордер это или вакансия
	responses, err := h.useCase.FindByPostID(id)
	if err != nil {
		//TODO: ошибка
		httputils.RespondError(w, reqId, New(err)) //Bad json
		return
	}

	httputils.Respond(w, reqId, http.StatusOK, responses)
}

func(h *Handlers) ChangePostResponse(w http.ResponseWriter, r *http.Request) {
	response := &models.Response{}
	reqID := r.Context().Value(ctxKeyReqID).(uint64)

	if err := json.NewDecoder(r.Body).Decode(response); err != nil {
		httputils.RespondError(w, reqID, InvalidJSON)
		return
	}
	params := mux.Vars(r)
	var err error
	response.PostID, err = strconv.ParseUint(params["id"], 10, 64)
	if err != nil {
		//TODO: ошибка
		httputils.RespondError(w, reqID, New(err))
		return
	}
	//TODO: не понятно откуда брать инфу об айдишнике
	response.UserID = r.Context().Value(ctxKeySession).(*model.Session).UserID
	responses, err := h.useCase.Change(*response)

	if err != nil {
		//TODO: ошибка
		httputils.RespondError(w, reqID, New(err))
		return
	}

	httputils.Respond(w, reqID, http.StatusOK, responses)
}

func(h *Handlers) DelPostResponse(w http.ResponseWriter, r *http.Request) {
	response := &models.Response{}
	params := mux.Vars(r)
	reqID := r.Context().Value(ctxKeyReqID).(uint64)
	var err error
	response.PostID, err = strconv.ParseUint(params["id"], 10, 64)
	if err != nil {
		//TODO: ошибка
		httputils.RespondError(w, reqID, New(err)) //Bad json
		return
	}
	//TODO: откуда-то брать этот чертов айдишник
	response.UserID = r.Context().Value(ctxKeySession).(*model.Session).UserID
	err = h.useCase.Delete(*response)

	if err != nil {
		//TODO: ошибка
		httputils.RespondError(w, reqID, New(err)) //Bad json
		return
	}
	var emptyInterface interface{}

	httputils.Respond(w, reqID, http.StatusOK, emptyInterface)
}
