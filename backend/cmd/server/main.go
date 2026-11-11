// Command server is the Go port of backend/app/main.py: it wires up the
// same HTTP routes, static file serving, CORS-allow-all policy, and the
// now-playing detector/update loop as the Python FastAPI app, listening on
// the same default port 9863.
package main

import (
	"context"
	"io/fs"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	webview "github.com/jchv/go-webview2"

	"miniplayer-backend/internal/api"
	"miniplayer-backend/internal/config"
	"miniplayer-backend/internal/detectors"
	"miniplayer-backend/internal/services"
	"miniplayer-backend/internal/staticassets"
	"miniplayer-backend/internal/ws"
)

func main() {
	cfg := config.Load()
	setupLogFile(cfg)
	log.Printf("miniplayer-backend starting: assets=%s settings=%s port=%d",
		cfg.AssetsDir, cfg.SettingsDir, cfg.ServerPort)

	dist, err := fs.Sub(staticassets.Dist(), "dist")
	if err != nil {
		log.Fatalf("embedded frontend assets unavailable: %v", err)
	}

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
	(&api.PageHandlers{Dist: dist}).Register(r)
	(&api.PublicHandlers{}).Register(r)
	(&api.WidgetHandlers{Dist: dist, Settings: settingsSvc}).Register(r)
	(&api.QueryHandlers{NowPlaying: nowPlayingSvc, Cover: coverSvc}).Register(r)
	(&api.SettingsHandlers{Settings: settingsSvc}).Register(r)
	(&api.AudioHandlers{}).Register(r)
	(&api.SystemHandlers{}).Register(r)
	(&api.ServerHandlers{}).Register(r)
	(&api.PluginsHandlers{}).Register(r)
	(&api.OutputHandlers{Settings: settingsSvc}).Register(r)
	(&api.ProfileHandlers{}).Register(r)
	(&api.WebsocketLyricHandlers{Manager: connManager}).Register(r)

	// Static file serving of the embedded frontend, matching main.py's
	// `app.mount("/static", StaticFiles(...))` (and its secondary mount at
	// "/" which the explicit routes above already take precedence over).
	fileServer := http.FileServer(http.FS(dist))
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
			log.Printf("server error: %v", err)
		}
	}()

	// The webview window owns the process lifetime: Run() blocks until the
	// window is closed, which is this app's only "quit" signal (there's no
	// console to Ctrl+C since the binary is built with -H=windowsgui).
	w := webview.New(false)
	defer w.Destroy()
	w.SetTitle("MiniPlayer")
	w.SetSize(420, 640, webview.HintNone)
	w.Navigate("http://127.0.0.1:" + strconv.Itoa(cfg.ServerPort) + "/")
	w.Run()

	log.Println("shutting down")
	detector.Stop()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_ = srv.Shutdown(ctx)
}

// setupLogFile redirects the standard logger to a file next to
// SettingsDir, since the binary is built with -H=windowsgui and has no
// console to print to.
func setupLogFile(cfg config.Settings) {
	_ = os.MkdirAll(cfg.SettingsDir, 0o755)
	f, err := os.OpenFile(filepath.Join(cfg.SettingsDir, "miniplayer.log"), os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0o644)
	if err != nil {
		return
	}
	log.SetOutput(f)
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

