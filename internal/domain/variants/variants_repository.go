package variants

import (
	"github.com/evermos/boilerplate-go/infras"
	"github.com/evermos/boilerplate-go/shared/failure"
	"github.com/evermos/boilerplate-go/shared/logger"
	"github.com/gofrs/uuid"
)

var (
	variantsQueries = struct {
		selectVariants string
		insertVariants string
	}{
		selectVariants: `
			SELECT
			FROM variant v`,
		insertVariants: `INSERT INTO variant 
				(variantId, variantName, brandId, price, createdAt, createdBy)
				VALUES
				(:variantId, :variantName, :brandId, :price, NOW(), :createdBy)`,
	}
)

type VariantRepository interface {
	Create(variants Variants) (err error)
	ExistsByID(id uuid.UUID) (exists bool, err error)
}

type VariantRepositoryMySQL struct {
	DB *infras.MySQLConn
}

func ProvideVariantRepositoryMySQl(db *infras.MySQLConn) *VariantRepositoryMySQL {
	return &VariantRepositoryMySQL{
		DB: db,
	}
}

func (v *VariantRepositoryMySQL) Create(variants Variants) (err error) {
	exists, err := v.ExistsByID(variants.VariantId)
	if err != nil {
		logger.ErrorWithStack(err)
		return
	}
	if exists {
		err = failure.Conflict("create", "variants", "already exists")
		logger.ErrorWithStack(err)
		return
	}
	stmt, err := v.DB.Write.PrepareNamed(variantsQueries.insertVariants)
	if err != nil {
		logger.ErrorWithStack(err)
	}
	defer stmt.Close()

	_, err = stmt.Exec(variants)
	if err != nil {
		logger.ErrorWithStack(err)
	}
	return
}

func (v *VariantRepositoryMySQL) ExistsByID(id uuid.UUID) (exists bool, err error) {
	err = v.DB.Read.Get(
		&exists,
		"SELECT COUNT(variantId) FROM variant WHERE variant.variantId = ?",
		id.String())
	if err != nil {
		logger.ErrorWithStack(err)
	}

	return
}
