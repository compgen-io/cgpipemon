package config

import (
    "flag"
)

type Config struct {
    Listen string
    DBConnect string
    Command string
    Asset func(string) ([]byte, error)
}

func Init() *Config {
    cfg := &Config{}

    flag.StringVar(&cfg.Listen, "listen", "localhost:3000", "HTTP listen spec")
    flag.StringVar(&cfg.DBConnect, "db-connect", "host=/var/run/postgresql dbname=gowebapp sslmode=disable", "DB Connect String")
    flag.Parse()

    if (len(flag.Args()) > 0) {
        cfg.Command = flag.Args()[0]
    } else {
        cfg.Command = "";
    }

    cfg.Asset = Asset

    return cfg
}

