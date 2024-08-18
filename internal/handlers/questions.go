package handlers

import (
    "net/http"
    "strconv"
    "github.com/gin-gonic/gin"
    "github.com/Pranav-1100/quiz-app-backend/internal/services"
)

func GetQuestions(quizService *services.QuizService) gin.HandlerFunc {
    return func(c *gin.Context) {
        levelID, err := strconv.Atoi(c.Param("level"))
        if err != nil {
            c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid level ID"})
            return
        }

        questions, err := quizService.GetQuestionsByLevel(levelID)
        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
            return
        }

        c.JSON(http.StatusOK, questions)
    }
}

func SubmitAnswer(quizService *services.QuizService, userService *services.UserService) gin.HandlerFunc {
    return func(c *gin.Context) {
        var request struct {
            UserID     int    `json:"user_id"`
            QuestionID int    `json:"question_id"`
            Answer     string `json:"answer"`
        }

        if err := c.BindJSON(&request); err != nil {
            c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
            return
        }

        correct, explanation, err := quizService.CheckAnswer(request.QuestionID, request.Answer)
        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
            return
        }

        var responseEmoji, responseImage string
        var coinsEarned int

        if correct {
            responseEmoji = "🎉"
            responseImage = "https://static.vecteezy.com/system/resources/previews/011/422/127/non_2x/correct-answer-text-button-correct-answer-speech-bubble-correct-answer-banner-label-template-illustration-vector.jpg"
            coinsEarned = 10
        } else {
            responseEmoji = "😢"
            responseImage = "https://static.vecteezy.com/system/resources/previews/011/422/127/non_2x/correct-answer-text-button-correct-answer-speech-bubble-correct-answer-banner-label-template-illustration-vector.jpg"
            coinsEarned = 0
        }

        // Update user progress and add coins
        if err := userService.UpdateUserProgress(request.UserID, coinsEarned, false); err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
            return
        }

        // Get updated user information
        user, err := userService.GetUser(request.UserID)
        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
            return
        }

        c.JSON(http.StatusOK, gin.H{
            "correct":     correct,
            "explanation": explanation,
            "emoji":       responseEmoji,
            "image":       responseImage,
            "coins_earned": coinsEarned,
            "total_coins": user.Coins,
            "current_level": user.CurrentLevel,
            "streak":      user.Streak,
        })
    }
}