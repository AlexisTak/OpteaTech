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

type MessagesHandler struct {
	store *Store
	db    *pgxpool.Pool
}

func NewMessagesHandler(store *Store, db *pgxpool.Pool) *MessagesHandler {
	return &MessagesHandler{store: store, db: db}
}

func (h *MessagesHandler) Push(input models.CreateMessageInput) {
	if h.db != nil {
		_ = h.insertMessage(context.Background(), input)
		return
	}

	h.store.mu.Lock()
	defer h.store.mu.Unlock()

	h.store.messages = append(h.store.messages, models.ContactMessage{
		ID:              uuid.New(),
		Name:            input.Name,
		Email:           input.Email,
		Company:         input.Company,
		ServiceInterest: input.ServiceInterest,
		BudgetRange:     input.BudgetRange,
		Message:         input.Message,
		IsRead:          false,
		IsReplied:       false,
		IPAddress:       input.IPAddress,
		UserAgent:       input.UserAgent,
		CreatedAt:       time.Now().UTC(),
	})
}

func (h *MessagesHandler) List(c fiber.Ctx) error {
	p := parsePagination(c)
	if h.db != nil {
		messages, total, err := h.listMessages(c, p)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "unable to fetch messages"})
		}
		setTotalCountHeader(c, total)
		return c.JSON(messages)
	}

	h.store.mu.RLock()
	defer h.store.mu.RUnlock()
	setTotalCountHeader(c, len(h.store.messages))
	return c.JSON(applySlicePagination(h.store.messages, p))
}

func (h *MessagesHandler) MarkRead(c fiber.Ctx) error {
	messageID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid id"})
	}

	if h.db != nil {
		message, err := h.markMessageRead(c.Context(), messageID)
		if err != nil {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "message not found"})
		}
		return c.JSON(message)
	}

	h.store.mu.Lock()
	defer h.store.mu.Unlock()
	for index, message := range h.store.messages {
		if message.ID == messageID {
			message.IsRead = true
			h.store.messages[index] = message
			return c.JSON(message)
		}
	}

	return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "message not found"})
}

func (h *MessagesHandler) Delete(c fiber.Ctx) error {
	messageID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid id"})
	}

	if h.db != nil {
		deleted, err := h.deleteMessage(c.Context(), messageID)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "unable to delete message"})
		}
		if !deleted {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "message not found"})
		}
		return c.SendStatus(fiber.StatusNoContent)
	}

	h.store.mu.Lock()
	defer h.store.mu.Unlock()
	for index, message := range h.store.messages {
		if message.ID == messageID {
			h.store.messages = append(h.store.messages[:index], h.store.messages[index+1:]...)
			return c.SendStatus(fiber.StatusNoContent)
		}
	}

	return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "message not found"})
}

func (h *MessagesHandler) Dashboard(c fiber.Ctx) error {
	if h.db != nil {
		dashboard, err := h.dashboardCounts(c.Context())
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "unable to build dashboard"})
		}
		return c.JSON(dashboard)
	}

	h.store.mu.RLock()
	defer h.store.mu.RUnlock()

	unread := 0
	for _, message := range h.store.messages {
		if !message.IsRead {
			unread++
		}
	}

	return c.JSON(fiber.Map{
		"messages_unread": unread,
		"projects":        len(h.store.projects),
		"services":        len(h.store.services),
		"testimonials":    len(h.store.testimonials),
	})
}

func (h *MessagesHandler) insertMessage(ctx context.Context, input models.CreateMessageInput) error {
	_, err := h.db.Exec(ctx, `
		INSERT INTO contact_messages (name, email, company, service_interest, budget_range, message, ip_address, user_agent)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`, input.Name, input.Email, input.Company, input.ServiceInterest, input.BudgetRange, input.Message, input.IPAddress, input.UserAgent)
	return err
}

func (h *MessagesHandler) listMessages(c fiber.Ctx, p pagination) ([]models.ContactMessage, int, error) {
	ctx := c.Context()
	whereClause, whereArgs, nextIndex := buildWhereClause(c, map[string]string{
		"name":             "name",
		"email":            "email",
		"company":          "company",
		"service_interest": "service_interest",
		"is_read":          "is_read",
		"is_replied":       "is_replied",
		"createdAt":        "created_at",
	}, 1)
	orderClause := buildOrderClause(c, map[string]string{
		"createdAt": "created_at",
		"name":      "name",
		"email":     "email",
		"isRead":    "is_read",
	}, "created_at", "DESC")

	var total int
	countQuery := "SELECT COUNT(*) FROM contact_messages " + whereClause
	if err := h.db.QueryRow(ctx, countQuery, whereArgs...).Scan(&total); err != nil {
		return nil, 0, err
	}

	query := fmt.Sprintf(`
		SELECT id, name, email, company, service_interest, budget_range, message, is_read, is_replied, ip_address, user_agent, created_at
		FROM contact_messages
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

	messages := make([]models.ContactMessage, 0)
	for rows.Next() {
		var message models.ContactMessage
		if err := rows.Scan(
			&message.ID,
			&message.Name,
			&message.Email,
			&message.Company,
			&message.ServiceInterest,
			&message.BudgetRange,
			&message.Message,
			&message.IsRead,
			&message.IsReplied,
			&message.IPAddress,
			&message.UserAgent,
			&message.CreatedAt,
		); err != nil {
			return nil, 0, err
		}
		messages = append(messages, message)
	}

	return messages, total, rows.Err()
}

func (h *MessagesHandler) markMessageRead(ctx context.Context, messageID uuid.UUID) (models.ContactMessage, error) {
	var message models.ContactMessage
	err := h.db.QueryRow(ctx, `
		UPDATE contact_messages
		SET is_read = true
		WHERE id = $1
		RETURNING id, name, email, company, service_interest, budget_range, message, is_read, is_replied, ip_address, user_agent, created_at
	`, messageID).Scan(
		&message.ID,
		&message.Name,
		&message.Email,
		&message.Company,
		&message.ServiceInterest,
		&message.BudgetRange,
		&message.Message,
		&message.IsRead,
		&message.IsReplied,
		&message.IPAddress,
		&message.UserAgent,
		&message.CreatedAt,
	)
	return message, err
}

func (h *MessagesHandler) deleteMessage(ctx context.Context, messageID uuid.UUID) (bool, error) {
	result, err := h.db.Exec(ctx, `DELETE FROM contact_messages WHERE id = $1`, messageID)
	if err != nil {
		return false, err
	}
	return result.RowsAffected() > 0, nil
}

func (h *MessagesHandler) dashboardCounts(ctx context.Context) (fiber.Map, error) {
	var messagesUnread int
	if err := h.db.QueryRow(ctx, `SELECT COUNT(*) FROM contact_messages WHERE is_read = false`).Scan(&messagesUnread); err != nil {
		return nil, err
	}

	var projects int
	if err := h.db.QueryRow(ctx, `SELECT COUNT(*) FROM projects`).Scan(&projects); err != nil {
		return nil, err
	}

	var services int
	if err := h.db.QueryRow(ctx, `SELECT COUNT(*) FROM services`).Scan(&services); err != nil {
		return nil, err
	}

	var testimonials int
	if err := h.db.QueryRow(ctx, `SELECT COUNT(*) FROM testimonials`).Scan(&testimonials); err != nil {
		return nil, err
	}

	return fiber.Map{
		"messages_unread": messagesUnread,
		"projects":        projects,
		"services":        services,
		"testimonials":    testimonials,
	}, nil
}
