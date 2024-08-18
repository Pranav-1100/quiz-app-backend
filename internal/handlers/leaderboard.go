package handlers

import (
    "net/http"
    "github.com/gin-gonic/gin"
    "github.com/Pranav-1100/quiz-app-backend/internal/services"
)

func GetLeaderboard(leaderboardService *services.LeaderboardService) gin.HandlerFunc {
    return func(c *gin.Context) {
        leaderboard, err := leaderboardService.GetTopScores(10) // Get top 10 scores
        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
            return
        }

        c.JSON(http.StatusOK, leaderboard)
    }
}