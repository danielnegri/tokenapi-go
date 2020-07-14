package server

import (
	"context"
	"fmt"
	"net/http"
	"runtime"
	"strconv"
	"time"

	"github.com/danielnegri/adheretech/errors"
	"github.com/danielnegri/adheretech/ledger"
	"github.com/danielnegri/adheretech/log"
	"github.com/danielnegri/adheretech/net/httputil"
	"github.com/danielnegri/adheretech/sync"
	"github.com/danielnegri/adheretech/version"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

const Prefix = "/api/v1"

func (s *service) newHandler() http.Handler {
	router := gin.New()
	router.Use(gin.Recovery())
	router.Use(cors.Default())
	router.NoRoute(httputil.NotFoundHandler)

	router.GET("/", s.handleRoot())
	router.GET("/health/heartbeat", s.handleHeartbeat())

	api := router.Group(Prefix)
	api.POST("/tokens", s.handleInsert())
	return router
}

func (s *service) handleRoot() gin.HandlerFunc {
	root := gin.H{
		"service":         ledger.Description,
		"arch":            runtime.GOARCH,
		"build_time":      version.BuildTime,
		"commit":          version.CommitHash,
		"os":              runtime.GOOS,
		"runtime_version": runtime.Version(),
		"version":         version.Version,
	}

	return func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, root)
	}
}

func (s *service) handleHeartbeat() gin.HandlerFunc {
	heartbeat := gin.H{
		"startup_time": time.Now(),
		"current_time": time.Now(),
		"message":      http.StatusText(http.StatusOK),
		"service":      ledger.Description,
		"status":       http.StatusOK,
		"version":      version.Version,
	}

	return func(ctx *gin.Context) {
		now := time.Now()
		switch ctx.NegotiateFormat(gin.MIMEHTML, gin.MIMEJSON) {
		case gin.MIMEJSON:
			heartbeat["current_time"] = now
			ctx.JSON(http.StatusOK, heartbeat)
		default:
			ctx.String(http.StatusOK, "%s @ %s", ledger.Description, now.Format(time.RFC3339))
		}
	}
}

func (s *service) handleInsert() gin.HandlerFunc {
	const op errors.Op = "server/service.handleInsert"

	return func(ctx *gin.Context) {
		rawsize := ctx.DefaultQuery("size", "0")
		size, err := strconv.Atoi(rawsize)
		if err != nil {
			log.Error(errors.E(op, err))
			httputil.AbortWithError(ctx, errors.E(errors.Invalid, "size must be an integer"))
			return
		}

		tokens, err := s.source.Generate(ctx, size)
		if err != nil {
			log.Error(errors.E(op, err))
			httputil.AbortWithError(ctx, err)
			return
		}

		w := ctx.Writer
		h := w.Header()
		h.Set("Transfer-Encoding", "chunked")
		h.Set("Content-Type", gin.MIMEPlain)
		w.WriteHeader(http.StatusOK)

		count := 0
		lines := make(chan string)
		finished := make(chan struct{}, 1)
		go s.insert(ctx, tokens, lines, finished)
		for {
			select {
			case line := <-lines:
				w.Write([]byte(line + "\n"))
				w.Flush()
				count++
			case <-finished:
				log.Debugf("Processed %d tokens", count)
				return
			}
		}
	}
}

func (s *service) insert(ctx context.Context, tokens []ledger.Token, lines chan string, finished chan struct{}) {
	wg := sync.NewWaitGroup(s.cfg.Concurrency)
	for _, token := range tokens {
		wg.Add()
		go func(c context.Context, t ledger.Token) {
			defer wg.Done()

			line := fmt.Sprintf("OK : %v", t)
			if err := s.storage.Insert(c, t); err != nil {
				line = fmt.Sprintf("ERR: %v", t)
			}

			log.Debug(line)
			lines <- line
		}(ctx, token)
	}

	wg.Wait()
	finished <- struct{}{}
}
