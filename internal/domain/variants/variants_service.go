package variants

import (
	"github.com/evermos/boilerplate-go/configs"
	"github.com/evermos/boilerplate-go/event/producer"
	"github.com/evermos/boilerplate-go/shared/failure"
	"github.com/gofrs/uuid"
)

type VariantService interface {
	Create(requestFormat VariantRequestFormat, varintId uuid.UUID) (variant Variants, err error)
}

type VariantServiceImpl struct {
	VariantRepository VariantRepository
	Producer          producer.Producer
	Config            *configs.Config
}

func ProvideVariantServiceImpl(variantRepository VariantRepository, producer producer.Producer, config *configs.Config) *VariantServiceImpl {
	return &VariantServiceImpl{
		VariantRepository: variantRepository,
		Producer:          producer,
		Config:            config,
	}
}

func (v *VariantServiceImpl) Create(requestFormat VariantRequestFormat, variantId uuid.UUID) (variant Variants, err error) {
	variant, err = variant.NewFromRequestFormat(requestFormat, variantId)
	if err != nil {
		return
	}
	if err != nil {
		return variant, failure.BadRequest(err)
	}
	err = v.VariantRepository.Create(variant)
	if err != nil {
		return
	}
	return
}
