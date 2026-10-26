// Command rivals-of-aetherd answers Rivals of Aether's first-party web hosts for
// the Switch build (0100D2700DE7C000) when a console is redirected to OpenPak.
//
// Netplay itself is NEX 4.6.4 over secure PRUDP (auth via nncs), served by
// platforms/switch/games/switch-nex once the title profile lands (see
// docs/nex.md); this service is only the developer's small HTTP side:
//
//	/api/games/start    match telemetry, sent when a game starts
//	/api/games/finish   match telemetry, sent when a game ends
//	/files/ingame_assets/news           in-game news page
//	/files/ingame_assets/news_data.txt  in-game news feed
//
// The client dials stats.rivalsofaether.com over plain HTTP and
// www.rivalsofaether.com for news. The service listens on plain HTTP (-addr)
// behind the TLS terminator. Unknown routes are logged with their body and
// answered 404, so a capture shows what the client asks next.
//
// The contract was read out of the Switch executable's strings; see
// docs/protocol.md.
package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

// Config is the service's env-driven configuration.
type Config struct {
	// Addr serves the telemetry and news stubs (behind the TLS terminator).
	Addr string
	// HealthAddr serves GET /healthz.
	HealthAddr string
	// NexSecureURL is the switch-nex secure upstream for this title, once its
	// registry profile lands. Empty means unconfigured (see docs/nex.md).
	NexSecureURL string
	// NewsDir optionally serves the news CDN files from disk; empty serves 204.
	NewsDir string
	// NXInternalURL/Key point at nx-baas for Nintendo id-token checks.
	NXInternalURL string
	NXInternalKey string
}

func configFromEnv() Config {
	return Config{
		Addr:         envOr("RIVALS_ADDR", ":23570"),
		HealthAddr:   envOr("RIVALS_HEALTH_ADDR", ":23571"),
		NexSecureURL: os.Getenv("RIVALS_NEX_SECURE_URL"),
		NewsDir:      os.Getenv("RIVALS_NEWS_DIR"),
		NXInternalURL: envOr("NX_INTERNAL_URL", "http://127.0.0.1:20070"),
		NXInternalKey: firstNonEmpty(os.Getenv("NX_INTERNAL_KEY"), os.Getenv("OPENPAK_INTERNAL_KEY")),
	}
}

func envOr(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func firstNonEmpty(vs ...string) string {
	for _, v := range vs {
		if v != "" {
			return v
		}
	}
	return ""
}

func main() {
	cfg := configFromEnv()
	s := &server{newsDir: cfg.NewsDir, nex: newNexStub(cfg)}

	go func() {
		h := http.NewServeMux()
		h.HandleFunc("/healthz", func(w http.ResponseWriter, _ *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(map[string]string{"status": "ok", "service": "rivals-of-aetherd"})
		})
		log.Printf("[rivals] health on %s", cfg.HealthAddr)
		log.Fatal(http.ListenAndServe(cfg.HealthAddr, h))
	}()

	mux := http.NewServeMux()
	mux.HandleFunc("/api/games/start", s.telemetry)
	mux.HandleFunc("/api/games/finish", s.telemetry)
	mux.HandleFunc("/files/ingame_assets/news", s.news)
	mux.HandleFunc("/files/ingame_assets/news_data.txt", s.news)
	mux.HandleFunc("/v1/nex/status", s.nexStatus)
	mux.HandleFunc("/", s.unknown)
	log.Printf("[rivals] serving telemetry+news stub on %s (newsDir=%q nex=%q)",
		cfg.Addr, cfg.NewsDir, cfg.NexSecureURL)
	log.Fatal(http.ListenAndServe(cfg.Addr, mux))
}

type server struct {
	newsDir string
	nex     *nexStub
}

// telemetry accepts a match start/finish report and keeps nothing: it is the
// player's match record, not ours to store. The request method is unconfirmed
// (static strings only name the paths), so any method is accepted; the shape
// of a real report is logged by size only until a capture pins it down.
func (s *server) telemetry(w http.ResponseWriter, r *http.Request) {
	n, _ := io.Copy(io.Discard, io.LimitReader(r.Body, 1<<20))
	log.Printf("[rivals] telemetry %s %s body=%d contentType=%q", r.Method, r.URL.Path, n, r.Header.Get("Content-Type"))
	w.Header().Set("Content-Type", "application/json")
	io.WriteString(w, "{}\n")
}

// news serves the in-game news CDN files. With RIVALS_NEWS_DIR set, files are
// served from disk by base name; otherwise 204, so the client shows no news
// instead of a stale or error page.
func (s *server) news(w http.ResponseWriter, r *http.Request) {
	if s.newsDir != "" {
		p := filepath.Join(s.newsDir, filepath.Base(r.URL.Path))
		if fi, err := os.Stat(p); err == nil && !fi.IsDir() {
			log.Printf("[rivals] news %s -> %s (%d bytes)", r.URL.Path, p, fi.Size())
			http.ServeFile(w, r, p)
			return
		}
	}
	log.Printf("[rivals] news %s -> 204 (nothing staged)", r.URL.Path)
	w.WriteHeader(http.StatusNoContent)
}

// nexStatus reports the NEX middleware stub: 501 until the title's NEX key is
// paired and RIVALS_NEX_SECURE_URL points at the switch-nex secure listener.
func (s *server) nexStatus(w http.ResponseWriter, r *http.Request) {
	if ok, missing := s.nex.status(); !ok {
		log.Printf("[rivals] nex/status -> 501 (%s)", missing)
		writeStub(w, missing)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]string{"status": "ok", "secure": s.nex.secureURL})
}

func (s *server) unknown(w http.ResponseWriter, r *http.Request) {
	n, _ := io.Copy(io.Discard, io.LimitReader(r.Body, 1<<20))
	log.Printf("[rivals] !! unhandled %s %s body=%d", r.Method, r.URL.Path, n)
	http.NotFound(w, r)
}

func writeStub(w http.ResponseWriter, missing string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNotImplemented)
	_ = json.NewEncoder(w).Encode(map[string]string{
		"error":   "not_implemented",
		"missing": missing,
	})
}

