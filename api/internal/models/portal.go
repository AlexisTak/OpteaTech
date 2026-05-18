package models

import (
	"time"

	"github.com/google/uuid"
)

type MilestoneStatus string

const (
	MilestonePending    MilestoneStatus = "pending"
	MilestoneInProgress MilestoneStatus = "in_progress"
	MilestoneDone       MilestoneStatus = "done"
	MilestoneBlocked    MilestoneStatus = "blocked"
)

type SenderType string

const (
	SenderAdmin  SenderType = "admin"
	SenderClient SenderType = "client"
)

type ProjectMilestone struct {
	ID          uuid.UUID       `json:"id" db:"id"`
	RequestID   uuid.UUID       `json:"request_id" db:"request_id"`
	Title       string          `json:"title" db:"title"`
	Description *string         `json:"description,omitempty" db:"description"`
	Status      MilestoneStatus `json:"status" db:"status"`
	OrderIndex  int             `json:"order_index" db:"order_index"`
	DueDate     *time.Time      `json:"due_date,omitempty" db:"due_date"`
	CompletedAt *time.Time      `json:"completed_at,omitempty" db:"completed_at"`
	IsVisible   bool            `json:"is_visible" db:"is_visible"`
	CreatedAt   time.Time       `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time       `json:"updated_at" db:"updated_at"`
}

type ProjectMessage struct {
	ID          uuid.UUID  `json:"id" db:"id"`
	RequestID   uuid.UUID  `json:"request_id" db:"request_id"`
	SenderType  SenderType `json:"sender_type" db:"sender_type"`
	SenderName  string     `json:"sender_name" db:"sender_name"`
	Content     string     `json:"content" db:"content"`
	Attachments []string   `json:"attachments" db:"attachments"`
	IsRead      bool       `json:"is_read" db:"is_read"`
	ReadAt      *time.Time `json:"read_at,omitempty" db:"read_at"`
	CreatedAt   time.Time  `json:"created_at" db:"created_at"`
}

type ProjectDeliverable struct {
	ID            uuid.UUID `json:"id" db:"id"`
	RequestID     uuid.UUID `json:"request_id" db:"request_id"`
	Name          string    `json:"name" db:"name"`
	Description   *string   `json:"description,omitempty" db:"description"`
	FileURL       string    `json:"-" db:"file_url"`
	FileType      string    `json:"file_type" db:"file_type"`
	FileSize      *int64    `json:"file_size,omitempty" db:"file_size"`
	Version       *string   `json:"version,omitempty" db:"version"`
	IsVisible     bool      `json:"is_visible" db:"is_visible"`
	DownloadCount int       `json:"download_count" db:"download_count"`
	CreatedAt     time.Time `json:"created_at" db:"created_at"`
}

type DashboardData struct {
	Request      RequestPublicView `json:"request"`
	Milestones   []MilestoneView   `json:"milestones"`
	Messages     []MessageView     `json:"messages"`
	Deliverables []DeliverableView `json:"deliverables"`
	UnreadCount  int               `json:"unread_messages_count"`
}

type RequestPublicView struct {
	ID                uuid.UUID     `json:"id"`
	ClientName        string        `json:"client_name"`
	Title             string        `json:"title"`
	ServiceType       ServiceType   `json:"service_type"`
	Status            RequestStatus `json:"status"`
	StatusLabel       string        `json:"status_label"`
	Progress          int           `json:"progress"`
	QuoteAmount       *float64      `json:"quote_amount,omitempty"`
	QuotePDFURL       *string       `json:"quote_pdf_url,omitempty"`
	CreatedAt         time.Time     `json:"created_at"`
	ClientEmailMasked string        `json:"client_email_masked"`
}

type MilestoneView struct {
	ID          uuid.UUID  `json:"id"`
	Title       string     `json:"title"`
	Description *string    `json:"description,omitempty"`
	Status      string     `json:"status"`
	OrderIndex  int        `json:"order_index"`
	DueDate     *time.Time `json:"due_date,omitempty"`
	CompletedAt *time.Time `json:"completed_at,omitempty"`
}

type MessageView struct {
	ID          uuid.UUID `json:"id"`
	SenderType  string    `json:"sender_type"`
	SenderName  string    `json:"sender_name"`
	Content     string    `json:"content"`
	Attachments []string  `json:"attachments"`
	IsRead      bool      `json:"is_read"`
	CreatedAt   time.Time `json:"created_at"`
}

type DeliverableView struct {
	ID            uuid.UUID `json:"id"`
	Name          string    `json:"name"`
	Description   *string   `json:"description,omitempty"`
	FileType      string    `json:"file_type"`
	Version       *string   `json:"version,omitempty"`
	DownloadCount int       `json:"download_count"`
	CreatedAt     time.Time `json:"created_at"`
}

type SendMessageInput struct {
	Content     string   `json:"content" validate:"required,min=2,max=2000"`
	Attachments []string `json:"attachments,omitempty" validate:"omitempty,max=10,dive,url"`
}

type CreateMilestoneInput struct {
	Title       string          `json:"title" validate:"required,min=2,max=200"`
	Description *string         `json:"description,omitempty"`
	Status      MilestoneStatus `json:"status" validate:"required,oneof=pending in_progress done blocked"`
	OrderIndex  int             `json:"order_index"`
	DueDate     *time.Time      `json:"due_date,omitempty"`
	IsVisible   *bool           `json:"is_visible,omitempty"`
}

type UpdateMilestoneInput struct {
	Title       *string          `json:"title,omitempty" validate:"omitempty,min=2,max=200"`
	Description *string          `json:"description,omitempty"`
	Status      *MilestoneStatus `json:"status,omitempty" validate:"omitempty,oneof=pending in_progress done blocked"`
	OrderIndex  *int             `json:"order_index,omitempty"`
	DueDate     *time.Time       `json:"due_date,omitempty"`
	CompletedAt *time.Time       `json:"completed_at,omitempty"`
	IsVisible   *bool            `json:"is_visible,omitempty"`
}

type CreateAdminMessageInput struct {
	SenderName  string   `json:"sender_name" validate:"required,min=2,max=100"`
	Content     string   `json:"content" validate:"required,min=2,max=2000"`
	Attachments []string `json:"attachments,omitempty" validate:"omitempty,max=10,dive,url"`
}

type CreateDeliverableInput struct {
	Name        string  `json:"name" validate:"required,min=2,max=200"`
	Description *string `json:"description,omitempty"`
	FileURL     string  `json:"file_url" validate:"required,url"`
	FileType    string  `json:"file_type" validate:"omitempty,max=50"`
	FileSize    *int64  `json:"file_size,omitempty" validate:"omitempty,gte=0"`
	Version     *string `json:"version,omitempty" validate:"omitempty,max=20"`
	IsVisible   *bool   `json:"is_visible,omitempty"`
}
