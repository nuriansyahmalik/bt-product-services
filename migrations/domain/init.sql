CREATE TABLE user (
                      userId VARCHAR(36) PRIMARY KEY,
                      username VARCHAR(50) NOT NULL UNIQUE,
                      email VARCHAR(50) NOT NULL UNIQUE,
                      userType ENUM('admin', 'regular') NOT NULL,
                      createdAt TIMESTAMP NOT NULL,
                      createdBy VARCHAR(36) NOT NULL,
                      updatedAt TIMESTAMP,
                      updatedBy VARCHAR(36),
                      deletedAt TIMESTAMP,
                      deletedBy VARCHAR(36)
);


CREATE TABLE brand (
                       brandId VARCHAR(36) PRIMARY KEY,
                       brandName VARCHAR(100) NOT NULL,
                       createdAt TIMESTAMP NOT NULL,
                       createdBy VARCHAR(36) NOT NULL,
                       updatedAt TIMESTAMP,
                       updatedBy VARCHAR(36),
                       deletedAt TIMESTAMP,
                       deletedBy VARCHAR(36)
);

CREATE TABLE variant (
                         variantId VARCHAR(36) PRIMARY KEY,
                         variantName VARCHAR(100) NOT NULL,
                         brandId VARCHAR(36) NOT NULL,
                         price DECIMAL(10, 2),
                         createdAt TIMESTAMP NOT NULL,
                         createdBy VARCHAR(36) NOT NULL,
                         updatedAt TIMESTAMP,
                         updatedBy VARCHAR(36),
                         deletedAt TIMESTAMP,
                         deletedBy VARCHAR(36),
                         FOREIGN KEY (brandId) REFERENCES brand (brandId)
);

CREATE TABLE products (
                          productId VARCHAR(36) PRIMARY KEY,
                          productName VARCHAR(200) NOT NULL,
                          variantId VARCHAR(36) NOT NULL,
                          createdAt TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
                          createdBy VARCHAR(36) ,
                          updatedAt TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
                          updatedBy VARCHAR(36),
                          deletedAt TIMESTAMP,
                          deletedBy VARCHAR(36),
                          FOREIGN KEY (variantId) REFERENCES variant(variantId)
);

CREATE TABLE images (
                        imageId VARCHAR(36) PRIMARY KEY,
                        productId VARCHAR(36) NOT NULL,
                        imageUrl VARCHAR(200) NOT NULL,
                        createdAt TIMESTAMP NOT NULL,
                        createdBy VARCHAR(36) ,
                        FOREIGN KEY (productId) REFERENCES products (productId)
);

CREATE TABLE warehouses (
                            warehouseId VARCHAR(36) PRIMARY KEY,
                            warehouseName VARCHAR(100) NOT NULL,
                            createdAt TIMESTAMP NOT NULL,
                            createdBy VARCHAR(36)
);

CREATE TABLE quantity (
                          quantityId VARCHAR(36) PRIMARY KEY,
                          productId VARCHAR(36) NOT NULL,
                          warehouseId VARCHAR(36) NOT NULL,
                          quantity INT NOT NULL,
                          status VARCHAR(200) NOT NULL,
                          createdAt TIMESTAMP NOT NULL,
                          createdBy VARCHAR(36)  ,
                          updatedAt TIMESTAMP,
                          updatedBy VARCHAR(36),
                          FOREIGN KEY (productId) REFERENCES products(productId),
                          FOREIGN KEY (warehouseId) REFERENCES warehouses(warehouseId)
);

CREATE TABLE variantWarehouse (
                                  variantId VARCHAR(36),
                                  warehouseId VARCHAR(36),
                                  PRIMARY KEY (variantId, warehouseId),
                                  FOREIGN KEY (variantId) REFERENCES variant(variantId),
                                  FOREIGN KEY (warehouseId) REFERENCES warehouses(warehouseId)
);

CREATE TABLE userProducts (
                              userProductId VARCHAR(36) PRIMARY KEY,
                              userId VARCHAR(36) NOT NULL,
                              productId VARCHAR(36) NOT NULL,
                              createdAt TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
                              FOREIGN KEY (userId) REFERENCES user(userId),
                              FOREIGN KEY (productId) REFERENCES products(productId)
);

CREATE INDEX idx_product_name ON products (productName);
CREATE INDEX idx_brand_variant ON products (brandId, variantId);


CREATE TRIGGER update_product_updatedAt
    AFTER UPDATE ON Variants
    FOR EACH ROW
BEGIN
    UPDATE products
    SET updatedAt = NOW(),
        updatedBy = 1
    WHERE productId = NEW.productId;
END;


