package models

type User struct {
    ID           int    `json:"id"`
    Username     string `json:"username"`
    Coins        int    `json:"coins"`
    CurrentLevel int    `json:"current_level"`
    Streak       int    `json:"streak"`
}