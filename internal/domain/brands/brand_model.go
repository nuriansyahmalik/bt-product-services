package brands

import (
	"github.com/evermos/boilerplate-go/shared/nuuid"
	"github.com/gofrs/uuid"
	"github.com/guregu/null"
	"time"
)

type Brands struct {
	BrandId   uuid.UUID   `db:"brandId"`
	BrandName string      `db:"brandName"`
	CreatedAt time.Time   `db:"createdAt"`
	CreatedBy uuid.UUID   `db:"createdBy"`
	UpdatedAt null.Time   `db:"updatedAt"`
	UpdatedBy nuuid.NUUID `db:"updatedBy"`
	Deleted   null.Time   `db:"deletedAt"`
	DeletedBy nuuid.NUUID `db:"deletedBy"`
}

type BrandRequestFormat struct {
	BrandName string `json:"brandName"`
}

func (b Brands) NewFromRequestFormat(req BrandRequestFormat, brandId uuid.UUID) (newBrand Brands, err error) {
	brandId, _ = uuid.NewV4()
	newBrand = Brands{
		BrandId:   brandId,
		BrandName: req.BrandName,
		CreatedAt: time.Now(),
		CreatedBy: brandId,
	}
	brands := make([]Brands, 0)
	brands = append(brands, newBrand)
	return
}
