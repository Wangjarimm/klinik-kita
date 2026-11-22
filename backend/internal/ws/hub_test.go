package ws

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/gorilla/websocket"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func wsURL(server *httptest.Server) string {
	return "ws" + strings.TrimPrefix(server.URL, "http")
}

// Review Focus 5: broadcast room A tidak boleh sampai ke subscriber room B.
func TestBroadcastRoomIsolation(t *testing.T) {
	hub := NewHub()

	roomA, roomB := "antrean:1:2026-11-23", "antrean:2:2026-11-23"
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		room := r.URL.Query().Get("room")
		hub.ServeAntrean(w, r, room, Snapshot{DokterID: 0, Tanggal: room, SedangMelayani: 0, SisaMenunggu: 0})
	}))
	defer srv.Close()

	dial := func(room string) *websocket.Conn {
		conn, _, err := websocket.DefaultDialer.Dial(wsURL(srv)+"?room="+room, nil)
		require.NoError(t, err)
		return conn
	}
	clientA := dial(roomA)
	defer clientA.Close()
	clientB := dial(roomB)
	defer clientB.Close()

	// snapshot awal keduanya
	var snap Snapshot
	require.NoError(t, clientA.ReadJSON(&snap))
	require.NoError(t, clientB.ReadJSON(&snap))

	hub.Broadcast(roomA, Snapshot{DokterID: 1, Tanggal: "2026-11-23", SedangMelayani: 3, SisaMenunggu: 5})

	// A menerima update
	var got Snapshot
	clientA.SetReadDeadline(time.Now().Add(2 * time.Second))
	require.NoError(t, clientA.ReadJSON(&got))
	assert.Equal(t, uint(3), got.SedangMelayani)
	assert.Equal(t, 5, got.SisaMenunggu)

	// B TIDAK menerima apa pun (timeout = benar)
	clientB.SetReadDeadline(time.Now().Add(500 * time.Millisecond))
	_, _, err := clientB.ReadMessage()
	assert.Error(t, err, "room B tidak boleh menerima broadcast room A")
}

// I-1: broadcast bersamaan ke koneksi yang sama tidak boleh menjatuhkan pesan
// (gorilla/websocket hanya mengizinkan satu penulis; tanpa serialisasi beberapa
// write gagal dan subscriber hilang tanpa sadar).
func TestBroadcastConcurrentTidakHilang(t *testing.T) {
	hub := NewHub()
	room := "antrean:9:2026-11-30"
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		hub.ServeAntrean(w, r, room, Snapshot{})
	}))
	defer srv.Close()

	conn, _, err := websocket.DefaultDialer.Dial(wsURL(srv)+"?room="+room, nil)
	require.NoError(t, err)
	defer conn.Close()
	var awal Snapshot
	require.NoError(t, conn.ReadJSON(&awal)) // snapshot awal

	const N = 20
	var wg sync.WaitGroup
	for i := 0; i < N; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			hub.Broadcast(room, map[string]int{"n": i})
		}(i)
	}
	wg.Wait()

	terima := 0
	conn.SetReadDeadline(time.Now().Add(3 * time.Second))
	for i := 0; i < N; i++ {
		if _, _, err := conn.ReadMessage(); err == nil {
			terima++
		} else {
			break
		}
	}
	assert.Equal(t, N, terima, "semua broadcast bersamaan harus sampai ke subscriber")
}
