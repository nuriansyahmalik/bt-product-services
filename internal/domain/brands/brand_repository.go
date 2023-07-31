package brands

import (
	"github.com/evermos/boilerplate-go/infras"
	"github.com/evermos/boilerplate-go/shared/failure"
	"github.com/evermos/boilerplate-go/shared/logger"
	"github.com/gofrs/uuid"
)

var (
	brandQueries = struct {
		selectBrand string
		insertBrand string
	}{
		selectBrand: `
		SELECT
   			b.brandId,
			b.brandName,
			b.createdAt,
			b.createdBy,
			b.updatedAt,
			b.updatedBy,
			b.deletedAt,
			b.deletedBy
		FROM
			brands b;
`,
		insertBrand: `
			INSERT INTO brand (
			           brandId, brandName, createdAt,createdBy
			) VALUES (
			          :brandId, :brandName, NOW(),:createdBy)`,
	}
)

type BrandRepository interface {
	Create(brand Brands) (err error)
	ExistsByID(id uuid.UUID) (exists bool, err error)
	ResolveByID(id uuid.UUID) (brand Brands, err error)
}

type BrandRepositoryMySQL struct {
	DB *infras.MySQLConn
}

func ProvideBrandRepository(db *infras.MySQLConn) *BrandRepositoryMySQL {
	return &BrandRepositoryMySQL{
		DB: db}
}
func (b *BrandRepositoryMySQL) Create(brand Brands) (err error) {
	exists, err := b.ExistsByID(brand.BrandId)
	if err != nil {
		logger.ErrorWithStack(err)
		return
	}
	if exists {
		err = failure.Conflict("create", "bra d", "already exists")
		logger.ErrorWithStack(err)
		return
	}
	stmt, err := b.DB.Write.PrepareNamed(brandQueries.insertBrand)
	if err != nil {
		logger.ErrorWithStack(err)
		return
	}
	defer stmt.Close()

	_, err = stmt.Exec(brand)
	if err != nil {
		logger.ErrorWithStack(err)
	}
	return
}

func (b *BrandRepositoryMySQL) ExistsByID(id uuid.UUID) (exists bool, err error) {
	err = b.DB.Read.Get(
		&exists,
		"SELECT COUNT(brandId) FROM brand WHERE brand.brandId = ?",
		id.String())
	if err != nil {
		logger.ErrorWithStack(err)
	}

	return
}

func (b *BrandRepositoryMySQL) ResolveByID(id uuid.UUID) (brand Brands, err error) {
	err = b.DB.Read.Get(
		&brand,
		brandQueries.selectBrand+" WHERE b.brandId = ?",
		id.String())
	if err != nil {
		err = failure.NotFound("brand")
		logger.ErrorWithStack(err)
		return
	}
	return
}
