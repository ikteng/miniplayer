package api

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/gorilla/websocket"

	"miniplayer-backend/internal/ws"
)

// WebsocketLyricHandlers mirrors app/api/websocket_lyric.py and the
// app/ws/{lyric_ws,websocket_lyric}.py aliases: registers the lyric
// WebSocket endpoint under both /api/ws/lyric and /ws/lyric.
type WebsocketLyricHandlers struct {
	Manager *ws.ConnectionManager
}

var upgrader = websocket.Upgrader{
	// CORS-allow-all, matching main.py's CORSMiddleware(allow_origins=["*"]).
	CheckOrigin: func(r *http.Request) bool { return true },
}

func (h *WebsocketLyricHandlers) Register(r chi.Router) {
	r.Get("/api/ws/lyric", h.serve)
	r.Get("/ws/lyric", h.serve)
}

func (h *WebsocketLyricHandlers) serve(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}
	h.Manager.Connect(conn)
	defer func() {
		h.Manager.Disconnect(conn)
		_ = conn.Close()
	}()

	// Mirror the Python handler's loop: just keep the connection open,
	// reading (and discarding) any client messages until it closes. All
	// outbound lyric pushes come from Manager.Broadcast, called elsewhere
	// (e.g. from the now-playing update loop) rather than from this read loop.
	for {
		if _, _, err := conn.ReadMessage(); err != nil {
			return
		}
	}
}
