package main

import (
    "log"
    "net/http"
    "os"

    "github.com/Pranav-1100/quiz-app-backend/internal/database"
    "github.com/joho/godotenv"
    "github.com/Pranav-1100/quiz-app-backend/internal/router"
)

func main() {
    godotenv.Load()
    db, err := database.InitDB("quiz.db")
    if err != nil {
        log.Fatalf("Error initializing database: %v", err)
    }
    defer db.Close()

    // Check if data was inserted correctly
    var count int
    err = db.QueryRow("SELECT COUNT(*) FROM questions").Scan(&count)
    if err != nil {
        log.Printf("Error checking question count: %v", err)
    } else {
        log.Printf("Number of questions in database: %d", count)
    }

    r := router.SetupRouter(db)

    port := os.Getenv("PORT")
    if port == "" {
        port = "8080" // Default port if not specified in environment
    }

    log.Printf("Server starting on :%s", port)
    log.Fatal(http.ListenAndServe(":"+port, r))
}