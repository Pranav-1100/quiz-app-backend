package models

import "time"

type LeaderboardEntry struct {
    UserID    int       `json:"user_id"`
    Username  string    `json:"username"`
    Score     int       `json:"score"`
    Timestamp time.Time `json:"timestamp"`
}