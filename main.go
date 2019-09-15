package main

import (
	"github.com/button-tech/utils-price-tool/controllers"
	"github.com/button-tech/utils-price-tool/services"
	"github.com/button-tech/utils-price-tool/storage/storecrc"
	"github.com/button-tech/utils-price-tool/storage/storetoplist"
	"github.com/button-tech/utils-price-tool/storage/storetrustwallet"
	"github.com/button-tech/utils-price-tool/tasks"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"log"
	"time"
)

func main() {
	r := gin.Default()
	r.Use(cors.Default())

	store := storetrustwallet.NewInMemoryCMCStore()
	storeCRC := storecrc.NewInMemoryCRCStore()
	storeList := storetoplist.NewInMemoryListStore()
	serviceGetPrices := services.NewService()

	// Container for tasks
	toTask := tasks.DuiCont{
		TimeOut:   time.Minute * 7,
		Service:   serviceGetPrices,
		Store:     store,
		StoreCRC:  storeCRC,
		StoreList: storeList,
	}

	go tasks.NewGetGroupTask(&toTask)
	//fmt.Println(runtime.NumGoroutine())
	controllers.NewController(store, storeCRC, storeList).Mount(r)

	if err := r.Run(":5000"); err != nil {
		log.Fatal(err)
	}

}
