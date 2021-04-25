package handlers

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
	"post/internal/app/httputils"
	"post/internal/app/models"
	orderUseCase "post/internal/app/order/usecase"
	"strconv"
)

type ctxKey uint8

const (
	ctxKeySession ctxKey = iota
	ctxKeyReqID   ctxKey = 1
)

type Handlers struct {
	useCase orderUseCase.UseCase
}

func NewHandler(useCase orderUseCase.UseCase) *Handlers {
	return &Handlers{
		useCase: useCase,
	}
}

func(h *Handlers) CreateOrder(w http.ResponseWriter, r *http.Request) {
	reqId := r.Context().Value(ctxKeyReqID).(uint64)
	//TODO: не оч понятно как вытащить айдишник
	id := r.Context().Value(ctxKeySession).(*models.Session).UserID
	o := &models.Order{}
	if err := json.NewDecoder(r.Body).Decode(o); err != nil {
		httputils.RespondError(w, reqId, InvalidJSON)
		return
	}
	o.CustomerID = id
	var err error
	o, err = h.useCase.Create(*o)
	if err != nil {
		//TODO: ошибка
		httputils.RespondError(w, reqId, New(err))
		return
	}
	httputils.Respond(w, reqId, http.StatusCreated, o)
}

func(h *Handlers) GetActualOrder(w http.ResponseWriter, r *http.Request) {
	reqID := r.Context().Value(ctxKeyReqID).(uint64)
	o, err := h.useCase.GetActualOrders()
	if err != nil {
		//TODO: ошибка
		httputils.RespondError(w, reqID, New(err))
		return
	}
	httputils.Respond(w, reqID, http.StatusOK, o)
}

func(h *Handlers) GetOrder(w http.ResponseWriter, r *http.Request) {
	reqID := r.Context().Value(ctxKeyReqID).(uint64)
	params := mux.Vars(r)
	id, err := strconv.ParseUint(params["id"], 10, 64)
	if err != nil {
		//TODO: ошибка
		httputils.RespondError(w, reqID, New(err))
		return
	}
	o, err := h.useCase.FindByID(id)
	if err != nil {
		//TODO: ошибка
		httputils.RespondError(w, reqID, New(err))
		return
	}
	httputils.Respond(w, reqID, http.StatusOK, o)
}

func(h *Handlers) SelectExecutor(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	reqID := r.Context().Value(ctxKeyReqID).(uint64)

	order := models.Order{}
	if err := json.NewDecoder(r.Body).Decode(&order); err != nil {
		httputils.RespondError(w, reqID, InvalidJSON)
		return
	}
	var err error
	order.ID, err = strconv.ParseUint(params["id"], 10, 64)
	if err != nil {
		httputils.RespondError(w, http.StatusBadRequest, InvalidJSON)
		return
	}
	err = h.useCase.SelectExecutor(order)
	if err != nil {
		//TODO: ошибка
		httputils.RespondError(w, http.StatusInternalServerError, New(err))
		return
	}
	httputils.Respond(w, reqID, http.StatusOK, order)
}

func(h *Handlers) DeleteExecutor(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	reqId := r.Context().Value(ctxKeyReqID).(uint64)

	order := models.Order{}
	var err error
	order.ID, err = strconv.ParseUint(params["id"], 10, 64)
	if err != nil {
		httputils.RespondError(w, http.StatusBadRequest, InvalidJSON)
		return
	}
	err = h.useCase.DeleteExecutor(order)
	if err != nil {
		//TODO: ошибка
		httputils.RespondError(w, http.StatusInternalServerError, New(err))
		return
	}
	var emptyInterface interface{}
	httputils.Respond(w, reqId, http.StatusOK, emptyInterface)
}

func(h *Handlers) GetAllUserOrders(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	reqID := r.Context().Value(ctxKeyReqID).(uint64)
	userID, err := strconv.ParseUint(params["id"], 10, 64)
	if err != nil {
		httputils.RespondError(w, http.StatusBadRequest, InvalidJSON)
		return
	}

	o, err := h.useCase.FindByUserID(userID)
	if err != nil {
		//TODO: ошибка
		httputils.RespondError(w, http.StatusNotFound, New(err))
		return
	}
	httputils.Respond(w, reqID, http.StatusOK, o)
}