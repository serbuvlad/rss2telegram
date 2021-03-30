package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/adrg/xdg"
	"gopkg.in/yaml.v3"
)

var (
	configPath string
	dbPath     string
	timeout    time.Duration
	appname    string
)

func init() {
	flag.StringVar(&configPath, "config", "", "YAML configuation file (default $XDG_CONFIG_HOME/appname/config.yaml)")
	flag.StringVar(&dbPath, "db", "", "sqlite3 database to store persistant data in (default $XDG_DATA_HOME/appname/data.db)")
	flag.StringVar(&appname, "appname", "rss2telegram", "appname to use for XDG directories")
	flag.DurationVar(&timeout, "t", time.Minute, "time to wait between polling feeds")
}

type Config struct {
	Feeds    []Feed
	Telegram struct {
		Token  string
		ChatID int64
	}
}

func main() {
	flag.Parse()

	if configPath == "" {
		var err error

		configPath, err = xdg.ConfigFile(filepath.Join(appname, "config.yaml"))
		if err != nil {
			log.Fatal(fmt.Errorf("Coult not find XDG config path: %w", err))
		}
	}

	if dbPath == "" {
		var err error

		dbPath, err = xdg.DataFile(filepath.Join(appname, "data.db"))
		if err != nil {
			log.Fatal(fmt.Errorf("Coult not find XDG data path: %w", err))
		}
	}

	var c Config

	f, err := os.Open(configPath)

	if err != nil {
		log.Fatal(fmt.Errorf("Can not open config file, make sure you created one: %w", err))
	}

	err = yaml.NewDecoder(f).Decode(&c)

	if err != nil {
		log.Fatal(fmt.Errorf("Error parsing config file: %w", err))
	}

	f.Close()

	initDB(dbPath)

	initTelegram(c.Telegram.Token, c.Telegram.ChatID)

	ch := make(chan string)

	go serveTelegram(ch)
	go serveFeeds(c.Feeds, timeout, ch)

	stop := make(chan struct{})
	<-stop
}
