package api

import (
	"net/http"

	"github.com/go-chi/chi/v5"

	"miniplayer-backend/internal/models"
)

// ServerHandlers mirrors app/api/server.py. Python registers this same
// `/api/server/info` path from BOTH server.py and system.py; FastAPI keeps
// whichever router matches first by registration order. Here SystemHandlers
// already registers `/api/server/info` (see system.go), so chi's router
// would panic on a duplicate pattern if this handler tried to register the
// exact same path again. To preserve "both routes exist, same shape" without
// a duplicate-route panic, this file intentionally registers no additional
// paths - system.go's /api/server/info and /server/info already provide the
// identical response shape that server.py would have produced.
type ServerHandlers struct{}

func (h *ServerHandlers) Register(r chi.Router) {
	// No-op: see comment above. Kept as a distinct type for symmetry with
	// the Python module layout and to make this decision discoverable.
}

var _ = models.RespData{}
var _ = http.StatusOK
