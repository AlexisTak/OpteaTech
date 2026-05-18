package handlers

import (
	"context"
	"strings"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"

	"github.com/optea-tech/api/internal/models"
	"github.com/optea-tech/api/internal/repository"
	"github.com/optea-tech/api/internal/services"
)

type AdminRequestsHandler struct {
	requestsRepo *repository.RequestsRepo
	portalRepo   *repository.PortalRepo
	emailSvc     *services.EmailService
	validate     *validator.Validate
	baseURL      string
}

func NewAdminRequestsHandler(requestsRepo *repository.RequestsRepo, portalRepo *repository.PortalRepo, emailSvc *services.EmailService, baseURL string) *AdminRequestsHandler {
	return &AdminRequestsHandler{
		requestsRepo: requestsRepo,
		portalRepo:   portalRepo,
		emailSvc:     emailSvc,
		validate:     validator.New(),
		baseURL:      baseURL,
	}
}

func (h *AdminRequestsHandler) ensureReady(c fiber.Ctx) error {
	if h.requestsRepo == nil || !h.requestsRepo.Ready() || h.portalRepo == nil || !h.portalRepo.Ready() {
		return c.Status(fiber.StatusServiceUnavailable).JSON(fiber.Map{"error": "client requests require postgres"})
	}
	return nil
}

func (h *AdminRequestsHandler) List(c fiber.Ctx) error {
	if err := h.ensureReady(c); err != nil {
		return err
	}

	p := parsePagination(c)
	requests, total, err := h.requestsRepo.List(c.Context(), repository.RequestListFilter{
		Status:     c.Query("status"),
		Query:      c.Query("q"),
		Offset:     p.offset,
		Limit:      p.limit,
		HasPaging:  p.valid,
		SortColumn: c.Query("_sort"),
		SortOrder:  c.Query("_order"),
	})
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "unable to fetch requests"})
	}

	setTotalCountHeader(c, total)
	return c.JSON(requests)
}

func (h *AdminRequestsHandler) Get(c fiber.Ctx) error {
	if err := h.ensureReady(c); err != nil {
		return err
	}

	request, err := h.requestsRepo.FindByID(c.Context(), c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "request not found"})
	}
	return c.JSON(request)
}

func (h *AdminRequestsHandler) UpdateStatus(c fiber.Ctx) error {
	if err := h.ensureReady(c); err != nil {
		return err
	}

	requestID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid id"})
	}

	var input models.UpdateRequestStatusInput
	if err := c.Bind().Body(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid payload"})
	}
	if err := h.validate.Struct(input); err != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{"error": "validation failed", "fields": formatValidationErrors(err)})
	}

	request, err := h.requestsRepo.UpdateStatus(c.Context(), requestID, input.Status)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "request not found"})
	}
	return c.JSON(request)
}

func (h *AdminRequestsHandler) UpdateProgress(c fiber.Ctx) error {
	if err := h.ensureReady(c); err != nil {
		return err
	}

	requestID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid id"})
	}

	var input models.UpdateRequestProgressInput
	if err := c.Bind().Body(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid payload"})
	}
	if err := h.validate.Struct(input); err != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{"error": "validation failed", "fields": formatValidationErrors(err)})
	}

	request, err := h.requestsRepo.UpdateProgress(c.Context(), requestID, input.Progress)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "request not found"})
	}
	return c.JSON(request)
}

func (h *AdminRequestsHandler) CreateMilestone(c fiber.Ctx) error {
	if err := h.ensureReady(c); err != nil {
		return err
	}

	requestID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid id"})
	}

	var input models.CreateMilestoneInput
	if err := c.Bind().Body(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid payload"})
	}
	if err := h.validate.Struct(input); err != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{"error": "validation failed", "fields": formatValidationErrors(err)})
	}

	milestone, err := h.portalRepo.CreateMilestone(c.Context(), requestID, input)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "unable to create milestone"})
	}
	return c.Status(fiber.StatusCreated).JSON(milestone)
}

func (h *AdminRequestsHandler) UpdateMilestone(c fiber.Ctx) error {
	if err := h.ensureReady(c); err != nil {
		return err
	}

	requestID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid request id"})
	}
	milestoneID, err := uuid.Parse(c.Params("mid"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid milestone id"})
	}

	var input models.UpdateMilestoneInput
	if err := c.Bind().Body(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid payload"})
	}
	if err := h.validate.Struct(input); err != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{"error": "validation failed", "fields": formatValidationErrors(err)})
	}

	milestone, err := h.portalRepo.UpdateMilestone(c.Context(), requestID, milestoneID, input)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "milestone not found"})
	}
	return c.JSON(milestone)
}

func (h *AdminRequestsHandler) SendMessage(c fiber.Ctx) error {
	if err := h.ensureReady(c); err != nil {
		return err
	}

	requestID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid id"})
	}

	var input models.CreateAdminMessageInput
	if err := c.Bind().Body(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid payload"})
	}
	if err := h.validate.Struct(input); err != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{"error": "validation failed", "fields": formatValidationErrors(err)})
	}

	message, err := h.portalRepo.CreateMessage(c.Context(), repository.CreateMessageInput{
		RequestID:   requestID,
		SenderType:  models.SenderAdmin,
		SenderName:  normalizeAdminSenderName(input.SenderName),
		Content:     input.Content,
		Attachments: input.Attachments,
	})
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "unable to send message"})
	}
	return c.Status(fiber.StatusCreated).JSON(message)
}

func (h *AdminRequestsHandler) AddDeliverable(c fiber.Ctx) error {
	if err := h.ensureReady(c); err != nil {
		return err
	}

	requestID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid id"})
	}

	var input models.CreateDeliverableInput
	if err := c.Bind().Body(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid payload"})
	}
	if err := h.validate.Struct(input); err != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{"error": "validation failed", "fields": formatValidationErrors(err)})
	}

	deliverable, err := h.portalRepo.CreateDeliverable(c.Context(), requestID, input)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "unable to add deliverable"})
	}
	return c.Status(fiber.StatusCreated).JSON(deliverable)
}

func (h *AdminRequestsHandler) SetQuote(c fiber.Ctx) error {
	if err := h.ensureReady(c); err != nil {
		return err
	}

	requestID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid id"})
	}

	var input models.SetQuoteInput
	if err := c.Bind().Body(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid payload"})
	}
	if err := h.validate.Struct(input); err != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{"error": "validation failed", "fields": formatValidationErrors(err)})
	}

	request, err := h.requestsRepo.SetQuote(c.Context(), requestID, input.Amount, input.Currency, input.ValidUntil, input.PDFURL)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "request not found"})
	}
	return c.JSON(request)
}

func (h *AdminRequestsHandler) RevokeToken(c fiber.Ctx) error {
	if err := h.ensureReady(c); err != nil {
		return err
	}

	requestID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid id"})
	}

	if err := h.requestsRepo.RevokeToken(c.Context(), requestID); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "unable to revoke token"})
	}
	return c.JSON(fiber.Map{"message": "token revoked"})
}

func (h *AdminRequestsHandler) RegenerateToken(c fiber.Ctx) error {
	if err := h.ensureReady(c); err != nil {
		return err
	}

	request, err := h.requestsRepo.FindByID(c.Context(), c.Params("id"))
	if err != nil {
		if err == pgx.ErrNoRows {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "request not found"})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "unable to fetch request"})
	}

	rawToken, tokenHash, expiresAt, err := services.GenerateClientAccessToken()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "unable to generate token"})
	}

	if err := h.requestsRepo.RegenerateToken(c.Context(), request.ID, tokenHash, expiresAt); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "unable to regenerate token"})
	}

	magicLink := services.BuildClientMagicLink(h.baseURL, rawToken, request.ID.String())
	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		_ = h.emailSvc.SendAccessLink(ctx, services.AccessLinkEmailData{
			ClientName:   request.ClientName,
			ClientEmail:  request.ClientEmail,
			MagicLink:    magicLink,
			ServiceType:  string(request.ServiceType),
			RequestTitle: request.Title,
			ExpiresAt:    expiresAt,
			IsRenewal:    true,
		})
		_ = h.requestsRepo.MarkEmailSent(context.Background(), request.ID)
	}()

	return c.JSON(fiber.Map{"message": "new client access link queued"})
}

func normalizeAdminSenderName(name string) string {
	trimmed := strings.TrimSpace(name)
	if trimmed == "" {
		return "Optea Tech"
	}
	return trimmed
}
