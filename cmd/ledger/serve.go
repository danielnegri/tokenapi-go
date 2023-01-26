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
	"os"
	"runtime"
	"time"

	"github.com/danielnegri/tokenapi-go/errors"
	"github.com/danielnegri/tokenapi-go/log"
	"github.com/danielnegri/tokenapi-go/server"
	"github.com/danielnegri/tokenapi-go/source"
	"github.com/danielnegri/tokenapi-go/storage/postgres"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func commandServe() *cobra.Command {
	var (
		concurrency   int
		databaseURL   string
		logFormat     string
		logLevel      string
		port          int
		sourceRetry   int
		sourceTimeout time.Duration
		sourceURL     string
	)

	cmd := cobra.Command{
		Use:     "serve",
		Short:   "Start Ledger HTTP server",
		Example: fmt.Sprintf("%s serve", shortDescription),
		Run: func(cmd *cobra.Command, args []string) {
			serverCfg := newServerConfig()
			serverCfg.Source = newSourceConfig()
			serverCfg.Storage = newStorageConfig()

			log.SetLogger(newLogger())
			svr := server.New(serverCfg)
			errors.Separator = ":: "

			if err := svr.Run(); err != nil {
				_, _ = fmt.Fprintln(os.Stderr, err)
				os.Exit(1)
			}
		},
	}

	cmd.Flags().IntVar(&concurrency, "concurrency", runtime.NumCPU(), "number of concurrent workers")
	_ = viper.BindPFlag("concurrency", cmd.Flags().Lookup("concurrency"))

	cmd.Flags().StringVar(&databaseURL, "database-url", postgres.DefaultURL, "database connection string")
	_ = viper.BindPFlag("database_url", cmd.Flags().Lookup("database-url"))

	cmd.Flags().StringVar(&logFormat, "log-format", log.DefaultFormat, "logger format")
	_ = viper.BindPFlag("log_format", cmd.Flags().Lookup("log-format"))

	cmd.Flags().StringVar(&logLevel, "log-level", log.DefaultLevel, "logger level")
	_ = viper.BindPFlag("log_level", cmd.Flags().Lookup("log-level"))

	cmd.Flags().IntVar(&port, "port", server.DefaultPort, "HTTP server port")
	_ = viper.BindPFlag("port", cmd.Flags().Lookup("port"))

	cmd.Flags().IntVar(&sourceRetry, "source-retry", source.DefaultRetry, "token source max retries")
	_ = viper.BindPFlag("source_retry", cmd.Flags().Lookup("source-retry"))

	cmd.Flags().DurationVar(&sourceTimeout, "source-timeout", source.DefaultTimeout, "token source timeout")
	_ = viper.BindPFlag("source_timeout", cmd.Flags().Lookup("source-timeout"))

	cmd.Flags().StringVar(&sourceURL, "source-url", source.DefaultURL, "token source address")
	_ = viper.BindPFlag("source_url", cmd.Flags().Lookup("source-url"))

	return &cmd
}
