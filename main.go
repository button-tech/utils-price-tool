package main

import (
	"github.com/gin-gonic/gin"
	"github.com/utils-price-tool/handlers"
	"github.com/utils-price-tool/services"
	"github.com/utils-price-tool/storage"
	"github.com/utils-price-tool/tasks"

	"log"
	"time"
)

func main() {
	r := gin.Default()

	store := storage.NewInMemoryStore()
	serviceGetPrices := services.NewService()
	tasks.NewGetGroupTask(time.Second*5, serviceGetPrices, store)
	handlers.NewController(store).Mount(r)

	if err := r.Run(":5000"); err != nil {
		log.Fatal(err)
	}

}
