package products

import (
	"database/sql"
	"fmt"
	"github.com/evermos/boilerplate-go/infras"
	"github.com/evermos/boilerplate-go/shared/failure"
	"github.com/evermos/boilerplate-go/shared/logger"
	"github.com/gofrs/uuid"
	"github.com/jmoiron/sqlx"
)

var (
	productQueries = struct {
		selectProduct          string
		selectProducts         string
		insertProduct          string
		insertImage            string
		insertImagePlaceholder string
		updateProduct          string
	}{
		selectProduct: ``,
		selectProducts: `SELECT
    b.brandName,
    p.productName,
    v.variantName,
    i.imageUrl,
    v.price,
    q.quantity AS stock,
    p.status,
    p.updatedBy
FROM products p
JOIN variants v ON p.variantId = v.variantId
JOIN brands b ON v.brandId = b.brandId
LEFT JOIN images i ON p.productId = i.productId
LEFT JOIN (
    SELECT
        productId,
        SUM(quantity) AS quantity,
        status
    FROM quantity
    GROUP BY productId, status
) q ON p.productId = q.productId AND q.status = 'in_stock'
WHERE
    (:brandName IS NULL OR b.brandName LIKE CONCAT('%', :brandName, '%'))
    AND (:productName IS NULL OR p.productName LIKE CONCAT('%', :productName, '%'))
    AND (:variantName IS NULL OR v.variantName LIKE CONCAT('%', :variantName, '%'))
ORDER BY
    p.createdAt DESC, q.quantity DESC;`,
		insertProduct: `
			INSERT INTO products (
			          productId,
                      productName,
                      variantId,
                      createdAt,
                      createdBy,
			          updatedAt
			) VALUES (
			          :productId,
			          :productName,
			          :variantId,
			          :createdAt,
			          :createdBy,
			          :updatedAt)`,
		insertImage: `
			INSERT INTO images (
			          imageId, 
			          productId,
			          imageUrl,
			          createdAt,
			          createdBy
			) VALUES `,
		insertImagePlaceholder: `
					:imageId,
					:productId, 
					:imageUrl, 
					:createdAt,
					:createdBy)`,
		updateProduct: `
		UPDATE Products
		SET 
		    product_name = :product_name, 
		    variantId = :variantId,
		    updatedAt = NOW(),
		    updatedBy = :updatedBy
		WHERE productId = :productId`,
	}
)

type ProductRepository interface {
	CreateProduct(product Product) error
	UpdateProduct(product Product) error
	HardDeleteProduct(productID string) error
	ListProducts() ([]Product, error)
	SearchProducts(params ProductSearchParams) ([]Product, error)
	SortProducts(sortBy string) ([]Product, error)
	PaginateProducts(page, pageSize int) ([]Product, error)
	ExistsByID(id uuid.UUID) (exists bool, err error)
	ResolveByID(id uuid.UUID) (product Product, err error)
}

type ProductRepositoryMySQL struct {
	DB *infras.MySQLConn
}

func ProvideProductRepositoryMySQL(db *infras.MySQLConn) *ProductRepositoryMySQL {
	return &ProductRepositoryMySQL{DB: db}
}

func (p *ProductRepositoryMySQL) CreateProduct(product Product) error {
	exists, err := p.ExistsByID(product.ProductId)
	if err != nil {
		logger.ErrorWithStack(err)
		return err
	}

	if exists {
		err = failure.Conflict("create", "product", "already exists")
		logger.ErrorWithStack(err)
		return err
	}
	return p.DB.WithTransaction(func(tx *sqlx.Tx, e chan error) {
		if err := p.txCreate(tx, product); err != nil {
			e <- err
			return
		}

		e <- nil
	})
}
func (p *ProductRepositoryMySQL) UpdateProduct(product Product) error {
	exists, err := p.ExistsByID(product.ProductId)
	if err != nil {
		logger.ErrorWithStack(err)
		return err
	}

	if !exists {
		err = failure.NotFound("foo")
		logger.ErrorWithStack(err)
		return err
	}
	return p.DB.WithTransaction(func(tx *sqlx.Tx, e chan error) {
		if err := p.txUpdate(tx, product); err != nil {
			e <- err
			return
		}

		e <- nil
	})
}
func (p *ProductRepositoryMySQL) HardDeleteProduct(productID string) error {
	return nil
}
func (p *ProductRepositoryMySQL) ListProducts() ([]Product, error) {
	var products []Product
	err := p.DB.Read.Select(&products, productQueries.selectProducts)
	if err != nil {
		return nil, err
	}
	return products, nil
}
func (p *ProductRepositoryMySQL) SearchProducts(params ProductSearchParams) ([]Product, error) {
	query := `
		SELECT
			b.brandName,
			p.productName,
			v.variantName,
			i.imageUrl,
			v.price,
			q.quantity AS stock,
			q.status,
			p.updatedBy
		FROM products p
		JOIN variants v ON p.variantId = v.variantId
		JOIN brands b ON v.brandId = b.brandId
		LEFT JOIN images i ON p.productId = i.productId
		LEFT JOIN (
			SELECT
				productId,
				SUM(quantity) AS quantity,
				status
			FROM quantity
			GROUP BY productId, status
		) q ON p.productId = q.productId AND q.status = 'in_stock'
		WHERE 1 = 1
	`

	args := []interface{}{}

	if params.BrandName != "" {
		query += "AND b.brandName LIKE ? "
		args = append(args, "%"+params.BrandName+"%")
	}

	if params.ProductName != "" {
		query += "AND p.productName LIKE ? "
		args = append(args, "%"+params.ProductName+"%")
	}

	if params.VariantName != "" {
		query += "AND v.variantName LIKE ? "
		args = append(args, "%"+params.VariantName+"%")
	}

	if params.Status != "" {
		query += "AND q.status = ? "
		args = append(args, params.Status)
	}

	query += "ORDER BY p.createdAt DESC, q.quantity DESC"

	var products []Product
	err := p.DB.Read.Select(&products, query, args...)
	if err != nil {
		return nil, err
	}
	return products, nil
}

func (p *ProductRepositoryMySQL) SortProducts(sortBy string) ([]Product, error) {
	var sortField string
	switch sortBy {
	case "createdAt":
		sortField = "p.createdAt"
	case "stock":
		sortField = "q.quantity"
	default:
		return nil, fmt.Errorf("unsupported sorting field: %s", sortBy)
	}

	query := fmt.Sprintf("%s ORDER BY %s DESC", productQueries.selectProducts, sortField)
	var products []Product
	err := p.DB.Read.Select(&products, query)
	if err != nil {
		return nil, err
	}
	return products, nil
}
func (p *ProductRepositoryMySQL) PaginateProducts(page, pageSize int) ([]Product, error) {
	offset := (page - 1) * pageSize
	query := productQueries.selectProducts + " LIMIT ? OFFSET ?"
	products := []Product{}
	err := p.DB.Read.Select(&products, query, pageSize, offset)
	if err != nil {
		return nil, err
	}
	return products, nil
}
func (p *ProductRepositoryMySQL) ResolveByID(id uuid.UUID) (product Product, err error) {
	err = p.DB.Read.Get(
		&product,
		productQueries.selectProduct+" WHERE foo.entity_id = ?",
		id.String())
	if err != nil && err == sql.ErrNoRows {
		err = failure.NotFound("product")
		logger.ErrorWithStack(err)
		return
	}
	return
}
func (r *ProductRepositoryMySQL) ExistsByID(id uuid.UUID) (exists bool, err error) {
	err = r.DB.Read.Get(
		&exists,
		"SELECT COUNT(entity_id) FROM product p WHERE p.entity_id = ?",
		id.String())
	if err != nil {
		logger.ErrorWithStack(err)
	}

	return
}

func (p *ProductRepositoryMySQL) txCreate(tx *sqlx.Tx, product Product) (err error) {
	stmt, err := tx.PrepareNamed(productQueries.insertProduct)
	if err != nil {
		logger.ErrorWithStack(err)
		return
	}
	defer stmt.Close()

	_, err = stmt.Exec(product)
	if err != nil {
		logger.ErrorWithStack(err)
	}

	return
}
func (r *ProductRepositoryMySQL) txDeleteItems(tx *sqlx.Tx, productID uuid.UUID) (err error) {
	_, err = tx.Exec("DELETE FROM foo_item WHERE foo_id = ?", productID.String())
	return
}
func (r *ProductRepositoryMySQL) txUpdate(tx *sqlx.Tx, product Product) (err error) {
	stmt, err := tx.PrepareNamed(productQueries.updateProduct)
	if err != nil {
		logger.ErrorWithStack(err)
		return
	}
	defer stmt.Close()

	_, err = stmt.Exec(product)
	if err != nil {
		logger.ErrorWithStack(err)
	}

	return
}
