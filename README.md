# go-explore-api

Proyek pembelajaran untuk membangun **REST API** menggunakan **Go murni** (native stdlib) — tanpa framework seperti Echo, Gin, atau Fiber. Cocok untuk memahami cara kerja HTTP server, routing, middleware, dan JSON encoding dari dasar.

---

## Tujuan Repositori

Repositori ini dirancang sebagai **bahan eksplorasi dan pembelajaran bahasa Go**, khususnya:

- Cara kerja `net/http` dan `http.ServeMux` terbaru (Go 1.22+)
- Routing berbasis method dan path parameter tanpa library eksternal
- Penulisan middleware yang bisa di-chain (Logger, CORS)
- Structured logging dengan `log/slog` (Go 1.21+)
- Pattern penulisan HTTP handler yang bersih dan idiomatic
- Pengelolaan modul dengan `go.mod`

---

## Tech Stack

| Komponen        | Detail                                        |
|-----------------|-----------------------------------------------|
| Language        | Go 1.24                                       |
| HTTP Server     | `net/http` (stdlib)                           |
| Router          | `http.ServeMux` dengan method routing (Go 1.22+) |
| JSON            | `encoding/json` (stdlib)                      |
| Logging         | `log/slog` structured logger (stdlib, Go 1.21+) |
| Dependency      | Tidak ada dependency eksternal                |

---

## Struktur Proyek

```
go-explore-api/
├── main.go               # Entry point: konfigurasi server & middleware chain
├── go.mod                # Go module definition
├── .env.example          # Contoh environment variable
│
├── handler/
│   ├── routes.go         # Registrasi semua endpoint ke ServeMux
│   ├── info.go           # GET /  — info API dan daftar endpoint
│   ├── health.go         # GET /health — health check & uptime
│   ├── echo.go           # GET+POST /echo — mirror request kembali ke caller
│   └── response.go       # Helper writeJSON()
│
└── middleware/
    └── middleware.go     # Logger, CORS, dan Chain() helper
```

---

## Instalasi & Menjalankan

### Prasyarat

- [Go 1.24+](https://go.dev/dl/) terpasang di sistem

### Clone & Run

```bash
git clone https://github.com/jakabrajadenta/go-explore-api.git
cd go-explore-api

# Jalankan langsung (tanpa build)
go run .

# Atau build dulu, lalu jalankan binary
go build -o bin/api .
./bin/api
```

Server berjalan di `http://localhost:8080` secara default.

### Konfigurasi Port

```bash
# Via environment variable
PORT=9000 go run .

# Atau salin .env.example ke .env lalu edit
cp .env.example .env
```

---

## API Endpoints

### `GET /`
Menampilkan informasi API dan daftar endpoint yang tersedia.

```bash
curl http://localhost:8080/
```

```json
{
  "name": "go-explore-api",
  "version": "1.0.0",
  "description": "A learning project for building REST APIs with native Go — no framework, just stdlib.",
  "language": "Go 1.24",
  "endpoints": [
    { "method": "GET",  "path": "/",               "description": "API info and available endpoints" },
    { "method": "GET",  "path": "/health",          "description": "Health check with uptime" },
    { "method": "GET",  "path": "/echo",            "description": "Echo query params and request headers" },
    { "method": "POST", "path": "/echo",            "description": "Echo JSON request body back to caller" },
    { "method": "GET",  "path": "/echo/{message}",  "description": "Echo a path parameter as message" }
  ]
}
```

---

### `GET /health`
Health check endpoint. Mengembalikan status server dan berapa lama server sudah berjalan (uptime).

```bash
curl http://localhost:8080/health
```

```json
{
  "status": "ok",
  "timestamp": "2026-06-03T10:00:00Z",
  "uptime": "5m32s"
}
```

---

### `GET /echo`
Memantulkan kembali query parameter dan request header dari caller.

```bash
curl "http://localhost:8080/echo?nama=jaka&kota=jakarta"
```

```json
{
  "method": "GET",
  "path": "/echo",
  "headers": {
    "User-Agent": "curl/8.7.1",
    "Accept": "*/*"
  },
  "query": {
    "nama": "jaka",
    "kota": "jakarta"
  }
}
```

---

### `POST /echo`
Memantulkan kembali JSON body yang dikirim oleh caller.

```bash
curl -X POST http://localhost:8080/echo \
  -H "Content-Type: application/json" \
  -d '{"pesan": "halo", "angka": 42}'
```

```json
{
  "method": "POST",
  "path": "/echo",
  "headers": {
    "Content-Type": "application/json",
    "User-Agent": "curl/8.7.1"
  },
  "body": {
    "pesan": "halo",
    "angka": 42
  }
}
```

---

### `GET /echo/{message}`
Memantulkan kembali pesan dari path parameter URL.

```bash
curl http://localhost:8080/echo/halo-dunia
```

```json
{
  "method": "GET",
  "path": "/echo/halo-dunia",
  "message": "halo-dunia"
}
```

---

## Fitur

| Fitur                     | Keterangan                                                         |
|---------------------------|--------------------------------------------------------------------|
| Method-aware routing      | `GET /echo` dan `POST /echo` bisa ditangani handler berbeda        |
| Path parameters           | `r.PathValue("message")` native Go 1.22                           |
| Middleware chain          | Logger + CORS dikomposisi dengan `Chain()`                         |
| Structured logging        | `log/slog` — output JSON-friendly, siap production                 |
| CORS header               | Permissive CORS untuk kemudahan development & testing              |
| Timeout konfigurasi       | `ReadTimeout`, `WriteTimeout`, `IdleTimeout` terset di server      |
| Zero external dependency  | `go.mod` tanpa entry `require` — hanya stdlib                      |

---

## Konsep Go yang Dipelajari

```
net/http       → HTTP server, handler, ServeMux, method routing, path params
encoding/json  → Marshal/Unmarshal, Encoder, Decoder
log/slog       → Structured logging dengan key-value attributes
os             → Membaca environment variable
time           → Timeout, RFC3339, Duration formatting
```

---

## Pengembangan Lanjutan (Ide)

Repositori ini bisa dikembangkan lebih jauh sebagai eksperimen:

- [ ] Middleware autentikasi sederhana (API Key via header)
- [ ] Rate limiter sederhana (in-memory per IP)
- [ ] Endpoint `POST /validate` — validasi JSON schema secara manual
- [ ] Graceful shutdown dengan `os.Signal` dan `context`
- [ ] Unit test dengan `net/http/httptest`
- [ ] Integrasi database (SQLite via `database/sql`)
- [ ] Dockerize dengan multi-stage build

---

## Lisensi

MIT
