package models

import "time"

// Item represents an example data model.
// Add your own models in this package.
type Item struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// CreateItemRequest represents the request to create an item.
type CreateItemRequest struct {
	Name        string `json:"name" validate:"required"`
	Description string `json:"description,omitempty"`
}

// UpdateItemRequest represents the request to update an item.
type UpdateItemRequest struct {
	Name        string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`
}
