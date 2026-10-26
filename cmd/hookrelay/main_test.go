package main

import (
	"log/slog"
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
