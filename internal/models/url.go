package models

import "time"

type URL struct {
    ID          int       `json:"id"`   
    ShortCode   string    `json:"short_code"` 
    OriginalURL string    `json:"original_url"`
    CreatedAt   time.Time `json:"created_at"`
    Clicks      int       `json:"clicks"`
}