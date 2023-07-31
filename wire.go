//go:build wireinject
// +build wireinject

package main

import (
	"github.com/evermos/boilerplate-go/configs"
	"github.com/evermos/boilerplate-go/event/producer"
	"github.com/evermos/boilerplate-go/infras"
	"github.com/evermos/boilerplate-go/internal/domain/brands"
	"github.com/evermos/boilerplate-go/internal/domain/foobarbaz"
	"github.com/evermos/boilerplate-go/internal/domain/products"
	"github.com/evermos/boilerplate-go/internal/domain/users"
	"github.com/evermos/boilerplate-go/internal/domain/variants"
	"github.com/evermos/boilerplate-go/internal/domain/warehouse"
	"github.com/evermos/boilerplate-go/internal/handlers"
	"github.com/evermos/boilerplate-go/transport/http"
	"github.com/evermos/boilerplate-go/transport/http/middleware"
	"github.com/evermos/boilerplate-go/transport/http/router"
	"github.com/google/wire"
)

// Wiring for configurations.
var configurations = wire.NewSet(
	configs.Get,
)

// Wiring for persistences.
var persistences = wire.NewSet(
	infras.ProvideMySQLConn,
)

// Wiring for domain FooBarBaz.
var domainFooBarBaz = wire.NewSet(
	// FooService interface and implementation
	foobarbaz.ProvideFooServiceImpl,
	wire.Bind(new(foobarbaz.FooService), new(*foobarbaz.FooServiceImpl)),
	// FooRepository interface and implementation
	foobarbaz.ProvideFooRepositoryMySQL,
	wire.Bind(new(foobarbaz.FooRepository), new(*foobarbaz.FooRepositoryMySQL)),
	// Producer interface and implementation
	producer.NewSNSProducer,
	wire.Bind(new(producer.Producer), new(*producer.SNSProducer)),
)

// Wiring for domainUser
var domainUser = wire.NewSet(
	//Service interface and implement
	users.ProvideUserServiceImpl,
	wire.Bind(new(users.UserService), new(*users.UserSerivceImpl)),
	//Repository interface and implementation
	users.ProvideUserRepositoryMySQL,
	wire.Bind(new(users.UserRepository), new(*users.UserRepositoryMysql)),
)

// Wiring for domain Brand
var domainBrand = wire.NewSet(
	//Service interface and implement
	brands.ProvideBrandServiceImpl,
	wire.Bind(new(brands.BrandService), new(*brands.BrandServiceImpl)),
	//Repository interface and implement
	brands.ProvideBrandRepository,
	wire.Bind(new(brands.BrandRepository), new(*brands.BrandRepositoryMySQL)),
)

var domainProduct = wire.NewSet(
	//Service interface and implement
	products.ProvideProductServiceImpl,
	wire.Bind(new(products.ProductService), new(*products.ProductServiceImpl)),
	//Repository interface and implement
	products.ProvideProductRepositoryMySQL,
	wire.Bind(new(products.ProductRepository), new(*products.ProductRepositoryMySQL)),
)

var domainVariant = wire.NewSet(
	//Service interface and implement
	variants.ProvideVariantServiceImpl,
	wire.Bind(new(variants.VariantService), new(*variants.VariantServiceImpl)),
	//Repository interface and implement
	variants.ProvideVariantRepositoryMySQl,
	wire.Bind(new(variants.VariantRepository), new(*variants.VariantRepositoryMySQL)),
)
var domainWarehouse = wire.NewSet(
	//Service interface and implement
	warehouse.ProvideWarehouseServiceImpl,
	wire.Bind(new(warehouse.WarehouseService), new(*warehouse.WarehouseServiceImpl)),
	//Repository interface and implement
	warehouse.ProvideWarehouseRepositoryMySQL,
	wire.Bind(new(warehouse.WarehouseRepository), new(*warehouse.WarehouseRepositoryMySQL)),
)

// Wiring for all domains.
var domains = wire.NewSet(
	domainFooBarBaz,
	domainUser,
	domainBrand,
	domainProduct,
	domainVariant,
	domainWarehouse,
)

var authMiddleware = wire.NewSet(
	middleware.ProvideAuthentication,
)

// Wiring for HTTP routing.
var routing = wire.NewSet(
	wire.Struct(new(router.DomainHandlers), "FooBarBazHandler", "UserHandler", "BrandHandler", "ProductHandler", "VariantHandler", "WarehouseHandler"),
	handlers.ProvideFooBarBazHandler,
	handlers.ProvideUserHandler,
	handlers.ProvideBrandHandler,
	handlers.ProvideProductHandler,
	handlers.ProvideVariantHandler,
	handlers.ProvideWarehouseHandler,
	router.ProvideRouter,
)

// Wiring for all domains event consumer.
//var evco = wire.NewSet(
//	wire.Struct(new(event.Consumers), "FooBarBaz"),
//	fooBarBazEvent.ProvideConsumerImpl,
//)

// Wiring for everything.
func InitializeService() *http.HTTP {
	wire.Build(
		// configurations
		configurations,
		// persistences
		persistences,
		// middleware
		authMiddleware,
		// domains
		domains,
		// routing
		routing,
		// selected transport layer
		http.ProvideHTTP)
	return &http.HTTP{}
}

// Wiring the event needs.
//func InitializeEvent() event.Consumers {
//	wire.Build(
//		// configurations
//		configurations,
//		// persistences
//		persistences,
//		// domains
//		domains,
//		// event consumer
//		evco)
//
//	return event.Consumers{}
//}
