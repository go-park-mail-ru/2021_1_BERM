package handlers

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
	"post/internal/app/models"
	orderUseCase "post/internal/app/order/usecase"
	"post/pkg/httputils"
	"strconv"
)


const (
	ctxKeySession uint8 = 3
	ctxKeyReqID   uint8 = 1
	ctxUserInfo   uint8  = 2
)

type Handlers struct {
	useCase orderUseCase.UseCase
}

func NewHandler(useCase orderUseCase.UseCase) *Handlers {
	return &Handlers{
		useCase: useCase,
	}
}

func (h *Handlers) CreateOrder(w http.ResponseWriter, r *http.Request) {
	reqID := r.Context().Value(ctxKeyReqID).(uint64)
	id := r.Context().Value(ctxUserInfo).(uint64)
	o := &models.Order{}
	if err := json.NewDecoder(r.Body).Decode(o); err != nil {
		httputils.RespondError(w, reqID, err)
		return
	}
	o.CustomerID = id
	o, err := h.useCase.Create(*o)
	if err != nil {
		httputils.RespondError(w, reqID, err)
		return
	}
	httputils.Respond(w, reqID, http.StatusCreated, o)
}

func (h *Handlers) GetActualOrder(w http.ResponseWriter, r *http.Request) {
	reqID := r.Context().Value(ctxKeyReqID).(uint64)
	o, err := h.useCase.GetActualOrders()
	if err != nil {
		httputils.RespondError(w, reqID, err)
		return
	}
	httputils.Respond(w, reqID, http.StatusOK, o)
}

func (h *Handlers) GetOrder(w http.ResponseWriter, r *http.Request) {
	reqID := r.Context().Value(ctxKeyReqID).(uint64)
	params := mux.Vars(r)
	id, err := strconv.ParseUint(params["id"], 10, 64)
	if err != nil {
		httputils.RespondError(w, reqID, err)
		return
	}
	o, err := h.useCase.FindByID(id)
	if err != nil {
		httputils.RespondError(w, reqID, err)
		return
	}
	httputils.Respond(w, reqID, http.StatusOK, o)
}

func (h *Handlers) ChangeOrder(w http.ResponseWriter, r *http.Request) {
	reqID := r.Context().Value(ctxKeyReqID).(uint64)
	order := models.Order{}
	if err := json.NewDecoder(r.Body).Decode(&order); err != nil {
		httputils.RespondError(w, reqID, err)
		return
	}
	params := mux.Vars(r)
	var err error
	order.ID, err = strconv.ParseUint(params["id"], 10, 64)
	if err != nil {
		httputils.RespondError(w, reqID, err)
		return
	}
	order, err = h.useCase.ChangeOrder(order)
	if err != nil {
		httputils.RespondError(w, reqID, err)
		return
	}
	httputils.Respond(w, reqID, http.StatusOK, order)
}

func (h *Handlers) DeleteOrder(w http.ResponseWriter, r *http.Request) {
	reqID := r.Context().Value(ctxKeyReqID).(uint64)
	params := mux.Vars(r)
	id, err := strconv.ParseUint(params["id"], 10, 64)
	if err != nil {
		httputils.RespondError(w, reqID, err)
		return
	}
	err = h.useCase.DeleteOrder(id)
	if err != nil {
		httputils.RespondError(w, reqID, err)
		return
	}
	var emptyInterface interface{}
	httputils.Respond(w, reqID, http.StatusOK, emptyInterface)
}

func (h *Handlers) SelectExecutor(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	reqID := r.Context().Value(ctxKeyReqID).(uint64)
	order := models.Order{}
	if err := json.NewDecoder(r.Body).Decode(&order); err != nil {
		httputils.RespondError(w, reqID, err)
		return
	}
	var err error
	order.ID, err = strconv.ParseUint(params["id"], 10, 64)
	if err != nil {
		httputils.RespondError(w, reqID, err)
		return
	}
	err = h.useCase.SelectExecutor(order)
	if err != nil {
		httputils.RespondError(w, reqID, err)
		return
	}
	httputils.Respond(w, reqID, http.StatusOK, order)
}

func (h *Handlers) DeleteExecutor(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	reqID := r.Context().Value(ctxKeyReqID).(uint64)

	order := models.Order{}

	var err error
	order.ID, err = strconv.ParseUint(params["id"], 10, 64)
	if err != nil {
		httputils.RespondError(w, reqID, err)
		return
	}
	err = h.useCase.DeleteExecutor(order)
	if err != nil {
		httputils.RespondError(w, reqID, err)
		return
	}
	var emptyInterface interface{}
	httputils.Respond(w, reqID, http.StatusOK, emptyInterface)
}

func (h *Handlers) GetAllUserOrders(w http.ResponseWriter, r *http.Request) {
	reqID := r.Context().Value(ctxKeyReqID).(uint64)

	params := mux.Vars(r)
	userID, err := strconv.ParseUint(params["id"], 10, 64)
	if err != nil {
		httputils.RespondError(w, reqID, err)
		return
	}

	o, err := h.useCase.FindByUserID(userID)
	if err != nil {
		httputils.RespondError(w, reqID, err)
		return
	}
	httputils.Respond(w, reqID, http.StatusOK, o)
}
