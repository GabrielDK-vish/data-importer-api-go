CREATE TABLE IF NOT EXISTS partners (
    id SERIAL PRIMARY KEY,
    partner_id VARCHAR(255) UNIQUE NOT NULL,
    partner_name VARCHAR(255) NOT NULL,
    mpn_id VARCHAR(255),
    tier2_mpn_id VARCHAR(255),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_partners_partner_id ON partners(partner_id);
