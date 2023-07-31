package handlers

import (
	"encoding/json"
	"github.com/evermos/boilerplate-go/internal/domain/variants"
	"github.com/evermos/boilerplate-go/shared"
	"github.com/evermos/boilerplate-go/shared/failure"
	"github.com/evermos/boilerplate-go/transport/http/response"
	"github.com/go-chi/chi"
	"github.com/gofrs/uuid"
	"net/http"
)

type VariantHandler struct {
	VariantService variants.VariantService
}

func ProvideVariantHandler(VariantService variants.VariantService) VariantHandler {
	return VariantHandler{VariantService: VariantService}
}

func (h *VariantHandler) Router(r chi.Router) {
	r.Route("/variant", func(r chi.Router) {
		r.Post("/", h.CreateVariant)
	})
}
func (h *VariantHandler) CreateVariant(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var requestFormat variants.VariantRequestFormat
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

	variantID, _ := uuid.NewV4()
	variant, err := h.VariantService.Create(requestFormat, variantID)
	if err != nil {
		response.WithError(w, err)
		return
	}

	response.WithJSON(w, http.StatusCreated, variant)
}
