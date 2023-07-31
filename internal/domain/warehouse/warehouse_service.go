package warehouse

import (
	"github.com/evermos/boilerplate-go/configs"
	"github.com/evermos/boilerplate-go/event/producer"
	"github.com/evermos/boilerplate-go/shared/failure"
	"github.com/gofrs/uuid"
)

type WarehouseService interface {
	Create(requestFormat WarehouseRequestFormat, warehouseId uuid.UUID) (warehouse Warehouses, err error)
	CreateQuantity(requestFormat QuantityRequestFormat, quantityId uuid.UUID) (quantity Quantity, err error)
}

type WarehouseServiceImpl struct {
	WarehouseRepository WarehouseRepository
	Producer            producer.Producer
	Config              *configs.Config
}

func ProvideWarehouseServiceImpl(werehouseRepository WarehouseRepository, producer producer.Producer, config *configs.Config) *WarehouseServiceImpl {
	return &WarehouseServiceImpl{
		WarehouseRepository: werehouseRepository,
		Producer:            producer,
		Config:              config,
	}
}

func (w *WarehouseServiceImpl) Create(requestFormat WarehouseRequestFormat, warehouseId uuid.UUID) (warehouse Warehouses, err error) {
	warehouse, err = warehouse.NewFromRequestFormat(requestFormat, warehouseId)
	if err != nil {
		return
	}
	if err != nil {
		return warehouse, failure.BadRequest(err)
	}
	err = w.WarehouseRepository.Create(warehouse)
	if err != nil {
		return
	}
	return
}

func (w *WarehouseServiceImpl) CreateQuantity(requestFormat QuantityRequestFormat, quantityId uuid.UUID) (quantity Quantity, err error) {
	quantity, err = quantity.NewFromRequestFormat(requestFormat, quantityId)
	if err != nil {
		return
	}
	if err != nil {
		return quantity, failure.BadRequest(err)
	}
	err = w.WarehouseRepository.CreateQuantity(quantity)
	if err != nil {
		return
	}
	return
}
