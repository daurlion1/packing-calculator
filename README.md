# 📦 Packing Calculator

A web-based tool to calculate optimal product pack combinations based on user-defined pack sizes.

This app is built using **Go (Golang)** for the backend and **vanilla HTML/JS** for the frontend.  
It fulfills product orders using **whole packs only**, following specific optimization rules.

---

## 🚀 Features

- 🔧 **Configure Pack Sizes**  
  Easily set or update available pack sizes via `POST /pack-sizes`.

- 📥 **Calculate Optimal Packs**  
  Submit an order amount and receive the most efficient pack combination via `POST /calculate`.

- 📤 **View Current Pack Sizes**  
  Fetch current pack sizes via `GET /pack-sizes`.

---

## 📦 Packaging Rules

1. **Only full packs** can be shipped. No partial packs.
2. **Minimize the total number of items** shipped while fulfilling the order.
3. If multiple combinations result in the same total items, **minimize the number of packs**.

> ⚠️ Rule 2 has higher priority than Rule 3.

---

## 🛠️ Tech Stack

- **Golang** (backend API and core logic)
- **HTML + JS** (UI)
- **Docker & Docker Compose** (for containerized deployment)

---

## 🐳 Run with Docker Compose

> ✅ No need for Go or Node.js installed locally

```bash
docker-compose up --build
```

Once the containers are up, visit:  
🌐 [http://localhost:8080](http://localhost:8080)

---

## 🧪 Running Tests (Optional)

If you're running tests locally:

```bash
go test ./internal/packing
```

---

## 📁 Project Structure

```
.
├── cmd/server          # Go HTTP server entrypoint
├── internal/packing    # Core business logic for packing calculation
├── config              # JSON file with current pack sizes
├── ui                  # Static HTML/JS frontend
├── Dockerfile          # Multi-stage Docker build
├── docker-compose.yml  # Container orchestration
```

---

## ✅ Example Usage

### 1. Submit Pack Sizes

```http
POST /pack-sizes
Content-Type: application/json

{
  "sizes": [250, 500, 1000, 2000, 5000]
}
```

### 2. Calculate Packs

```http
POST /calculate
Content-Type: application/json

{
  "amount": 12001
}
```

**Response:**
```json
{
  "packs": {
    "5000": 2,
    "2000": 1,
    "250": 1
  }
}
```

---

Happy packing! 🎁