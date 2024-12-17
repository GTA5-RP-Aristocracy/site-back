package product

import "time"

type (
	// Product model
	Product struct {
		ID          int       `json:"id"`
		Name        string    `json:"name"`
		Price       float64   `json:"price"`
		Description string    `json:"description"`
		Picture     string    `json:"picture"`
		Created     time.Time `json:"created"`
		Updated     time.Time `json:"updated"`
	}
)
