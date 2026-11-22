package ws

import (
	"net/http"
	"net/http/httptest"
	"strings"
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
