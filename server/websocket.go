package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

// Struktur pesan donasi
type DonationMessage struct {
	Message string `json:"message"`
	Amount  int    `json:"amount"`
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

// Daftar klien WebSocket dan mutex untuk menghindari race condition
var clients = make([]*websocket.Conn, 0)
var clientsMutex sync.Mutex // Mutex untuk melindungi akses ke clients

func startWebSocketServer() {
	// Menjalankan server WebSocket
	http.HandleFunc("/ws", handleWebSocket)
	fmt.Println("Server WebSocket berjalan di port 3000...")
	if err := http.ListenAndServe(":3000", nil); err != nil {
		fmt.Println("Gagal memulai server WebSocket:", err)
	}
}

func handleWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println("Gagal meng-upgrade koneksi:", err)
		return
	}
	defer conn.Close()

	// Tambahkan klien baru ke daftar clients dengan proteksi mutex
	// clientsMutex.Lock()
	// clients = append(clients, conn)
	// clientsMutex.Unlock()

	// fmt.Println("Client WebSocket terhubung")

	// Tambahkan klien baru ke daftar klien
	addClient(conn)
	fmt.Println("Client WebSocket terhubung")

	// Hapus klien dari daftar saat koneksi ditutup
	defer removeClient(conn)

	// Hapus klien dari daftar ketika koneksi ditutup
	for {
		_, _, err := conn.ReadMessage()
		if err != nil {
			fmt.Println("Client WebSocket terputus:", err)
			break
		}
	}
}

// Fungsi untuk menambahkan klien baru ke daftar klien dengan proteksi mutex
func addClient(conn *websocket.Conn) {
	clientsMutex.Lock()
	defer clientsMutex.Unlock()
	clients = append(clients, conn)
}

// Fungsi untuk menghapus klien dari daftar klien dengan proteksi mutex
func removeClient(conn *websocket.Conn) {
	clientsMutex.Lock()
	defer clientsMutex.Unlock()
	for i := 0; i < len(clients); i++ {
		if clients[i] == conn {
			clients = append(clients[:i], clients[i+1:]...)
			break
		}
	}
}

// clientsMutex.Lock()
// for i := 0; i < len(clients); i++ {
// 	if clients[i] == conn {
// 		clients = append(clients[:i], clients[i+1:]...)
// 		break
// 	}
// }
// clientsMutex.Unlock()

// Fungsi untuk mem-broadcast pesan donasi ke semua klien WebSocket
func broadcastMessage(message string, amount int) {
	donation := DonationMessage{
		Message: message,
		Amount:  amount,
	}

	data, err := json.Marshal(donation)
	if err != nil {
		fmt.Println("Error encoding donation message:", err)
		return
	}

	// Broadcast ke setiap klien
	clientsMutex.Lock()
	defer clientsMutex.Unlock()

	for i := 0; i < len(clients); i++ {
		conn := clients[i]
		err := conn.WriteMessage(websocket.TextMessage, data)
		if err != nil {
			fmt.Println("Error broadcasting message:", err)
			conn.Close()
			// Hapus klien yang terputus dari daftar
			clients = append(clients[:i], clients[i+1:]...)
			i-- // Kurangi indeks karena slice telah bergeser
		}
	}
}
