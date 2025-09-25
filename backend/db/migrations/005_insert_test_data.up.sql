-- Inserir dados de teste para partners
INSERT INTO partners (partner_id, partner_name, mpn_id, tier2_mpn_id) VALUES
('PARTNER001', 'Microsoft Corporation', 'MPN001', 'T2MPN001'),
('PARTNER002', 'Amazon Web Services', 'MPN002', 'T2MPN002'),
('PARTNER003', 'Google Cloud Platform', 'MPN003', 'T2MPN003')
ON CONFLICT (partner_id) DO NOTHING;

-- Inserir dados de teste para customers
INSERT INTO customers (customer_id, customer_name, customer_domain_name, country) VALUES
('CUST001', 'TechCorp Solutions', 'techcorp.com', 'Brazil'),
('CUST002', 'DataFlow Inc', 'dataflow.com', 'United States'),
('CUST003', 'CloudTech Ltd', 'cloudtech.co.uk', 'United Kingdom'),
('CUST004', 'InnovateSoft', 'innovatesoft.com', 'Canada')
ON CONFLICT (customer_id) DO NOTHING;

-- Inserir dados de teste para products
INSERT INTO products (product_id, sku_id, sku_name, product_name, meter_type, category, sub_category, unit_type) VALUES
('PROD001', 'SKU001', 'Azure Virtual Machine', 'Azure Compute', 'Compute', 'Virtual Machines', 'Standard', 'Hours'),
('PROD002', 'SKU002', 'AWS EC2 Instance', 'AWS Compute', 'Compute', 'EC2', 'General Purpose', 'Hours'),
('PROD003', 'SKU003', 'Google Cloud Storage', 'GCP Storage', 'Storage', 'Cloud Storage', 'Standard', 'GB-Month'),
('PROD004', 'SKU004', 'Azure SQL Database', 'Azure Database', 'Database', 'SQL Database', 'Standard', 'DTU-Hours')
ON CONFLICT (product_id) DO NOTHING;
