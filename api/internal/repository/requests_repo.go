package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/optea-tech/api/internal/models"
)

type scanner interface {
	Scan(dest ...any) error
}

type RequestListFilter struct {
	Status     string
	Query      string
	Offset     int
	Limit      int
	SortColumn string
	SortOrder  string
	HasPaging  bool
}

type RequestsRepo struct {
	db *pgxpool.Pool
}

func NewRequestsRepo(db *pgxpool.Pool) *RequestsRepo {
	return &RequestsRepo{db: db}
}

func (r *RequestsRepo) Ready() bool {
	return r != nil && r.db != nil
}

func (r *RequestsRepo) CreateRequest(ctx context.Context, input models.CreateRequestInput, tokenHash string, tokenExpiresAt time.Time, ipAddress string, userAgent string) (*models.ClientRequest, error) {
	attachments, err := json.Marshal(defaultStringSlice(input.Attachments))
	if err != nil {
		return nil, fmt.Errorf("marshal attachments: %w", err)
	}

	metadata, err := json.Marshal(defaultMetadata(input.Metadata))
	if err != nil {
		return nil, fmt.Errorf("marshal metadata: %w", err)
	}

	row := r.db.QueryRow(ctx, `
		INSERT INTO client_requests (
			client_name, client_email, client_company, client_phone,
			service_type, title, description, budget_range, deadline,
			attachments, metadata, access_token_hash, token_expires_at,
			ip_address, user_agent
		) VALUES (
			$1, $2, $3, $4,
			$5, $6, $7, $8, $9,
			$10, $11, $12, $13,
			$14, $15
		)
		RETURNING
			id, client_name, client_email, client_company, client_phone,
			service_type, title, description, budget_range, deadline,
			attachments, metadata, status, progress, quote_amount,
			quote_currency, quote_valid_until, quote_accepted_at, quote_pdf_url,
			access_token_hash, token_created_at, token_expires_at, token_last_used,
			token_use_count, email_verified, email_sent_at, ip_address, user_agent,
			created_at, updated_at
	`,
		input.ClientName,
		input.ClientEmail,
		nullableString(input.ClientCompany),
		nullableString(input.ClientPhone),
		input.ServiceType,
		input.Title,
		input.Description,
		nullableString(input.BudgetRange),
		input.Deadline,
		attachments,
		metadata,
		tokenHash,
		tokenExpiresAt,
		nullableString(ipAddress),
		nullableString(userAgent),
	)

	request, err := scanClientRequest(row)
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}

	return request, nil
}

func (r *RequestsRepo) FindByTokenHash(ctx context.Context, tokenHash string) (*models.ClientRequest, error) {
	row := r.db.QueryRow(ctx, `
		SELECT
			id, client_name, client_email, client_company, client_phone,
			service_type, title, description, budget_range, deadline,
			attachments, metadata, status, progress, quote_amount,
			quote_currency, quote_valid_until, quote_accepted_at, quote_pdf_url,
			access_token_hash, token_created_at, token_expires_at, token_last_used,
			token_use_count, email_verified, email_sent_at, ip_address, user_agent,
			created_at, updated_at
		FROM client_requests
		WHERE access_token_hash = $1
		  AND token_expires_at > NOW()
		  AND status != 'annule'
		LIMIT 1
	`, tokenHash)

	return scanClientRequest(row)
}

func (r *RequestsRepo) FindByID(ctx context.Context, requestID string) (*models.ClientRequest, error) {
	row := r.db.QueryRow(ctx, `
		SELECT
			id, client_name, client_email, client_company, client_phone,
			service_type, title, description, budget_range, deadline,
			attachments, metadata, status, progress, quote_amount,
			quote_currency, quote_valid_until, quote_accepted_at, quote_pdf_url,
			access_token_hash, token_created_at, token_expires_at, token_last_used,
			token_use_count, email_verified, email_sent_at, ip_address, user_agent,
			created_at, updated_at
		FROM client_requests
		WHERE id = $1
	`, requestID)

	return scanClientRequest(row)
}

func (r *RequestsRepo) FindByIDAndEmail(ctx context.Context, requestID string, email string) (*models.ClientRequest, error) {
	row := r.db.QueryRow(ctx, `
		SELECT
			id, client_name, client_email, client_company, client_phone,
			service_type, title, description, budget_range, deadline,
			attachments, metadata, status, progress, quote_amount,
			quote_currency, quote_valid_until, quote_accepted_at, quote_pdf_url,
			access_token_hash, token_created_at, token_expires_at, token_last_used,
			token_use_count, email_verified, email_sent_at, ip_address, user_agent,
			created_at, updated_at
		FROM client_requests
		WHERE id = $1 AND LOWER(client_email) = LOWER($2)
	`, requestID, email)

	return scanClientRequest(row)
}

func (r *RequestsRepo) List(ctx context.Context, filter RequestListFilter) ([]models.ClientRequest, int, error) {
	conditions := make([]string, 0, 2)
	args := make([]any, 0, 4)
	index := 1

	if strings.TrimSpace(filter.Status) != "" {
		conditions = append(conditions, fmt.Sprintf("status = $%d", index))
		args = append(args, filter.Status)
		index++
	}

	if strings.TrimSpace(filter.Query) != "" {
		conditions = append(conditions, fmt.Sprintf("(client_name ILIKE $%d OR client_email ILIKE $%d OR title ILIKE $%d)", index, index, index))
		args = append(args, "%"+strings.TrimSpace(filter.Query)+"%")
		index++
	}

	whereClause := ""
	if len(conditions) > 0 {
		whereClause = " WHERE " + strings.Join(conditions, " AND ")
	}

	var total int
	if err := r.db.QueryRow(ctx, "SELECT COUNT(*) FROM client_requests"+whereClause, args...).Scan(&total); err != nil {
		return nil, 0, err
	}

	sortColumn := "created_at"
	if filter.SortColumn == "updated_at" || filter.SortColumn == "progress" || filter.SortColumn == "status" {
		sortColumn = filter.SortColumn
	}

	sortOrder := "DESC"
	if strings.EqualFold(filter.SortOrder, "ASC") {
		sortOrder = "ASC"
	}

	query := `
		SELECT
			id, client_name, client_email, client_company, client_phone,
			service_type, title, description, budget_range, deadline,
			attachments, metadata, status, progress, quote_amount,
			quote_currency, quote_valid_until, quote_accepted_at, quote_pdf_url,
			access_token_hash, token_created_at, token_expires_at, token_last_used,
			token_use_count, email_verified, email_sent_at, ip_address, user_agent,
			created_at, updated_at
		FROM client_requests` + whereClause + fmt.Sprintf(" ORDER BY %s %s", sortColumn, sortOrder)

	queryArgs := append([]any{}, args...)
	if filter.HasPaging {
		query += fmt.Sprintf(" LIMIT $%d OFFSET $%d", index, index+1)
		queryArgs = append(queryArgs, filter.Limit, filter.Offset)
	}

	rows, err := r.db.Query(ctx, query, queryArgs...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	requests := make([]models.ClientRequest, 0)
	for rows.Next() {
		request, scanErr := scanClientRequest(rows)
		if scanErr != nil {
			return nil, 0, scanErr
		}
		requests = append(requests, *request)
	}

	return requests, total, rows.Err()
}

func (r *RequestsRepo) UpdateTokenUsage(ctx context.Context, requestID uuid.UUID) error {
	_, err := r.db.Exec(ctx, `
		UPDATE client_requests
		SET token_last_used = NOW(),
		    token_use_count = token_use_count + 1,
		    email_verified = true
		WHERE id = $1
	`, requestID)
	return err
}

func (r *RequestsRepo) MarkEmailSent(ctx context.Context, requestID uuid.UUID) error {
	_, err := r.db.Exec(ctx, `UPDATE client_requests SET email_sent_at = NOW() WHERE id = $1`, requestID)
	return err
}

func (r *RequestsRepo) RevokeToken(ctx context.Context, requestID uuid.UUID) error {
	_, err := r.db.Exec(ctx, `UPDATE client_requests SET token_expires_at = NOW() - INTERVAL '1 second' WHERE id = $1`, requestID)
	return err
}

func (r *RequestsRepo) RegenerateToken(ctx context.Context, requestID uuid.UUID, newHash string, expiresAt time.Time) error {
	_, err := r.db.Exec(ctx, `
		UPDATE client_requests
		SET access_token_hash = $1,
		    token_expires_at = $2,
		    token_created_at = NOW(),
		    token_last_used = NULL,
		    token_use_count = 0,
		    email_verified = false
		WHERE id = $3
	`, newHash, expiresAt, requestID)
	return err
}

func (r *RequestsRepo) AcceptQuote(ctx context.Context, requestID uuid.UUID) error {
	_, err := r.db.Exec(ctx, `
		UPDATE client_requests
		SET quote_accepted_at = COALESCE(quote_accepted_at, NOW()),
		    status = CASE WHEN status IN ('nouveau', 'en_etude', 'devis_envoye') THEN 'accepte' ELSE status END,
		    updated_at = NOW()
		WHERE id = $1
	`, requestID)
	return err
}

func (r *RequestsRepo) UpdateStatus(ctx context.Context, requestID uuid.UUID, status models.RequestStatus) (*models.ClientRequest, error) {
	row := r.db.QueryRow(ctx, `
		UPDATE client_requests
		SET status = $1
		WHERE id = $2
		RETURNING
			id, client_name, client_email, client_company, client_phone,
			service_type, title, description, budget_range, deadline,
			attachments, metadata, status, progress, quote_amount,
			quote_currency, quote_valid_until, quote_accepted_at, quote_pdf_url,
			access_token_hash, token_created_at, token_expires_at, token_last_used,
			token_use_count, email_verified, email_sent_at, ip_address, user_agent,
			created_at, updated_at
	`, status, requestID)

	return scanClientRequest(row)
}

func (r *RequestsRepo) UpdateProgress(ctx context.Context, requestID uuid.UUID, progress int) (*models.ClientRequest, error) {
	row := r.db.QueryRow(ctx, `
		UPDATE client_requests
		SET progress = $1
		WHERE id = $2
		RETURNING
			id, client_name, client_email, client_company, client_phone,
			service_type, title, description, budget_range, deadline,
			attachments, metadata, status, progress, quote_amount,
			quote_currency, quote_valid_until, quote_accepted_at, quote_pdf_url,
			access_token_hash, token_created_at, token_expires_at, token_last_used,
			token_use_count, email_verified, email_sent_at, ip_address, user_agent,
			created_at, updated_at
	`, progress, requestID)

	return scanClientRequest(row)
}

func (r *RequestsRepo) SetQuote(ctx context.Context, requestID uuid.UUID, amount float64, currency string, validUntil *time.Time, pdfURL *string) (*models.ClientRequest, error) {
	if strings.TrimSpace(currency) == "" {
		currency = "EUR"
	}

	row := r.db.QueryRow(ctx, `
		UPDATE client_requests
		SET quote_amount = $1,
		    quote_currency = $2,
		    quote_valid_until = $3,
		    quote_pdf_url = $4,
		    status = CASE WHEN status IN ('nouveau', 'en_etude') THEN 'devis_envoye' ELSE status END
		WHERE id = $5
		RETURNING
			id, client_name, client_email, client_company, client_phone,
			service_type, title, description, budget_range, deadline,
			attachments, metadata, status, progress, quote_amount,
			quote_currency, quote_valid_until, quote_accepted_at, quote_pdf_url,
			access_token_hash, token_created_at, token_expires_at, token_last_used,
			token_use_count, email_verified, email_sent_at, ip_address, user_agent,
			created_at, updated_at
	`, amount, strings.ToUpper(currency), validUntil, pdfURL, requestID)

	return scanClientRequest(row)
}

func scanClientRequest(row scanner) (*models.ClientRequest, error) {
	var request models.ClientRequest
	var attachmentsRaw []byte
	var metadataRaw []byte

	err := row.Scan(
		&request.ID,
		&request.ClientName,
		&request.ClientEmail,
		&request.ClientCompany,
		&request.ClientPhone,
		&request.ServiceType,
		&request.Title,
		&request.Description,
		&request.BudgetRange,
		&request.Deadline,
		&attachmentsRaw,
		&metadataRaw,
		&request.Status,
		&request.Progress,
		&request.QuoteAmount,
		&request.QuoteCurrency,
		&request.QuoteValidUntil,
		&request.QuoteAcceptedAt,
		&request.QuotePDFURL,
		&request.AccessTokenHash,
		&request.TokenCreatedAt,
		&request.TokenExpiresAt,
		&request.TokenLastUsed,
		&request.TokenUseCount,
		&request.EmailVerified,
		&request.EmailSentAt,
		&request.IPAddress,
		&request.UserAgent,
		&request.CreatedAt,
		&request.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	request.Attachments = decodeStringSlice(attachmentsRaw)
	request.Metadata = decodeMetadata(metadataRaw)

	return &request, nil
}

func decodeStringSlice(raw []byte) []string {
	if len(raw) == 0 {
		return []string{}
	}

	var items []string
	if err := json.Unmarshal(raw, &items); err != nil {
		return []string{}
	}
	if items == nil {
		return []string{}
	}
	return items
}

func decodeMetadata(raw []byte) map[string]any {
	if len(raw) == 0 {
		return map[string]any{}
	}

	var value map[string]any
	if err := json.Unmarshal(raw, &value); err != nil {
		return map[string]any{}
	}
	if value == nil {
		return map[string]any{}
	}
	return value
}

func defaultStringSlice(items []string) []string {
	if items == nil {
		return []string{}
	}
	return items
}

func defaultMetadata(metadata map[string]any) map[string]any {
	if metadata == nil {
		return map[string]any{}
	}
	return metadata
}

func nullableString(value string) any {
	trimmed := strings.TrimSpace(value)
	if trimmed == "" {
		return nil
	}
	return trimmed
}

var ErrNotReady = fmt.Errorf("repository unavailable")

var ErrNoRows = pgx.ErrNoRows
