package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestTelemetryAccepted(t *testing.T) {
	s := &server{nex: &nexStub{}}
	for _, target := range []string{"/api/games/start", "/api/games/finish"} {
		for _, method := range []string{http.MethodGet, http.MethodPost} {
			req := httptest.NewRequest(method, target, strings.NewReader(`{"game":"x"}`))
			w := httptest.NewRecorder()
			s.telemetry(w, req)
			if w.Code != 200 {
				t.Errorf("%s %s: got %d, want 200", method, target, w.Code)
			}
		}
	}
}

func TestNewsEmptyWhenNothingStaged(t *testing.T) {
	s := &server{nex: &nexStub{}}
	for _, target := range []string{"/files/ingame_assets/news", "/files/ingame_assets/news_data.txt"} {
		w := httptest.NewRecorder()
		s.news(w, httptest.NewRequest(http.MethodGet, target, nil))
		if w.Code != 204 {
			t.Errorf("%s: got %d, want 204", target, w.Code)
		}
	}
}

func TestNewsServedFromDir(t *testing.T) {
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "news_data.txt"), []byte("hello"), 0o644); err != nil {
		t.Fatal(err)
	}
	s := &server{newsDir: dir, nex: &nexStub{}}
	w := httptest.NewRecorder()
	s.news(w, httptest.NewRequest(http.MethodGet, "/files/ingame_assets/news_data.txt", nil))
	if w.Code != 200 || w.Body.String() != "hello" {
		t.Fatalf("news_data.txt: got %d %q, want 200 %q", w.Code, w.Body.String(), "hello")
	}
}

func TestNexStatusBlockedOnKey(t *testing.T) {
	s := &server{nex: &nexStub{secureURL: "https://127.0.0.1:23161"}}
	w := httptest.NewRecorder()
	s.nexStatus(w, httptest.NewRequest(http.MethodGet, "/v1/nex/status", nil))
	if w.Code != 501 {
		t.Fatalf("nex/status: got %d, want 501 while TODO-KEY is open", w.Code)
	}
	var body map[string]string
	if err := json.Unmarshal(w.Body.Bytes(), &body); err != nil {
		t.Fatal(err)
	}
	if body["missing"] != "nex-key" {
		t.Fatalf("nex/status missing=%q, want %q", body["missing"], "nex-key")
	}
}

func TestConfigDefaults(t *testing.T) {
	t.Setenv("RIVALS_ADDR", "")
	t.Setenv("RIVALS_HEALTH_ADDR", "")
	cfg := configFromEnv()
	if cfg.Addr != ":23570" || cfg.HealthAddr != ":23571" {
		t.Fatalf("defaults: addr=%q health=%q", cfg.Addr, cfg.HealthAddr)
	}
	t.Setenv("RIVALS_ADDR", ":9999")
	if cfg := configFromEnv(); cfg.Addr != ":9999" {
		t.Fatalf("RIVALS_ADDR override: got %q", cfg.Addr)
	}
}
