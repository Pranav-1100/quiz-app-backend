package main

import (
    "log"
    "net/http"

    "github.com/Pranav-1100/quiz-app-backend/internal/database"
    "github.com/Pranav-1100/quiz-app-backend/internal/router"
)

func main() {
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

    log.Println("Server starting on :8080")
    log.Fatal(http.ListenAndServe(":8080", r))
}