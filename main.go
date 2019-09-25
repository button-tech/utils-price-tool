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

//todo: add packages response, request, error

func main() {
	r := gin.Default()
	r.Use(cors.Default())

	store := storage.NewCache()
	serviceGetPrices := services.New()

	// Container for tasks
	toTask := tasks.DuiCont{
		TimeOut: time.Minute * 7,
		Service: serviceGetPrices,
		Store:   store,
	}

	go tasks.NewGetGroup(&toTask)
	controllers.New(store).Mount(r)

	if err := r.Run(":5000"); err != nil {
		log.Fatal(err)
	}

}
