CREATE TABLE IF NOT EXISTS usages (
    id SERIAL PRIMARY KEY,
    invoice_number VARCHAR(255),
    charge_start_date DATE,
    usage_date DATE NOT NULL,
    quantity DECIMAL(15,6) NOT NULL,
    unit_price DECIMAL(15,6) NOT NULL,
    billing_pre_tax_total DECIMAL(15,2) NOT NULL,
    resource_location VARCHAR(255),
    tags TEXT,
    benefit_type VARCHAR(100),
    partner_id INTEGER REFERENCES partners(id) ON DELETE CASCADE,
    customer_id INTEGER REFERENCES customers(id) ON DELETE CASCADE,
    product_id INTEGER REFERENCES products(id) ON DELETE CASCADE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_usages_partner_id ON usages(partner_id);
CREATE INDEX idx_usages_customer_id ON usages(customer_id);
CREATE INDEX idx_usages_product_id ON usages(product_id);
CREATE INDEX idx_usages_usage_date ON usages(usage_date);
CREATE INDEX idx_usages_charge_start_date ON usages(charge_start_date);
CREATE INDEX idx_usages_invoice_number ON usages(invoice_number);
