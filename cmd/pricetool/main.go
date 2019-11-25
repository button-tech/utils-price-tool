package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/button-tech/utils-price-tool/api"
	"github.com/button-tech/utils-price-tool/services"
	"github.com/button-tech/utils-price-tool/storage"
	"github.com/button-tech/utils-price-tool/tasks"
	"github.com/valyala/fasthttp"
)

func main() {
	store := storage.NewCache()
	go tasks.NewGetGroup(services.New(), store)

	s := api.NewServer(store)
	server := fasthttp.Server{
		Handler:      s.R.HandleRequest,
		ReadTimeout:  time.Second * 30,
		WriteTimeout: time.Second * 30,
	}

	signalEx := make(chan os.Signal, 1)
	defer close(signalEx)

	signal.Notify(signalEx,
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT)

	go func() {
		if err := server.ListenAndServe(":5000"); err != nil {
			log.Fatal(err)
		}
	}()
	defer func() {
		if err := server.Shutdown(); err != nil {
			log.Fatal(err)
		}
	}()

	stop := <-signalEx
	log.Println("Received", stop)
	log.Println("Waiting for all jobs to stop")

}
