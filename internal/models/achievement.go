package models

import "time"

type Achievement struct {
    ID          int       `json:"id"`
    Name        string    `json:"name"`
    Description string    `json:"description"`
    UnlockedAt  time.Time `json:"unlocked_at"`
}