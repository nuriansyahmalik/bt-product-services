package warehouse

import (
	"github.com/evermos/boilerplate-go/shared/nuuid"
	"github.com/gofrs/uuid"
	"github.com/guregu/null"
	"time"
)

type Warehouses struct {
	WarehouseId   uuid.UUID   `db:"warehouseId"`
	WarehouseName string      `db:"warehouseName"`
	CreatedAt     time.Time   `db:"createdAt"`
	CreatedBy     uuid.UUID   `db:"createdBy"`
	UpdatedAt     null.Time   `db:"updatedAt"`
	UpdatedBy     nuuid.NUUID `db:"updatedBy"`
	Deleted       null.Time   `db:"deletedAt"`
	DeletedBy     nuuid.NUUID `db:"deletedBy"`
}

type Quantity struct {
	QuantityId  uuid.UUID `db:"quantityId"`
	ProductId   uuid.UUID `db:"productId"`
	WarehouseId uuid.UUID `db:"warehouseId"`
	Quantity    int       `db:"quantity"`
	Status      string    `db:"status"`
	CreatedAt   time.Time `db:"createdAt"`
	CreatedBy   uuid.UUID `db:"createdBy"`
	UpdatedAt   null.Time `db:"updatedAt"`
	UpdatedBy   uuid.UUID `db:"updatedBy"`
}

type WarehouseRequestFormat struct {
	WarehouseName string `json:"warehouseName"`
}
type QuantityRequestFormat struct {
	ProductId   uuid.UUID `json:"productId"`
	WarehouseId uuid.UUID `json:"warehouseId"`
	Quantity    int       `json:"quantity"`
	Status      string    `json:"status"`
}

func (w Warehouses) NewFromRequestFormat(req WarehouseRequestFormat, warehouseId uuid.UUID) (newWarehouse Warehouses, err error) {
	warehouseId, _ = uuid.NewV4()
	newWarehouse = Warehouses{
		WarehouseId:   warehouseId,
		WarehouseName: req.WarehouseName,
		CreatedAt:     time.Now(),
	}
	warehouses := make([]Warehouses, 0)
	warehouses = append(warehouses, newWarehouse)
	return
}

func (q Quantity) NewFromRequestFormat(req QuantityRequestFormat, quantityId uuid.UUID) (newQuantity Quantity, err error) {
	quantityId, _ = uuid.NewV4()
	newQuantity = Quantity{
		QuantityId:  quantityId,
		ProductId:   req.ProductId,
		WarehouseId: req.WarehouseId,
		Quantity:    req.Quantity,
		Status:      req.Status,
		CreatedAt:   time.Now(),
	}
	quantities := make([]Quantity, 0)
	quantities = append(quantities, newQuantity)
	return
}
