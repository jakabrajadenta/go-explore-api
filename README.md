# go-explore-api

Proyek pembelajaran untuk membangun **REST API** menggunakan **Go murni** (native stdlib + satu driver DB) — tanpa framework seperti Echo, Gin, atau Fiber. Mencakup arsitektur berlapis ala Spring Boot: handler → service → repository → database.

---

## Tujuan Repositori

Repositori ini dirancang sebagai **bahan eksplorasi dan pembelajaran bahasa Go**, mencakup:

- Cara kerja `net/http` dan `http.ServeMux` terbaru (Go 1.22+) dengan method routing & path parameter
- Arsitektur berlapis: **Handler → Service → Repository → Model/DTO**
- Standarisasi response API dengan envelope konsisten (mirip Spring `ResponseBody` / `@RestControllerAdvice`)
- Penulisan middleware yang bisa di-chain: Logger & CORS
- Structured logging dengan `log/slog` (built-in Go 1.21+) + **Trace ID** untuk pelacakan request end-to-end
- Propagasi konteks (`context.Context`) dari HTTP request hingga query database
- Koneksi PostgreSQL via `pgxpool` (connection pool)
- Validasi input sederhana tanpa library eksternal
- Pengelolaan modul dengan `go.mod`

---

## Tech Stack

| Komponen        | Detail                                                |
|-----------------|-------------------------------------------------------|
| Language        | Go 1.25.0                                             |
| HTTP Server     | `net/http` stdlib                                     |
| Router          | `http.ServeMux` method routing (Go 1.22+)             |
| JSON            | `encoding/json` stdlib                                |
| Logging         | `log/slog` structured logger + trace ID per request   |
| Database        | PostgreSQL 14+                                        |
| DB Driver       | `github.com/jackc/pgx/v5` v5.10.0 (pgxpool)          |
| Validation      | Custom — stdlib `strings` + `regexp`                 |

---

## Perbandingan Struktur: Spring Boot vs Go

| Spring Boot                | Go (project ini)                   | Peran                              |
|----------------------------|------------------------------------|------------------------------------|
| `@RestController`          | `internal/handler/`                | Menerima HTTP request & kirim response |
| `@Service`                 | `internal/service/`                | Business logic & orchestration     |
| `@Repository` / JPA        | `internal/repository/`             | Akses data ke database             |
| `@Entity`                  | `internal/model/`                  | Representasi row database          |
| `RequestDTO` / `ResponseDTO` | `internal/dto/`                  | Kontrak request & response         |
| `ResponseEntity` / `@ControllerAdvice` | `pkg/response/`        | Standarisasi bentuk response       |
| `application.properties`   | `.env` + `config/`                 | Konfigurasi environment            |

---

## Struktur Proyek

```
go-explore-api/
├── main.go                          # Entry point: wiring & server startup
├── go.mod / go.sum                  # Module & dependency lock
├── .env.example                     # Template environment variable
│
├── config/
│   └── database.go                  # DB config & pgxpool factory
│
├── internal/                        # Kode bisnis, tidak boleh diimport dari luar
│   ├── handler/                     # = Controller: terima request, delegasi ke service
│   │   ├── routes.go                # Registrasi semua route ke ServeMux
│   │   ├── user_handler.go          # CRUD /api/v1/users
│   │   ├── health_handler.go        # GET /health
│   │   ├── info_handler.go          # GET /
│   │   └── echo_handler.go          # GET+POST /echo (endpoint belajar)
│   │
│   ├── service/                     # Business logic
│   │   └── user_service.go          # Interface + implementasi UserService
│   │
│   ├── repository/                  # Data access layer
│   │   └── user_repository.go       # Interface + implementasi UserRepository (pgx)
│   │
│   ├── model/                       # = Entity: struct yang merepresentasikan tabel DB
│   │   └── user.go
│   │
│   └── dto/                         # Data Transfer Objects: request & response
│       └── user_dto.go              # CreateUserRequest, UpdateUserRequest, UserResponse
│
├── middleware/
│   └── middleware.go                # Logger (inject trace ID), CORS, Chain() helper
│
├── pkg/                             # Package reusable (bisa dipakai modul lain)
│   ├── logger/
│   │   └── logger.go                # Trace ID: generate, simpan/baca context, bungkus slog
│   └── response/
│       └── response.go              # Standard response envelope (OK, Created, BadRequest, dst.)
│
└── sql/
    ├── ddl.sql                      # Schema, table, index, trigger definition
    └── dml.sql                      # Seed / default data
```

---

## Prasyarat

- [Go 1.25.0+](https://go.dev/dl/)
- PostgreSQL 14+ berjalan di localhost

---

## Setup Database

### 1. Buat database

```sql
CREATE DATABASE go_explore;
```

### 2. Jalankan DDL (buat schema, tabel, trigger)

```bash
psql -U postgres -d go_explore -f sql/ddl.sql
```

Atau buka `sql/ddl.sql` di DBeaver dan eksekusi.

### 3. Jalankan DML (isi data default)

```bash
psql -U postgres -d go_explore -f sql/dml.sql
```

---

## Instalasi & Menjalankan

```bash
git clone https://github.com/jakabrajadenta/go-explore-api.git
cd go-explore-api

# Salin dan sesuaikan env
cp .env.example .env
# Edit .env: isi DB_PASSWORD

# Download dependency
go mod download

# Jalankan langsung
go run .

# Atau build binary dulu
go build -o bin/api .
./bin/api
```

Server berjalan di `http://localhost:8080`.

---

## Konfigurasi (`.env`)

| Variable      | Default          | Keterangan                                          |
|---------------|------------------|-----------------------------------------------------|
| `PORT`        | `8080`           | Port HTTP server                                    |
| `DB_HOST`     | `localhost`      | Host PostgreSQL                                     |
| `DB_PORT`     | `5432`           | Port PostgreSQL                                     |
| `DB_USER`     | `postgres`       | Username DB                                         |
| `DB_PASSWORD` | _(wajib diisi)_  | Password DB                                         |
| `DB_NAME`     | `go_explore`     | Nama database                                       |
| `DB_SCHEMA`   | `user_management`| Search path / schema aktif                          |
| `DB_SSLMODE`  | `disable`        | SSL mode (`disable`/`require`)                      |
| `LOG_FORMAT`  | _(text)_         | Set `json` untuk output JSON (cocok untuk produksi) |
| `LOG_LEVEL`   | `info`           | Set `debug` untuk melihat log query DB              |

---

## Standard Response Format

Semua endpoint mengembalikan envelope JSON yang konsisten:

### Response tunggal
```json
{
  "success": true,
  "message": "User retrieved successfully",
  "data": {
    "id": 1,
    "username": "john_doe",
    "email": "john@example.com",
    "full_name": "John Doe",
    "phone": "081200000002",
    "is_active": true,
    "created_at": "2026-06-03T10:00:00Z",
    "updated_at": "2026-06-03T10:00:00Z"
  },
  "meta": {
    "trace_id": "a3f2b1c4d5e6f7a8",
    "timestamp": "2026-06-03T10:05:00Z",
    "path": "/api/v1/users/1"
  }
}
```

### Response list (dengan pagination)
```json
{
  "success": true,
  "message": "Users retrieved successfully",
  "data": [...],
  "meta": {
    "trace_id": "a3f2b1c4d5e6f7a8",
    "timestamp": "2026-06-03T10:05:00Z",
    "path": "/api/v1/users",
    "pagination": {
      "page": 1,
      "per_page": 10,
      "total": 5,
      "total_pages": 1
    }
  }
}
```

### Response error validasi
```json
{
  "success": false,
  "message": "Validation failed",
  "errors": {
    "email": "email format is invalid",
    "full_name": "full_name is required"
  },
  "meta": {
    "trace_id": "a3f2b1c4d5e6f7a8",
    "timestamp": "2026-06-03T10:05:00Z",
    "path": "/api/v1/users"
  }
}
```

---

## Logging & Trace ID

Setiap HTTP request mendapat **trace ID** unik — string hex 16 karakter yang mengalir dari middleware hingga database query, sehingga seluruh proses satu request bisa dibaca end-to-end di log.

### Alur trace ID

```
[Middleware] request received   → generate trace ID, simpan ke context
[Service]    service.GetByID    → baca trace ID dari context, log operasi
[Repository] repo.FindByID      → baca trace ID dari context, log query (debug)
[Middleware] request completed  → log status & latency
```

### Contoh output log (format teks, development)

```
time=2026-06-04T10:00:00Z level=INFO  trace_id=a3f2b1c4d5e6f7a8 msg="request received"   method=GET path=/api/v1/users/1 remote=127.0.0.1:54321
time=2026-06-04T10:00:00Z level=INFO  trace_id=a3f2b1c4d5e6f7a8 msg="service.GetByID"    user_id=1
time=2026-06-04T10:00:00Z level=DEBUG trace_id=a3f2b1c4d5e6f7a8 msg="repo.FindByID"      user_id=1
time=2026-06-04T10:00:00Z level=INFO  trace_id=a3f2b1c4d5e6f7a8 msg="request completed"  method=GET path=/api/v1/users/1 status=200 latency=3ms
```

Untuk menelusuri satu request, filter berdasarkan `trace_id`:

```bash
grep "a3f2b1c4d5e6f7a8" app.log
```

### Trace ID juga tersedia di:

- **Response header** `X-Trace-Id` — client bisa membacanya untuk keperluan laporan bug
- **Response body** `meta.trace_id` — tercantum di setiap JSON response

```bash
# Lihat trace ID dari header response
curl -I http://localhost:8080/api/v1/users/1
# X-Trace-Id: a3f2b1c4d5e6f7a8
```

### Konfigurasi log

| Mode | Perintah | Kapan digunakan |
|------|----------|-----------------|
| Text (default) | `go run .` | Development — mudah dibaca manusia |
| JSON | `LOG_FORMAT=json go run .` | Produksi — mudah di-parse log aggregator |
| Debug | `LOG_LEVEL=debug go run .` | Debugging — tampilkan juga log query DB |

### Log levels per layer

| Layer | Level | Contoh |
|-------|-------|--------|
| Middleware | `INFO` | `request received`, `request completed` |
| Service | `INFO` | `service.Create`, `service.Delete` |
| Service (error) | `ERROR` | `service.Create failed error=...` |
| Repository | `DEBUG` | `repo.FindByID`, `repo.Update` |

---

## API Endpoints

### Utility

| Method | Path               | Deskripsi                              |
|--------|--------------------|----------------------------------------|
| GET    | `/`                | Info API & daftar endpoint             |
| GET    | `/health`          | Health check & uptime server           |
| GET    | `/echo`            | Mirror query params & headers          |
| POST   | `/echo`            | Mirror JSON body                       |
| GET    | `/echo/{message}`  | Mirror path parameter                  |

### User Management

| Method | Path                    | Deskripsi                      | Status Sukses |
|--------|-------------------------|--------------------------------|---------------|
| GET    | `/api/v1/users`         | List semua user (paginasi)     | 200           |
| GET    | `/api/v1/users/{id}`    | Detail user berdasarkan ID     | 200           |
| POST   | `/api/v1/users`         | Buat user baru                 | 201           |
| PUT    | `/api/v1/users/{id}`    | Update seluruh data user       | 200           |
| DELETE | `/api/v1/users/{id}`    | Hapus user                     | 200           |

---

## Contoh Request

### List Users
```bash
curl "http://localhost:8080/api/v1/users?page=1&per_page=5"
```

### Create User
```bash
curl -X POST http://localhost:8080/api/v1/users \
  -H "Content-Type: application/json" \
  -d '{
    "username": "rafi_a",
    "email": "rafi@example.com",
    "full_name": "Rafi Ahmad",
    "phone": "081200000010"
  }'
```

### Update User
```bash
curl -X PUT http://localhost:8080/api/v1/users/1 \
  -H "Content-Type: application/json" \
  -d '{
    "username": "admin_v2",
    "email": "admin@example.com",
    "full_name": "Administrator V2",
    "phone": "081200000001",
    "is_active": true
  }'
```

### Delete User
```bash
curl -X DELETE http://localhost:8080/api/v1/users/5
```

---

## Validasi

Field wajib yang divalidasi:

| Field       | Create | Update | Aturan                        |
|-------------|--------|--------|-------------------------------|
| `username`  | ✓      | ✓      | Tidak boleh kosong, unik      |
| `email`     | ✓      | ✓      | Tidak boleh kosong, format valid, unik |
| `full_name` | ✓      | ✓      | Tidak boleh kosong            |
| `phone`     | —      | —      | Opsional                      |
| `is_active` | —      | opsional | Boolean, default `true` saat create |

---

## Konsep Go yang Dipelajari

```
net/http         → HTTP server, ServeMux, method routing, path params (Go 1.22)
encoding/json    → NewDecoder / NewEncoder
log/slog         → Structured logging dengan attributes key-value, JSON/text handler
context          → Propagasi nilai (trace ID) dari HTTP layer hingga database query
crypto/rand      → Generate trace ID yang aman secara kriptografi
encoding/hex     → Encode bytes menjadi string hex yang readable
os               → Membaca environment variable
time             → Timeout, RFC3339, Duration
errors           → errors.Is() untuk error sentinel
regexp           → Validasi format email
pgx/v5           → PostgreSQL driver, pgxpool, QueryRow, Query, Exec
```

---

## Pengembangan Lanjutan (Ide)

- [x] Trace ID logging end-to-end per request (`pkg/logger`)
- [ ] Graceful shutdown dengan `os.Signal` + `context.WithCancel`
- [ ] Unit test dengan `net/http/httptest` dan mock repository
- [ ] Middleware autentikasi sederhana (API Key via header)
- [ ] Rate limiter sederhana (in-memory per IP)
- [ ] Soft delete (`deleted_at TIMESTAMPTZ`)
- [ ] Search & filter di endpoint list users
- [ ] Dockerize dengan multi-stage build
- [ ] Migration tool (`golang-migrate`)

---

## Lisensi

MIT
