package main

import (
	"log"

	"github.com/abik1221/bdu-dormitory-backend/config"
	"github.com/abik1221/bdu-dormitory-backend/routes"
	"github.com/gin-gonic/gin"
)

func main() {

	db, err := config.InitDB()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	r := gin.Default()

	routes.SetupRoutes(r, db)

	if err := r.Run(":8080"); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
