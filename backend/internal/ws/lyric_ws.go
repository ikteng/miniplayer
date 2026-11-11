// Package ws mirrors app/ws/lyric_ws.py and app/ws/websocket_lyric.py
// (which were themselves aliases of app/api/websocket_lyric.py's
// ConnectionManager): a simple fan-out broadcaster for the lyric WebSocket
// endpoint.
package ws

import (
	"sync"

	"github.com/gorilla/websocket"
)

// ConnectionManager mirrors the Python ConnectionManager: tracks active
// connections and broadcasts JSON messages to all of them, silently
// dropping any connection that errors on send (same best-effort behavior as
// the Python `except Exception: pass`).
type ConnectionManager struct {
	mu    sync.Mutex
	conns map[*websocket.Conn]struct{}
}

func NewConnectionManager() *ConnectionManager {
	return &ConnectionManager{conns: make(map[*websocket.Conn]struct{})}
}

func (m *ConnectionManager) Connect(c *websocket.Conn) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.conns[c] = struct{}{}
}

func (m *ConnectionManager) Disconnect(c *websocket.Conn) {
	m.mu.Lock()
	defer m.mu.Unlock()
	delete(m.conns, c)
}

// Broadcast sends a JSON message to every connected client, mirroring
// ConnectionManager.broadcast.
func (m *ConnectionManager) Broadcast(message interface{}) {
	m.mu.Lock()
	conns := make([]*websocket.Conn, 0, len(m.conns))
	for c := range m.conns {
		conns = append(conns, c)
	}
	m.mu.Unlock()

	for _, c := range conns {
		if err := c.WriteJSON(message); err != nil {
			m.Disconnect(c)
			_ = c.Close()
		}
	}
}
