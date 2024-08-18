package handlers

import (
    "net/http"
    "strconv"
    "github.com/gin-gonic/gin"
    "github.com/Pranav-1100/quiz-app-backend/internal/services"
)

func GetAllQuizzes(quizService *services.QuizService, userService *services.UserService) gin.HandlerFunc {
    return func(c *gin.Context) {
        userID, err := getUserIDFromContext(c)
        if err != nil {
            c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
            return
        }

        quizzes, err := quizService.GetAllQuizzes()
        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch quizzes"})
            return
        }

        completedQuizzes, err := userService.GetCompletedQuizzes(userID)
        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch user progress"})
            return
        }

        response := make([]gin.H, len(quizzes))
        for i, quiz := range quizzes {
            completed := false
            for _, completedQuizID := range completedQuizzes {
                if completedQuizID == quiz.ID {
                    completed = true
                    break
                }
            }

            response[i] = gin.H{
                "id":          quiz.ID,
                "name":        quiz.Name,
                "description": quiz.Description,
                "difficulty":  quiz.Difficulty,
                "threshold":   quiz.UnlockThreshold,
                "completed":   completed,
            }
        }

        c.JSON(http.StatusOK, gin.H{
            "quizzes": response,
            "user_id": userID,
        })
    }
}

func getUserIDFromContext(c *gin.Context) (int, error) {
    userID, err := strconv.Atoi(c.Query("user_id"))
    if err != nil {
        return 0, err
    }
    return userID, nil
}