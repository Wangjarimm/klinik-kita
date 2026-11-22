package ws

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"sync"

	"github.com/gorilla/websocket"
)

// Snapshot kondisi antrean satu dokter pada satu tanggal.
type Snapshot struct {
	DokterID       uint   `json:"dokter_id"`
	Tanggal        string `json:"tanggal"`
	SedangMelayani uint   `json:"sedang_melayani"`
	SisaMenunggu   int    `json:"sisa_menunggu"`
}

// Pesan notifikasi personal untuk satu user.
type NotifikasiMsg struct {
	Judul string `json:"judul"`
	Pesan string `json:"pesan"`
	Tipe  string `json:"tipe"` // info | sukses | peringatan
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // shortcut: origin dicek sama seperti cookie di middleware saat upgrade dilakukan handler
	},
}

// Hub mengelola koneksi per room. Room antrean: "antrean:<dokter>:<tanggal>".
// Room notifikasi: "user:<id>".
type Hub struct {
	mu    sync.RWMutex
	rooms map[string]map[*websocket.Conn]struct{}
}

func NewHub() *Hub {
	return &Hub{rooms: map[string]map[*websocket.Conn]struct{}{}}
}

func (h *Hub) Join(room string, conn *websocket.Conn) {
	h.mu.Lock()
	defer h.mu.Unlock()
	if h.rooms[room] == nil {
		h.rooms[room] = map[*websocket.Conn]struct{}{}
	}
	h.rooms[room][conn] = struct{}{}
}

func (h *Hub) Leave(room string, conn *websocket.Conn) {
	h.mu.Lock()
	defer h.mu.Unlock()
	if set, ok := h.rooms[room]; ok {
		delete(set, conn)
		if len(set) == 0 {
			delete(h.rooms, room)
		}
	}
}

func (h *Hub) Broadcast(room string, payload any) {
	data, err := json.Marshal(payload)
	if err != nil {
		return
	}
	h.mu.RLock()
	peers := make([]*websocket.Conn, 0, len(h.rooms[room]))
	for conn := range h.rooms[room] {
		peers = append(peers, conn)
	}
	h.mu.RUnlock()

	for _, conn := range peers {
		if err := conn.WriteMessage(websocket.TextMessage, data); err != nil {
			log.Printf("ws write gagal: %v", err)
			h.Leave(room, conn)
		}
	}
}

// ServeAntrean upgrade + join room, lalu kirim snapshot awal. Blok sampai koneksi tutup.
func (h *Hub) ServeAntrean(w http.ResponseWriter, r *http.Request, room string, awal Snapshot) {
	h.serveRoom(w, r, room, func(conn *websocket.Conn) error {
		return conn.WriteJSON(awal)
	})
}

// ServeNotifikasi upgrade + join room user, tanpa pesan awal.
func (h *Hub) ServeNotifikasi(w http.ResponseWriter, r *http.Request, room string) {
	h.serveRoom(w, r, room, func(*websocket.Conn) error { return nil })
}

func (h *Hub) serveRoom(w http.ResponseWriter, r *http.Request, room string, onJoin func(*websocket.Conn) error) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}
	defer conn.Close()
	h.Join(room, conn)
	defer h.Leave(room, conn)

	if err := onJoin(conn); err != nil {
		return
	}
	for {
		if _, _, err := conn.ReadMessage(); err != nil {
			return
		}
	}
}

func RoomUser(userID uint) string {
	return "user:" + strconv.FormatUint(uint64(userID), 10)
}

// NotifyUser mengirim pesan ke semua koneksi milik satu user.
func (h *Hub) NotifyUser(userID uint, msg NotifikasiMsg) {
	h.Broadcast(RoomUser(userID), map[string]any{
		"tipe":  "notifikasi",
		"judul": msg.Judul,
		"pesan": msg.Pesan,
		"label": msg.Tipe,
	})
}
