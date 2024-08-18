package services

import (
    "database/sql"
    "github.com/Pranav-1100/quiz-app-backend/internal/models"
)

type LeaderboardService struct {
    DB *sql.DB
}

func NewLeaderboardService(db *sql.DB) *LeaderboardService {
    return &LeaderboardService{DB: db}
}

func (s *LeaderboardService) GetTopScores(limit int) ([]models.LeaderboardEntry, error) {
    rows, err := s.DB.Query(`
        SELECT u.id, u.username, l.score, l.timestamp
        FROM leaderboard l
        JOIN users u ON l.user_id = u.id
        ORDER BY l.score DESC
        LIMIT ?
    `, limit)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var entries []models.LeaderboardEntry
    for rows.Next() {
        var e models.LeaderboardEntry
        if err := rows.Scan(&e.UserID, &e.Username, &e.Score, &e.Timestamp); err != nil {
            return nil, err
        }
        entries = append(entries, e)
    }

    return entries, nil
}

func (s *LeaderboardService) UpdateScore(userID int, score int) error {
    _, err := s.DB.Exec("INSERT INTO leaderboard (user_id, score) VALUES (?, ?)", userID, score)
    return err
}