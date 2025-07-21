# 🚀 PT XYZ Multifinance Backend Service

Sistem backend untuk perusahaan pembiayaan **PT XYZ Multifinance** yang menangani data konsumen, limit pembiayaan, dan transaksi pembelian asset (white goods, motor, mobil).

---

## 🧾 Studi Kasus

PT XYZ Multifinance adalah perusahaan pembiayaan yang ingin mentransformasi sistem monolith menjadi scalable, secure, dan maintainable system. Sistem ini mencakup:

- Data pelanggan
- Limit pembiayaan berdasarkan tenor
- Transaksi
- Login berbasis JWT

---

## 🔧 Fitur

- ✅ Manajemen data konsumen
- ✅ Limit pembiayaan
- ✅ Transaksi pembiayaan (OTR, bunga, cicilan)
- ✅ Login + JWT Auth
- ✅ Pencegahan OWASP Top 10 (SQLi, file abuse, auth)
- ✅ Unit test & Docker support

---
## 🧪 Cara Menjalankan (Local Dev)
### 1. Clone Repo

```bash
git clone https://github.com/dianerrwansyah/xyz-multifinance.git
cd xyz-multifinance
```

### 2. Setup .env
APP_NAME=xyz-mulfinance
APP_PORT=8080

DB_HOST=127.0.0.1 // host.docker.internal
DB_PORT=3306
DB_USER=root
DB_PASSWORD=
DB_NAME=xyz_db

JWT_SECRET=1234
MAX_UPLOAD_SIZE_MB=5

### 3. Setup Database
Jalankan file SQL berikut ke dalam MySQL:
```bash
database/migration.sql
```

### 4. Jalankan Aplikasi
```bash
go run main.go
``` 
🐳 Jalankan via Docker (Optional)
docker-compose up --build


🧱 Teknologi
    Golang
    Gin Web Framework
    MySQL
    Docker
    JWT
    Clean Architecture
    Git Flow

📂 SQL Migration
Simpan di database/migration.sql:
<details> <summary>📄 Klik untuk lihat SQL</summary>

```bash
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

```
</details> 


## Base URL
```
/api/v1
```

---

## 1. Health Check

### GET /health
Cek status API.

**Response**

```json
{
  "status": "ok"
}
```

---

## 2. Authentication

### POST /login
Login user untuk mendapatkan JWT token.

**Request Body**

```json
{
  "username": "string",
  "password": "string"
}
```

**Response Success (200 OK)**

```json
{
  "token": "jwt_token_string"
}
```

**Response Failure (401 Unauthorized)**

```json
{
  "error": "Unauthorized: invalid credentials"
}
```

---

## 3. Customer APIs (Protected)

> Semua endpoint di sini butuh header `Authorization: Bearer <token>`

### Register User (Only Admin)
### POST /api/v1/users
📌 Hanya bisa diakses oleh user dengan role admin
**Response Success (200 OK)**
```json
{
  "id": 2,
  "username": "admin"
}
```

### GET /customers/:nik
Ambil data customer berdasarkan NIK.

**Response Success (200 OK)**

```json
{
  "id": 1,
  "nik": "123456",
  "full_name": "Dian Erwansyah",
  "legal_name": "Dian Erwansyah",
  "place_of_birth": "Tarakan",
  "date_of_birth": "1996-05-03T00:00:00Z",
  "salary": 7000000,
  "ktp_photo": "",
  "selfie_photo": "",
  "created_at": "",
  "updated_at": ""
}
```

### POST /customers
Buat customer baru.

**Request Body**

```json
{
  "user_id": 2,
  "nik": "1234567",
  "full_name": "Agustiansyah",
  "legal_name": "Agus",
  "place_of_birth": "Bunyu",
  "date_of_birth": "1972-08-13",
  "salary": 5000000,
  "ktp_photo": "",
  "selfie_photo": "",
}
```

**Response Success (201 Created)**

```json
{
  "message": "Customer created successfully"
}
```

---

## 4. Limit APIs (Protected)

### POST /limits
Buat limit baru.

**Request Body**

```json
{
  "customer_id": 1,
  "tenor": 12,
  "limit": 10000000
}
```

**Response Success (201 Created)**

```json
{
  "message": "Limit created successfully"
}
```

### GET /limits/:id
Ambil limit berdasarkan ID.

**Response Success (200 OK)**

```json
{
  "id": 1,
  "customer_id": 1,
  "tenor": 12,
  "limit": 10000000,
  "created_at": "...",
  "updated_at": "..."
}
```

### PUT /limits/:id
Update limit berdasarkan ID.

**Request Body**

```json
{
  "customer_id": 1,
  "tenor": 12,
  "limit": 15000000
}
```

**Response Success (200 OK)**

```json
{
  "message": "Limit updated successfully"
}
```

### DELETE /limits/:id
Hapus limit berdasarkan ID.

**Response Success (200 OK)**

```json
{
  "message": "Limit deleted successfully"
}
```

### GET /limits/customer/:customer_id
Ambil semua limit untuk customer tertentu.

**Response Success (200 OK)**

```json
[
  {
    "id": 1,
    "customer_id": 1,
    "tenor": 12,
    "limit": 10000000,
    "created_at": "...",
    "updated_at": "..."
  },
  {
    "id": 2,
    "customer_id": 1,
    "tenor": 24,
    "limit": 20000000,
    "created_at": "...",
    "updated_at": "..."
  }
]
```

---

## 5. Transaction APIs (Protected)

### POST /transactions
Buat transaksi baru.

**Request Body**

```json
{
  "contract_number": "CN123456",
  "customer_id": 1,
  "tenor": 12,
  "amount": 5000000,
  "otr": 5500000,
  "admin_fee": 50000,
  "installment": 450000,
  "interest": 5,
  "asset_name": "Motorcycle"
}
```

**Response Success (201 Created)**

```json
{
  "message": "Transaction created successfully"
}
```

### GET /transactions/:nik
Ambil semua transaksi berdasarkan NIK customer.

**Response Success (200 OK)**

```json
[
  {
    "id": 1,
    "contract_number": "CN123456",
    "customer_id": 1,
    "tenor": 12,
    "amount": 5000000,
    "otr": 5500000,
    "admin_fee": 50000,
    "installment": 450000,
    "interest": 5,
    "asset_name": "Motorcycle",
    "created_at": "...",
    "updated_at": "..."
  }
]
```

### PUT /transactions/:id
Update transaksi berdasarkan ID.

**Request Body**

```json
{
  "contract_number": "CN123456",
  "customer_id": 1,
  "tenor": 12,
  "amount": 6000000,
  "otr": 6500000,
  "admin_fee": 60000,
  "installment": 500000,
  "interest": 5,
  "asset_name": "Motorcycle"
}
```

**Response Success (200 OK)**

```json
{
  "message": "Transaction updated successfully"
}
```

### DELETE /transactions/:id
Hapus transaksi berdasarkan ID.

**Response Success (200 OK)**

```json
{
  "message": "Transaction deleted successfully"
}
```

---

## Notes
- Semua endpoint kecuali `/login` dan `/health` membutuhkan header Authorization Bearer token.
- Pastikan JWT token valid dan belum expired.
- Gunakan NIK sebagai identifier unik untuk customer pada beberapa endpoint.