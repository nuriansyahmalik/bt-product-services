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
	"strconv"
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
		r.Get("/search", h.SearchProducts)
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
func (h *ProductHandler) SearchProducts(w http.ResponseWriter, r *http.Request) {
	page, err := strconv.Atoi(r.URL.Query().Get("page"))
	if err != nil {
		response.WithMessage(w, http.StatusBadRequest, "Missing Param Query Limit")
		response.WithError(w, err)
		return
	}
	pageSize, err := strconv.Atoi(r.URL.Query().Get("page_size"))
	if err != nil {
		response.WithMessage(w, http.StatusBadRequest, "Missing Param Query Page")
		response.WithError(w, err)
		return
	}
	params := products.ProductSearchParams{
		BrandName:   r.URL.Query().Get("brand_name"),
		ProductName: r.URL.Query().Get("product_name"),
		VariantName: r.URL.Query().Get("variant_name"),
		Status:      r.URL.Query().Get("status"),
		SortBy:      r.URL.Query().Get("sort_by"),
		Page:        page - 1,
		PageSize:    pageSize,
	}

	searchResult, err := h.ProductService.SearchProducts(params)
	if err != nil {
		response.WithMessage(w, http.StatusBadRequest, "Error searching products")
		response.WithError(w, err)
		return
	}
	response.WithJSON(w, http.StatusOK, searchResult)
}
