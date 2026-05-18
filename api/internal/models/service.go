package models

import (
	"time"

	"github.com/google/uuid"
)

type Service struct {
	ID              uuid.UUID `json:"id" db:"id"`
	Name            string    `json:"name" db:"name"`
	Slug            string    `json:"slug" db:"slug"`
	Description     *string   `json:"description,omitempty" db:"description"`
	LongDescription *string   `json:"longDescription,omitempty" db:"long_description"`
	Icon            *string   `json:"icon,omitempty" db:"icon"`
	Color           *string   `json:"color,omitempty" db:"color"`
	Features        []string  `json:"features" db:"features"`
	StartingPrice   *int      `json:"startingPrice,omitempty" db:"starting_price"`
	OrderIndex      int       `json:"orderIndex" db:"order_index"`
	IsActive        bool      `json:"isActive" db:"is_active"`
	CreatedAt       time.Time `json:"createdAt" db:"created_at"`
}

type CreateServiceInput struct {
	Name            string   `json:"name" validate:"required,min=2,max=100"`
	Slug            string   `json:"slug" validate:"required,min=2,max=100"`
	Description     *string  `json:"description,omitempty"`
	LongDescription *string  `json:"longDescription,omitempty"`
	Icon            *string  `json:"icon,omitempty"`
	Color           *string  `json:"color,omitempty"`
	Features        []string `json:"features,omitempty"`
	StartingPrice   *int     `json:"startingPrice,omitempty"`
	OrderIndex      int      `json:"orderIndex"`
	IsActive        bool     `json:"isActive"`
}

type UpdateServiceInput struct {
	Name            *string   `json:"name,omitempty" validate:"omitempty,min=2,max=100"`
	Slug            *string   `json:"slug,omitempty" validate:"omitempty,min=2,max=100"`
	Description     *string   `json:"description,omitempty"`
	LongDescription *string   `json:"longDescription,omitempty"`
	Icon            *string   `json:"icon,omitempty"`
	Color           *string   `json:"color,omitempty"`
	Features        *[]string `json:"features,omitempty"`
	StartingPrice   *int      `json:"startingPrice,omitempty"`
	OrderIndex      *int      `json:"orderIndex,omitempty"`
	IsActive        *bool     `json:"isActive,omitempty"`
}
