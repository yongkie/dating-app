## Dating App

### Struktur Proyek

```plaintext
dating-app/
├── main.go              # Main application entry point
├── .env                 # Environment variables
├── go.mod               # Go module file
├── go.sum               # Go dependencies
├── README.md            # Documentation
├── service/
│   ├── UserService.go      # User service interface
│   └── UserServiceImpl.go  # User service implementation
├── models/
│   └── User.go             # User model
├── db/
│   └── db.go               # Database connection and queries
```

---

## API Endpoints

### **1. User Signup**
**Endpoint:** `/signup`  
**Method:** `POST`  
**Description:** Mendaftarkan pengguna baru.  

**Request Body:**
```json
{
  "username": "johndoe",
  "password": "securepassword"
}
```

**Responses:**
- `201 Created` – Signup successful
- `400 Bad Request` – Invalid request payload

---

### **2. User Login**
**Endpoint:** `/login`  
**Method:** `POST`  
**Description:** Autentikasi pengguna dan mengembalikan token sesi.  

**Request Body:**
```json
{
  "username": "johndoe",
  "password": "securepassword"
}
```

**Responses:**
- `200 OK` – Login successful (Token dikembalikan)
- `400 Bad Request` – Invalid request payload

---

### **3. Swipe Action**
**Endpoint:** `/swipe`  
**Method:** `POST`  
**Description:** Merekam tindakan swipe dari pengguna.  

**Request Body:**
```json
{
  "user_id": 1,
  "target_id": 2,
  "action": "right"
}
```

**Responses:**
- `200 OK` – Swipe action recorded
- `400 Bad Request` – Invalid request payload
- `500 Internal Server Error` – Database error

---

### **4. Purchase Action**
**Endpoint:** `/purchase`  
**Method:** `POST`  
**Description:** Menangani transaksi untuk menghapus batasan swipe dan menambahkan label verifikasi.  

**Request Body:**
```json
{
  "user_id": 1,
  "amount": 9.99
}
```

**Responses:**
- `200 OK` – Purchase action completed
- `400 Bad Request` – Invalid request payload
- `500 Internal Server Error` – Database error

---

## **Database Schema**

### **Users Table**
| Column       | Type         | Description |
|-------------|-------------|-------------|
| `id`        | INT (PK)     | Primary key |
| `username`  | VARCHAR(50)  | Unique username |
| `password`  | TEXT         | Hashed password |
| `swipes`    | INT          | Number of swipes made today |
| `last_swipe` | TIMESTAMP   | Timestamp of last swipe |
| `verified`  | BOOLEAN      | User verification status |

### **Swipes Table**
| Column       | Type         | Description |
|-------------|-------------|-------------|
| `id`        | INT (PK)     | Primary key |
| `user_id`   | INT (FK)     | Foreign key to Users table |
| `target_id` | INT          | ID of the swiped user |
| `action`    | VARCHAR(10)  | Swipe action (left/right) |
| `created_at` | TIMESTAMP   | Timestamp of the swipe |

---

## **Cara Menjalankan Service**

1. Clone repository:
   ```bash
   git clone https://github.com/your-repo/dating-app.git
   cd dating-app
   ```

2. Install dependencies:
   ```bash
   go mod tidy
   ```

3. Konfigurasi database:
   - Buat database sesuai dengan schema di atas.
   - Perbarui koneksi database di `db/db.go`.

4. Jalankan service:
   ```bash
   go run main.go
   ```

---

## **Lisensi**
Proyek ini dilisensikan di bawah **MIT License**.
