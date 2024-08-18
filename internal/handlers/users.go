package handlers

import (
    "net/http"
	"fmt"
    "strconv"
    "github.com/gin-gonic/gin"
    "github.com/Pranav-1100/quiz-app-backend/internal/services"
)

func GetUser(userService *services.UserService) gin.HandlerFunc {
    return func(c *gin.Context) {
        userID, err := strconv.Atoi(c.Param("id"))
        if err != nil {
            c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
            return
        }

        user, err := userService.GetUser(userID)
        if err != nil {
            if err.Error() == "user not found" {
                // Create a new user with a default username
                newUser, err := userService.CreateUser(fmt.Sprintf("User%d", userID))
                if err != nil {
                    c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
                    return
                }
                c.JSON(http.StatusOK, newUser)
                return
            }
            c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
            return
        }

        c.JSON(http.StatusOK, user)
    }
}

func GetUserCoins(userService *services.UserService) gin.HandlerFunc {
    return func(c *gin.Context) {
        userID, err := strconv.Atoi(c.Param("id"))
        if err != nil {
            c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
            return
        }

        user, err := userService.GetUser(userID)
        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
            return
        }

        c.JSON(http.StatusOK, gin.H{"coins": user.Coins})
    }
}

func CreateUser(userService *services.UserService) gin.HandlerFunc {
    return func(c *gin.Context) {
        var request struct {
            Username string `json:"username" binding:"required"`
        }

        if err := c.ShouldBindJSON(&request); err != nil {
            c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
            return
        }

        user, err := userService.CreateUser(request.Username)
        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
            return
        }

        c.JSON(http.StatusCreated, user)
    }
}