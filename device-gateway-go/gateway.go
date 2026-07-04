package devicegateway

import (
	"net/http"
	"time"

	"github.com/lobehub/lobehub/apps/device-gateway-go/internal/gateway"
)

type Config = gateway.Config

type Server = gateway.Server

func ConfigFromEnv() Config { return gateway.ConfigFromEnv() }

func NewServer(cfg Config) *Server { return gateway.NewServer(cfg) }

func Routes(s *Server) http.Handler { return s.Routes() }

func ReadTimeout(cfg Config) time.Duration { return cfg.ReadTimeout }

func WriteTimeout(cfg Config) time.Duration { return cfg.WriteTimeout }

func ShutdownTimeout(cfg Config) time.Duration { return cfg.ShutdownTimeout }
