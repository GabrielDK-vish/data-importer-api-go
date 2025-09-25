-- Remover dados de teste
DELETE FROM usages WHERE partner_id IN (SELECT id FROM partners WHERE partner_id IN ('PARTNER001', 'PARTNER002', 'PARTNER003'));
DELETE FROM partners WHERE partner_id IN ('PARTNER001', 'PARTNER002', 'PARTNER003');
DELETE FROM customers WHERE customer_id IN ('CUST001', 'CUST002', 'CUST003', 'CUST004');
DELETE FROM products WHERE product_id IN ('PROD001', 'PROD002', 'PROD003', 'PROD004');
