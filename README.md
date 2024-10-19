## Dokumentasi API
## https://quiz-sanbercode-go-production.up.railway.app/
Proyek ini menggunakan **JWT Authentication**, di mana pengguna harus melakukan login terlebih dahulu untuk
mengakses endpoint API. Login dapat dilakkukan dengan data pengguna yang telah di-*seed* saat migrasi database. Pengguna yang sudah terdaftar adalah:

- **Username**: `admin`
- **Password**: `admin`

### Autentikasi

Pastikan untuk memberikan JWT token yang tersimpan di Cookies yang benar melalui header Bearer Token saat mengakses
API. Server akan memverifikasi kredensial tersebut dan memberikan akses jika cocok.

Jika Authentikasi tidak dilakukan dan mencoba untuk mengakses API Kategori, Kriteria, Produk dan Value, maka akan memberikan response sebagai berikut:

```json
  {
    "details": "token is malformed: token contains an invalid number of segments",
    "error": "Invalid token"
  }
 ```


## API User

### 1. Login

**Endpoint**: `POST /api/login`

Endpoint ini digunakan untuk mendapatkan sebuah token jwt.
- **Response** (jika berhasil login):
  ```json
  {
    "message": "login berhasil"
  }
  ```

- **Response** (jika memasukkan username atau password yang salah):
  ```json
  {
    "message": "Username atau password salah"
  }
  ```

- **Response** (jika memasukkan username atau password yang tidak sesuai dengan tipe data):
  ```json
  {
    "error": "failed to read body"
  }
  ```
  
### 2. Logout

**Endpoint**: `GET /api/logout`

Endpoint ini digunakan untuk menghapus token jwt yang tersimpan
- **Response** (jika berhasil logout):
  ```json
  {
    "message": "logout berhasil"
  }
  ```


## API Kategori

### 1. Melihat Semua Kategori

**Endpoint**: `GET /api/categories`

Endpoint ini digunakan untuk mendapatkan daftar semua kategori yang tersedia.

- **Response** (jika terdapat data):
  ```json
  {
    "categories": categories
  }
  ```

- **Response** (jika tidak ada data):
  ```json
  {
    "message": "Data kategori masih kosong
  }
  ```

### 2. Menambahkan Kategori Baru

**Endpoint**: `POST /api/categories`

Endpoint ini digunakan untuk menambahkan kategori baru.

- **Body Request**:
  ```json
  {
    "name": "Testing"
  }
  ```

- **Response** (jika berhasil):
  ```json
  {
    categories
  }
  ```

- **Response** (jika terjadi mengisi string kosong (""):
  ```json
  {
    "error": "tolong masukan nama kategori"
  }
  ```

- **Response** (jika terjadi duplikasi nama kategori):
  ```json
  {
    "message": "nama kategori sudah ada"
  }
  ```

### 3. Melihat Detail Kategori Berdasarkan ID

**Endpoint**: `GET /api/categories/:id`

Endpoint ini digunakan untuk melihat detail kategori berdasarkan ID tertentu.

- **Response** (jika berhasil):
  ```json
  {
    categories
  }
  ```

- **Response** (jika ID tidak ditemukan):
  ```json
  {
    "error": "Kategori dengan ID:%d tidak ditemukan"
  }
  ```

- **Response** (jika ID tidak valid):
  ```json
  {
    "error": "ID tidak sesuai"
  }
  ```

### 4. Memperbarui Kategori Berdasarkan ID

**Endpoint**: `PUT /api/categories/:id`

Endpoint ini digunakan untuk mengubah data kategori berdasarkan ID tertentu.

- **Body Request**:
  ```json
  {
    "name": "Update kategori testing"
  }
  ```

- **Response** (jika berhasil):
  ```json
  {
    "message": "Data kategori berhasil diperbarui"
  }
  ```

- **Response** (jika ID tidak valid):
  ```json
  {
    "error": "ID tidak sesuai"
  }
  ```

- **Response** (jika ID tidak ditemukan):
  ```json
  {
    "error": "Kategori dengan ID:%d tidak ditemukan"
  }
  ```

- **Response** (jika data yang dimasukkan sama seperti data sebelumnya):
  ```json
  {
    "message": "masukkan data nama kategori yang baru"
  }
  ```

### 5. Menghapus Kategori Berdasarkan ID

**Endpoint**: `DELETE /api/categories/:id`

Endpoint ini digunakan untuk menghapus data kategori berdasarkan ID tertentu.

- **Response** (jika berhasil):
  ```json
  {
    "message": "Kategori dengan ID:%d berhasil dihapus"
  }
  ```

- **Response** (jika ID tidak valid):
  ```json
  {
    "error": "ID tidak valid"
  }
  ```
  
- **Response** (jika ID tidak ditemukan):
  ```json
  {
    "error": "Kategori dengan ID:%d tidak ditemukan"
  }
  ```


## API Criteria

### 1. Melihat Semua Kriteria

**Endpoint**: `GET /api/criterias`

Endpoint ini digunakan untuk melihat daftar semua kriteria yang tersedia.

- **Response** (jika terdapat data):
  ```json
  {
    "criterias": criterias
  }
  ```

- **Response** (jika tidak ada data):
  ```json
  {
    "message": "data kriteria masih kosong"
  }
  ```

### 2. Menambahkan Kriteria Baru

**Endpoint**: `POST /api/criterias`

Endpoint ini digunakan untuk menambahkan kriteria baru.

- **Body Request**:
  ```json
  {
    "name": "Gross Profit Margin",
    "weight": 5,
    "type": "benefit"
  }
  ```

- **Response** (jika berhasil):
  ```json
  {
    criterias
  }
  ```

- **Response** (jika input nilai kosong atau null):
  ```json
  {
    "error": "inputan tidak valid"
  }
  ```
  
- **Response** (jika nama kriteria sudah ada):
  ```json
  {
    "message": "nama kriteria sudah ada"
  }
  ```

### 3. Melihat Detail Kriteria Berdasarkan ID

**Endpoint**: `GET /api/criterias/:id`

Endpoint ini digunakan untuk melihat detail kriteria berdasarkan ID tertentu.

- **Response** (jika berhasil):
  ```json
  {
    criterias
  }
  ```

- **Response** (jika ID tidak valid):
  ```json
  {
    "error": "ID tidak sesuai"
  }
  ```
  
- **Response** (jika ID tidak ditemukan):
  ```json
  {
    "error": "Kriteria dengan ID:%d tidak ditemukan"
  }
  ```


### 4. Memperbarui Kriteria Berdasarkan ID

**Endpoint**: `PUT /api/criterias/:id`

Endpoint ini digunakan untuk memperbarui data kriteria berdasarkan ID tertentu.

- **Body Request**:
  ```json
  {
    "name": "ROI",
    "weight": 5,
    "type": "benefit"
  }
  ```

- **Response** (jika berhasil):
  ```json
  {
    "message": "Data kriteria berhasil diperbarui"
  }
  ```

- **Response** (jika ID tidak valid):
  ```json
  {
    "error": "ID tidak sesuai"
  }
  ```

- **Response** (jika ID tidak ditemukan):
  ```json
  {
    "error": "Kriteria dengan ID:%d tidak ditemukan"
  }
  ```

- **Response** (jika gagal membaca request json):
  ```json
  {
    "message": "failed to read json"
  }
  ```

- **Response** (jika nama kriteria tidak diubah):
  ```json
  {
    "error": "masukkan data nama kriteria yang baru"
  }
  ```

- **Response** (jika nama sudah ada pada data dengan ID yang lain):
  ```json
  {
    "message": "nama kriteria sudah ada"
  }
  ```

### 5. Menghapus Kriteria Berdasarkan ID

**Endpoint**: `DELETE /api/criterias/:id`

Endpoint ini digunakan untuk menghapus kriteria berdasarkan ID tertentu.

- **Response** (jika berhasil):
  ```json
  {
    "message": "Kriteria dengan ID:%d berhasil dihapus"
  }
  ```

- **Response** (jika ID tidak valid):
  ```json
  {
    "error": "ID tidak valid"
  }
  ```
  
- **Response** (jika ID tidak ditemukan):
  ```json
  {
    "error": "Kriteria dengan ID:%d tidak ditemukan"
  }
  ```


## API Product

### 1. Melihat Semua Produk

**Endpoint**: `GET /api/products`

Endpoint ini digunakan untuk melihat daftar semua produk yang tersedia.

- **Response** (jika terdapat data):
  ```json
  {
    "products": products
  }
  ```

- **Response** (jika tidak ada data):
  ```json
  {
    "message": "data produk masih kosong"
  }
  ```

### 2. Menambahkan Produk Baru

**Endpoint**: `POST /api/products`

Endpoint ini digunakan untuk menambahkan produk baru.

- **Body Request**:
  ```json
  {
    "name": "Post Data 1",
    "net_profit": 20000,
    "gross_profit": 80000,
    "gross_sale": 200000,
    "purchase_cost": 150000,
    "initial_stock": 100,
    "final_stock": 50,
    "category_id": 3
  }
  ```

- **Response** (jika berhasil):
  ```json
  {
    products
  }
  ```

- **Response** (jika input nilai kosong atau null):
  ```json
  {
    "error": "semua field harus diisi dengan nilai yang valid"
  }
  ```
  
- **Response** (jika nama kriteria sudah ada):
  ```json
  {
    "message": "nama produk sudah ada"
  }
  ```

### 3. Melihat Detail Produk Berdasarkan ID

**Endpoint**: `GET /api/products/:id`

Endpoint ini digunakan untuk melihat detail produk berdasarkan ID tertentu.

- **Response** (jika berhasil):
  ```json
  {
    products
  }
  ```

- **Response** (jika ID tidak valid):
  ```json
  {
    "error": "ID tidak sesuai"
  }
  ```
  
- **Response** (jika ID tidak ditemukan):
  ```json
  {
    "error": "Produk dengan ID:%d tidak ditemukan"
  }
  ```


### 4. Memperbarui Kriteria Berdasarkan ID

**Endpoint**: `PUT /api/criterias/:id`

Endpoint ini digunakan untuk memperbarui data kriteria berdasarkan ID tertentu.

- **Body Request**:
  ```json
  {
    "name": "Update Data 1",
    "net_profit": 20000,
    "gross_profit": 80000,
    "gross_sale": 200000,
    "purchase_cost": 150000,
    "initial_stock": 100,
    "final_stock": 50,
    "category_id": 3
  }
  ```

- **Response** (jika berhasil):
  ```json
  {
    "message": "Data produk berhasil diperbarui"
  }
  ```

- **Response** (jika ID tidak valid):
  ```json
  {
    "error": "ID tidak sesuai"
  }
  ```

- **Response** (jika ID tidak ditemukan):
  ```json
  {
    "error": "Produk dengan ID:%d tidak ditemukan"
  }
  ```

- **Response** (jika gagal membaca request json):
  ```json
  {
    "message": "failed to read json"
  }
  ```

- **Response** (jika input nilai kosong atau null):
  ```json
  {
    "message": "semua field harus diisi dengan nilai yang valid"
  }
  ```

- **Response** (jika tidak ada data yang diubah):
  ```json
  {
    "error": "masukkan minimal satu data yang baru"
  }
  ```

- **Response** (jika nama sudah ada pada data dengan ID yang lain):
  ```json
  {
    "message": "nama produk sudah ada"
  }
  ```

### 5. Menghapus Produk Berdasarkan ID

**Endpoint**: `DELETE /api/products/:id`

Endpoint ini digunakan untuk menghapus produk berdasarkan ID tertentu.

- **Response** (jika berhasil):
  ```json
  {
    "message": "Produk dengan ID:%d berhasil dihapus"
  }
  ```

- **Response** (jika ID tidak valid):
  ```json
  {
    "error": "ID tidak valid"
  }
  ```
  
- **Response** (jika ID tidak ditemukan):
  ```json
  {
    "error": "Produk dengan ID:%d tidak ditemukan"
  }
  ```

## API Value

### 1. Melihat Semua Nilai dari Produk

**Endpoint**: `GET /api/values`

Endpoint ini digunakan untuk melihat daftar semua nilai produk yang tersedia.

- **Response** (jika terdapat data):
  ```json
  {
    "values": values
  }
  ```

- **Response** (jika tidak ada data):
  ```json
  {
    "message": "data nilai masih kosong"
  }
  ```

### 2. Menambahkan Nilai dari semua Produk yang ada

**Endpoint**: `POST /api/values`

Endpoint ini digunakan untuk menambahkan nilai dari semua produk yang ada.

- **Response** (jika berhasil):
  ```json
  {
    "message": "nilai produk berhasil dihitung untuk semua kriteria dan disimpan"
  }
  ```
  
- **Response** (jika berhasil):
  ```json
  {
    "message": "nilai produk berhasil dihitung untuk semua kriteria dan disimpan"
  }
  ```

- **Response** (jika gagal mengambil data kriteria):
  ```json
  {
    "error": "gagal mengambil data kriteria"
  }
  ```

- **Response** (jika gagal mengambil data produk):
  ```json
  {
    "error": "gagal mengambil data produk"
  }
  ```
  
- **Response** (jika mendapatkan data kriteria yang tidak dikenali):
  ```json
  {
    "error": "kriteria tidak dikenali"
  }
  ```


### 3. Menghapus Semua Nilai Produk

**Endpoint**: `DELETE /api/values`

Endpoint ini digunakan untuk menghapus seluruh nilai produk.

- **Response** (jika berhasil):
  ```json
  {
    "message": "semua data nilai berhasil dihapus"
  }
  ```
