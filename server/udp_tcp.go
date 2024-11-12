package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net"
	"sync"
)

// Struktur User untuk menyimpan saldo setiap pengguna
type User struct {
	Wallet Wallet
}

// Struktur Wallet untuk saldo dan mutex
type Wallet struct {
	Balance int
	Mutex   sync.Mutex
}

// Menambahkan saldo ke dompet
func (w *Wallet) TopUp(amount int) {
	w.Mutex.Lock()
	defer w.Mutex.Unlock()
	w.Balance += amount
}

// Mendapatkan saldo saat ini
func (w *Wallet) GetBalance() int {
	w.Mutex.Lock()
	defer w.Mutex.Unlock()
	return w.Balance
}

// Struktur Request untuk menerima permintaan
type Request struct {
	Type     string `json:"type"`     // "saldo", "topup", atau "donasi"
	Username string `json:"username"` // username pengguna
	Amount   int    `json:"amount"`   // nominal uang untuk topup atau donasi
	Message  string `json:"message"`  // pesan donasi
}

// Variabel peta untuk menyimpan data pengguna
var users = map[string]*User{}
var usersMutex sync.Mutex // Mutex untuk melindungi akses ke peta pengguna

func getUser(username string) *User {
	usersMutex.Lock()
	defer usersMutex.Unlock()
	// Jika pengguna belum ada, buat entri baru
	if _, exists := users[username]; !exists {
		users[username] = &User{Wallet: Wallet{}}
	}
	return users[username]
}

func main() {
	go startUDPServer()
	go startWebSocketServer() // Memulai server WebSocket
	startTCPServer()
}

// Fungsi untuk server UDP (pengelolaan dompet)
func startUDPServer() {
	addr, err := net.ResolveUDPAddr("udp", ":8080")
	if err != nil {
		fmt.Println("Error resolving address:", err)
		return
	}

	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		fmt.Println("Error listening on UDP:", err)
		return
	}
	defer conn.Close()

	fmt.Println("Server UDP berjalan di port 8080...")

	for {
		buf := make([]byte, 1024)
		n, clientAddr, err := conn.ReadFromUDP(buf)
		if err != nil {
			fmt.Println("Error reading from UDP:", err)
			continue
		}

		var request Request
		err = json.Unmarshal(buf[:n], &request)
		if err != nil {
			fmt.Println("Error decoding request:", err)
			continue
		}

		user := getUser(request.Username) // Mendapatkan data pengguna berdasarkan username

		var response string
		if request.Type == "saldo" {
			response = fmt.Sprintf("Saldo saat ini untuk %s: %d", request.Username, user.Wallet.GetBalance())
		} else if request.Type == "topup" {
			user.Wallet.TopUp(request.Amount)
			response = fmt.Sprintf("Top-up berhasil. Saldo saat ini untuk %s: %d", request.Username, user.Wallet.GetBalance())
		} else {
			response = "Permintaan tidak dikenal."
		}

		_, err = conn.WriteToUDP([]byte(response), clientAddr)
		if err != nil {
			fmt.Println("Error sending response:", err)
		}
	}
}

// Fungsi untuk memulai server TCP (pesan donasi)
func startTCPServer() {
	listener, err := net.Listen("tcp", ":9090")
	if err != nil {
		fmt.Println("Error starting TCP server:", err)
		return
	}
	defer listener.Close()

	fmt.Println("Server TCP berjalan di port 9090...")

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting TCP connection:", err)
			continue
		}

		go handleTCPConnection(conn)
	}
}

// Menghandle koneksi TCP dan mem-broadcast pesan donasi ke WebSocket
func handleTCPConnection(conn net.Conn) {
	defer conn.Close()
	buf := make([]byte, 1024)

	n, err := conn.Read(buf)
	if err != nil && err != io.EOF {
		fmt.Println("Error membaca data TCP:", err)
		return
	}

	var request Request
	err = json.Unmarshal(buf[:n], &request)
	if err != nil {
		fmt.Println("Error decoding request:", err)
		return
	}

	user := getUser(request.Username) // Mendapatkan data pengguna berdasarkan username

	if request.Type == "donasi" {
		user.Wallet.Mutex.Lock()
		if user.Wallet.Balance < request.Amount {
			user.Wallet.Mutex.Unlock()
			conn.Write([]byte("Saldo tidak mencukupi"))
			return
		}
		user.Wallet.Balance -= request.Amount
		user.Wallet.Mutex.Unlock()

		message := fmt.Sprintf("Pesan donasi dari %s: %s", request.Username, request.Message)
		broadcastMessage(message, request.Amount) // Memanggil broadcastMessage dari websocket.go

		conn.Write([]byte("Pesan donasi berhasil diterima"))
	}
}
