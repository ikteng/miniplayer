// Command server is the Go port of backend/app/main.py: it wires up the
// same HTTP routes, static file serving, CORS-allow-all policy, and the
// now-playing detector/update loop as the Python FastAPI app, listening on
// the same default port 9863.
package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"miniplayer-backend/internal/api"
	"miniplayer-backend/internal/config"
	"miniplayer-backend/internal/detectors"
	"miniplayer-backend/internal/services"
	"miniplayer-backend/internal/ws"
)

func main() {
	cfg := config.Load()
	log.Printf("miniplayer-backend starting: assets=%s settings=%s static=%s port=%d",
		cfg.AssetsDir, cfg.SettingsDir, cfg.StaticDir, cfg.ServerPort)

	settingsSvc := services.NewSettingsService(cfg.SettingsDir)
	coverSvc := services.NewCoverService(cfg.AssetsDir)
	outputSvc := services.NewOutputService(settingsSvc)
	nowPlayingSvc := services.NewNowPlayingService()

	detector := detectors.NewAudioDetector(cfg.AssetsDir)
	updater := services.NewAppUpdater(detector, nowPlayingSvc, outputSvc)

	connManager := ws.NewConnectionManager()

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(corsAllowAll)

	// Route registration order mirrors main.py's router-include order: the
	// explicit page/widget routes are registered before the StaticFiles
	// mount equivalent, so they take precedence over the generic static
	// handler for the exact paths they claim ("/", "/index", etc.).
	(&api.PageHandlers{StaticDir: cfg.StaticDir}).Register(r)
	(&api.PublicHandlers{}).Register(r)
	(&api.WidgetHandlers{StaticDir: cfg.StaticDir, Settings: settingsSvc}).Register(r)
	(&api.QueryHandlers{NowPlaying: nowPlayingSvc, Cover: coverSvc}).Register(r)
	(&api.SettingsHandlers{Settings: settingsSvc}).Register(r)
	(&api.AudioHandlers{}).Register(r)
	(&api.SystemHandlers{}).Register(r)
	(&api.ServerHandlers{}).Register(r)
	(&api.PluginsHandlers{}).Register(r)
	(&api.OutputHandlers{Settings: settingsSvc}).Register(r)
	(&api.ProfileHandlers{}).Register(r)
	(&api.WebsocketLyricHandlers{Manager: connManager}).Register(r)

	// Static file serving of frontend/dist, matching main.py's
	// `app.mount("/static", StaticFiles(...))` (and its secondary mount at
	// "/" which the explicit routes above already take precedence over).
	fileServer := http.FileServer(http.Dir(cfg.StaticDir))
	r.Handle("/static/*", http.StripPrefix("/static", fileServer))
	r.Handle("/assets/*", fileServer)

	detector.Start(updater.OnUpdate)

	srv := &http.Server{
		Addr:    ":" + strconv.Itoa(cfg.ServerPort),
		Handler: r,
	}

	go func() {
		log.Printf("listening on %s", srv.Addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("server error: %v", err)
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	<-stop

	log.Println("shutting down")
	detector.Stop()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_ = srv.Shutdown(ctx)
}

// corsAllowAll mirrors main.py's CORSMiddleware(allow_origins=["*"],
// allow_credentials=True, allow_methods=["*"], allow_headers=["*"]).
func corsAllowAll(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "*")
		w.Header().Set("Access-Control-Allow-Headers", "*")
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}
		next.ServeHTTP(w, r)
	})
}

