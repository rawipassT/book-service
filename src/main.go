package main

import (
	"log"

	"github.com/rawipassT/book-service/config"
	"github.com/rawipassT/book-service/internal/http"
	"github.com/rawipassT/book-service/routes"
)

func main() {

	// Init Config
	config.InitConfig()
	
	// Init Database
	config.ConnectDatabase()

	// Setup Route
	r := routes.SetupRoutes(&http.BookHandler{})
	if err := r.Run(":3000"); err != nil {
		log.Fatalf("Could not run server: %v\n", err)
	}
}
