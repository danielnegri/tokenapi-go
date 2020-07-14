package main

import (
	"fmt"
	"os"
	"runtime"
	"time"

	"github.com/danielnegri/adheretech/errors"
	"github.com/danielnegri/adheretech/log"
	"github.com/danielnegri/adheretech/server"
	"github.com/danielnegri/adheretech/source"
	"github.com/danielnegri/adheretech/storage/postgres"
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
