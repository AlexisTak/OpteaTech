package handlers

import (
	"context"
	"log/slog"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v3"

	"github.com/optea-tech/api/internal/models"
	"github.com/optea-tech/api/internal/repository"
	"github.com/optea-tech/api/internal/services"
)

type RequestsHandler struct {
	repo     *repository.RequestsRepo
	emailSvc *services.EmailService
	validate *validator.Validate
	baseURL  string
}

func NewRequestsHandler(repo *repository.RequestsRepo, emailSvc *services.EmailService, baseURL string) *RequestsHandler {
	return &RequestsHandler{
		repo:     repo,
		emailSvc: emailSvc,
		validate: validator.New(),
		baseURL:  baseURL,
	}
}

func (h *RequestsHandler) CreateRequest(c fiber.Ctx) error {
	if h.repo == nil || !h.repo.Ready() {
		return c.Status(fiber.StatusServiceUnavailable).JSON(fiber.Map{"error": "client requests require postgres"})
	}

	var input models.CreateRequestInput
	if err := c.Bind().Body(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Format de requete invalide.", "code": "INVALID_JSON"})
	}

	if input.Website != "" {
		slog.Info("client request honeypot triggered", "ip", c.IP(), "email", input.ClientEmail)
		return c.Status(fiber.StatusCreated).JSON(models.CreateRequestResponse{
			Message:   "Votre demande a ete envoyee. Verifiez votre email.",
			RequestID: "00000000-0000-0000-0000-000000000000",
		})
	}

	if err := h.validate.Struct(input); err != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
			"error":  "Donnees invalides.",
			"code":   "VALIDATION_ERROR",
			"fields": formatValidationErrors(err),
		})
	}

	rawToken, tokenHash, expiresAt, err := services.GenerateClientAccessToken()
	if err != nil {
		slog.Error("generate client access token failed", "error", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Erreur interne. Reessayez plus tard."})
	}

	request, err := h.repo.CreateRequest(c.Context(), input, tokenHash, expiresAt, c.IP(), c.Get("User-Agent"))
	if err != nil {
		slog.Error("create client request failed", "error", err, "email", input.ClientEmail)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Impossible d'enregistrer votre demande."})
	}

	magicLink := services.BuildClientMagicLink(h.baseURL, rawToken, request.ID.String())
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := h.emailSvc.SendAccessLink(ctx, services.AccessLinkEmailData{
		ClientName:   request.ClientName,
		ClientEmail:  request.ClientEmail,
		MagicLink:    magicLink,
		ServiceType:  string(request.ServiceType),
		RequestTitle: request.Title,
		ExpiresAt:    expiresAt,
	}); err != nil {
		slog.Warn("send client access link failed", "error", err, "request_id", request.ID)
	} else if err := h.repo.MarkEmailSent(context.Background(), request.ID); err != nil {
		slog.Warn("mark client access link as sent failed", "error", err, "request_id", request.ID)
	}

	go func() {
		if err := h.emailSvc.NotifyAdminNewRequest(context.Background(), request); err != nil {
			slog.Warn("notify admin about new client request failed", "error", err, "request_id", request.ID)
		}
	}()

	return c.Status(fiber.StatusCreated).JSON(models.CreateRequestResponse{
		Message:   "Demande recue. Un email avec votre lien de suivi vous a ete envoye.",
		RequestID: request.ID.String(),
	})
}
