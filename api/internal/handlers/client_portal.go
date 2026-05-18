package handlers

import (
	"context"
	"strings"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v3"

	"github.com/optea-tech/api/internal/middleware"
	"github.com/optea-tech/api/internal/models"
	"github.com/optea-tech/api/internal/repository"
	"github.com/optea-tech/api/internal/services"
)

type ClientPortalHandler struct {
	requestsRepo  *repository.RequestsRepo
	portalRepo    *repository.PortalRepo
	accessLogRepo *repository.AccessLogRepo
	emailSvc      *services.EmailService
	validate      *validator.Validate
	baseURL       string
}

func NewClientPortalHandler(requestsRepo *repository.RequestsRepo, portalRepo *repository.PortalRepo, accessLogRepo *repository.AccessLogRepo, emailSvc *services.EmailService, baseURL string) *ClientPortalHandler {
	return &ClientPortalHandler{
		requestsRepo:  requestsRepo,
		portalRepo:    portalRepo,
		accessLogRepo: accessLogRepo,
		emailSvc:      emailSvc,
		validate:      validator.New(),
		baseURL:       baseURL,
	}
}

func (h *ClientPortalHandler) ensureReady(c fiber.Ctx) error {
	if h.requestsRepo == nil || !h.requestsRepo.Ready() || h.portalRepo == nil || !h.portalRepo.Ready() {
		return c.Status(fiber.StatusServiceUnavailable).JSON(fiber.Map{"error": "client portal requires postgres"})
	}
	return nil
}

func (h *ClientPortalHandler) GetDashboard(c fiber.Ctx) error {
	if err := h.ensureReady(c); err != nil {
		return err
	}

	request := middleware.GetClientRequest(c)
	if request == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "acces refuse"})
	}

	milestones, _ := h.portalRepo.GetMilestones(c.Context(), request.ID)
	messages, _ := h.portalRepo.GetMessages(c.Context(), request.ID)
	deliverables, _ := h.portalRepo.GetDeliverables(c.Context(), request.ID)
	unreadCount, _ := h.portalRepo.CountUnreadMessages(c.Context(), request.ID)

	return c.JSON(models.DashboardData{
		Request:      buildPublicView(request),
		Milestones:   milestones,
		Messages:     messages,
		Deliverables: deliverables,
		UnreadCount:  unreadCount,
	})
}

func (h *ClientPortalHandler) SendMessage(c fiber.Ctx) error {
	if err := h.ensureReady(c); err != nil {
		return err
	}

	request := middleware.GetClientRequest(c)
	if request == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "acces refuse"})
	}

	var input models.SendMessageInput
	if err := c.Bind().Body(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Format invalide."})
	}

	if err := h.validate.Struct(input); err != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{"error": "Message invalide.", "fields": formatValidationErrors(err)})
	}

	message, err := h.portalRepo.CreateMessage(c.Context(), repository.CreateMessageInput{
		RequestID:   request.ID,
		SenderType:  models.SenderClient,
		SenderName:  request.ClientName,
		Content:     input.Content,
		Attachments: input.Attachments,
	})
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Impossible d'envoyer le message."})
	}

	go h.emailSvc.NotifyAdminNewClientMessage(context.Background(), request, message)
	return c.Status(fiber.StatusCreated).JSON(message)
}

func (h *ClientPortalHandler) DownloadDeliverable(c fiber.Ctx) error {
	if err := h.ensureReady(c); err != nil {
		return err
	}

	request := middleware.GetClientRequest(c)
	if request == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "acces refuse"})
	}

	deliverable, err := h.portalRepo.GetDeliverable(c.Context(), c.Params("id"), request.ID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Fichier introuvable."})
	}

	go h.portalRepo.IncrementDownloadCount(context.Background(), deliverable.ID)
	return c.JSON(fiber.Map{"url": deliverable.FileURL, "expires_in": 0})
}

func (h *ClientPortalHandler) RequestNewToken(c fiber.Ctx) error {
	if err := h.ensureReady(c); err != nil {
		return err
	}

	var input struct {
		Email     string `json:"email" validate:"required,email"`
		RequestID string `json:"request_id" validate:"required,uuid"`
	}
	if err := c.Bind().Body(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Format invalide."})
	}

	if err := h.validate.Struct(input); err != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{"error": "Donnees invalides.", "fields": formatValidationErrors(err)})
	}

	request, err := h.requestsRepo.FindByIDAndEmail(c.Context(), input.RequestID, input.Email)
	if err != nil {
		return c.JSON(fiber.Map{"message": "Si votre email est associe a cette demande, vous recevrez un nouveau lien sous peu."})
	}

	rawToken, tokenHash, expiresAt, err := services.GenerateClientAccessToken()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Erreur interne."})
	}

	if err := h.requestsRepo.RegenerateToken(c.Context(), request.ID, tokenHash, expiresAt); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Erreur interne."})
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

	return c.JSON(fiber.Map{"message": "Si votre email est associe a cette demande, vous recevrez un nouveau lien sous peu."})
}

func (h *ClientPortalHandler) AcceptQuote(c fiber.Ctx) error {
	if err := h.ensureReady(c); err != nil {
		return err
	}

	request := middleware.GetClientRequest(c)
	if request == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "acces refuse"})
	}

	if request.QuoteAmount == nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Aucun devis disponible pour cette demande."})
	}
	if request.QuoteAcceptedAt != nil {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{"error": "Le devis a deja ete accepte."})
	}

	if err := h.requestsRepo.AcceptQuote(c.Context(), request.ID); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Impossible d'enregistrer votre accord."})
	}

	go h.emailSvc.NotifyAdminQuoteAccepted(context.Background(), request)
	return c.JSON(fiber.Map{"message": "Devis accepte. Nous vous contacterons tres prochainement pour demarrer le projet."})
}

func buildPublicView(request *models.ClientRequest) models.RequestPublicView {
	return models.RequestPublicView{
		ID:                request.ID,
		ClientName:        request.ClientName,
		Title:             request.Title,
		ServiceType:       request.ServiceType,
		Status:            request.Status,
		StatusLabel:       statusLabels[request.Status],
		Progress:          request.Progress,
		QuoteAmount:       request.QuoteAmount,
		QuotePDFURL:       request.QuotePDFURL,
		CreatedAt:         request.CreatedAt,
		ClientEmailMasked: maskEmail(request.ClientEmail),
	}
}

var statusLabels = map[models.RequestStatus]string{
	models.StatusNouveau:     "Demande recue",
	models.StatusEnEtude:     "En cours d'analyse",
	models.StatusDevisEnvoye: "Devis envoye - en attente de votre accord",
	models.StatusAccepte:     "Projet lance",
	models.StatusEnCours:     "En cours de developpement",
	models.StatusEnRevision:  "Phase de revision",
	models.StatusLivre:       "Livre - en attente de validation",
	models.StatusTermine:     "Projet termine",
	models.StatusAnnule:      "Annule",
}

func maskEmail(email string) string {
	parts := strings.Split(email, "@")
	if len(parts) != 2 {
		return "****"
	}
	if len(parts[0]) < 2 {
		return "**@" + parts[1]
	}
	return parts[0][:2] + "**@" + parts[1]
}
