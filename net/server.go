// Copyright 2020 The Ledger Authors
//
// Licensed under the AGPL, Version 3.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     https://www.gnu.org/licenses/agpl-3.0.en.html
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package net

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/danielnegri/adheretech/log"
)

// ServerConfig holds info required to configure a Server.Server.
type ServerConfig struct {
	// MaxHeaderBytes can be used to override the default of 1<<20.
	MaxHeaderBytes int

	// ReadTimeout can be used to override the default http Server timeout of 20s.
	// The string should be formatted like a time.Duration string.
	ReadTimeout time.Duration

	// WriteTimeout can be used to override the default http Server timeout of 20s.
	// The string should be formatted like a time.Duration string.
	WriteTimeout time.Duration

	// IdleTimeout can be used to override the default http Server timeout of 120s.
	// The string should be formatted like a time.Duration string.
	IdleTimeout time.Duration

	// ShutdownTimeout can be used to override the default http Server Shutdown timeout
	// of 5m.
	ShutdownTimeout time.Duration

	// HTTPPort is the port the Server implementation will serve HTTP over.
	// The default is 8080
	HTTPPort int
}

type Server interface {
	Run() error
}

// server encapsulates all logic for registering and running a Server.
type server struct {
	cfg *ServerConfig

	httpServer *http.Server
	Shutdown   func()

	// Exit chan for graceful Shutdown
	Exit chan chan error
}

func NewServer(cfg *ServerConfig, handler http.Handler) *server {
	if cfg == nil {
		cfg = &ServerConfig{}
	}

	if cfg.MaxHeaderBytes == 0 {
		cfg.MaxHeaderBytes = 1 << 20
	}

	if cfg.ReadTimeout == 0 {
		cfg.ReadTimeout = 20 * time.Second
	}

	if cfg.WriteTimeout == 0 {
		cfg.WriteTimeout = 20 * time.Second
	}

	if cfg.IdleTimeout == 0 {
		cfg.IdleTimeout = 120 * time.Second
	}

	if cfg.ShutdownTimeout == 0 {
		cfg.ShutdownTimeout = 5 * time.Minute
	}

	if cfg.HTTPPort == 0 {
		cfg.HTTPPort = 8080
	}

	httpServer := &http.Server{
		Handler:        handler,
		Addr:           fmt.Sprintf(":%d", cfg.HTTPPort),
		MaxHeaderBytes: cfg.MaxHeaderBytes,
		ReadTimeout:    cfg.ReadTimeout,
		WriteTimeout:   cfg.WriteTimeout,
		IdleTimeout:    cfg.IdleTimeout,
	}

	return &server{
		cfg:        cfg,
		httpServer: httpServer,
		Exit:       make(chan chan error),
	}
}

func (s *server) start() error {
	go func() {
		log.Infof("Listening and serving HTTP on %s", s.httpServer.Addr)
		err := s.httpServer.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			log.Errorf("HTTP Server error - initiating shutting down: %v", err)
			s.stop()
			return
		}
	}()

	go func() {
		exit := <-s.Exit

		// Stop listener with timeout
		ctx, cancel := context.WithTimeout(context.Background(), s.cfg.ShutdownTimeout)
		defer cancel()

		// Stop service
		if s.Shutdown != nil {
			s.Shutdown()
		}

		// Stop HTTP Server
		if s.httpServer != nil {
			log.Infof("Stopping HTTP Server on %s", s.httpServer.Addr)
			exit <- s.httpServer.Shutdown(ctx)
			return
		}

		exit <- nil
	}()

	return nil
}

func (s *server) stop() error {
	ch := make(chan error)
	s.Exit <- ch
	return <-ch
}

// Run will create a new Server and register the given
// Service and start up the Server(s).
// This will block until the Server shuts down.
func (s *server) Run() error {
	if err := s.start(); err != nil {
		return err
	}

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGTERM, syscall.SIGINT, syscall.SIGKILL)
	log.Info("Received signal ", <-ch)
	return s.stop()
}

func NewHTTPServer(cfg ServerConfig, handler http.Handler) *http.Server {
	return &http.Server{
		Handler:        handler,
		Addr:           fmt.Sprintf(":%d", cfg.HTTPPort),
		MaxHeaderBytes: cfg.MaxHeaderBytes,
		ReadTimeout:    cfg.ReadTimeout,
		WriteTimeout:   cfg.WriteTimeout,
		IdleTimeout:    cfg.IdleTimeout,
	}
}
