package main

import (
	"context"
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	config "github.com/Aureys96/hw12_13_14_15_calendar/internal/config" //nolint
	storage "github.com/Aureys96/hw12_13_14_15_calendar/internal/storage"
	sqlstorage "github.com/Aureys96/hw12_13_14_15_calendar/internal/storage/sql"

	"github.com/Aureys96/hw12_13_14_15_calendar/internal/app"
	"github.com/Aureys96/hw12_13_14_15_calendar/internal/logger"
	internalhttp "github.com/Aureys96/hw12_13_14_15_calendar/internal/server/http"
	memorystorage "github.com/Aureys96/hw12_13_14_15_calendar/internal/storage/memory"
)

var configFile string

func init() {
	flag.StringVar(&configFile, "config", "/configs/config.toml", "Path to configuration file")
}

func main() {
	flag.Parse()

	if flag.Arg(0) == "version" {
		printVersion()
		return
	}

	cfg, err := config.NewConfig(configFile)
	if err != nil {
		log.Fatalln("Error while reading cfg file", err)
	}
	logg := logger.New(cfg)

	var store storage.IStorage
	if cfg.DBConfig.Inmemory {
		store = memorystorage.New()
	} else {
		store = sqlstorage.New(cfg.DBConfig)
	}

	calendar := app.New(logg, store)

	server := internalhttp.NewServer(logg, calendar, cfg.Server)

	ctx, cancel := signal.NotifyContext(context.Background(),
		syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	defer cancel()

	go func() {
		<-ctx.Done()

		ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
		defer cancel()

		if err := server.Stop(ctx); err != nil {
			logg.Error("failed to stop http server: " + err.Error())
		}
	}()

	logg.Info("calendar is running...")

	if err := server.Start(); err != nil {
		os.Exit(1) //nolint:gocritic
	}
}
