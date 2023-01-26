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

package httputil

import (
	"net/http"
	"regexp"
	"runtime"
	"strings"
	"time"

	"github.com/danielnegri/tokenapi-go/errors"
	"github.com/danielnegri/tokenapi-go/version"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func RootHandler(msg string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{
			"message":         msg,
			"arch":            runtime.GOARCH,
			"build_time":      version.BuildTime,
			"commit":          version.CommitHash,
			"os":              runtime.GOOS,
			"runtime_version": runtime.Version(),
			"version":         version.Version,
		})
	}

}

// NotFoundHandler is a helper function that calls Server.Abort.
func NotFoundHandler(c *gin.Context) {
	Abort(c, http.StatusNotFound, http.StatusText(http.StatusNotFound))
}

// LoggerHandler returns a gin.HandlerFunc (middleware) that logs requests using logrus.
//
// Requests with errors are logged using logrus.Error().
// Requests without errors are logged using logrus.Info().
//
// It receives:
//  1. A time package format string (e.g. time.RFC3339).
//  2. A boolean stating whether to use UTC time zone or local.
func LoggerHandler(logger logrus.FieldLogger, timeFormat string, utc bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		// some evil middlewares modify this values
		path := c.Request.URL.Path
		c.Next()

		end := time.Now()
		latency := end.Sub(start)
		if utc {
			end = end.UTC()
		}

		entry := logger.WithFields(logrus.Fields{
			"status":       c.Writer.Status(),
			"method":       c.Request.Method,
			"uri":          c.Request.RequestURI,
			"path":         path,
			"content_type": c.ContentType(),
			"remote-addr":  c.ClientIP(),
			"user-agent":   c.Request.UserAgent(),
			"x-request-id": c.GetHeader("X-Request-Id"),
			"latency":      latency,
			"time":         end.Format(timeFormat),
		})

		if len(c.Errors) > 0 {
			// Append error field if this is an erroneous request.
			entry.Error(c.Errors.String())
		} else {
			entry.Info()
		}
	}
}

var newLine = regexp.MustCompile(`\r?\n?\t`)

func AbortWithError(ctx *gin.Context, err error) {
	code := http.StatusInternalServerError
	msg := newLine.ReplaceAllString(err.Error(), " ")
	e, ok := err.(*errors.Error)
	if ok {
		if index := strings.Index(msg, ":"); len(msg) > index+1 {
			msg = strings.TrimSpace(msg[index+1:])
		}

		switch e.Kind {
		case errors.Duplicate:
			code = http.StatusBadRequest
		case errors.Invalid:
			code = http.StatusBadRequest
		case errors.NotFound:
			code = http.StatusNotFound
		case errors.Permission:
			code = http.StatusForbidden
		}
	}

	ctx.AbortWithStatusJSON(code, &ErrorResponse{
		Code:    code,
		Message: msg,
	})
}
