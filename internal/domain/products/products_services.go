package products

import (
	"github.com/evermos/boilerplate-go/configs"
	"github.com/evermos/boilerplate-go/shared/failure"
	"github.com/gofrs/uuid"
)

type ProductService interface {
	Create(requestFormat ProductRequestFormat, variantID uuid.UUID) (product Product, err error)
	Update(id uuid.UUID) (product Product, err error)
	SearchProducts(params ProductSearchParams) ([]Product, error)
}

type ProductServiceImpl struct {
	ProductRepository ProductRepository
	Config            *configs.Config
}

func ProvideProductServiceImpl(productRepository ProductRepository, config *configs.Config) *ProductServiceImpl {
	return &ProductServiceImpl{ProductRepository: productRepository, Config: config}
}

func (p *ProductServiceImpl) Create(requestFormat ProductRequestFormat, productID uuid.UUID) (product Product, err error) {
	product, err = product.NewFromRequestFormat(requestFormat, productID)
	if err != nil {
		return
	}
	if err != nil {
		return product, failure.BadRequest(err)
	}

	err = p.ProductRepository.CreateProduct(product)
	if err != nil {
		return
	}
	return
}

func (p *ProductServiceImpl) Update(id uuid.UUID) (product Product, err error) {
	product, err = p.ProductRepository.ResolveByID(id)
	if err != nil {
		return
	}
	err = product.Update(id)
	if err != nil {
		return
	}
	err = p.ProductRepository.UpdateProduct(product)
	return
}

func (s *ProductServiceImpl) ListProducts() ([]Product, error) {
	products, err := s.ProductRepository.ListProducts()
	if err != nil {
		return nil, err
	}
	return products, nil
}

func (s *ProductServiceImpl) SearchProducts(params ProductSearchParams) ([]Product, error) {
	products, err := s.ProductRepository.SearchProducts(params)
	if err != nil {
		return nil, err
	}
	return products, nil
}
