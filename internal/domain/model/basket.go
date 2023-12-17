package model

import "time"

type Basket struct {
	Id        uint      `json:"id, omitempty"`
	CreatedAt time.Time `json:"created_at, omitempty"`
	UpdatedAt time.Time `json:"updated_at, omitempty"`
	Data      string    `json:"data, omitempty"`
	State     string
}
