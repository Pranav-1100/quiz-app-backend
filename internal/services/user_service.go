package services

import (
    "database/sql"
	"fmt"
    "github.com/Pranav-1100/quiz-app-backend/internal/models"
)

type UserService struct {
    DB *sql.DB
}

func NewUserService(db *sql.DB) *UserService {
    return &UserService{DB: db}
}

func (s *UserService) GetUser(userID int) (*models.User, error) {
    var user models.User
    err := s.DB.QueryRow("SELECT id, username, coins, current_level, streak FROM users WHERE id = ?", userID).Scan(&user.ID, &user.Username, &user.Coins, &user.CurrentLevel, &user.Streak)
    if err == sql.ErrNoRows {
        return nil, fmt.Errorf("user not found")
    }
    if err != nil {
        return nil, err
    }
    return &user, nil
}

func (s *UserService) UpdateUserProgress(userID int, coinsEarned int, levelCompleted bool) error {
    tx, err := s.DB.Begin()
    if err != nil {
        return err
    }
    defer tx.Rollback()

    _, err = tx.Exec("UPDATE users SET coins = coins + ?, streak = streak + 1 WHERE id = ?", coinsEarned, userID)
    if err != nil {
        return err
    }

    if levelCompleted {
        _, err = tx.Exec("UPDATE users SET current_level = current_level + 1 WHERE id = ?", userID)
        if err != nil {
            return err
        }
    }

    return tx.Commit()
}

func (s *UserService) GetUserAchievements(userID int) ([]models.Achievement, error) {
    rows, err := s.DB.Query(`
        SELECT a.id, a.name, a.description, ua.unlocked_at
        FROM achievements a
        JOIN user_achievements ua ON a.id = ua.achievement_id
        WHERE ua.user_id = ?
    `, userID)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var achievements []models.Achievement
    for rows.Next() {
        var a models.Achievement
        if err := rows.Scan(&a.ID, &a.Name, &a.Description, &a.UnlockedAt); err != nil {
            return nil, err
        }
        achievements = append(achievements, a)
    }

    return achievements, nil
}

func (s *UserService) AddCoins(userID int, coins int) error {
    _, err := s.DB.Exec("UPDATE users SET coins = coins + ? WHERE id = ?", coins, userID)
    return err
}

func (s *UserService) UseLifeline(userID int, lifelineCost int) error {
    result, err := s.DB.Exec("UPDATE users SET coins = coins - ? WHERE id = ? AND coins >= ?", lifelineCost, userID, lifelineCost)
    if err != nil {
        return err
    }
    rowsAffected, err := result.RowsAffected()
    if err != nil {
        return err
    }
    if rowsAffected == 0 {
        return fmt.Errorf("not enough coins to use lifeline")
    }
    return nil
}

func (s *UserService) CreateUser(username string) (*models.User, error) {
    result, err := s.DB.Exec("INSERT INTO users (username, coins, current_level, streak) VALUES (?, ?, ?, ?)", username, 0, 1, 0)
    if err != nil {
        return nil, err
    }

    id, err := result.LastInsertId()
    if err != nil {
        return nil, err
    }

    return &models.User{
        ID:           int(id),
        Username:     username,
        Coins:        0,
        CurrentLevel: 1,
        Streak:       0,
    }, nil
}

func (s *UserService) GetUserProgress(userID int) (*models.UserProgress, error) {
    var progress models.UserProgress
    progress.UserID = userID

    rows, err := s.DB.Query("SELECT level_id FROM user_progress WHERE user_id = ?", userID)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    for rows.Next() {
        var levelID int
        if err := rows.Scan(&levelID); err != nil {
            return nil, err
        }
        progress.CompletedQuizzes = append(progress.CompletedQuizzes, levelID)
    }

    return &progress, nil
}

func (s *UserService) GetCompletedQuizzes(userID int) ([]int, error) {
    rows, err := s.DB.Query("SELECT level_id FROM user_progress WHERE user_id = ?", userID)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var completedQuizzes []int
    for rows.Next() {
        var levelID int
        if err := rows.Scan(&levelID); err != nil {
            return nil, err
        }
        completedQuizzes = append(completedQuizzes, levelID)
    }

    return completedQuizzes, nil
}