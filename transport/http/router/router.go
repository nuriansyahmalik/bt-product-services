package router

import (
	"github.com/evermos/boilerplate-go/internal/handlers"
	"github.com/go-chi/chi"
)

// DomainHandlers is a struct that contains all domain-specific handlers.
type DomainHandlers struct {
	FooBarBazHandler handlers.FooBarBazHandler
	UserHandler      handlers.UserHandler
	BrandHandler     handlers.BrandHandler
	ProductHandler   handlers.ProductHandler
	VariantHandler   handlers.VariantHandler
	WarehouseHandler handlers.WarehouseHandler
}

// Router is the router struct containing handlers.
type Router struct {
	DomainHandlers DomainHandlers
}

// ProvideRouter is the provider function for this router.
func ProvideRouter(domainHandlers DomainHandlers) Router {
	return Router{
		DomainHandlers: domainHandlers,
	}
}

// SetupRoutes sets up all routing for this server.
func (r *Router) SetupRoutes(mux *chi.Mux) {
	mux.Route("/v1", func(rc chi.Router) {
		r.DomainHandlers.FooBarBazHandler.Router(rc)
		r.DomainHandlers.UserHandler.Router(rc)
		r.DomainHandlers.BrandHandler.Router(rc)
		r.DomainHandlers.ProductHandler.Router(rc)
		r.DomainHandlers.VariantHandler.Router(rc)
		r.DomainHandlers.WarehouseHandler.Router(rc)
	})
}
