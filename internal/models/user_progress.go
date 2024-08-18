package models

type UserProgress struct {
    UserID           int   `json:"user_id"`
    CompletedQuizzes []int `json:"completed_quizzes"`
}