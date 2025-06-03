package main

import (
    "bdu-dormitory-backend/config"
    "bdu-dormitory-backend/routes"
    "github.com/gin-gonic/gin"
)

func main() {
    // Initialize database connection
    db, err := config.InitDB()
    if err != nil {
        panic("Failed to connect to database: " + err.Error())
    }
    defer db.Close()

    // Set up Gin router
    r := gin.Default()

    // Initialize routes
    routes.SetupRoutes(r, db)

    // Start server
    r.Run(":8080") // Runs on http://localhost:8080
}