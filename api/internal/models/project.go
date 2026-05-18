package models

import (
	"time"

	"github.com/google/uuid"
)

type ProjectCategory string

type ProjectStatus string

const (
	ProjectCategoryWeb      ProjectCategory = "web"
	ProjectCategoryLogiciel ProjectCategory = "logiciel"
	ProjectCategoryIA       ProjectCategory = "ia"
	ProjectCategoryConseil  ProjectCategory = "conseil"
)

const (
	ProjectStatusDraft     ProjectStatus = "draft"
	ProjectStatusPublished ProjectStatus = "published"
	ProjectStatusArchived  ProjectStatus = "archived"
)

type Project struct {
	ID               uuid.UUID       `json:"id" db:"id"`
	Title            string          `json:"title" db:"title"`
	Slug             string          `json:"slug" db:"slug"`
	ShortDescription *string         `json:"shortDescription,omitempty" db:"short_description"`
	FullDescription  *string         `json:"fullDescription,omitempty" db:"full_description"`
	CoverImageURL    *string         `json:"coverImageUrl,omitempty" db:"cover_image_url"`
	Images           []string        `json:"images" db:"images"`
	Tags             []string        `json:"tags" db:"tags"`
	Category         ProjectCategory `json:"category" db:"category"`
	ClientName       *string         `json:"clientName,omitempty" db:"client_name"`
	ProjectURL       *string         `json:"projectUrl,omitempty" db:"project_url"`
	GitHubURL        *string         `json:"githubUrl,omitempty" db:"github_url"`
	Status           ProjectStatus   `json:"status" db:"status"`
	Featured         bool            `json:"featured" db:"featured"`
	OrderIndex       int             `json:"orderIndex" db:"order_index"`
	CreatedAt        time.Time       `json:"createdAt" db:"created_at"`
	UpdatedAt        time.Time       `json:"updatedAt" db:"updated_at"`
}

type CreateProjectInput struct {
	Title            string          `json:"title" validate:"required,min=2,max=200"`
	Slug             string          `json:"slug" validate:"required,min=2,max=200"`
	ShortDescription *string         `json:"shortDescription,omitempty"`
	FullDescription  *string         `json:"fullDescription,omitempty"`
	CoverImageURL    *string         `json:"coverImageUrl,omitempty"`
	Images           []string        `json:"images,omitempty"`
	Tags             []string        `json:"tags,omitempty"`
	Category         ProjectCategory `json:"category" validate:"required,oneof=web logiciel ia conseil"`
	ClientName       *string         `json:"clientName,omitempty"`
	ProjectURL       *string         `json:"projectUrl,omitempty"`
	GitHubURL        *string         `json:"githubUrl,omitempty"`
	Status           ProjectStatus   `json:"status" validate:"omitempty,oneof=draft published archived"`
	Featured         bool            `json:"featured"`
	OrderIndex       int             `json:"orderIndex"`
}

type UpdateProjectInput struct {
	Title            *string          `json:"title,omitempty" validate:"omitempty,min=2,max=200"`
	Slug             *string          `json:"slug,omitempty" validate:"omitempty,min=2,max=200"`
	ShortDescription *string          `json:"shortDescription,omitempty"`
	FullDescription  *string          `json:"fullDescription,omitempty"`
	CoverImageURL    *string          `json:"coverImageUrl,omitempty"`
	Images           *[]string        `json:"images,omitempty"`
	Tags             *[]string        `json:"tags,omitempty"`
	Category         *ProjectCategory `json:"category,omitempty" validate:"omitempty,oneof=web logiciel ia conseil"`
	ClientName       *string          `json:"clientName,omitempty"`
	ProjectURL       *string          `json:"projectUrl,omitempty"`
	GitHubURL        *string          `json:"githubUrl,omitempty"`
	Status           *ProjectStatus   `json:"status,omitempty" validate:"omitempty,oneof=draft published archived"`
	Featured         *bool            `json:"featured,omitempty"`
	OrderIndex       *int             `json:"orderIndex,omitempty"`
}
