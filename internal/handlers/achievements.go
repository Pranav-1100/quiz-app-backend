package handlers

import (
    "net/http"
	"strconv"
    "github.com/gin-gonic/gin"
    "github.com/Pranav-1100/quiz-app-backend/internal/services"
)

func GetAchievements(userService *services.UserService) gin.HandlerFunc {
    return func(c *gin.Context) {
        userID, err := strconv.Atoi(c.Param("id"))
        if err != nil {
            c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
            return
        }

        achievements, err := userService.GetUserAchievements(userID)
        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
            return
        }

        if len(achievements) == 0 {
            c.JSON(http.StatusOK, gin.H{"message": "No achievements yet"})
            return
        }

        c.JSON(http.StatusOK, achievements)
    }
}