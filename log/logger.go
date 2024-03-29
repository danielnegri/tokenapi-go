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

package log

import (
	"fmt"
	"os"
	"strings"

	"github.com/sirupsen/logrus"
)

var (
	logLevels  = []string{"panic", "fatal", "error", "warn", "info", "debug"}
	logFormats = []string{"json", "text"}
)

type utcFormatter struct {
	f logrus.Formatter
}

// Format log entries to UTC location.
func (f *utcFormatter) Format(e *logrus.Entry) ([]byte, error) {
	e.Time = e.Time.UTC()
	return f.f.Format(e)
}

// New creates a new Logger. Configuration should be set by changing level (eg.: panic, fatal, error, warn, info, debug)
// format (eg.: text, json).
func New(level string, format string) logrus.FieldLogger {
	var logLevel logrus.Level
	switch strings.ToLower(level) {
	case "panic":
		logLevel = logrus.PanicLevel
	case "fatal":
		logLevel = logrus.FatalLevel
	case "error":
		logLevel = logrus.ErrorLevel
	case "warn":
		logLevel = logrus.WarnLevel
	case "info":
		logLevel = logrus.InfoLevel
	case "debug":
		logLevel = logrus.DebugLevel
	default:
		panic(fmt.Sprintf("log level is not one of the supported values (%s): %s", strings.Join(logLevels, ", "), level))
	}

	var formatter utcFormatter
	switch strings.ToLower(format) {
	case "text":
		formatter.f = &logrus.TextFormatter{DisableColors: true}
	case "json":
		formatter.f = &logrus.JSONFormatter{}
	default:
		panic(fmt.Sprintf("log format is not one of the supported values (%s), falling back to json: %s", strings.Join(logFormats, ", "), format))
	}

	return &logrus.Logger{
		Out:       os.Stderr,
		Formatter: &formatter,
		Level:     logLevel,
	}
}
