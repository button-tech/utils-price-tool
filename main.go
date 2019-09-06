package main

import (
	"github.com/gin-gonic/gin"
	"github.com/utils-tool_prices/handlers"
	"github.com/utils-tool_prices/services"
	"github.com/utils-tool_prices/storage"
	"github.com/utils-tool_prices/tasks"
	"log"
	"time"
)

func main() {
	r := gin.Default()

	store := storage.NewInMemoryStore()
	serviceGetPrices := services.NewService()
	tasks.NewGetGroupTask(time.Second * 5, serviceGetPrices, store)
	handlers.NewController(store).Mount(r)

	if err := r.Run(":5000"); err != nil {
		log.Fatal(err)
	}

}
