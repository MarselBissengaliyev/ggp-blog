package main

import (
	"log"

	"github.com/MarselBissengaliyev/ggp-blog/config"
	"github.com/MarselBissengaliyev/ggp-blog/routes"
	"github.com/MarselBissengaliyev/ggp-blog/storage"
	"github.com/gin-gonic/gin"
)

func main() {
	config, err := config.LoadConfig("./")
	if err != nil {
		log.Fatal(err)
	}

	db, err := storage.NewConnection(&config)
	if err != nil {
		log.Fatalf("Could not load the database: %s", err.Error())
	}

	r := gin.Default()
	routes.SetupRoutes(r, db)
	r.Run()
}
