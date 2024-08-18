package router

import (
    "database/sql"
    "time"
    "github.com/gin-contrib/cors"
    "github.com/gin-gonic/gin"
    "github.com/Pranav-1100/quiz-app-backend/internal/handlers"
    "github.com/Pranav-1100/quiz-app-backend/internal/services"
)

func SetupRouter(db *sql.DB) *gin.Engine {
    r := gin.Default()

    // Apply CORS middleware with a permissive configuration
    r.Use(cors.New(cors.Config{
        AllowOrigins:     []string{"*"},
        AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
        AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization", "X-Requested-With"},
        ExposeHeaders:    []string{"Content-Length"},
        AllowCredentials: true,
        MaxAge:           12 * time.Hour,
    }))

    quizService := services.NewQuizService(db)
    userService := services.NewUserService(db)
    leaderboardService := services.NewLeaderboardService(db)

    r.GET("/questions/:level", handlers.GetQuestions(quizService))
    r.GET("/quizzes", handlers.GetAllQuizzes(quizService, userService))
    r.POST("/answer", handlers.SubmitAnswer(quizService, userService))
    r.POST("/user", handlers.CreateUser(userService))
    r.GET("/user/:id", handlers.GetUser(userService))
    r.POST("/lifeline", handlers.UseLifeline(quizService, userService))
    r.GET("/user/:id/coins", handlers.GetUserCoins(userService))
    r.GET("/achievements/:id", handlers.GetAchievements(userService))
    r.GET("/leaderboard", handlers.GetLeaderboard(leaderboardService))
    r.GET("/next-level", handlers.GetNextLevel(quizService, userService))

    return r
}