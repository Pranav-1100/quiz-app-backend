package models

type Quiz struct {
    ID              int    `json:"id"`
    Name            string `json:"name"`
    Description     string `json:"description"`
    Difficulty      int    `json:"difficulty"`
    UnlockThreshold int    `json:"unlock_threshold"`
}