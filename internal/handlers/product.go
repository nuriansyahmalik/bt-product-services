package handlers

import (
	"encoding/json"
	"github.com/evermos/boilerplate-go/internal/domain/products"
	"github.com/evermos/boilerplate-go/shared"
	"github.com/evermos/boilerplate-go/shared/failure"
	"github.com/evermos/boilerplate-go/transport/http/response"
	"github.com/go-chi/chi"
	"github.com/gofrs/uuid"
	"net/http"
)

type ProductHandler struct {
	ProductService products.ProductService
}

func ProvideProductHandler(ProductService products.ProductService) ProductHandler {
	return ProductHandler{ProductService: ProductService}
}

func (h *ProductHandler) Router(r chi.Router) {
	r.Route("/product", func(r chi.Router) {
		r.Post("/", h.CreateProduct)
	})
}
func (h *ProductHandler) CreateProduct(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var requestFormat products.ProductRequestFormat
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

	variantID, _ := uuid.NewV4() // TODO: read from context
	product, err := h.ProductService.Create(requestFormat, variantID)
	if err != nil {
		response.WithError(w, err)
		return
	}

	response.WithJSON(w, http.StatusCreated, product)
}
