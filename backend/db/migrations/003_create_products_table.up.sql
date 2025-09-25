CREATE TABLE IF NOT EXISTS products (
    id SERIAL PRIMARY KEY,
    product_id VARCHAR(255) UNIQUE NOT NULL,
    sku_id VARCHAR(255) NOT NULL,
    sku_name VARCHAR(255) NOT NULL,
    product_name VARCHAR(255) NOT NULL,
    meter_type VARCHAR(100),
    category VARCHAR(100),
    sub_category VARCHAR(100),
    unit_type VARCHAR(50),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_products_product_id ON products(product_id);
CREATE INDEX idx_products_sku_id ON products(sku_id);
CREATE INDEX idx_products_category ON products(category);
