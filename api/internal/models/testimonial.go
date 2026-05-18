package models

import (
	"time"

	"github.com/google/uuid"
)

type Testimonial struct {
	ID              uuid.UUID  `json:"id" db:"id"`
	ClientName      string     `json:"clientName" db:"client_name"`
	ClientRole      *string    `json:"clientRole,omitempty" db:"client_role"`
	ClientCompany   *string    `json:"clientCompany,omitempty" db:"client_company"`
	ClientAvatarURL *string    `json:"clientAvatarUrl,omitempty" db:"client_avatar_url"`
	Content         string     `json:"content" db:"content"`
	Rating          int        `json:"rating" db:"rating"`
	ProjectID       *uuid.UUID `json:"projectId,omitempty" db:"project_id"`
	IsActive        bool       `json:"isActive" db:"is_active"`
	OrderIndex      int        `json:"orderIndex" db:"order_index"`
	CreatedAt       time.Time  `json:"createdAt" db:"created_at"`
}

type CreateTestimonialInput struct {
	ClientName      string     `json:"clientName" validate:"required,min=2,max=100"`
	ClientRole      *string    `json:"clientRole,omitempty"`
	ClientCompany   *string    `json:"clientCompany,omitempty"`
	ClientAvatarURL *string    `json:"clientAvatarUrl,omitempty"`
	Content         string     `json:"content" validate:"required,min=10"`
	Rating          int        `json:"rating" validate:"required,min=1,max=5"`
	ProjectID       *uuid.UUID `json:"projectId,omitempty"`
	IsActive        bool       `json:"isActive"`
	OrderIndex      int        `json:"orderIndex"`
}

type UpdateTestimonialInput struct {
	ClientName      *string    `json:"clientName,omitempty" validate:"omitempty,min=2,max=100"`
	ClientRole      *string    `json:"clientRole,omitempty"`
	ClientCompany   *string    `json:"clientCompany,omitempty"`
	ClientAvatarURL *string    `json:"clientAvatarUrl,omitempty"`
	Content         *string    `json:"content,omitempty" validate:"omitempty,min=10"`
	Rating          *int       `json:"rating,omitempty" validate:"omitempty,min=1,max=5"`
	ProjectID       *uuid.UUID `json:"projectId,omitempty"`
	IsActive        *bool      `json:"isActive,omitempty"`
	OrderIndex      *int       `json:"orderIndex,omitempty"`
}
