package models

type Level struct {
    ID              int    `json:"id"`
    Name            string `json:"name"`
    Difficulty      int    `json:"difficulty"`
    UnlockThreshold int    `json:"unlock_threshold"`
}