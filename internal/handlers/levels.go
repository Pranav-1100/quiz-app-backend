package handlers

import (
    "net/http"
    "github.com/gin-gonic/gin"
    "github.com/Pranav-1100/quiz-app-backend/internal/services"
)

func GetNextLevel(quizService *services.QuizService, userService *services.UserService) gin.HandlerFunc {
    return func(c *gin.Context) {
        userID, _ := c.Get("userID") // Assuming you have middleware to set this
        user, err := userService.GetUser(userID.(int))
        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
            return
        }

        nextLevel, err := quizService.GetNextLevel(user.CurrentLevel)
        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
            return
        }

        c.JSON(http.StatusOK, nextLevel)
    }
}