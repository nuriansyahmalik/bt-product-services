package products

import (
	"encoding/json"
	"fmt"
	"github.com/evermos/boilerplate-go/shared/failure"
	"github.com/evermos/boilerplate-go/shared/nuuid"
	"github.com/gofrs/uuid"
	"github.com/guregu/null"
	"time"
)

type Product struct {
	ProductId   uuid.UUID   `db:"productId"`
	ProductName string      `db:"productName"`
	VariantId   uuid.UUID   `db:"variantId"`
	BrandId     uuid.UUID   `db:"brandId"`
	Price       float64     `db:"price"`
	Stock       int         `db:"stock"`
	Status      string      `db:"status"`
	ImageURL    string      `db:"imageUrl"`
	CreatedAt   time.Time   `db:"createdAt"`
	CreatedBy   uuid.UUID   `db:"createdBy"`
	UpdatedAt   null.Time   `db:"updatedAt"`
	UpdatedBy   nuuid.NUUID `db:"updatedBy"`
	Deleted     null.Time   `db:"deletedAt"`
	DeletedBy   nuuid.NUUID `db:"deletedBy"`
}

type FilterOptions struct {
	BrandName   string
	ProductName string
	VariantName string
	Status      string
}

type SortingOptions struct {
	By        string
	Ascending bool
}

type PaginationOptions struct {
	Page     int
	PageSize int
}

type Image struct {
	ImageId   uuid.UUID `db:"imageId"`
	ProductId uuid.UUID `db:"productId"`
	ImageURL  string    `db:"imageUrl"`
	CreatedAt time.Time `db:"createdAt"`
	CreatedBy uuid.UUID `db:"createdBy"`
}

type ImageRequestFormat struct {
	ImageURL  string    `json:"imageURL"`
	CreatedAt time.Time `json:"created"`
}

type ProductRequestFormat struct {
	ProductName string    `json:"productName" validate:"required"`
	VariantId   uuid.UUID `json:"variantId"`
}

type ProductResponseFormat struct {
	ID          uuid.UUID  `json:"id"`
	VariantId   uuid.UUID  `json:"variantId"`
	ProductName string     `json:"productName"`
	Created     time.Time  `json:"created"`
	CreatedBy   uuid.UUID  `json:"createdBy"`
	Updated     null.Time  `json:"updated,omitempty"`
	UpdatedBy   *uuid.UUID `json:"updatedBy,omitempty"`
	Deleted     null.Time  `json:"deleted,omitempty"`
	DeletedBy   *uuid.UUID `json:"deletedBy,omitempty"`
}

func (p Product) MarshalJSON() ([]byte, error) {
	return json.Marshal(p.ToResponseFormat())
}

func (p *Product) SoftDelete(id uuid.UUID) (err error) {
	if p.IsDeleted() {
		return failure.Conflict("softDelete", "foo", "already marked as deleted")
	}

	p.Deleted = null.TimeFrom(time.Now())
	p.DeletedBy = nuuid.From(id)
	return
}
func (p *Product) IsDeleted() (deleted bool) {
	return p.Deleted.Valid && p.DeletedBy.Valid
}

func (p *Product) Update(id uuid.UUID) (err error) {
	fmt.Println(id)
	return err
}
func (p Product) NewFromRequestFormat(req ProductRequestFormat, productID uuid.UUID) (newProduct Product, err error) {
	productID, _ = uuid.NewV4()
	newProduct = Product{
		ProductId:   productID,
		ProductName: req.ProductName,
		VariantId:   req.VariantId,
		CreatedAt:   time.Now(),
		CreatedBy:   productID,
	}
	products := make([]Product, 0)
	products = append(products, newProduct)
	return
}

func (p Image) NewFromRequestFormat(format ImageRequestFormat, productId uuid.UUID) (newImage Image) {
	imageId, _ := uuid.NewV4()
	newImage = Image{
		ImageId:   imageId,
		ProductId: productId,
		ImageURL:  format.ImageURL,
		CreatedAt: time.Now(),
		CreatedBy: imageId,
	}

	return
}

func (i *Image) ToResponseFormat() ImageResponseFormat {
	return ImageResponseFormat{
		ImageId:   i.ImageId,
		ProductId: i.ProductId,
		ImageURL:  i.ImageURL,
	}
}

type ImageResponseFormat struct {
	ImageId   uuid.UUID `json:"imageId"`
	ProductId uuid.UUID `json:"productId"`
	ImageURL  string    `json:"imageURL"`
}

func (p *Product) ToResponseFormat() ProductResponseFormat {
	return ProductResponseFormat{
		ID:          p.ProductId,
		VariantId:   p.VariantId,
		ProductName: p.ProductName,
		Created:     p.CreatedAt,
		CreatedBy:   p.CreatedBy,
		Updated:     p.UpdatedAt,
		Deleted:     p.Deleted,
	}
}
