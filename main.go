package main

import (
	"github.com/button-tech/utils-price-tool/handlers"
	"github.com/button-tech/utils-price-tool/services"
	"github.com/button-tech/utils-price-tool/storage"
	"github.com/button-tech/utils-price-tool/storage/storecrc"
	"github.com/button-tech/utils-price-tool/storage/storetoplist"
	"github.com/button-tech/utils-price-tool/tasks"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"log"
	"time"
)

func main() {
	r := gin.Default()
	r.Use(cors.Default())

	store := storage.NewInMemoryStore()
	storeCRC := storecrc.NewInMemoryCRCStore()
	storeList := storetoplist.NewInMemoryListStore()
	serviceGetPrices := services.NewService()

	// container
	toTask := tasks.DuiCont{
		TimeOut:   time.Minute*7,
		Service:   serviceGetPrices,
		Store:     store,
		StoreCRC:  storeCRC,
		StoreList: storeList,
	}

	go tasks.NewGetGroupTask(&toTask)
	handlers.NewController(store, storeCRC, storeList).Mount(r)

	if err := r.Run(":5001"); err != nil {
		log.Fatal(err)
	}

}
