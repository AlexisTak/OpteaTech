package handlers

import (
	"sync"
	"time"

	"github.com/google/uuid"

	"github.com/optea-tech/api/internal/models"
)

type Store struct {
	mu           sync.RWMutex
	projects     []models.Project
	services     []models.Service
	testimonials []models.Testimonial
	messages     []models.ContactMessage
	refresh      map[string]time.Time
}

func NewStore() *Store {
	now := time.Now().UTC()
	projectID := uuid.New()

	short := "Refonte complete d'un portail citoyen pour une collectivite."
	cover := "/project-1.svg"
	projectURL := "https://example.com"
	serviceDescription := "Sites performants, visibles et maintenables"
	icon := "Globe"
	color := "#4f8ef7"
	clientRole := "CEO"
	clientCompany := "FlowScale"

	return &Store{
		projects: []models.Project{
			{
				ID:               projectID,
				Title:            "Portail de services municipaux",
				Slug:             "portail-services-municipaux",
				ShortDescription: &short,
				CoverImageURL:    &cover,
				Tags:             []string{"Next.js", "Go", "Supabase"},
				Category:         models.ProjectCategoryWeb,
				ProjectURL:       &projectURL,
				Status:           models.ProjectStatusPublished,
				Featured:         true,
				OrderIndex:       0,
				CreatedAt:        now,
				UpdatedAt:        now,
			},
		},
		services: []models.Service{
			{
				ID:          uuid.New(),
				Name:        "Sites Web et Vitrines",
				Slug:        "sites-web-vitrines",
				Description: &serviceDescription,
				Icon:        &icon,
				Color:       &color,
				Features:    []string{"Next.js", "Design responsive", "SEO"},
				OrderIndex:  0,
				IsActive:    true,
				CreatedAt:   now,
			},
		},
		testimonials: []models.Testimonial{
			{
				ID:            uuid.New(),
				ClientName:    "Nicolas Rey",
				ClientRole:    &clientRole,
				ClientCompany: &clientCompany,
				Content:       "Execution nette, livrables propres et communication fluide.",
				Rating:        5,
				ProjectID:     &projectID,
				IsActive:      true,
				OrderIndex:    0,
				CreatedAt:     now,
			},
		},
		messages: []models.ContactMessage{},
		refresh:  make(map[string]time.Time),
	}
}
