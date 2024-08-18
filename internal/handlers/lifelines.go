package handlers

import (
    "net/http"
    "github.com/gin-gonic/gin"
    "github.com/Pranav-1100/quiz-app-backend/internal/services"
)

func UseLifeline(quizService *services.QuizService, userService *services.UserService) gin.HandlerFunc {
    return func(c *gin.Context) {
        var request struct {
            UserID       int    `json:"user_id"`
            QuestionID   int    `json:"question_id"`
            LifelineType string `json:"lifeline_type"`
        }

        if err := c.BindJSON(&request); err != nil {
            c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
            return
        }

        cost := getLifelineCost(request.LifelineType)
        err := userService.UseLifeline(request.UserID, cost)
        if err != nil {
            c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
            return
        }

        var lifelineResult interface{}
        var lifelineError error

        switch request.LifelineType {
        case "50-50":
            lifelineResult, lifelineError = quizService.GetFiftyFiftyOptions(request.QuestionID)
        case "hint":
            lifelineResult, lifelineError = quizService.GetHint(request.QuestionID)
        case "expert":
            lifelineResult, lifelineError = quizService.GetExpertAdvice(request.QuestionID)
        default:
            c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid lifeline type"})
            return
        }

        if lifelineError != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": lifelineError.Error()})
            return
        }

        // Get updated user information
        user, err := userService.GetUser(request.UserID)
        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
            return
        }

        c.JSON(http.StatusOK, gin.H{
            "message": "Lifeline used successfully",
            "lifeline_result": lifelineResult,
            "remaining_coins": user.Coins,
        })
    }
}

func getLifelineCost(lifelineType string) int {
    switch lifelineType {
    case "50-50":
        return 20
    case "hint":
        return 10
    case "expert":
        return 30
    default:
        return 0
    }
}
