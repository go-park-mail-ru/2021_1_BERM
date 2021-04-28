package handlers

import (
	"encoding/json"
	"errors"
	"github.com/gorilla/mux"
	"net/http"
	"post/internal/app/models"
	orderUseCase "post/internal/app/order/usecase"
	"post/pkg/Error"
	"post/pkg/httputils"
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
	reqID, err := strconv.ParseUint(r.Header.Get("X_Request_Id"), 10, 64)
	if err != nil{
		httputils.RespondError(w, reqID, err, http.StatusInternalServerError)
	}
	id, err := strconv.ParseUint(r.Header.Get("X_Id"), 10, 64)
	if err != nil{
		httputils.RespondError(w, reqID, err, http.StatusInternalServerError)
	}
	if err != nil{
		httputils.RespondError(w, reqID, err, http.StatusInternalServerError)
	}
	o := &models.Order{}
	if err = json.NewDecoder(r.Body).Decode(o); err != nil {
		httputils.RespondError(w, reqID, err, http.StatusInternalServerError)
		return
	}
	o.CustomerID = id
	o, err = h.useCase.Create(*o)
	if err != nil {
		httpErr := &Error.Error{}
		errors.As(err, &httpErr)
		if httpErr.InternalError {
			httputils.RespondError(w, reqID, err, http.StatusInternalServerError)
		} else {
			httputils.RespondError(w, reqID, err, http.StatusBadRequest)
		}
		return
	}
	httputils.Respond(w, reqID, http.StatusCreated, o)
}

func(h *Handlers) GetActualOrder(w http.ResponseWriter, r *http.Request) {
	reqID, err := strconv.ParseUint(r.Header.Get("X_Request_Id"), 10, 64)
	if err != nil{
		httputils.RespondError(w, reqID, err, http.StatusInternalServerError)
	}
	o, err := h.useCase.GetActualOrders()
	if err != nil {
		httpErr := &Error.Error{}
		errors.As(err, &httpErr)
		if httpErr.InternalError {
			httputils.RespondError(w, reqID, err, http.StatusInternalServerError)
		} else {
			httputils.RespondError(w, reqID, err, http.StatusBadRequest)
		}
		return
	}
	httputils.Respond(w, reqID, http.StatusOK, o)
}

func(h *Handlers) GetOrder(w http.ResponseWriter, r *http.Request) {
	reqID, err := strconv.ParseUint(r.Header.Get("X_Request_Id"), 10, 64)
	if err != nil{
		httputils.RespondError(w, reqID, err, http.StatusInternalServerError)
	}
	params := mux.Vars(r)
	id, err := strconv.ParseUint(params["id"], 10, 64)
	if err != nil {
		httputils.RespondError(w, reqID, err, http.StatusInternalServerError)
		return
	}
	o, err := h.useCase.FindByID(id)
	if err != nil {
		httpErr := &Error.Error{}
		errors.As(err, &httpErr)
		if httpErr.InternalError {
			httputils.RespondError(w, reqID, err, http.StatusInternalServerError)
		} else {
			httputils.RespondError(w, reqID, err, http.StatusBadRequest)
		}
		return
	}
	httputils.Respond(w, reqID, http.StatusOK, o)
}

func (h *Handlers) ChangeOrder(w http.ResponseWriter, r *http.Request) {
	reqID, err := strconv.ParseUint(r.Header.Get("X_Request_Id"), 10, 64)
	if err != nil{
		httputils.RespondError(w, reqID, err, http.StatusInternalServerError)
	}
	order := models.Order{}
	if err = json.NewDecoder(r.Body).Decode(&order); err != nil {
		httputils.RespondError(w, reqID, err, http.StatusInternalServerError)
		return
	}
	params := mux.Vars(r)
	order.ID, err = strconv.ParseUint(params["id"], 10, 64)
	if err != nil {
		httputils.RespondError(w, reqID, err, http.StatusInternalServerError)
		return
	}
	order, err = h.useCase.ChangeOrder(order)
	if err != nil {
		httputils.RespondError(w, reqID, err, http.StatusInternalServerError)
		return
	}
	httputils.Respond(w, reqID, http.StatusOK, order)
}

func (h *Handlers) DeleteOrder(w http.ResponseWriter, r *http.Request) {
	reqID, err := strconv.ParseUint(r.Header.Get("X_Request_Id"), 10, 64)
	if err != nil{
		httputils.RespondError(w, reqID, err, http.StatusInternalServerError)
	}
	params := mux.Vars(r)
	id, err := strconv.ParseUint(params["id"], 10, 64)
	if err != nil {
		httputils.RespondError(w, reqID, err, http.StatusInternalServerError)
		return
	}
	err = h.useCase.DeleteOrder(id)
	if err != nil {
		httputils.RespondError(w, reqID, err, http.StatusInternalServerError)
		return
	}
	var emptyInterface interface{}
	httputils.Respond(w, reqID, http.StatusOK, emptyInterface)
}

func(h *Handlers) SelectExecutor(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	reqID, err := strconv.ParseUint(r.Header.Get("X_Request_Id"), 10, 64)
	if err != nil{
		httputils.RespondError(w, reqID, err, http.StatusInternalServerError)
	}

	order := models.Order{}
	if err = json.NewDecoder(r.Body).Decode(&order); err != nil {
		httputils.RespondError(w, reqID, err, http.StatusInternalServerError)
		return
	}
	order.ID, err = strconv.ParseUint(params["id"], 10, 64)
	if err != nil {
		httputils.RespondError(w, reqID, err, http.StatusInternalServerError)
		return
	}
	err = h.useCase.SelectExecutor(order)
	if err != nil {
		httpErr := &Error.Error{}
		errors.As(err, &httpErr)
		if httpErr.InternalError {
			httputils.RespondError(w, reqID, err, http.StatusInternalServerError)
		} else {
			httputils.RespondError(w, reqID, err, http.StatusBadRequest)
		}
		return
	}
	httputils.Respond(w, reqID, http.StatusOK, order)
}

func(h *Handlers) DeleteExecutor(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	reqID, err := strconv.ParseUint(r.Header.Get("X_Request_Id"), 10, 64)
	if err != nil{
		httputils.RespondError(w, reqID, err, http.StatusInternalServerError)
	}

	order := models.Order{}

	order.ID, err = strconv.ParseUint(params["id"], 10, 64)
	if err != nil {
		httputils.RespondError(w, reqID, err, http.StatusInternalServerError)
		return
	}
	err = h.useCase.DeleteExecutor(order)
	if err != nil {
		httpErr := &Error.Error{}
		errors.As(err, &httpErr)
		if httpErr.InternalError {
			httputils.RespondError(w, reqID, err, http.StatusInternalServerError)
		} else {
			httputils.RespondError(w, reqID, err, http.StatusBadRequest)
		}
		return
	}
	var emptyInterface interface{}
	httputils.Respond(w, reqID, http.StatusOK, emptyInterface)
}

func(h *Handlers) GetAllUserOrders(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	reqID, err := strconv.ParseUint(r.Header.Get("X_Request_Id"), 10, 64)
	if err != nil{
		httputils.RespondError(w, reqID, err, http.StatusInternalServerError)
	}
	userID, err := strconv.ParseUint(params["id"], 10, 64)
	if err != nil {
		httputils.RespondError(w, reqID, err, http.StatusInternalServerError)
		return
	}

	o, err := h.useCase.FindByUserID(userID)
	if err != nil {
		httpErr := &Error.Error{}
		errors.As(err, &httpErr)
		if httpErr.InternalError {
			httputils.RespondError(w, reqID, err, http.StatusInternalServerError)
		} else {
			httputils.RespondError(w, reqID, err, http.StatusBadRequest)
		}
		return
	}
	httputils.Respond(w, reqID, http.StatusOK, o)
}
