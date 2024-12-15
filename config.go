package main

import (
	"errors"
	"fmt"
	"io"
	"os"
)

type config struct {
	Prefix string
	Host   string
	Port   string
	PathDB string
	LogTo  io.Writer
}

func NewConfig() *config {
	return &config{
		Prefix: "METRICS_SERVICE",
	}
}

func (c *config) Load() error {
	c.Host = os.Getenv(fmt.Sprintf("%s_HOST", c.Prefix))
	c.Port = os.Getenv(fmt.Sprintf("%s_PORT", c.Prefix))
	c.PathDB = os.Getenv(fmt.Sprintf("%s_PATH_DB", c.Prefix))

	logfile := os.Getenv(fmt.Sprintf("%s_LOG_FILE", c.Prefix))

	if c.Host == "" {
		c.Host = "127.0.0.1"
	}
	if c.Port == "" {
		c.Port = "8081"
	}
	if c.PathDB == "" {
		return errors.New("Config: Unknown DB path")
	}
	if logfile != "" {
		file, err := os.OpenFile(logfile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
		if err != nil {
			return err
		}
		c.LogTo = file
	} else {
		c.LogTo = os.Stdout
	}

	return nil
}

func (c *config) Addr() string {
	return fmt.Sprintf("%s:%s", c.Host, c.Port)
}
