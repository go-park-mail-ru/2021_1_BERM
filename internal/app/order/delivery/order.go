package delivery

import (
	"encoding/json"
	"ff/internal/app/models"
	"ff/internal/app/order"
	"ff/internal/app/user"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
)

type ctxKey uint8

const (
	ctxKeySession ctxKey = iota
	ctxKeyReqID   ctxKey = 1
)

type OrderHandler struct {
	orderUseCase order.OrderUseCase
}


func (o *OrderHandler) HandleGetOrder(w http.ResponseWriter, r *http.Request) {
	reqId := r.Context().Value(ctxKeyReqID).(uint64)
	params := mux.Vars(r)
	id, err := strconv.ParseUint(params["id"], 10, 64)
	if err != nil {
		o.error(w, reqId, New(err))
		return
	}
	ord, err := o.orderUseCase.FindByID(id)
	if err != nil {
		o.error(w, reqId, New(err))
		return
	}
	o.respond(w, reqId, http.StatusOK, ord)
}

func (o *OrderHandler) HandleGetActualOrder(w http.ResponseWriter, r *http.Request) {
	reqId := r.Context().Value(ctxKeyReqID).(uint64)
	ord, err := order.GetActualOrders()
	if err != nil {
		o.error(w, reqId, New(err))
		return
	}
	o.respond(w, reqId, http.StatusOK, ord)
}

func (o *OrderHandler) handleCreateOrderResponse(w http.ResponseWriter, r *http.Request) {
	response := &models.ResponseOrder{}
	reqId := r.Context().Value(ctxKeyReqID).(uint64)

	if err := json.NewDecoder(r.Body).Decode(response); err != nil {
		o.error(w, reqId, InvalidJSON) //Bad json
		return
	}
	params := mux.Vars(r)
	id, err := strconv.ParseUint(params["id"], 10, 64)
	if err != nil {
		o.error(w, reqId, InvalidJSON) //Bad json
		return
	}
	response.OrderID = id
	response, err = o.useCase.ResponseOrder().Create(*response)
	if err != nil {
		o.error(w, reqId, New(err))
		return
	}
	o.respond(w, reqId, http.StatusCreated, response)
}

func (o *OrderHandler) handleGetAllOrderResponses(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	reqId := r.Context().Value(ctxKeyReqID).(uint64)
	id, err := strconv.ParseUint(params["id"], 10, 64)
	if err != nil {
		o.error(w, reqId, New(err)) //Bad json
		return
	}
	responses, err := o.orderUseCase.FindByVacancyID(id)
	if err != nil {
		o.error(w, reqId, New(err)) //Bad json
		return
	}

	o.respond(w, reqId, http.StatusOK, responses)
}

func (o *OrderHandler) handleChangeOrderResponse(w http.ResponseWriter, r *http.Request) {
	response := &models.ResponseOrder{}
	reqId := r.Context().Value(ctxKeyReqID).(uint64)

	if err := json.NewDecoder(r.Body).Decode(response); err != nil {
		o.error(w, reqId, InvalidJSON)
		return
	}
	params := mux.Vars(r)
	var err error
	response.OrderID, err = strconv.ParseUint(params["id"], 10, 64)
	if err != nil {
		o.error(w, reqId, New(err))
		return
	}
	response.UserID = r.Context().Value(ctxKeySession).(*models.Session).UserID
	responses, err := s.useCase.ResponseOrder().Change(*response)

	if err != nil {
		o.error(w, reqId, New(err))
		return
	}

	o.respond(w, reqId, http.StatusOK, responses)
}

func (o *OrderHandler) handleDeleteOrderResponse(w http.ResponseWriter, r *http.Request) {
	response := &models.ResponseOrder{}
	params := mux.Vars(r)
	reqId := r.Context().Value(ctxKeyReqID).(uint64)
	var err error
	response.OrderID, err = strconv.ParseUint(params["id"], 10, 64)
	if err != nil {
		o.error(w, reqId, New(err)) //Bad json
		return
	}
	response.UserID = r.Context().Value(ctxKeySession).(*models.Session).UserID
	err = o.orderUseCase.Delete(*response)

	if err != nil {
		o.error(w, reqId, New(err)) //Bad json
		return
	}
	var emptyInterface interface{}

	o.respond(w, reqId, http.StatusOK, emptyInterface)
}

func (o *OrderHandler) handleSelectExecutor(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	reqId := r.Context().Value(ctxKeyReqID).(uint64)

	order := models.Order{}
	if err := json.NewDecoder(r.Body).Decode(&order); err != nil {
		o.error(w, reqId, InvalidJSON)
		return
	}
	var err error
	order.ID, err = strconv.ParseUint(params["id"], 10, 64)
	if err != nil {
		o.error(w, http.StatusBadRequest, InvalidJSON)
		return
	}
	err = o.orderUseCase.SelectExecutor(order)
	if err != nil {
		o.error(w, http.StatusInternalServerError, New(err))
		return
	}
	o.respond(w, reqId, http.StatusOK, order)
}

func (o *OrderHandler) handleDeleteExecutor(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	reqId := r.Context().Value(ctxKeyReqID).(uint64)

	order := models.Order{}
	var err error
	order.ID, err = strconv.ParseUint(params["id"], 10, 64)
	if err != nil {
		o.error(w, http.StatusBadRequest, InvalidJSON)
		return
	}
	err = o.orderUseCase.DeleteExecutor(order)
	if err != nil {
		o.error(w, http.StatusInternalServerError, New(err))
		return
	}
	var emptyInterface interface{}
	s.respond(w, reqId, http.StatusOK, emptyInterface)
}

func (o *OrderHandler) handleGetAllUserOrders(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	reqId := r.Context().Value(ctxKeyReqID).(uint64)
	userID, err := strconv.ParseUint(params["id"], 10, 64)
	if err != nil {
		o.error(w, http.StatusBadRequest, InvalidJSON)
		return
	}
	userReq, err := user.UserUseCase.FindByID(userID)
	if err != nil {
		o.error(w, http.StatusInternalServerError, New(err))
		return
	}
	isExecutor := userReq.Executor
	var ord []models.Order
	if isExecutor {
		ord, err = o.use.FindByExecutorID(userID)
	} else {
		ord, err = o.orderUseCase.FindByCustomerID(userID)
	}
	if err != nil {
		o.error(w, http.StatusNotFound, New(err))
		return
	}
	o.respond(w, reqId, http.StatusOK, ord)
}
