package models

type Question struct {
    ID            int      `json:"id"`
    LevelID       int      `json:"level_id"`
    Question      string   `json:"question"`
    Options       []string `json:"options"`
    CorrectAnswer string   `json:"-"` // Don't send this to the client
    Difficulty    int      `json:"difficulty"`
    ImageURL      string   `json:"image_url"`
    Explanation   string   `json:"-"` // Don't send this to the client initially
}