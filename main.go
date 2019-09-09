package main

import (
	"github.com/button-tech/utils-price-tool/handlers"
	"github.com/button-tech/utils-price-tool/services"
	"github.com/button-tech/utils-price-tool/storage"
	"github.com/button-tech/utils-price-tool/storage/storecrc"
	"github.com/button-tech/utils-price-tool/tasks"
	"github.com/gin-gonic/gin"
<<<<<<< HEAD
=======
	"github.com/button-tech/utils-price-tool/handlers"
	"github.com/button-tech/utils-price-tool/services"
	"github.com/button-tech/utils-price-tool/storage"
	"github.com/button-tech/utils-price-tool/storage/storecrc"
	"github.com/button-tech/utils-price-tool/tasks"
>>>>>>> 65de444f817d6aa0bf31207504f19ae7961ace98

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
		TimeOut:  time.Second * 5,
		Service:  serviceGetPrices,
		Store:    store,
		StoreCRC: storeCRC,
	}

	tasks.NewGetGroupTask(&toTask)
	handlers.NewController(store, storeCRC).Mount(r)

	if err := r.Run(":5000"); err != nil {
		log.Fatal(err)
	}

}
