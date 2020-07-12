package server

import (
	"context"
	"fmt"

	"github.com/danielnegri/adheretech/log"

	"github.com/danielnegri/adheretech/errors"
	"github.com/danielnegri/adheretech/ledger"
	"github.com/danielnegri/adheretech/net"
	"github.com/danielnegri/adheretech/source"
	"github.com/danielnegri/adheretech/version"
	"github.com/gin-gonic/gin"
)

const (
	DefaultPort = 8080
)

type Server interface {
	//ledger.Ledger

	Run() error
	Shutdown()
}

type service struct {
	cfg    *Config
	server net.Server
	source source.Source
	debug  bool
}

var _ Server = (*service)(nil)

type Config struct {
	Concurrency int
	Debug       bool
	HTTPServer  *net.ServerConfig
	Source      *source.Config
}

func New(cfg *Config) *service {
	if !cfg.Debug {
		gin.SetMode("release")
	}

	svc := &service{
		cfg:    cfg,
		source: source.New(cfg.Source),
	}

	server := net.NewServer(cfg.HTTPServer, svc.newHandler())
	server.Shutdown = svc.Shutdown
	svc.server = server

	return svc
}

func (s *service) Run() error {
	log.Infof("%s: Starting Ledger service (%s)", ledger.Description, version.Version)
	ctx := context.Background()

	if s.cfg.Source != nil && s.source != nil {
		if err := s.source.Check(ctx); err != nil {
			log.Errorf("error while connecting to Token source: %v", err)
		} else {
			log.Infof("Connected to Token source at %v", s.cfg.Source.URL)
		}
	} else {
		return errors.E(errors.Internal, "invalid Token source configuration")
	}

	// Start Server
	if err := s.server.Run(); err != nil {
		return fmt.Errorf("failed to start server: %v", err)
	}

	return nil
}

func (s *service) Shutdown() {
	log.Infof("%s: Stopping Ledger service", ledger.Description)

	// TODO: Stop storage
}
