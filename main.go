package main

import (
	"github.com/gin-gonic/gin"
	"github.com/utils-price-tool/handlers"
	"github.com/utils-price-tool/services"
	"github.com/utils-price-tool/storage"
	"github.com/utils-price-tool/storage/storecrc"
	"github.com/utils-price-tool/tasks"

	"log"
	"time"
)

func main() {
	r := gin.Default()

	store := storage.NewInMemoryStore()
	storeCRC := storecrc.NewInMemoryCRCStore()
	serviceGetPrices := services.NewService()

	// container
	toTask := tasks.DuiCont{
		TimeOut:   time.Second*5,
		Service:   serviceGetPrices,
		Store:     store,
		StoreCRC:  storeCRC,
	}

	tasks.NewGetGroupTask(&toTask)
	handlers.NewController(store, storeCRC).Mount(r)

	if err := r.Run(":5000"); err != nil {
		log.Fatal(err)
	}

}
