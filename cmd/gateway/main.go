package main

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	agentgateway "github.com/lobehub/lobehub/apps/agent-gateway-go"
	devicegateway "github.com/lobehub/lobehub/apps/device-gateway-go"
)

func main() {
	agentCfg := agentgateway.ConfigFromEnv()
	if err := agentCfg.Validate(); err != nil {
		slog.Error("invalid agent configuration", "error", err)
		os.Exit(1)
	}
	deviceCfg := devicegateway.ConfigFromEnv()
	if err := deviceCfg.Validate(); err != nil {
		slog.Error("invalid device configuration", "error", err)
		os.Exit(1)
	}

	agentServer := agentgateway.NewServer(agentCfg)
	deviceServer := devicegateway.NewServer(deviceCfg)

	agentHTTP := &http.Server{
		Addr:         ":" + agentCfg.Port,
		Handler:      agentgateway.Routes(agentServer),
		ReadTimeout:  agentCfg.ReadTimeout,
		WriteTimeout: agentCfg.WriteTimeout,
	}
	deviceHTTP := &http.Server{
		Addr:         ":" + deviceCfg.Port,
		Handler:      devicegateway.Routes(deviceServer),
		ReadTimeout:  deviceCfg.ReadTimeout,
		WriteTimeout: deviceCfg.WriteTimeout,
	}

	errCh := make(chan error, 2)
	go func() {
		slog.Info("agent gateway listening", "addr", agentHTTP.Addr)
		errCh <- agentHTTP.ListenAndServe()
	}()
	go func() {
		slog.Info("device gateway listening", "addr", deviceHTTP.Addr)
		errCh <- deviceHTTP.ListenAndServe()
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	select {
	case <-stop:
		ctx, cancel := context.WithTimeout(context.Background(), agentCfg.ShutdownTimeout)
		defer cancel()
		_ = agentHTTP.Shutdown(ctx)
		_ = deviceHTTP.Shutdown(ctx)
	case err := <-errCh:
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			slog.Error("server failed", "error", err)
			ctx, cancel := context.WithTimeout(context.Background(), agentCfg.ShutdownTimeout)
			defer cancel()
			_ = agentHTTP.Shutdown(ctx)
			_ = deviceHTTP.Shutdown(ctx)
			os.Exit(1)
		}
	}
}
