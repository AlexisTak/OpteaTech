package handlers

import (
	"context"
	"errors"
	"net/mail"
	"os"
	"strings"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v3"
	"github.com/resend/resend-go/v2"

	"github.com/optea-tech/api/internal/models"
)

type ContactHandler struct {
	validator *validator.Validate
	resend    *resend.Client
	messages  *MessagesHandler
}

func NewContactHandler(messages *MessagesHandler) *ContactHandler {
	apiKey := os.Getenv("RESEND_API_KEY")
	client := resend.NewClient(apiKey)

	return &ContactHandler{
		validator: validator.New(),
		resend:    client,
		messages:  messages,
	}
}

func (h *ContactHandler) Submit(c fiber.Ctx) error {
	var input models.CreateMessageInput
	if err := c.Bind().Body(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid payload"})
	}

	if strings.TrimSpace(input.Honeypot) != "" {
		return c.Status(fiber.StatusOK).JSON(fiber.Map{"status": "accepted"})
	}

	if err := h.validator.Struct(input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "validation failed"})
	}

	ipAddr := c.IP()
	userAgent := c.Get("User-Agent")
	input.IPAddress = &ipAddr
	input.UserAgent = &userAgent
	h.messages.Push(input)

	if _, err := mail.ParseAddress(input.Email); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid email"})
	}

	if err := h.sendEmails(input); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "unable to send email"})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"status": "message received"})
}

func (h *ContactHandler) sendEmails(input models.CreateMessageInput) error {
	from := os.Getenv("CONTACT_EMAIL_FROM")
	to := os.Getenv("CONTACT_EMAIL_TO")
	if from == "" || to == "" {
		return errors.New("email addresses are not configured")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	adminParams := &resend.SendEmailRequest{
		From:    from,
		To:      []string{to},
		Subject: "Nouveau message optea.tech",
		Html: "<h2>Nouveau lead</h2>" +
			"<p><strong>Nom:</strong> " + input.Name + "</p>" +
			"<p><strong>Email:</strong> " + input.Email + "</p>" +
			"<p><strong>Message:</strong> " + input.Message + "</p>",
	}

	_, err := h.resend.Emails.SendWithContext(ctx, adminParams)
	if err != nil {
		return err
	}

	confirmParams := &resend.SendEmailRequest{
		From:    from,
		To:      []string{input.Email},
		Subject: "Bien recu - optea.tech",
		Html:    "<p>Merci pour votre message. Nous revenons vers vous sous 48h.</p>",
	}

	_, err = h.resend.Emails.SendWithContext(ctx, confirmParams)
	return err
}
