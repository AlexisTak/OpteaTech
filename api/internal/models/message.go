package models

import (
	"time"

	"github.com/google/uuid"
)

type ContactMessage struct {
	ID              uuid.UUID `json:"id" db:"id"`
	Name            string    `json:"name" db:"name"`
	Email           string    `json:"email" db:"email"`
	Company         *string   `json:"company,omitempty" db:"company"`
	ServiceInterest *string   `json:"serviceInterest,omitempty" db:"service_interest"`
	BudgetRange     *string   `json:"budgetRange,omitempty" db:"budget_range"`
	Message         string    `json:"message" db:"message"`
	IsRead          bool      `json:"isRead" db:"is_read"`
	IsReplied       bool      `json:"isReplied" db:"is_replied"`
	IPAddress       *string   `json:"ipAddress,omitempty" db:"ip_address"`
	UserAgent       *string   `json:"userAgent,omitempty" db:"user_agent"`
	CreatedAt       time.Time `json:"createdAt" db:"created_at"`
}

type CreateMessageInput struct {
	Name            string  `json:"name" validate:"required,min=1,max=100"`
	Email           string  `json:"email" validate:"required,email"`
	Company         *string `json:"company,omitempty"`
	ServiceInterest *string `json:"serviceInterest,omitempty"`
	BudgetRange     *string `json:"budgetRange,omitempty"`
	Message         string  `json:"message" validate:"required,min=20"`
	Honeypot        string  `json:"website"`
	IPAddress       *string `json:"ipAddress,omitempty"`
	UserAgent       *string `json:"userAgent,omitempty"`
}
