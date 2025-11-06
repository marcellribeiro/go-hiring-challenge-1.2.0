-- Add category_id column to products table
ALTER TABLE products ADD COLUMN IF NOT EXISTS category_id INTEGER REFERENCES categories(id);

-- Assign products to categories
-- PROD001, PROD004, PROD007 -> Clothing
UPDATE products SET category_id = (SELECT id FROM categories WHERE code = 'CLOTHING')
WHERE code IN ('PROD001', 'PROD004', 'PROD007');

-- PROD002, PROD006 -> Shoes
UPDATE products SET category_id = (SELECT id FROM categories WHERE code = 'SHOES')
WHERE code IN ('PROD002', 'PROD006');

-- PROD003, PROD005, PROD008 -> Accessories
UPDATE products SET category_id = (SELECT id FROM categories WHERE code = 'ACCESSORIES')
WHERE code IN ('PROD003', 'PROD005', 'PROD008');
CREATE TABLE IF NOT EXISTS categories (
    id SERIAL PRIMARY KEY,
    code VARCHAR(32) UNIQUE NOT NULL,
    name VARCHAR(256) NOT NULL,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

