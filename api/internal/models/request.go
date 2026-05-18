package models

import (
	"time"

	"github.com/google/uuid"
)

type ServiceType string

const (
	ServiceSiteWeb  ServiceType = "site_web"
	ServiceLogiciel ServiceType = "logiciel"
	ServiceIA       ServiceType = "ia"
	ServiceConseil  ServiceType = "conseil"
	ServiceAutre    ServiceType = "autre"
)

type RequestStatus string

const (
	StatusNouveau     RequestStatus = "nouveau"
	StatusEnEtude     RequestStatus = "en_etude"
	StatusDevisEnvoye RequestStatus = "devis_envoye"
	StatusAccepte     RequestStatus = "accepte"
	StatusEnCours     RequestStatus = "en_cours"
	StatusEnRevision  RequestStatus = "en_revision"
	StatusLivre       RequestStatus = "livre"
	StatusTermine     RequestStatus = "termine"
	StatusAnnule      RequestStatus = "annule"
)

type ClientRequest struct {
	ID              uuid.UUID      `json:"id" db:"id"`
	ClientName      string         `json:"client_name" db:"client_name"`
	ClientEmail     string         `json:"client_email" db:"client_email"`
	ClientCompany   *string        `json:"client_company,omitempty" db:"client_company"`
	ClientPhone     *string        `json:"client_phone,omitempty" db:"client_phone"`
	ServiceType     ServiceType    `json:"service_type" db:"service_type"`
	Title           string         `json:"title" db:"title"`
	Description     string         `json:"description" db:"description"`
	BudgetRange     *string        `json:"budget_range,omitempty" db:"budget_range"`
	Deadline        *time.Time     `json:"deadline,omitempty" db:"deadline"`
	Attachments     []string       `json:"attachments,omitempty" db:"attachments"`
	Metadata        map[string]any `json:"metadata,omitempty" db:"metadata"`
	Status          RequestStatus  `json:"status" db:"status"`
	Progress        int            `json:"progress" db:"progress"`
	QuoteAmount     *float64       `json:"quote_amount,omitempty" db:"quote_amount"`
	QuoteCurrency   string         `json:"quote_currency,omitempty" db:"quote_currency"`
	QuoteValidUntil *time.Time     `json:"quote_valid_until,omitempty" db:"quote_valid_until"`
	QuoteAcceptedAt *time.Time     `json:"quote_accepted_at,omitempty" db:"quote_accepted_at"`
	QuotePDFURL     *string        `json:"quote_pdf_url,omitempty" db:"quote_pdf_url"`
	AccessTokenHash string         `json:"-" db:"access_token_hash"`
	TokenCreatedAt  time.Time      `json:"-" db:"token_created_at"`
	TokenExpiresAt  time.Time      `json:"-" db:"token_expires_at"`
	TokenLastUsed   *time.Time     `json:"-" db:"token_last_used"`
	TokenUseCount   int            `json:"-" db:"token_use_count"`
	EmailVerified   bool           `json:"email_verified" db:"email_verified"`
	EmailSentAt     *time.Time     `json:"-" db:"email_sent_at"`
	IPAddress       *string        `json:"-" db:"ip_address"`
	UserAgent       *string        `json:"-" db:"user_agent"`
	CreatedAt       time.Time      `json:"created_at" db:"created_at"`
	UpdatedAt       time.Time      `json:"updated_at" db:"updated_at"`
}

type CreateRequestInput struct {
	ClientName    string         `json:"client_name" validate:"required,min=2,max=100"`
	ClientEmail   string         `json:"client_email" validate:"required,email,max=200"`
	ClientCompany string         `json:"client_company" validate:"omitempty,max=100"`
	ClientPhone   string         `json:"client_phone" validate:"omitempty,max=30"`
	ServiceType   ServiceType    `json:"service_type" validate:"required,oneof=site_web logiciel ia conseil autre"`
	Title         string         `json:"title" validate:"required,min=5,max=200"`
	Description   string         `json:"description" validate:"required,min=20,max=5000"`
	BudgetRange   string         `json:"budget_range" validate:"omitempty,oneof=moins_2k 2k_5k 5k_15k 15k_plus"`
	Deadline      *time.Time     `json:"deadline,omitempty"`
	Attachments   []string       `json:"attachments,omitempty" validate:"omitempty,max=10,dive,url"`
	Metadata      map[string]any `json:"metadata,omitempty"`
	Website       string         `json:"website" validate:"omitempty,max=0"`
}

type CreateRequestResponse struct {
	Message   string `json:"message"`
	RequestID string `json:"request_id"`
}

type UpdateRequestStatusInput struct {
	Status RequestStatus `json:"status" validate:"required,oneof=nouveau en_etude devis_envoye accepte en_cours en_revision livre termine annule"`
}

type UpdateRequestProgressInput struct {
	Progress int `json:"progress" validate:"min=0,max=100"`
}

type SetQuoteInput struct {
	Amount     float64    `json:"amount" validate:"required,gt=0"`
	Currency   string     `json:"currency" validate:"omitempty,len=3"`
	ValidUntil *time.Time `json:"valid_until,omitempty"`
	PDFURL     *string    `json:"pdf_url,omitempty" validate:"omitempty,url"`
}
