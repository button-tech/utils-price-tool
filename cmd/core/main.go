package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/button-tech/logger"
	core "github.com/button-tech/utils-price-tool/core/server"
	"github.com/button-tech/utils-price-tool/pkg/storage"
	"github.com/button-tech/utils-price-tool/services"
	"github.com/button-tech/utils-price-tool/tasks"
)

func main() {
	store := storage.NewCache()
	getPrices := services.New(store)
	go tasks.FetchGroup(getPrices, store)

	if err := logger.InitLogger(os.Getenv("DSN")); err != nil {
		log.Fatal(err)
	}

	c := core.New(store, getPrices)
	signalEx := make(chan os.Signal, 1)
	defer close(signalEx)

	signal.Notify(signalEx,
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

	stop := <-signalEx
	logger.Info("Received", stop)
	logger.Info("Waiting for all jobs to stop")
}
