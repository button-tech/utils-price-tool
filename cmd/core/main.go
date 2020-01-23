package main

import (
	"github.com/button-tech/utils-price-tool/pkg/storage/cache"
	"log"
	"os"
	"syscall"

	"github.com/button-tech/logger"
	core "github.com/button-tech/utils-price-tool/core/server"
	"github.com/button-tech/utils-price-tool/services/update"
	"os/signal"
)

func main() {
	storage := cache.NewCache()

	go update.Start(storage)

	if err := logger.InitLogger(os.Getenv("DSN")); err != nil {
		log.Fatal(err)
	}

	c := core.New(storage)
	signalForExit := make(chan os.Signal, 1)
	defer close(signalForExit)

	signal.Notify(signalForExit,
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT)

	go func() {
		if err := c.S.ListenAndServe(":5000"); err != nil {
			logger.Fatal(err)
		}
	}()
	defer func() {
		if err := c.S.Shutdown(); err != nil {
			logger.Fatal(err)
		}
	}()

	stop := <-signalForExit
	logger.Info("Received", stop)
	logger.Info("Waiting for all jobs to stop")
}
