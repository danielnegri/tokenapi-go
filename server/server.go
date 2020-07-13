package server

import (
	"context"
	"fmt"

	"github.com/danielnegri/adheretech/storage/postgres"

	"github.com/danielnegri/adheretech/errors"
	"github.com/danielnegri/adheretech/ledger"
	"github.com/danielnegri/adheretech/log"
	"github.com/danielnegri/adheretech/net"
	"github.com/danielnegri/adheretech/source"
	"github.com/danielnegri/adheretech/storage"
	"github.com/danielnegri/adheretech/version"
	"github.com/gin-gonic/gin"
	"github.com/go-pg/pg/v10"
)

const (
	DefaultPort = 8080
)

type Server interface {
	Run() error
	Shutdown()
}

type service struct {
	cfg     *Config
	server  net.Server
	source  source.Source
	storage storage.Storage
	debug   bool

	quit chan struct{}
}

var _ Server = (*service)(nil)

type Config struct {
	Concurrency int
	Debug       bool
	HTTPServer  *net.ServerConfig
	Source      *source.Config
	Storage     *pg.Options
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

	cfg := s.cfg
	if cfg == nil {
		return errors.E(errors.Internal, "invalid server configuration")
	}

	if cfg.Source != nil && s.source != nil {
		if err := s.source.Check(ctx); err != nil {
			log.Errorf("error while connecting to Token source: %v", err)
		} else {
			log.Infof("Connected to Token source at %v", cfg.Source.URL)
		}
	} else {
		return errors.E(errors.Internal, "invalid Token source configuration")
	}

	if cfg.Storage != nil {
		storage, err := postgres.Connect(cfg.Storage)
		if err != nil {
			log.Errorf("error while connecting to Postgres: %v", err)
			return err
		}

		s.storage = storage
		if err := s.storage.Check(ctx); err != nil {
			log.Errorf("error while checking connection with storage: %v", err)
		} else {
			log.Infof("Connected to Storage at %v", cfg.Source.URL)
		}
	} else {
		return errors.E(errors.Internal, "invalid storage configuration")
	}

	// Start Server
	if err := s.server.Run(); err != nil {
		return fmt.Errorf("failed to start server: %v", err)
	}

	return nil
}

func (s *service) Shutdown() {
	log.Infof("%s: Stopping Ledger service", ledger.Description)
	s.quit <- struct{}{}
}
