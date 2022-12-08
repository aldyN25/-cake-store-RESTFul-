package cake

import (
	"time"
)

type CakeResponse struct {
	ID          int        `json:"id"`
	Title       string     `json:"title"`
	Description string     `json:"description"`
	Image       string     `json:"image"`
	Rating      float32    `json:"rating"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   *time.Time `json:"updated_at"`
}

type CakesResponse []CakeResponse
