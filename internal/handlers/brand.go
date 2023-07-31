package handlers

import (
	"encoding/json"
	"github.com/evermos/boilerplate-go/internal/domain/brands"
	"github.com/evermos/boilerplate-go/shared"
	"github.com/evermos/boilerplate-go/shared/failure"
	"github.com/evermos/boilerplate-go/transport/http/response"
	"github.com/go-chi/chi"
	"github.com/gofrs/uuid"
	"net/http"
)

type BrandHandler struct {
	BrandService brands.BrandService
}

func ProvideBrandHandler(BrandService brands.BrandService) BrandHandler {
	return BrandHandler{BrandService: BrandService}
}

func (h *BrandHandler) Router(r chi.Router) {
	r.Route("/brand", func(r chi.Router) {
		r.Post("/", h.CreateBrand)
	})
}

func (h *BrandHandler) CreateBrand(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var requestFormat brands.BrandRequestFormat
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
	brandID, _ := uuid.NewV4()
	brand, err := h.BrandService.Create(requestFormat, brandID)
	if err != nil {
		return
	}
	response.WithJSON(w, http.StatusCreated, brand)
}
