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

package main

import (
	"fmt"
	"net/url"
	"os"

	"github.com/go-pg/pg/v10"

	"github.com/danielnegri/adheretech/log"
	"github.com/danielnegri/adheretech/net"
	"github.com/danielnegri/adheretech/server"
	"github.com/danielnegri/adheretech/source"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func newLogger() logrus.FieldLogger {
	return log.New(viper.GetString("log_level"), viper.GetString("log_format"))
}

func newServerConfig() *server.Config {
	cfg := &server.Config{}
	cfg.Concurrency = viper.GetInt("concurrency")
	cfg.Debug = viper.GetString("log_level") == "debug"
	cfg.HTTPServer = &net.ServerConfig{}
	cfg.HTTPServer.HTTPPort = viper.GetInt("port")
	return cfg
}

func newSourceConfig() *source.Config {
	cfg := &source.Config{}
	cfg.Retry = viper.GetInt("source_retry")
	cfg.Timeout = viper.GetDuration("source_timeout")

	rawurl := viper.GetString("source_url")
	srcURL, err := url.Parse(rawurl)
	if err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)
		os.Exit(2)
	}

	cfg.URL = srcURL

	return cfg
}

func newStorageConfig() *pg.Options {
	dburl := viper.GetString("database_url")
	cfg, err := pg.ParseURL(dburl)
	if err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)
		os.Exit(2)
	}

	return cfg
}
