package main

import (
	"cmp"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"
)

type Config struct {
	Port     string
	LogLevel string
}

func main() {
	cfg, err := ConfigLoad()
	if err != nil {
		slog.Error("config failed", "error", err)
		os.Exit(1)
	}

	if err := LoggerSetup(cfg.LogLevel); err != nil {
		slog.Error("logger failed", "error", err)
		os.Exit(1)
	}

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
	port := fmt.Sprintf(":"+"%s", portInput)
	cfg := &Config{
		Port:     port,
		LogLevel: cmp.Or(os.Getenv("LOG_LEVEL"), "INFO"),
	}
	return cfg, nil
}

func LoggerSetup(level string) error {

	lv := new(slog.LevelVar)
	if err := lv.UnmarshalText([]byte(level)); err != nil {
		return fmt.Errorf("parse LOG_LEVEL: %w", err)
	}
	h := slog.NewJSONHandler(
		os.Stdout, &slog.HandlerOptions{
			Level: lv,
		},
	)
	slog.SetDefault(slog.New(h))
	return nil
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

func ServerRun(port string) error {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /healthz", healthz)

	srv := &http.Server{
		Addr:              port,
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
