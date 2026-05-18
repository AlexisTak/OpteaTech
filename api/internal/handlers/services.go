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

type ServicesHandler struct {
	store *Store
	db    *pgxpool.Pool
}

func NewServicesHandler(store *Store, db *pgxpool.Pool) *ServicesHandler {
	return &ServicesHandler{store: store, db: db}
}

func (h *ServicesHandler) ListPublic(c fiber.Ctx) error {
	if h.db != nil {
		services, err := h.listServicesDB(c.Context(), true)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "unable to fetch services"})
		}
		return c.JSON(services)
	}

	h.store.mu.RLock()
	defer h.store.mu.RUnlock()

	result := make([]models.Service, 0, len(h.store.services))
	for _, service := range h.store.services {
		if service.IsActive {
			result = append(result, service)
		}
	}
	return c.JSON(result)
}

func (h *ServicesHandler) ListAdmin(c fiber.Ctx) error {
	p := parsePagination(c)
	if h.db != nil {
		services, total, err := h.listServicesAdminDB(c, p)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "unable to fetch services"})
		}
		setTotalCountHeader(c, total)
		return c.JSON(services)
	}

	h.store.mu.RLock()
	defer h.store.mu.RUnlock()
	setTotalCountHeader(c, len(h.store.services))
	return c.JSON(applySlicePagination(h.store.services, p))
}

func (h *ServicesHandler) Create(c fiber.Ctx) error {
	var input models.CreateServiceInput
	if err := c.Bind().Body(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid payload"})
	}

	service := models.Service{
		ID:              uuid.New(),
		Name:            input.Name,
		Slug:            input.Slug,
		Description:     input.Description,
		LongDescription: input.LongDescription,
		Icon:            input.Icon,
		Color:           input.Color,
		Features:        input.Features,
		StartingPrice:   input.StartingPrice,
		OrderIndex:      input.OrderIndex,
		IsActive:        input.IsActive,
		CreatedAt:       time.Now().UTC(),
	}

	if h.db != nil {
		created, err := h.insertServiceDB(c.Context(), service)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "unable to create service"})
		}
		return c.Status(fiber.StatusCreated).JSON(created)
	}

	h.store.mu.Lock()
	h.store.services = append(h.store.services, service)
	h.store.mu.Unlock()

	return c.Status(fiber.StatusCreated).JSON(service)
}

func (h *ServicesHandler) Update(c fiber.Ctx) error {
	serviceID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid id"})
	}

	var input models.UpdateServiceInput
	if err := c.Bind().Body(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid payload"})
	}

	if h.db != nil {
		service, fetchErr := h.getServiceByIDDB(c.Context(), serviceID)
		if fetchErr != nil {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "service not found"})
		}
		if input.Name != nil {
			service.Name = *input.Name
		}
		if input.Slug != nil {
			service.Slug = *input.Slug
		}
		if input.Description != nil {
			service.Description = input.Description
		}
		if input.LongDescription != nil {
			service.LongDescription = input.LongDescription
		}
		if input.Icon != nil {
			service.Icon = input.Icon
		}
		if input.Color != nil {
			service.Color = input.Color
		}
		if input.Features != nil {
			service.Features = *input.Features
		}
		if input.StartingPrice != nil {
			service.StartingPrice = input.StartingPrice
		}
		if input.OrderIndex != nil {
			service.OrderIndex = *input.OrderIndex
		}
		if input.IsActive != nil {
			service.IsActive = *input.IsActive
		}

		updated, updateErr := h.updateServiceDB(c.Context(), service)
		if updateErr != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "unable to update service"})
		}
		return c.JSON(updated)
	}

	h.store.mu.Lock()
	defer h.store.mu.Unlock()

	for index, service := range h.store.services {
		if service.ID != serviceID {
			continue
		}
		if input.Name != nil {
			service.Name = *input.Name
		}
		if input.Slug != nil {
			service.Slug = *input.Slug
		}
		if input.Description != nil {
			service.Description = input.Description
		}
		if input.LongDescription != nil {
			service.LongDescription = input.LongDescription
		}
		if input.Icon != nil {
			service.Icon = input.Icon
		}
		if input.Color != nil {
			service.Color = input.Color
		}
		if input.Features != nil {
			service.Features = *input.Features
		}
		if input.StartingPrice != nil {
			service.StartingPrice = input.StartingPrice
		}
		if input.OrderIndex != nil {
			service.OrderIndex = *input.OrderIndex
		}
		if input.IsActive != nil {
			service.IsActive = *input.IsActive
		}
		h.store.services[index] = service
		return c.JSON(service)
	}

	return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "service not found"})
}

func (h *ServicesHandler) Delete(c fiber.Ctx) error {
	serviceID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid id"})
	}

	if h.db != nil {
		result, delErr := h.db.Exec(c.Context(), `DELETE FROM services WHERE id = $1`, serviceID)
		if delErr != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "unable to delete service"})
		}
		if result.RowsAffected() == 0 {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "service not found"})
		}
		return c.SendStatus(fiber.StatusNoContent)
	}

	h.store.mu.Lock()
	defer h.store.mu.Unlock()

	for index, service := range h.store.services {
		if service.ID == serviceID {
			h.store.services = append(h.store.services[:index], h.store.services[index+1:]...)
			return c.SendStatus(fiber.StatusNoContent)
		}
	}

	return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "service not found"})
}

func (h *ServicesHandler) listServicesDB(ctx context.Context, activeOnly bool) ([]models.Service, error) {
	rows, err := h.db.Query(ctx, `
		SELECT id, name, slug, description, long_description, icon, color, features, starting_price, order_index, is_active, created_at
		FROM services
		WHERE ($1::boolean = false OR is_active = true)
		ORDER BY order_index ASC, created_at DESC
	`, activeOnly)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	services := make([]models.Service, 0)
	for rows.Next() {
		service, scanErr := scanService(rows)
		if scanErr != nil {
			return nil, scanErr
		}
		services = append(services, service)
	}
	return services, rows.Err()
}

func (h *ServicesHandler) listServicesAdminDB(c fiber.Ctx, p pagination) ([]models.Service, int, error) {
	ctx := c.Context()
	whereClause, whereArgs, nextIndex := buildWhereClause(c, map[string]string{
		"name":       "name",
		"slug":       "slug",
		"isActive":   "is_active",
		"orderIndex": "order_index",
		"createdAt":  "created_at",
	}, 1)
	orderClause := buildOrderClause(c, map[string]string{
		"name":       "name",
		"isActive":   "is_active",
		"orderIndex": "order_index",
		"createdAt":  "created_at",
	}, "order_index", "ASC")

	var total int
	countQuery := "SELECT COUNT(*) FROM services " + whereClause
	if err := h.db.QueryRow(ctx, countQuery, whereArgs...).Scan(&total); err != nil {
		return nil, 0, err
	}

	query := fmt.Sprintf(`
		SELECT id, name, slug, description, long_description, icon, color, features, starting_price, order_index, is_active, created_at
		FROM services
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

	services := make([]models.Service, 0)
	for rows.Next() {
		service, scanErr := scanService(rows)
		if scanErr != nil {
			return nil, 0, scanErr
		}
		services = append(services, service)
	}

	return services, total, rows.Err()
}

func (h *ServicesHandler) getServiceByIDDB(ctx context.Context, id uuid.UUID) (models.Service, error) {
	row := h.db.QueryRow(ctx, `
		SELECT id, name, slug, description, long_description, icon, color, features, starting_price, order_index, is_active, created_at
		FROM services
		WHERE id = $1
		LIMIT 1
	`, id)
	return scanService(row)
}

func (h *ServicesHandler) insertServiceDB(ctx context.Context, service models.Service) (models.Service, error) {
	row := h.db.QueryRow(ctx, `
		INSERT INTO services (id, name, slug, description, long_description, icon, color, features, starting_price, order_index, is_active, created_at)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12)
		RETURNING id, name, slug, description, long_description, icon, color, features, starting_price, order_index, is_active, created_at
	`, service.ID, service.Name, service.Slug, service.Description, service.LongDescription, service.Icon, service.Color,
		service.Features, service.StartingPrice, service.OrderIndex, service.IsActive, service.CreatedAt)
	return scanService(row)
}

func (h *ServicesHandler) updateServiceDB(ctx context.Context, service models.Service) (models.Service, error) {
	row := h.db.QueryRow(ctx, `
		UPDATE services
		SET name = $2,
		    slug = $3,
		    description = $4,
		    long_description = $5,
		    icon = $6,
		    color = $7,
		    features = $8,
		    starting_price = $9,
		    order_index = $10,
		    is_active = $11
		WHERE id = $1
		RETURNING id, name, slug, description, long_description, icon, color, features, starting_price, order_index, is_active, created_at
	`, service.ID, service.Name, service.Slug, service.Description, service.LongDescription, service.Icon, service.Color,
		service.Features, service.StartingPrice, service.OrderIndex, service.IsActive)
	return scanService(row)
}

type serviceScanner interface {
	Scan(dest ...interface{}) error
}

func scanService(scanner serviceScanner) (models.Service, error) {
	var service models.Service
	if err := scanner.Scan(
		&service.ID,
		&service.Name,
		&service.Slug,
		&service.Description,
		&service.LongDescription,
		&service.Icon,
		&service.Color,
		&service.Features,
		&service.StartingPrice,
		&service.OrderIndex,
		&service.IsActive,
		&service.CreatedAt,
	); err != nil {
		return models.Service{}, err
	}
	if service.Features == nil {
		service.Features = []string{}
	}
	return service, nil
}
