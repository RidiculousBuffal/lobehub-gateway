package main

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	agentgateway "github.com/lobehub/lobehub/apps/agent-gateway-go"
	devicegateway "github.com/lobehub/lobehub/apps/device-gateway-go"
)

type service struct {
	name            string
	server          *http.Server
	shutdownTimeout time.Duration
}

func main() {
	os.Exit(run())
}

func run() int {
	// The unified binary ignores PORT: each service gets its own variable so
	// the two listeners cannot collide in a single process.
	agentCfg := agentgateway.ConfigFromEnv()
	agentCfg.Port = envOrDefault("AGENT_PORT", "8787")
	if err := agentCfg.Validate(); err != nil {
		slog.Error("invalid agent configuration", "error", err)
		return 1
	}
	deviceCfg := devicegateway.ConfigFromEnv()
	deviceCfg.Port = envOrDefault("DEVICE_PORT", "8788")
	if err := deviceCfg.Validate(); err != nil {
		slog.Error("invalid device configuration", "error", err)
		return 1
	}

	services := []service{
		{
			name:            "agent gateway",
			server:          newHTTPServer(agentCfg.Port, agentgateway.NewServer(agentCfg).Routes(), agentCfg.ReadTimeout, agentCfg.WriteTimeout),
			shutdownTimeout: agentCfg.ShutdownTimeout,
		},
		{
			name:            "device gateway",
			server:          newHTTPServer(deviceCfg.Port, devicegateway.NewServer(deviceCfg).Routes(), deviceCfg.ReadTimeout, deviceCfg.WriteTimeout),
			shutdownTimeout: deviceCfg.ShutdownTimeout,
		},
	}

	errCh := make(chan error, len(services))
	for _, svc := range services {
		go func() {
			slog.Info(svc.name+" listening", "addr", svc.server.Addr)
			errCh <- svc.server.ListenAndServe()
		}()
	}

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	exitCode := 0
	select {
	case <-stop:
	case err := <-errCh:
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			slog.Error("server failed", "error", err)
			exitCode = 1
		}
	}

	for _, svc := range services {
		ctx, cancel := context.WithTimeout(context.Background(), svc.shutdownTimeout)
		if err := svc.server.Shutdown(ctx); err != nil {
			slog.Error(svc.name+" shutdown failed", "error", err)
			exitCode = 1
		}
		cancel()
	}
	return exitCode
}

func newHTTPServer(port string, handler http.Handler, readTimeout, writeTimeout time.Duration) *http.Server {
	return &http.Server{
		Addr:         ":" + port,
		Handler:      handler,
		ReadTimeout:  readTimeout,
		WriteTimeout: writeTimeout,
	}
}

func envOrDefault(key string, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}
