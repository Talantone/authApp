package main

import (
	"github.com/Talantone/authApp"
	"github.com/Talantone/authApp/pkg/handler"
	"log"
)

func main() {
	handlers := new(handler.Handler)
	srv := new(authApp.Server)
	if err := srv.Run("8080", handlers.InitRoutes()); err != nil {
		log.Fatalf("Error occured while running http server: %s", err.Error())
	}

}
