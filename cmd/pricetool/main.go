package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/button-tech/logger"
	"github.com/button-tech/utils-price-tool/api"
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

	s := api.NewServer(store)
	signalEx := make(chan os.Signal, 1)
	defer close(signalEx)

	signal.Notify(signalEx,
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT)

	go func() {
		if err := s.Core.ListenAndServe(":5000"); err != nil {
			logger.Fatal(err)
		}
	}()
	defer func() {
		if err := s.Core.Shutdown(); err != nil {
			logger.Fatal(err)
		}
	}()

	stop := <-signalEx
	logger.Info("Received", stop)
	logger.Info("Waiting for all jobs to stop")

}
