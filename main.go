package main

import (
	"github.com/button-tech/utils-price-tool/controllers"
	"github.com/button-tech/utils-price-tool/services"
	"github.com/button-tech/utils-price-tool/storage"
	"github.com/button-tech/utils-price-tool/tasks"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"log"
	"time"
)

//1.todo: debug fetch-task

//todo: rewrite at fast-http

//todo: make one storage

//todo: complete hexed functional

//todo: add packages response, request, error

//todo: from handlers make one func for mapping

func main() {
	r := gin.Default()
	r.Use(cors.Default())

	store := storage.NewCache()
	serviceGetPrices := services.NewService()

	// Container for tasks
	toTask := tasks.DuiCont{
		TimeOut:   time.Second * 3,
		Service:   serviceGetPrices,
		Store:     store,
	}

	go tasks.NewGetGroupTask(&toTask)
	controllers.NewController(store).Mount(r)

	if err := r.Run(":5000"); err != nil {
		log.Fatal(err)
	}

}
