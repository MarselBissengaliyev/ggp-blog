package main

import (
	"log"

	"github.com/MarselBissengaliyev/ggp-blog/config"
	"github.com/MarselBissengaliyev/ggp-blog/migrations"
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

	if err := migrations.Migrate(db); err != nil {
		log.Fatalf("Could not migrate database: %s", err.Error())
	}

	r := gin.Default()

	r.Use(gin.Logger())

	routes.SetupRoutes(r, db, &config)

	r.Run()
}
