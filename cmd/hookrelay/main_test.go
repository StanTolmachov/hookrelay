package main

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestConfigLoad(t *testing.T) {
	tests := []struct {
		name     string
		port     string
		logLevel string
		want     Config
		wantErr  bool
	}{
		{
			name: "all default",
			want: Config{Port: "8080", LogLevel: slog.LevelInfo},
		},
		{
			name:     "all changed",
			port:     "3333",
			logLevel: "warn",
			want:     Config{Port: "3333", LogLevel: slog.LevelWarn},
		},
		{
			name: "port ok",
			port: "3333",
			want: Config{Port: "3333", LogLevel: slog.LevelInfo},
		},
		{
			name:    "port abc",
			port:    "abc",
			wantErr: true,
		},
		{
			name: "port 1",
			port: "1",
			want: Config{Port: "1", LogLevel: slog.LevelInfo},
		},
		{
			name: "port 65535",
			port: "65535",
			want: Config{Port: "65535", LogLevel: slog.LevelInfo},
		},
		{
			name:    "port out of range zero",
			port:    "0",
			wantErr: true,
		},
		{
			name:    "port out of range more 65535",
			port:    "65536",
			wantErr: true,
		},
		{
			name:     "LOG_LEVEL abc",
			logLevel: "abc",
			wantErr:  true,
		},
		{
			name:     "LOG_LEVEL lowcase",
			logLevel: "debug",
			want:     Config{Port: "8080", LogLevel: slog.LevelDebug},
		},
	}

	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				t.Setenv("PORT", tt.port)
				t.Setenv("LOG_LEVEL", tt.logLevel)

				got, err := ConfigLoad()

				if (err != nil) != tt.wantErr {
					t.Fatalf("want = %v, got = %v", tt.wantErr, err)
				}
				if tt.wantErr {
					return
				}

				if *got != tt.want {
					t.Fatalf("want = %+v, got = %+v", tt.want, *got)
				}
			},
		)
	}
}

func TestNewRouter(t *testing.T) {
	mux := NewRouter()

	tests := []struct {
		name       string
		target     string
		method     string
		wantStatus int
		wantHeader map[string]string
		wantBody   *string
	}{
		{
			name:       "ok",
			target:     "/healthz",
			method:     http.MethodGet,
			wantStatus: http.StatusOK,
			wantHeader: map[string]string{
				"Content-Type": "application/json",
			},
			wantBody: new("ok"),
		},
		{
			name:       "405",
			target:     "/healthz",
			method:     http.MethodPost,
			wantStatus: http.StatusMethodNotAllowed,
		},
		{
			name:       "404",
			target:     "/healthzzz",
			method:     http.MethodGet,
			wantStatus: http.StatusNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				req := httptest.NewRequest(tt.method, tt.target, nil)
				rec := httptest.NewRecorder()

				mux.ServeHTTP(rec, req)

				if rec.Code != tt.wantStatus {
					t.Fatalf("status = %d, want = %d", rec.Code, tt.wantStatus)
				}

				for k, want := range tt.wantHeader {
					if got := rec.Result().Header.Get(k); got != want {
						t.Fatalf("header %s = %q, want = %q", k, got, want)
					}
				}
				if tt.wantBody == nil {
					return
				}

				type healthResp struct {
					Status string `json:"status"`
				}

				var gotBody healthResp

				if err := json.NewDecoder(rec.Body).Decode(&gotBody); err != nil {
					t.Fatalf("decode body: %v", err)
				}

				if *tt.wantBody != gotBody.Status {
					t.Fatalf("body got = %q, want = %q", gotBody.Status, *tt.wantBody)
				}

			},
		)
	}
}
