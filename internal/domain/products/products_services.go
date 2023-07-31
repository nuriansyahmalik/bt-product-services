package products

import (
	"github.com/evermos/boilerplate-go/configs"
	"github.com/evermos/boilerplate-go/event/producer"
	"github.com/evermos/boilerplate-go/shared/failure"
	"github.com/gofrs/uuid"
)

type ProductService interface {
	Create(requestFormat ProductRequestFormat, variantID uuid.UUID) (product Product, err error)
}

type ProductServiceImpl struct {
	ProductRepository ProductRepository
	Producer          producer.Producer
	Config            *configs.Config
}

func ProvideProductServiceImpl(productRepository ProductRepository, producer producer.Producer, config *configs.Config) *ProductServiceImpl {
	return &ProductServiceImpl{ProductRepository: productRepository, Producer: producer, Config: config}
}

func (p *ProductServiceImpl) Create(requestFormat ProductRequestFormat, productID uuid.UUID) (product Product, err error) {
	product, err = product.NewFromRequestFormat(requestFormat, productID)
	if err != nil {
		return
	}
	if err != nil {
		return product, failure.BadRequest(err)
	}

	err = p.ProductRepository.Create(product)
	if err != nil {
		return
	}
	return
}

func (p *ProductServiceImpl) GetAll() {
	return
}
