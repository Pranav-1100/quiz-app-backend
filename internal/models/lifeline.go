package models

type Lifeline struct {
    ID   int    `json:"id"`
    Type string `json:"type"`
    Cost int    `json:"cost"`
}