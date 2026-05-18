package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/optea-tech/api/internal/models"
)

type PortalRepo struct {
	db *pgxpool.Pool
}

type CreateMessageInput struct {
	RequestID   uuid.UUID
	SenderType  models.SenderType
	SenderName  string
	Content     string
	Attachments []string
}

type AccessLogEntry struct {
	RequestID     *uuid.UUID
	Action        string
	IPAddress     string
	UserAgent     string
	Success       bool
	FailureReason string
}

type AccessLogRepo struct {
	db *pgxpool.Pool
}

func NewPortalRepo(db *pgxpool.Pool) *PortalRepo {
	return &PortalRepo{db: db}
}

func NewAccessLogRepo(db *pgxpool.Pool) *AccessLogRepo {
	return &AccessLogRepo{db: db}
}

func (r *PortalRepo) Ready() bool {
	return r != nil && r.db != nil
}

func (r *PortalRepo) GetMilestones(ctx context.Context, requestID uuid.UUID) ([]models.MilestoneView, error) {
	rows, err := r.db.Query(ctx, `
		SELECT id, title, description, status, order_index, due_date, completed_at
		FROM project_milestones
		WHERE request_id = $1 AND is_visible = true
		ORDER BY order_index ASC, created_at ASC
	`, requestID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := make([]models.MilestoneView, 0)
	for rows.Next() {
		var item models.MilestoneView
		if err := rows.Scan(&item.ID, &item.Title, &item.Description, &item.Status, &item.OrderIndex, &item.DueDate, &item.CompletedAt); err != nil {
			return nil, err
		}
		items = append(items, item)
	}

	return items, rows.Err()
}

func (r *PortalRepo) GetMessages(ctx context.Context, requestID uuid.UUID) ([]models.MessageView, error) {
	rows, err := r.db.Query(ctx, `
		SELECT id, sender_type, sender_name, content, attachments, is_read, created_at
		FROM project_messages
		WHERE request_id = $1
		ORDER BY created_at ASC
	`, requestID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := make([]models.MessageView, 0)
	for rows.Next() {
		var item models.MessageView
		var attachmentsRaw []byte
		if err := rows.Scan(&item.ID, &item.SenderType, &item.SenderName, &item.Content, &attachmentsRaw, &item.IsRead, &item.CreatedAt); err != nil {
			return nil, err
		}
		item.Attachments = decodeStringSlice(attachmentsRaw)
		items = append(items, item)
	}

	return items, rows.Err()
}

func (r *PortalRepo) GetDeliverables(ctx context.Context, requestID uuid.UUID) ([]models.DeliverableView, error) {
	rows, err := r.db.Query(ctx, `
		SELECT id, name, description, file_type, version, download_count, created_at
		FROM project_deliverables
		WHERE request_id = $1 AND is_visible = true
		ORDER BY created_at DESC
	`, requestID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := make([]models.DeliverableView, 0)
	for rows.Next() {
		var item models.DeliverableView
		if err := rows.Scan(&item.ID, &item.Name, &item.Description, &item.FileType, &item.Version, &item.DownloadCount, &item.CreatedAt); err != nil {
			return nil, err
		}
		items = append(items, item)
	}

	return items, rows.Err()
}

func (r *PortalRepo) CountUnreadMessages(ctx context.Context, requestID uuid.UUID) (int, error) {
	var count int
	err := r.db.QueryRow(ctx, `SELECT COUNT(*) FROM project_messages WHERE request_id = $1 AND sender_type = 'admin' AND is_read = false`, requestID).Scan(&count)
	return count, err
}

func (r *PortalRepo) CreateMessage(ctx context.Context, input CreateMessageInput) (*models.ProjectMessage, error) {
	attachments, err := json.Marshal(defaultStringSlice(input.Attachments))
	if err != nil {
		return nil, fmt.Errorf("marshal message attachments: %w", err)
	}

	row := r.db.QueryRow(ctx, `
		INSERT INTO project_messages (request_id, sender_type, sender_name, content, attachments)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, request_id, sender_type, sender_name, content, attachments, is_read, read_at, created_at
	`, input.RequestID, input.SenderType, input.SenderName, input.Content, attachments)

	var message models.ProjectMessage
	var attachmentsRaw []byte
	if err := row.Scan(&message.ID, &message.RequestID, &message.SenderType, &message.SenderName, &message.Content, &attachmentsRaw, &message.IsRead, &message.ReadAt, &message.CreatedAt); err != nil {
		return nil, err
	}
	message.Attachments = decodeStringSlice(attachmentsRaw)

	return &message, nil
}

func (r *PortalRepo) CreateMilestone(ctx context.Context, requestID uuid.UUID, input models.CreateMilestoneInput) (*models.ProjectMilestone, error) {
	isVisible := true
	if input.IsVisible != nil {
		isVisible = *input.IsVisible
	}

	row := r.db.QueryRow(ctx, `
		INSERT INTO project_milestones (request_id, title, description, status, order_index, due_date, is_visible)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id, request_id, title, description, status, order_index, due_date, completed_at, is_visible, created_at, updated_at
	`, requestID, input.Title, input.Description, input.Status, input.OrderIndex, input.DueDate, isVisible)

	var milestone models.ProjectMilestone
	err := row.Scan(&milestone.ID, &milestone.RequestID, &milestone.Title, &milestone.Description, &milestone.Status, &milestone.OrderIndex, &milestone.DueDate, &milestone.CompletedAt, &milestone.IsVisible, &milestone.CreatedAt, &milestone.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &milestone, nil
}

func (r *PortalRepo) UpdateMilestone(ctx context.Context, requestID uuid.UUID, milestoneID uuid.UUID, input models.UpdateMilestoneInput) (*models.ProjectMilestone, error) {
	row := r.db.QueryRow(ctx, `
		UPDATE project_milestones
		SET title = COALESCE($1, title),
		    description = COALESCE($2, description),
		    status = COALESCE($3, status),
		    order_index = COALESCE($4, order_index),
		    due_date = COALESCE($5, due_date),
		    completed_at = COALESCE($6, completed_at),
		    is_visible = COALESCE($7, is_visible)
		WHERE id = $8 AND request_id = $9
		RETURNING id, request_id, title, description, status, order_index, due_date, completed_at, is_visible, created_at, updated_at
	`, input.Title, input.Description, input.Status, input.OrderIndex, input.DueDate, input.CompletedAt, input.IsVisible, milestoneID, requestID)

	var milestone models.ProjectMilestone
	err := row.Scan(&milestone.ID, &milestone.RequestID, &milestone.Title, &milestone.Description, &milestone.Status, &milestone.OrderIndex, &milestone.DueDate, &milestone.CompletedAt, &milestone.IsVisible, &milestone.CreatedAt, &milestone.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &milestone, nil
}

func (r *PortalRepo) CreateDeliverable(ctx context.Context, requestID uuid.UUID, input models.CreateDeliverableInput) (*models.ProjectDeliverable, error) {
	isVisible := true
	if input.IsVisible != nil {
		isVisible = *input.IsVisible
	}

	fileType := strings.TrimSpace(input.FileType)
	if fileType == "" {
		fileType = "file"
	}

	row := r.db.QueryRow(ctx, `
		INSERT INTO project_deliverables (request_id, name, description, file_url, file_type, file_size, version, is_visible)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING id, request_id, name, description, file_url, file_type, file_size, version, is_visible, download_count, created_at
	`, requestID, input.Name, input.Description, input.FileURL, fileType, input.FileSize, input.Version, isVisible)

	var deliverable models.ProjectDeliverable
	err := row.Scan(&deliverable.ID, &deliverable.RequestID, &deliverable.Name, &deliverable.Description, &deliverable.FileURL, &deliverable.FileType, &deliverable.FileSize, &deliverable.Version, &deliverable.IsVisible, &deliverable.DownloadCount, &deliverable.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &deliverable, nil
}

func (r *PortalRepo) GetDeliverable(ctx context.Context, deliverableID string, requestID uuid.UUID) (*models.ProjectDeliverable, error) {
	row := r.db.QueryRow(ctx, `
		SELECT id, request_id, name, description, file_url, file_type, file_size, version, is_visible, download_count, created_at
		FROM project_deliverables
		WHERE id = $1 AND request_id = $2 AND is_visible = true
	`, deliverableID, requestID)

	var deliverable models.ProjectDeliverable
	err := row.Scan(&deliverable.ID, &deliverable.RequestID, &deliverable.Name, &deliverable.Description, &deliverable.FileURL, &deliverable.FileType, &deliverable.FileSize, &deliverable.Version, &deliverable.IsVisible, &deliverable.DownloadCount, &deliverable.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &deliverable, nil
}

func (r *PortalRepo) IncrementDownloadCount(ctx context.Context, deliverableID uuid.UUID) error {
	_, err := r.db.Exec(ctx, `UPDATE project_deliverables SET download_count = download_count + 1 WHERE id = $1`, deliverableID)
	return err
}

func (r *AccessLogRepo) Log(ctx context.Context, entry AccessLogEntry) error {
	if r == nil || r.db == nil {
		return nil
	}

	_, err := r.db.Exec(ctx, `
		INSERT INTO client_access_logs (request_id, action, ip_address, user_agent, success, failure_reason)
		VALUES ($1, $2, $3, $4, $5, $6)
	`, entry.RequestID, entry.Action, nullableString(entry.IPAddress), nullableString(entry.UserAgent), entry.Success, nullableString(entry.FailureReason))
	return err
}
