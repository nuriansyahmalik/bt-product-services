package brands

import (
	"github.com/evermos/boilerplate-go/configs"
	"github.com/evermos/boilerplate-go/event/producer"
	"github.com/evermos/boilerplate-go/shared/failure"
	"github.com/gofrs/uuid"
)

type BrandService interface {
	Create(requestFormat BrandRequestFormat, brandId uuid.UUID) (brand Brands, err error)
	ResolveByID(id uuid.UUID) (brand Brands, err error)
}

type BrandServiceImpl struct {
	BrandRepository BrandRepository
	Producer        producer.Producer
	Config          *configs.Config
}

func ProvideBrandServiceImpl(brandRepository BrandRepository, producer producer.Producer, config *configs.Config) *BrandServiceImpl {
	return &BrandServiceImpl{
		BrandRepository: brandRepository,
		Producer:        producer,
		Config:          config,
	}
}

func (b *BrandServiceImpl) Create(requestFormat BrandRequestFormat, brandId uuid.UUID) (brand Brands, err error) {
	brand, err = brand.NewFromRequestFormat(requestFormat, brandId)
	if err != nil {
		return
	}
	if err != nil {
		return brand, failure.BadRequest(err)
	}
	err = b.BrandRepository.Create(brand)
	if err != nil {
		return
	}
	return
}

func (b *BrandServiceImpl) ResolveByID(id uuid.UUID) (brand Brands, err error) {
	brand, err = b.BrandRepository.ResolveByID(id)
	if err != nil {
		return brand, failure.NotFound("brand")
	}
	return
}
