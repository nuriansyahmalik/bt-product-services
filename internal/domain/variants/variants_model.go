package variants

import (
	"github.com/evermos/boilerplate-go/shared/nuuid"
	"github.com/gofrs/uuid"
	"github.com/guregu/null"
	"time"
)

type Variants struct {
	VariantId   uuid.UUID   `db:"variantId"`
	VariantName string      `db:"variantName"`
	BrandId     uuid.UUID   `db:"brandId"`
	Price       float64     `db:"price"`
	CreatedAt   time.Time   `db:"createdAt"`
	CreatedBy   uuid.UUID   `db:"createdBy"`
	UpdatedAt   null.Time   `db:"updatedAt"`
	UpdatedBy   nuuid.NUUID `db:"updatedBy"`
	Deleted     null.Time   `db:"deletedAt"`
	DeletedBy   nuuid.NUUID `db:"deletedBy"`
}

type VariantRequestFormat struct {
	VariantName string    `json:"variantName"`
	BrandId     uuid.UUID `json:"brandId"`
	Price       float64   `json:"price"`
}

func (v Variants) NewFromRequestFormat(req VariantRequestFormat, variantId uuid.UUID) (newVariant Variants, err error) {
	variantId, _ = uuid.NewV4()
	newVariant = Variants{
		VariantId:   variantId,
		VariantName: req.VariantName,
		BrandId:     req.BrandId,
		Price:       req.Price,
		CreatedAt:   time.Now(),
		CreatedBy:   variantId,
	}
	variants := make([]Variants, 0)
	variants = append(variants, newVariant)
	return
}
