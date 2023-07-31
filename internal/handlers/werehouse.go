package handlers

import (
	"encoding/json"
	"github.com/evermos/boilerplate-go/internal/domain/warehouse"
	"github.com/evermos/boilerplate-go/shared"
	"github.com/evermos/boilerplate-go/shared/failure"
	"github.com/evermos/boilerplate-go/transport/http/response"
	"github.com/go-chi/chi"
	"github.com/gofrs/uuid"
	"net/http"
)

type WarehouseHandler struct {
	WarehouseService warehouse.WarehouseService
}

func ProvideWarehouseHandler(WarehouseService warehouse.WarehouseService) WarehouseHandler {
	return WarehouseHandler{WarehouseService: WarehouseService}
}

func (h *WarehouseHandler) Router(r chi.Router) {
	r.Route("/warehouse", func(r chi.Router) {
		r.Post("/", h.CreateWarehouse)
		r.Post("/quantity", h.CreateQuantity)
	})
}

func (h *WarehouseHandler) CreateWarehouse(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var requestFormat warehouse.WarehouseRequestFormat
	err := decoder.Decode(&requestFormat)
	if err != nil {
		response.WithError(w, failure.BadRequest(err))
		return
	}

	err = shared.GetValidator().Struct(requestFormat)
	if err != nil {
		response.WithError(w, failure.BadRequest(err))
		return
	}
	warehouseID, _ := uuid.NewV4()
	warehouse, err := h.WarehouseService.Create(requestFormat, warehouseID)
	if err != nil {
		response.WithError(w, err)
		return
	}

	response.WithJSON(w, http.StatusCreated, warehouse)
}

func (h *WarehouseHandler) CreateQuantity(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var requestFormat warehouse.QuantityRequestFormat
	err := decoder.Decode(&requestFormat)
	if err != nil {
		response.WithError(w, failure.BadRequest(err))
		return
	}

	err = shared.GetValidator().Struct(requestFormat)
	if err != nil {
		response.WithError(w, failure.BadRequest(err))
		return
	}
	quantityId, _ := uuid.NewV4()
	quantity, err := h.WarehouseService.CreateQuantity(requestFormat, quantityId)
	if err != nil {
		response.WithError(w, err)
		return
	}

	response.WithJSON(w, http.StatusCreated, quantity)
}
