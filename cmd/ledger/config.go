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
