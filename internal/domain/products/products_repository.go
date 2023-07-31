package products

import (
	"fmt"
	"github.com/evermos/boilerplate-go/infras"
	"github.com/evermos/boilerplate-go/shared/failure"
	"github.com/evermos/boilerplate-go/shared/logger"
	"github.com/gofrs/uuid"
	"github.com/jmoiron/sqlx"
	"strings"
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
			q.status,
			p.updatedBy
		FROM products p
		INNER JOIN brands b ON p.brandId = b.brandId
		INNER JOIN variants v ON p.variantId = v.variantId
		LEFT JOIN images i ON p.productId = i.productId
		LEFT JOIN (
			SELECT
				productId,
				SUM(quantity) AS quantity,
				status
			FROM quantity
			GROUP BY productId, status
		) q ON p.productId = q.productId AND q.status = 'In Stock'`,
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
	Create(product Product) (err error)
	ExistsByID(id uuid.UUID) (exists bool, err error)
	ResolveByID(id uuid.UUID) (product Product, err error)
	Update(product Product) (err error)
}

type ProductRepositoryMySQL struct {
	DB *infras.MySQLConn
}

func ProvideProductRepositoryMySQL(db *infras.MySQLConn) *ProductRepositoryMySQL {
	return &ProductRepositoryMySQL{DB: db}
}
func (p *ProductRepositoryMySQL) Create(product Product) (err error) {
	exists, err := p.ExistsByID(product.ProductId)
	if err != nil {
		logger.ErrorWithStack(err)
		return
	}
	if exists {
		err = failure.Conflict("create", "products", "already exists")
		logger.ErrorWithStack(err)
		return
	}
	return p.DB.WithTransaction(func(tx *sqlx.Tx, e chan error) {
		if err := p.txCreate(tx, product); err != nil {
			e <- err
			return
		}

		//if err := p.txCreateItems(tx, product.); err != nil {
		//	e <- err
		//	return
		//}

		e <- nil
	})
}
func (p *ProductRepositoryMySQL) ExistsByID(id uuid.UUID) (exists bool, err error) {
	err = p.DB.Read.Get(
		&exists,
		"SELECT COUNT(productId) FROM products p WHERE p.productId = ?",
		id.String())
	if err != nil {
		logger.ErrorWithStack(err)
	}

	return
}
func (p *ProductRepositoryMySQL) ResolveByID(id uuid.UUID) (product Product, err error) {
	return
}
func (p *ProductRepositoryMySQL) Update(product Product) (err error) {
	exists, err := p.ExistsByID(product.ProductId)
	if err != nil {
		logger.ErrorWithStack(err)
		return
	}
	if !exists {
		err = failure.NotFound("product")
		logger.ErrorWithStack(err)
		return
	}

	return nil
}
func (p *ProductRepositoryMySQL) composeBulkInsertProductQuery(images []Image) (query string, params []interface{}, err error) {
	values := []string{}
	for _, img := range images {
		param := map[string]interface{}{
			"ImageId":   img.ImageId,
			"productId": img.ProductId,
			"ImageUrl":  img.ImageURL,
		}
		q, args, err := sqlx.Named(productQueries.insertImagePlaceholder, param)
		if err != nil {
			return query, params, err
		}
		values = append(values, q)
		params = append(params, args...)
	}
	query = fmt.Sprintf("%v,%v", productQueries.insertImage, strings.Join(values, ","))
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
func (p *ProductRepositoryMySQL) txCreateImage(tx *sqlx.Tx, images []Image) (err error) {
	if len(images) == 0 {
		return
	}
	query, args, err := p.composeBulkInsertProductQuery(images)
	if err != nil {
		return
	}
	stmt, err := tx.Preparex(query)
	if err != nil {
		return
	}
	defer stmt.Close()
	_, err = stmt.Stmt.Exec(args...)
	if err != nil {
		logger.ErrorWithStack(err)
	}

	return
}

func (r *ProductRepositoryMySQL) ListProducts(filter FilterOptions, sorting SortingOptions, pagination PaginationOptions) (products []Product, err error) {
	selectProducts := productQueries.selectProducts

	args := make(map[string]interface{})

	if filter.BrandName != "" {
		selectProducts += " AND b.brandName LIKE :brandName"
		args["brandName"] = "%" + filter.BrandName + "%"
	}
	if filter.ProductName != "" {
		selectProducts += " AND p.productName ILIKE :productName"
		args["productName"] = "%" + filter.ProductName + "%"
	}
	if filter.VariantName != "" {
		selectProducts += " AND v.variantName ILIKE :variantName"
		args["variantName"] = "%" + filter.VariantName + "%"
	}
	if filter.Status != "" {
		selectProducts += " AND q.status = :status"
		args["status"] = filter.Status
	}

	switch sorting.By {
	case "createdAt":
		selectProducts += " ORDER BY p.createdAt"
	case "stock":
		selectProducts += " ORDER BY q.quantity"
	}

	if sorting.Ascending {
		selectProducts += " ASC"
	} else {
		selectProducts += " DESC"
	}

	selectProducts += " LIMIT :pageSize OFFSET :offset"
	args["pageSize"] = pagination.PageSize
	args["offset"] = (pagination.Page - 1) * pagination.PageSize

	err = r.DB.Read.Select(&products, r.DB.Read.Rebind(selectProducts), args)
	if err != nil {
		logger.ErrorWithStack(err)
		return nil, err
	}

	return products, nil
}
