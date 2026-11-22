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

// koneksi pembungkus: gorilla/websocket hanya mengizinkan satu penulis pada satu
// waktu, jadi setiap koneksi punya mutex tulis sendiri (I-1).
type koneksi struct {
	conn *websocket.Conn
	wmu  sync.Mutex
}

func (k *koneksi) writeJSON(v any) error {
	data, err := json.Marshal(v)
	if err != nil {
		return err
	}
	k.wmu.Lock()
	defer k.wmu.Unlock()
	return k.conn.WriteMessage(websocket.TextMessage, data)
}

// Hub mengelola koneksi per room. Room antrean: "antrean:<dokter>:<tanggal>".
// Room notifikasi: "user:<id>".
type Hub struct {
	mu    sync.RWMutex
	rooms map[string]map[*koneksi]struct{}
}

func NewHub() *Hub {
	return &Hub{rooms: map[string]map[*koneksi]struct{}{}}
}

func (h *Hub) Join(room string, k *koneksi) {
	h.mu.Lock()
	defer h.mu.Unlock()
	if h.rooms[room] == nil {
		h.rooms[room] = map[*koneksi]struct{}{}
	}
	h.rooms[room][k] = struct{}{}
}

func (h *Hub) Leave(room string, k *koneksi) {
	h.mu.Lock()
	defer h.mu.Unlock()
	if set, ok := h.rooms[room]; ok {
		delete(set, k)
		if len(set) == 0 {
			delete(h.rooms, room)
		}
	}
}

func (h *Hub) Broadcast(room string, payload any) {
	h.mu.RLock()
	peers := make([]*koneksi, 0, len(h.rooms[room]))
	for k := range h.rooms[room] {
		peers = append(peers, k)
	}
	h.mu.RUnlock()

	for _, k := range peers {
		if err := k.writeJSON(payload); err != nil {
			log.Printf("ws write gagal: %v", err)
			h.Leave(room, k)
			_ = k.conn.Close()
		}
	}
}

func (h *Hub) serveRoom(w http.ResponseWriter, r *http.Request, room string, onJoin func(*koneksi) error) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}
	k := &koneksi{conn: conn}
	defer conn.Close()
	h.Join(room, k)
	defer h.Leave(room, k)

	if err := onJoin(k); err != nil {
		return
	}
	for {
		if _, _, err := conn.ReadMessage(); err != nil {
			return
		}
	}
}

// ServeAntrean upgrade + join room, lalu kirim snapshot awal. Blok sampai koneksi tutup.
func (h *Hub) ServeAntrean(w http.ResponseWriter, r *http.Request, room string, awal Snapshot) {
	h.serveRoom(w, r, room, func(k *koneksi) error {
		return k.writeJSON(awal)
	})
}

// ServeNotifikasi upgrade + join room user, tanpa pesan awal.
func (h *Hub) ServeNotifikasi(w http.ResponseWriter, r *http.Request, room string) {
	h.serveRoom(w, r, room, func(*koneksi) error { return nil })
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
