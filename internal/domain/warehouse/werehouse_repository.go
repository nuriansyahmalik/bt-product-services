package warehouse

import (
	"github.com/evermos/boilerplate-go/infras"
	"github.com/evermos/boilerplate-go/shared/failure"
	"github.com/evermos/boilerplate-go/shared/logger"
	"github.com/gofrs/uuid"
)

var (
	warehouseQueries = struct {
		selectWarehouse string
		insertWarehouse string
		insertQuantity  string
	}{
		selectWarehouse: `
SELECT
FROM warehouses w`,
		insertWarehouse: `INSERT INTO warehouse
				(warehouseId, warehouseName, createdAt)
				VALUES
				(:warehouseId, :warehouseName, NOW())`,
		insertQuantity: `INSERT INTO quantity 
				(quantityId, productId, warehouseId, quantity, status, createdAt)
				VALUES
				(:quantityId, :productId, :warehouseId, :quantity, :status, :createdAt)`,
	}
)

type WarehouseRepository interface {
	Create(warehouse Warehouses) (err error)
	CreateQuantity(quantity Quantity) (err error)
	ExistsByID(id uuid.UUID) (exists bool, err error)
}

type WarehouseRepositoryMySQL struct {
	DB *infras.MySQLConn
}

func ProvideWarehouseRepositoryMySQL(db *infras.MySQLConn) *WarehouseRepositoryMySQL {
	return &WarehouseRepositoryMySQL{DB: db}
}

func (w *WarehouseRepositoryMySQL) Create(warehouse Warehouses) (err error) {
	exists, err := w.ExistsByID(warehouse.WarehouseId)
	if err != nil {
		logger.ErrorWithStack(err)
		return
	}

	if exists {
		err = failure.Conflict("create", "warehouse", "already exists")
		logger.ErrorWithStack(err)
		return
	}

	stmt, err := w.DB.Write.PrepareNamed(warehouseQueries.insertWarehouse)
	if err != nil {
		logger.ErrorWithStack(err)
		return
	}
	defer stmt.Close()

	_, err = stmt.Exec(warehouse)
	if err != nil {
		logger.ErrorWithStack(err)
	}
	return
}
func (w *WarehouseRepositoryMySQL) CreateQuantity(quantity Quantity) (err error) {
	stmt, err := w.DB.Write.PrepareNamed(warehouseQueries.insertQuantity)
	if err != nil {
		logger.ErrorWithStack(err)
		return
	}
	defer stmt.Close()

	_, err = stmt.Exec(quantity)
	if err != nil {
		logger.ErrorWithStack(err)
	}
	return
}

func (w *WarehouseRepositoryMySQL) ExistsByID(id uuid.UUID) (exists bool, err error) {
	err = w.DB.Read.Get(
		&exists,
		"SELECT COUNT(warehouseId) FROM warehouse w WHERE w.warehouseId = ?",
		id.String())
	if err != nil {
		logger.ErrorWithStack(err)
	}

	return
}
