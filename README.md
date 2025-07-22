
# Aplikasi Donasi Saweria (Versi Lokal)
Aplikasi ini adalah implementasi sistem donasi seperti *Saweria* dengan dukungan client-server berbasis TCP, UDP, dan WebSocket. Aplikasi memungkinkan pengguna untuk melihat saldo, melakukan top-up, dan mengirim donasi dengan pesan yang akan disiarkan melalui WebSocket.

## Fitur Utama
- **Login sederhana** dengan username dan password
- **Cek saldo** menggunakan protokol UDP
- **Top-up saldo** menggunakan protokol UDP
- **Kirim donasi** menggunakan protokol TCP
- **Broadcast donasi** ke client WebSocket dalam format real-time

## Teknologi yang Digunakan
- Golang (Go)
- Protokol UDP & TCP
- WebSocket (menggunakan `gorilla/websocket`)
- JSON untuk komunikasi data

## Cara Menjalankan

### 1. Jalankan Server
Server terdiri dari tiga bagian dan akan dijalankan secara paralel dari satu binary:
```
go run main.go
```
Server akan berjalan di:
- UDP (dompet): `localhost:8080`
- TCP (donasi): `localhost:9090`
- WebSocket: `localhost:3000`

### 2. Jalankan Client
```
go run client.go
```
Masukkan username dan password, lalu pilih opsi:
- **1**: Cek saldo
- **2**: Top-up saldo
- **3**: Kirim pesan donasi

### 3. Jalankan Client WebSocket (Opsional)
Gunakan tools seperti Postman, browser, atau aplikasi frontend dengan koneksi ke:
```
ws://localhost:3000/ws
```
Akan menerima pesan donasi secara real-time dalam format JSON.

## Contoh Format Pesan JSON
```json
{
  "message": "Pesan donasi dari user123: Terima kasih!",
  "amount": 10000
}
```

## Struktur File
- `main.go` — berisi logic server
- `client.go` — berisi logic client
- `websocket.go` — berisi handler WebSocket dan fungsi broadcast

## Lisensi
MIT License

Copyright (c) 2025

Dengan menggunakan aplikasi ini, Anda bebas mengubah, menggunakan, dan menyebarkan kode sumber dengan tetap menyertakan lisensi MIT ini.
