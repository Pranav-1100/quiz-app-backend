package router

import (
    "database/sql"
    "github.com/gin-gonic/gin"
    "github.com/Pranav-1100/quiz-app-backend/internal/handlers"
    "github.com/Pranav-1100/quiz-app-backend/internal/services"
)

func SetupRouter(db *sql.DB) *gin.Engine {
    r := gin.Default()

    r.Use(CORSMiddleware())

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
    r.GET("/next-level", handlers.GetNextLevel(quizService, userService))  // Updated this line

    return r
}

func CORSMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
        c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
        c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
        c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

        if c.Request.Method == "OPTIONS" {
            c.AbortWithStatus(204)
            return
        }

        c.Next()
    }
}