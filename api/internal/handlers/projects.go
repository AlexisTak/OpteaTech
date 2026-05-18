package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/optea-tech/api/internal/models"
)

type ProjectsHandler struct {
	store *Store
	db    *pgxpool.Pool
}

func NewProjectsHandler(store *Store, db *pgxpool.Pool) *ProjectsHandler {
	return &ProjectsHandler{store: store, db: db}
}

func (h *ProjectsHandler) ListPublic(c fiber.Ctx) error {
	category := c.Query("category")
	featured := c.Query("featured")
	if h.db != nil {
		projects, err := h.listProjectsDB(c.Context(), true, category, featured == "true")
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "unable to fetch projects"})
		}
		return c.JSON(projects)
	}

	h.store.mu.RLock()
	defer h.store.mu.RUnlock()

	result := make([]models.Project, 0, len(h.store.projects))
	for _, project := range h.store.projects {
		if project.Status != models.ProjectStatusPublished {
			continue
		}
		if category != "" && string(project.Category) != category {
			continue
		}
		if featured == "true" && !project.Featured {
			continue
		}
		result = append(result, project)
	}

	return c.JSON(result)
}

func (h *ProjectsHandler) GetBySlug(c fiber.Ctx) error {
	slug := c.Params("slug")
	if h.db != nil {
		project, err := h.getProjectBySlugDB(c.Context(), slug)
		if err != nil {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "project not found"})
		}
		return c.JSON(project)
	}

	h.store.mu.RLock()
	defer h.store.mu.RUnlock()

	for _, project := range h.store.projects {
		if project.Slug == slug && project.Status == models.ProjectStatusPublished {
			return c.JSON(project)
		}
	}

	return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "project not found"})
}

func (h *ProjectsHandler) ListAdmin(c fiber.Ctx) error {
	p := parsePagination(c)
	if h.db != nil {
		projects, total, err := h.listProjectsAdminDB(c, p)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "unable to fetch projects"})
		}
		setTotalCountHeader(c, total)
		return c.JSON(projects)
	}

	h.store.mu.RLock()
	defer h.store.mu.RUnlock()
	setTotalCountHeader(c, len(h.store.projects))
	return c.JSON(applySlicePagination(h.store.projects, p))
}

func (h *ProjectsHandler) Create(c fiber.Ctx) error {
	var input models.CreateProjectInput
	if err := c.Bind().Body(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid payload"})
	}

	now := time.Now().UTC()
	project := models.Project{
		ID:               uuid.New(),
		Title:            input.Title,
		Slug:             input.Slug,
		ShortDescription: input.ShortDescription,
		FullDescription:  input.FullDescription,
		CoverImageURL:    input.CoverImageURL,
		Images:           input.Images,
		Tags:             input.Tags,
		Category:         input.Category,
		ClientName:       input.ClientName,
		ProjectURL:       input.ProjectURL,
		GitHubURL:        input.GitHubURL,
		Status:           input.Status,
		Featured:         input.Featured,
		OrderIndex:       input.OrderIndex,
		CreatedAt:        now,
		UpdatedAt:        now,
	}
	if project.Status == "" {
		project.Status = models.ProjectStatusDraft
	}

	if h.db != nil {
		created, err := h.insertProjectDB(c.Context(), project)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "unable to create project"})
		}
		return c.Status(fiber.StatusCreated).JSON(created)
	}

	h.store.mu.Lock()
	h.store.projects = append(h.store.projects, project)
	h.store.mu.Unlock()

	return c.Status(fiber.StatusCreated).JSON(project)
}

func (h *ProjectsHandler) Update(c fiber.Ctx) error {
	projectID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid id"})
	}

	var input models.UpdateProjectInput
	if err := c.Bind().Body(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid payload"})
	}

	if h.db != nil {
		project, err := h.getProjectByIDDB(c.Context(), projectID)
		if err != nil {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "project not found"})
		}

		if input.Title != nil {
			project.Title = *input.Title
		}
		if input.Slug != nil {
			project.Slug = *input.Slug
		}
		if input.ShortDescription != nil {
			project.ShortDescription = input.ShortDescription
		}
		if input.FullDescription != nil {
			project.FullDescription = input.FullDescription
		}
		if input.CoverImageURL != nil {
			project.CoverImageURL = input.CoverImageURL
		}
		if input.Images != nil {
			project.Images = *input.Images
		}
		if input.Tags != nil {
			project.Tags = *input.Tags
		}
		if input.Category != nil {
			project.Category = *input.Category
		}
		if input.ClientName != nil {
			project.ClientName = input.ClientName
		}
		if input.ProjectURL != nil {
			project.ProjectURL = input.ProjectURL
		}
		if input.GitHubURL != nil {
			project.GitHubURL = input.GitHubURL
		}
		if input.Status != nil {
			project.Status = *input.Status
		}
		if input.Featured != nil {
			project.Featured = *input.Featured
		}
		if input.OrderIndex != nil {
			project.OrderIndex = *input.OrderIndex
		}
		project.UpdatedAt = time.Now().UTC()

		updated, err := h.updateProjectDB(c.Context(), project)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "unable to update project"})
		}
		return c.JSON(updated)
	}

	h.store.mu.Lock()
	defer h.store.mu.Unlock()

	for index, project := range h.store.projects {
		if project.ID != projectID {
			continue
		}
		if input.Title != nil {
			project.Title = *input.Title
		}
		if input.Slug != nil {
			project.Slug = *input.Slug
		}
		if input.ShortDescription != nil {
			project.ShortDescription = input.ShortDescription
		}
		if input.FullDescription != nil {
			project.FullDescription = input.FullDescription
		}
		if input.CoverImageURL != nil {
			project.CoverImageURL = input.CoverImageURL
		}
		if input.Images != nil {
			project.Images = *input.Images
		}
		if input.Tags != nil {
			project.Tags = *input.Tags
		}
		if input.Category != nil {
			project.Category = *input.Category
		}
		if input.ClientName != nil {
			project.ClientName = input.ClientName
		}
		if input.ProjectURL != nil {
			project.ProjectURL = input.ProjectURL
		}
		if input.GitHubURL != nil {
			project.GitHubURL = input.GitHubURL
		}
		if input.Status != nil {
			project.Status = *input.Status
		}
		if input.Featured != nil {
			project.Featured = *input.Featured
		}
		if input.OrderIndex != nil {
			project.OrderIndex = *input.OrderIndex
		}
		project.UpdatedAt = time.Now().UTC()
		h.store.projects[index] = project
		return c.JSON(project)
	}

	return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "project not found"})
}

func (h *ProjectsHandler) Delete(c fiber.Ctx) error {
	projectID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid id"})
	}

	if h.db != nil {
		result, err := h.db.Exec(c.Context(), `DELETE FROM projects WHERE id = $1`, projectID)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "unable to delete project"})
		}
		if result.RowsAffected() == 0 {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "project not found"})
		}
		return c.SendStatus(fiber.StatusNoContent)
	}

	h.store.mu.Lock()
	defer h.store.mu.Unlock()

	for index, project := range h.store.projects {
		if project.ID == projectID {
			h.store.projects = append(h.store.projects[:index], h.store.projects[index+1:]...)
			return c.SendStatus(fiber.StatusNoContent)
		}
	}

	return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "project not found"})
}

func (h *ProjectsHandler) listProjectsDB(ctx context.Context, publishedOnly bool, category string, featuredOnly bool) ([]models.Project, error) {
	query := `
		SELECT id, title, slug, short_description, full_description, cover_image_url, images, tags, category,
		       client_name, project_url, github_url, status, featured, order_index, created_at, updated_at
		FROM projects
		WHERE ($1::boolean = false OR status = 'published')
		  AND ($2::text = '' OR category = $2)
		  AND ($3::boolean = false OR featured = true)
		ORDER BY order_index ASC, created_at DESC
	`
	rows, err := h.db.Query(ctx, query, publishedOnly, category, featuredOnly)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	projects := make([]models.Project, 0)
	for rows.Next() {
		project, scanErr := scanProject(rows)
		if scanErr != nil {
			return nil, scanErr
		}
		projects = append(projects, project)
	}
	return projects, rows.Err()
}

func (h *ProjectsHandler) listProjectsAdminDB(c fiber.Ctx, p pagination) ([]models.Project, int, error) {
	ctx := c.Context()
	whereClause, whereArgs, nextIndex := buildWhereClause(c, map[string]string{
		"title":      "title",
		"slug":       "slug",
		"status":     "status",
		"category":   "category",
		"featured":   "featured",
		"clientName": "client_name",
		"createdAt":  "created_at",
		"updatedAt":  "updated_at",
	}, 1)
	orderClause := buildOrderClause(c, map[string]string{
		"title":      "title",
		"status":     "status",
		"category":   "category",
		"orderIndex": "order_index",
		"createdAt":  "created_at",
		"updatedAt":  "updated_at",
	}, "created_at", "DESC")

	var total int
	countQuery := "SELECT COUNT(*) FROM projects " + whereClause
	if err := h.db.QueryRow(ctx, countQuery, whereArgs...).Scan(&total); err != nil {
		return nil, 0, err
	}

	query := fmt.Sprintf(`
		SELECT id, title, slug, short_description, full_description, cover_image_url, images, tags, category,
		       client_name, project_url, github_url, status, featured, order_index, created_at, updated_at
		FROM projects
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

	projects := make([]models.Project, 0)
	for rows.Next() {
		project, scanErr := scanProject(rows)
		if scanErr != nil {
			return nil, 0, scanErr
		}
		projects = append(projects, project)
	}

	return projects, total, rows.Err()
}

func (h *ProjectsHandler) getProjectBySlugDB(ctx context.Context, slug string) (models.Project, error) {
	row := h.db.QueryRow(ctx, `
		SELECT id, title, slug, short_description, full_description, cover_image_url, images, tags, category,
		       client_name, project_url, github_url, status, featured, order_index, created_at, updated_at
		FROM projects
		WHERE slug = $1 AND status = 'published'
		LIMIT 1
	`, slug)
	return scanProject(row)
}

func (h *ProjectsHandler) getProjectByIDDB(ctx context.Context, id uuid.UUID) (models.Project, error) {
	row := h.db.QueryRow(ctx, `
		SELECT id, title, slug, short_description, full_description, cover_image_url, images, tags, category,
		       client_name, project_url, github_url, status, featured, order_index, created_at, updated_at
		FROM projects
		WHERE id = $1
		LIMIT 1
	`, id)
	return scanProject(row)
}

func (h *ProjectsHandler) insertProjectDB(ctx context.Context, project models.Project) (models.Project, error) {
	imagesJSON, err := json.Marshal(project.Images)
	if err != nil {
		return models.Project{}, err
	}

	row := h.db.QueryRow(ctx, `
		INSERT INTO projects (id, title, slug, short_description, full_description, cover_image_url, images, tags, category,
		                      client_name, project_url, github_url, status, featured, order_index, created_at, updated_at)
		VALUES ($1,$2,$3,$4,$5,$6,$7::jsonb,$8,$9,$10,$11,$12,$13,$14,$15,$16,$17)
		RETURNING id, title, slug, short_description, full_description, cover_image_url, images, tags, category,
		          client_name, project_url, github_url, status, featured, order_index, created_at, updated_at
	`, project.ID, project.Title, project.Slug, project.ShortDescription, project.FullDescription, project.CoverImageURL,
		string(imagesJSON), project.Tags, string(project.Category), project.ClientName, project.ProjectURL, project.GitHubURL,
		string(project.Status), project.Featured, project.OrderIndex, project.CreatedAt, project.UpdatedAt)

	return scanProject(row)
}

func (h *ProjectsHandler) updateProjectDB(ctx context.Context, project models.Project) (models.Project, error) {
	imagesJSON, err := json.Marshal(project.Images)
	if err != nil {
		return models.Project{}, err
	}

	row := h.db.QueryRow(ctx, `
		UPDATE projects
		SET title = $2,
		    slug = $3,
		    short_description = $4,
		    full_description = $5,
		    cover_image_url = $6,
		    images = $7::jsonb,
		    tags = $8,
		    category = $9,
		    client_name = $10,
		    project_url = $11,
		    github_url = $12,
		    status = $13,
		    featured = $14,
		    order_index = $15,
		    updated_at = $16
		WHERE id = $1
		RETURNING id, title, slug, short_description, full_description, cover_image_url, images, tags, category,
		          client_name, project_url, github_url, status, featured, order_index, created_at, updated_at
	`, project.ID, project.Title, project.Slug, project.ShortDescription, project.FullDescription, project.CoverImageURL,
		string(imagesJSON), project.Tags, string(project.Category), project.ClientName, project.ProjectURL, project.GitHubURL,
		string(project.Status), project.Featured, project.OrderIndex, project.UpdatedAt)

	return scanProject(row)
}

type projectScanner interface {
	Scan(dest ...interface{}) error
}

func scanProject(scanner projectScanner) (models.Project, error) {
	var project models.Project
	var category string
	var status string
	var imagesJSON []byte
	if err := scanner.Scan(
		&project.ID,
		&project.Title,
		&project.Slug,
		&project.ShortDescription,
		&project.FullDescription,
		&project.CoverImageURL,
		&imagesJSON,
		&project.Tags,
		&category,
		&project.ClientName,
		&project.ProjectURL,
		&project.GitHubURL,
		&status,
		&project.Featured,
		&project.OrderIndex,
		&project.CreatedAt,
		&project.UpdatedAt,
	); err != nil {
		return models.Project{}, err
	}

	project.Category = models.ProjectCategory(category)
	project.Status = models.ProjectStatus(status)
	if len(imagesJSON) == 0 {
		project.Images = []string{}
	} else {
		if err := json.Unmarshal(imagesJSON, &project.Images); err != nil {
			project.Images = []string{}
		}
	}
	if project.Tags == nil {
		project.Tags = []string{}
	}
	return project, nil
}
