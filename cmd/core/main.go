package main

import (
	"github.com/button-tech/utils-price-tool/pkg/storage/cache"
	"log"
	"os"
	"syscall"

	"github.com/button-tech/logger"
	"github.com/button-tech/utils-price-tool/core/prices"
	core "github.com/button-tech/utils-price-tool/core/server"
	"github.com/button-tech/utils-price-tool/services/update"
	"os/signal"
)

func main() {
	store := cache.NewCache()

	p := prices.New(store)

	go update.Start(p)

	if err := logger.InitLogger(os.Getenv("DSN")); err != nil {
		log.Fatal(err)
	}

	c := core.New(store, p)
	signalforExit := make(chan os.Signal, 1)
	defer close(signalforExit)

	signal.Notify(signalforExit,
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

	stop := <-signalforExit
	logger.Info("Received", stop)
	logger.Info("Waiting for all jobs to stop")
}
