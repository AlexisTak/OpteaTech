package handlers

import (
	"context"
	"fmt"
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/optea-tech/api/internal/models"
)

type TestimonialsHandler struct {
	store *Store
	db    *pgxpool.Pool
}

func NewTestimonialsHandler(store *Store, db *pgxpool.Pool) *TestimonialsHandler {
	return &TestimonialsHandler{store: store, db: db}
}

func (h *TestimonialsHandler) ListPublic(c fiber.Ctx) error {
	if h.db != nil {
		testimonials, err := h.listTestimonialsDB(c.Context(), true)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "unable to fetch testimonials"})
		}
		return c.JSON(testimonials)
	}

	h.store.mu.RLock()
	defer h.store.mu.RUnlock()

	result := make([]models.Testimonial, 0, len(h.store.testimonials))
	for _, testimonial := range h.store.testimonials {
		if testimonial.IsActive {
			result = append(result, testimonial)
		}
	}

	return c.JSON(result)
}

func (h *TestimonialsHandler) ListAdmin(c fiber.Ctx) error {
	p := parsePagination(c)
	if h.db != nil {
		testimonials, total, err := h.listTestimonialsAdminDB(c, p)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "unable to fetch testimonials"})
		}
		setTotalCountHeader(c, total)
		return c.JSON(testimonials)
	}

	h.store.mu.RLock()
	defer h.store.mu.RUnlock()
	setTotalCountHeader(c, len(h.store.testimonials))
	return c.JSON(applySlicePagination(h.store.testimonials, p))
}

func (h *TestimonialsHandler) Create(c fiber.Ctx) error {
	var input models.CreateTestimonialInput
	if err := c.Bind().Body(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid payload"})
	}

	testimonial := models.Testimonial{
		ID:              uuid.New(),
		ClientName:      input.ClientName,
		ClientRole:      input.ClientRole,
		ClientCompany:   input.ClientCompany,
		ClientAvatarURL: input.ClientAvatarURL,
		Content:         input.Content,
		Rating:          input.Rating,
		ProjectID:       input.ProjectID,
		IsActive:        input.IsActive,
		OrderIndex:      input.OrderIndex,
		CreatedAt:       time.Now().UTC(),
	}

	if h.db != nil {
		created, err := h.insertTestimonialDB(c.Context(), testimonial)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "unable to create testimonial"})
		}
		return c.Status(fiber.StatusCreated).JSON(created)
	}

	h.store.mu.Lock()
	h.store.testimonials = append(h.store.testimonials, testimonial)
	h.store.mu.Unlock()

	return c.Status(fiber.StatusCreated).JSON(testimonial)
}

func (h *TestimonialsHandler) Update(c fiber.Ctx) error {
	testimonialID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid id"})
	}

	var input models.UpdateTestimonialInput
	if err := c.Bind().Body(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid payload"})
	}

	if h.db != nil {
		testimonial, fetchErr := h.getTestimonialByIDDB(c.Context(), testimonialID)
		if fetchErr != nil {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "testimonial not found"})
		}
		if input.ClientName != nil {
			testimonial.ClientName = *input.ClientName
		}
		if input.ClientRole != nil {
			testimonial.ClientRole = input.ClientRole
		}
		if input.ClientCompany != nil {
			testimonial.ClientCompany = input.ClientCompany
		}
		if input.ClientAvatarURL != nil {
			testimonial.ClientAvatarURL = input.ClientAvatarURL
		}
		if input.Content != nil {
			testimonial.Content = *input.Content
		}
		if input.Rating != nil {
			testimonial.Rating = *input.Rating
		}
		if input.ProjectID != nil {
			testimonial.ProjectID = input.ProjectID
		}
		if input.IsActive != nil {
			testimonial.IsActive = *input.IsActive
		}
		if input.OrderIndex != nil {
			testimonial.OrderIndex = *input.OrderIndex
		}

		updated, updateErr := h.updateTestimonialDB(c.Context(), testimonial)
		if updateErr != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "unable to update testimonial"})
		}
		return c.JSON(updated)
	}

	h.store.mu.Lock()
	defer h.store.mu.Unlock()

	for index, testimonial := range h.store.testimonials {
		if testimonial.ID != testimonialID {
			continue
		}
		if input.ClientName != nil {
			testimonial.ClientName = *input.ClientName
		}
		if input.ClientRole != nil {
			testimonial.ClientRole = input.ClientRole
		}
		if input.ClientCompany != nil {
			testimonial.ClientCompany = input.ClientCompany
		}
		if input.ClientAvatarURL != nil {
			testimonial.ClientAvatarURL = input.ClientAvatarURL
		}
		if input.Content != nil {
			testimonial.Content = *input.Content
		}
		if input.Rating != nil {
			testimonial.Rating = *input.Rating
		}
		if input.ProjectID != nil {
			testimonial.ProjectID = input.ProjectID
		}
		if input.IsActive != nil {
			testimonial.IsActive = *input.IsActive
		}
		if input.OrderIndex != nil {
			testimonial.OrderIndex = *input.OrderIndex
		}
		h.store.testimonials[index] = testimonial
		return c.JSON(testimonial)
	}

	return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "testimonial not found"})
}

func (h *TestimonialsHandler) Delete(c fiber.Ctx) error {
	testimonialID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid id"})
	}

	if h.db != nil {
		result, delErr := h.db.Exec(c.Context(), `DELETE FROM testimonials WHERE id = $1`, testimonialID)
		if delErr != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "unable to delete testimonial"})
		}
		if result.RowsAffected() == 0 {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "testimonial not found"})
		}
		return c.SendStatus(fiber.StatusNoContent)
	}

	h.store.mu.Lock()
	defer h.store.mu.Unlock()

	for index, testimonial := range h.store.testimonials {
		if testimonial.ID == testimonialID {
			h.store.testimonials = append(h.store.testimonials[:index], h.store.testimonials[index+1:]...)
			return c.SendStatus(fiber.StatusNoContent)
		}
	}

	return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "testimonial not found"})
}

func (h *TestimonialsHandler) listTestimonialsDB(ctx context.Context, activeOnly bool) ([]models.Testimonial, error) {
	rows, err := h.db.Query(ctx, `
		SELECT id, client_name, client_role, client_company, client_avatar_url, content, rating, project_id, is_active, order_index, created_at
		FROM testimonials
		WHERE ($1::boolean = false OR is_active = true)
		ORDER BY order_index ASC, created_at DESC
	`, activeOnly)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	testimonials := make([]models.Testimonial, 0)
	for rows.Next() {
		testimonial, scanErr := scanTestimonial(rows)
		if scanErr != nil {
			return nil, scanErr
		}
		testimonials = append(testimonials, testimonial)
	}
	return testimonials, rows.Err()
}

func (h *TestimonialsHandler) listTestimonialsAdminDB(c fiber.Ctx, p pagination) ([]models.Testimonial, int, error) {
	ctx := c.Context()
	whereClause, whereArgs, nextIndex := buildWhereClause(c, map[string]string{
		"clientName":    "client_name",
		"clientRole":    "client_role",
		"clientCompany": "client_company",
		"rating":        "rating",
		"isActive":      "is_active",
		"orderIndex":    "order_index",
		"createdAt":     "created_at",
	}, 1)
	orderClause := buildOrderClause(c, map[string]string{
		"clientName": "client_name",
		"rating":     "rating",
		"isActive":   "is_active",
		"orderIndex": "order_index",
		"createdAt":  "created_at",
	}, "order_index", "ASC")

	var total int
	countQuery := "SELECT COUNT(*) FROM testimonials " + whereClause
	if err := h.db.QueryRow(ctx, countQuery, whereArgs...).Scan(&total); err != nil {
		return nil, 0, err
	}

	query := fmt.Sprintf(`
		SELECT id, client_name, client_role, client_company, client_avatar_url, content, rating, project_id, is_active, order_index, created_at
		FROM testimonials
		%s
		ORDER BY %s
	`, whereClause, orderClause)
	args := append([]interface{}{}, whereArgs...)
	if p.valid {
		query += fmt.Sprintf(" LIMIT $%d OFFSET $%d", nextIndex, nextIndex+1)
		args = append(args, p.limit, p.offset)
	}

	rows, err := h.db.Query(ctx, query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	testimonials := make([]models.Testimonial, 0)
	for rows.Next() {
		testimonial, scanErr := scanTestimonial(rows)
		if scanErr != nil {
			return nil, 0, scanErr
		}
		testimonials = append(testimonials, testimonial)
	}

	return testimonials, total, rows.Err()
}

func (h *TestimonialsHandler) getTestimonialByIDDB(ctx context.Context, id uuid.UUID) (models.Testimonial, error) {
	row := h.db.QueryRow(ctx, `
		SELECT id, client_name, client_role, client_company, client_avatar_url, content, rating, project_id, is_active, order_index, created_at
		FROM testimonials
		WHERE id = $1
		LIMIT 1
	`, id)
	return scanTestimonial(row)
}

func (h *TestimonialsHandler) insertTestimonialDB(ctx context.Context, testimonial models.Testimonial) (models.Testimonial, error) {
	row := h.db.QueryRow(ctx, `
		INSERT INTO testimonials (id, client_name, client_role, client_company, client_avatar_url, content, rating, project_id, is_active, order_index, created_at)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11)
		RETURNING id, client_name, client_role, client_company, client_avatar_url, content, rating, project_id, is_active, order_index, created_at
	`, testimonial.ID, testimonial.ClientName, testimonial.ClientRole, testimonial.ClientCompany, testimonial.ClientAvatarURL,
		testimonial.Content, testimonial.Rating, testimonial.ProjectID, testimonial.IsActive, testimonial.OrderIndex, testimonial.CreatedAt)
	return scanTestimonial(row)
}

func (h *TestimonialsHandler) updateTestimonialDB(ctx context.Context, testimonial models.Testimonial) (models.Testimonial, error) {
	row := h.db.QueryRow(ctx, `
		UPDATE testimonials
		SET client_name = $2,
		    client_role = $3,
		    client_company = $4,
		    client_avatar_url = $5,
		    content = $6,
		    rating = $7,
		    project_id = $8,
		    is_active = $9,
		    order_index = $10
		WHERE id = $1
		RETURNING id, client_name, client_role, client_company, client_avatar_url, content, rating, project_id, is_active, order_index, created_at
	`, testimonial.ID, testimonial.ClientName, testimonial.ClientRole, testimonial.ClientCompany, testimonial.ClientAvatarURL,
		testimonial.Content, testimonial.Rating, testimonial.ProjectID, testimonial.IsActive, testimonial.OrderIndex)
	return scanTestimonial(row)
}

type testimonialScanner interface {
	Scan(dest ...interface{}) error
}

func scanTestimonial(scanner testimonialScanner) (models.Testimonial, error) {
	var testimonial models.Testimonial
	if err := scanner.Scan(
		&testimonial.ID,
		&testimonial.ClientName,
		&testimonial.ClientRole,
		&testimonial.ClientCompany,
		&testimonial.ClientAvatarURL,
		&testimonial.Content,
		&testimonial.Rating,
		&testimonial.ProjectID,
		&testimonial.IsActive,
		&testimonial.OrderIndex,
		&testimonial.CreatedAt,
	); err != nil {
		return models.Testimonial{}, err
	}
	return testimonial, nil
}
