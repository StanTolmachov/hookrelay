package main

import (
	"cmp"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"
)

type Config struct {
	Port     string
	LogLevel slog.Level
}

func main() {
	cfg, err := ConfigLoad()
	if err != nil {
		slog.Error("config failed", "error", err)
		os.Exit(1)
	}

	LoggerSetup(cfg.LogLevel)

	if err := ServerRun(cfg.Port); err != nil {
		slog.Error("server failed", "error", err)
		os.Exit(1)
	}

}

func ConfigLoad() (*Config, error) {
	portInput := cmp.Or(os.Getenv("PORT"), "8080")
	portInt, err := strconv.Atoi(portInput)
	if err != nil {
		return nil, fmt.Errorf("parse PORT: %w", err)
	}
	if portInt < 1 || portInt > 65535 {
		return nil, fmt.Errorf("parse PORT: %d out of range 1..65535", portInt)
	}

	level := cmp.Or(os.Getenv("LOG_LEVEL"), "INFO")

	var lvl slog.Level
	if err := lvl.UnmarshalText([]byte(level)); err != nil {
		return nil, fmt.Errorf("parse LOG_LEVEL: %w", err)
	}

	cfg := &Config{
		Port:     strconv.Itoa(portInt),
		LogLevel: lvl,
	}
	return cfg, nil
}

func LoggerSetup(level slog.Level) {
	h := slog.NewJSONHandler(
		os.Stdout, &slog.HandlerOptions{
			Level: level,
		},
	)
	slog.SetDefault(slog.New(h))
}
func healthz(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	type healthResp struct {
		Status string `json:"status"`
	}
	resp := healthResp{Status: "ok"}
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		slog.Error("encode health response", "error", err)
	}
}

func NewRouter() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /healthz", healthz)
	return mux
}
func ServerRun(port string) error {

	addr := net.JoinHostPort("", port)

	mux := NewRouter()

	srv := &http.Server{
		Addr:              addr,
		Handler:           mux,
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      10 * time.Second,
		IdleTimeout:       60 * time.Second,
	}

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	errChan := make(chan error, 1)
	go func() {
		slog.Info("starting server on port", "port", port)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			errChan <- err
		}
	}()

	select {
	case err := <-errChan:
		return err
	case sig := <-stop:
		slog.Info("stopping server with signal", "signal", sig.String())
	}

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	return srv.Shutdown(shutdownCtx)
}
