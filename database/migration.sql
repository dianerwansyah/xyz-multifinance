-- Buat database
CREATE DATABASE IF NOT EXISTS xyz_db;
USE xyz_db;

-- Tabel Users
CREATE TABLE IF NOT EXISTS users (
    id INT AUTO_INCREMENT PRIMARY KEY,
    username VARCHAR(50) NOT NULL UNIQUE,
    password_hash VARCHAR(255) NOT NULL,
    role VARCHAR(20) NOT NULL DEFAULT 'user',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL DEFAULT NULL
);

-- Tabel Customers
CREATE TABLE IF NOT EXISTS customers (
    id INT AUTO_INCREMENT PRIMARY KEY,
    user_id INT NOT NULL,
    nik VARCHAR(20) NOT NULL UNIQUE,
    full_name VARCHAR(100) NOT NULL,
    legal_name VARCHAR(100) NOT NULL,
    birth_place VARCHAR(100),
    birth_date DATE,
    salary BIGINT,
    photo_ktp VARCHAR(255),
    photo_selfie VARCHAR(255),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL DEFAULT NULL,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

-- Tabel Limits
CREATE TABLE IF NOT EXISTS limits (
    id INT AUTO_INCREMENT PRIMARY KEY,
    customer_id INT NOT NULL,
    tenor_month INT NOT NULL,
    limit_amount BIGINT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL DEFAULT NULL,
    FOREIGN KEY (customer_id) REFERENCES customers(id) ON DELETE CASCADE
);

-- Tabel Transactions
CREATE TABLE IF NOT EXISTS transactions (
    id INT AUTO_INCREMENT PRIMARY KEY,
    customer_id INT NOT NULL,
    contract_number VARCHAR(50) NOT NULL UNIQUE,
    tenor INT NOT NULL,
    otr BIGINT NOT NULL,
    admin_fee BIGINT NOT NULL,
    installment_amount BIGINT NOT NULL,
    interest_amount BIGINT NOT NULL,
    asset_name VARCHAR(100) NOT NULL,
    status VARCHAR(50) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL DEFAULT NULL,
    FOREIGN KEY (customer_id) REFERENCES customers(id) ON DELETE CASCADE
);

-- INSERT dummy customer for development

-- Dummy Admin
INSERT INTO users (username, password_hash, role, created_at)
VALUES (
    'dian',
    '$2a$10$7uRtHKrklb8BPzmGWAVtJuWPfKebl.eoCveYKAH4elXQpniGAA8IW', //password123
    'admin',
    NOW()
);

-- Ambil ID dari user 'cust001' untuk referensi ke customers
SET @cust_user_id = (SELECT id FROM users WHERE username = 'cust001' LIMIT 1);

INSERT INTO customers (
    user_id, nik, full_name, legal_name, birth_place, birth_date, salary,
    photo_ktp, photo_selfie, created_at, updated_at
) VALUES (
    1,
    '12345',
    'Dian Erwansyah',
    'Dian Erwansyah',
    'Tarakan',
    '1996-05-03',
    8000000,
    '',
    '',
    NOW(),
    NOW()
);

-- customer_id = 1 (Dian Erwansyah)
INSERT INTO limits (customer_id, tenor_month, limit_amount, created_at, updated_at)
VALUES
(1, 1, 5000000, NOW(), NOW()),
(1, 3, 15000000, NOW(), NOW());


-- Misal customer_id = 1
INSERT INTO transactions (
    customer_id, contract_number, tenor, otr, admin_fee,
    installment_amount, interest_amount, asset_name,
    created_at, updated_at
) VALUES (
    1,
    'CNTR202507001',
    3,
    20000000,
    1000000,
    5000000,
    2000000,
    'Yamaha NMAX 2024',
    'disbursed',
    NOW(),
    NOW()
);