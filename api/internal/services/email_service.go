package services

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"html/template"
	"strings"
	"time"

	"github.com/resend/resend-go/v2"

	"github.com/optea-tech/api/internal/models"
)

var ErrEmailServiceDisabled = errors.New("email service is not configured")

type EmailService struct {
	client     *resend.Client
	fromEmail  string
	adminEmail string
}

func NewEmailService(apiKey, fromEmail, adminEmail string) *EmailService {
	service := &EmailService{
		fromEmail:  strings.TrimSpace(fromEmail),
		adminEmail: strings.TrimSpace(adminEmail),
	}
	if strings.TrimSpace(apiKey) != "" {
		service.client = resend.NewClient(strings.TrimSpace(apiKey))
	}
	return service
}

type AccessLinkEmailData struct {
	ClientName   string
	ClientEmail  string
	MagicLink    string
	ServiceType  string
	RequestTitle string
	ExpiresAt    time.Time
	IsRenewal    bool
}

func (s *EmailService) SendAccessLink(ctx context.Context, data AccessLinkEmailData) error {
	if !s.enabled() {
		return ErrEmailServiceDisabled
	}

	subject := "Votre lien de suivi - optea.tech"
	if data.IsRenewal {
		subject = "Votre nouveau lien de suivi - optea.tech"
	}

	html, err := renderEmailTemplate(accessLinkTemplate, data)
	if err != nil {
		return err
	}

	_, err = s.client.Emails.SendWithContext(ctx, &resend.SendEmailRequest{
		From:    s.fromEmail,
		To:      []string{data.ClientEmail},
		Subject: subject,
		Html:    html,
	})
	return err
}

func (s *EmailService) NotifyAdminNewRequest(ctx context.Context, req *models.ClientRequest) error {
	if !s.enabled() || s.adminEmail == "" {
		return nil
	}

	body := fmt.Sprintf("<p>Nouvelle demande client</p><p><strong>%s</strong> (%s)</p><p>%s</p>", req.ClientName, req.ClientEmail, req.Title)
	_, err := s.client.Emails.SendWithContext(ctx, &resend.SendEmailRequest{
		From:    s.fromEmail,
		To:      []string{s.adminEmail},
		Subject: "Nouvelle demande client - optea.tech",
		Html:    body,
	})
	return err
}

func (s *EmailService) NotifyAdminNewClientMessage(ctx context.Context, req *models.ClientRequest, msg *models.ProjectMessage) error {
	if !s.enabled() || s.adminEmail == "" {
		return nil
	}
	body := fmt.Sprintf("<p>Nouveau message client pour <strong>%s</strong></p><p>%s</p>", req.Title, msg.Content)
	_, err := s.client.Emails.SendWithContext(ctx, &resend.SendEmailRequest{
		From:    s.fromEmail,
		To:      []string{s.adminEmail},
		Subject: "Nouveau message client - optea.tech",
		Html:    body,
	})
	return err
}

func (s *EmailService) NotifyAdminQuoteAccepted(ctx context.Context, req *models.ClientRequest) error {
	if !s.enabled() || s.adminEmail == "" {
		return nil
	}
	body := fmt.Sprintf("<p>Le devis a ete accepte pour <strong>%s</strong>.</p>", req.Title)
	_, err := s.client.Emails.SendWithContext(ctx, &resend.SendEmailRequest{
		From:    s.fromEmail,
		To:      []string{s.adminEmail},
		Subject: "Devis accepte - optea.tech",
		Html:    body,
	})
	return err
}

func (s *EmailService) enabled() bool {
	return s.client != nil && s.fromEmail != ""
}

func renderEmailTemplate(tmpl string, data any) (string, error) {
	t, err := template.New("email").Parse(tmpl)
	if err != nil {
		return "", err
	}

	var buffer bytes.Buffer
	if err := t.Execute(&buffer, data); err != nil {
		return "", err
	}

	return buffer.String(), nil
}

const accessLinkTemplate = `<!DOCTYPE html>
<html lang="fr">
<head>
  <meta charset="UTF-8" />
  <meta name="viewport" content="width=device-width, initial-scale=1.0" />
  <title>Votre suivi de projet - optea.tech</title>
</head>
<body style="margin:0;padding:0;background:#f7f5ef;font-family:Arial,sans-serif;color:#0f1722;">
  <table width="100%" cellpadding="0" cellspacing="0" style="padding:32px 0;">
    <tr>
      <td align="center">
        <table width="560" cellpadding="0" cellspacing="0" style="background:#ffffff;border-radius:18px;border:1px solid #e6e1d5;overflow:hidden;">
          <tr>
            <td style="padding:28px 36px;border-bottom:1px solid #ece6d7;font-size:18px;font-weight:700;letter-spacing:0.02em;">optea.tech</td>
          </tr>
          <tr>
            <td style="padding:36px;">
              <p style="margin:0 0 8px;font-size:12px;text-transform:uppercase;letter-spacing:0.14em;color:#7c6a41;">{{if .IsRenewal}}Nouveau lien d'acces{{else}}Espace client securise{{end}}</p>
              <h1 style="margin:0 0 18px;font-size:30px;line-height:1.15;">Bonjour {{ .ClientName }},</h1>
              <p style="margin:0 0 14px;font-size:15px;line-height:1.7;">Votre demande pour <strong>{{ .RequestTitle }}</strong> a bien ete enregistree. Votre espace client vous permet de suivre l'avancement du projet, les jalons, les messages et les livrables.</p>
              <p style="margin:0 0 26px;font-size:15px;line-height:1.7;">Service concerne: <strong>{{ .ServiceType }}</strong></p>
              <p style="margin:0 0 28px;">
                <a href="{{ .MagicLink }}" style="display:inline-block;padding:14px 24px;border-radius:999px;background:#0f1722;color:#ffffff;text-decoration:none;font-weight:600;">Acceder a mon espace</a>
              </p>
              <div style="padding:16px 18px;border-radius:12px;background:#f5f0e2;font-size:13px;line-height:1.6;color:#5d5135;">
                Ce lien est personnel et expire le <strong>{{ .ExpiresAt.Format "02/01/2006" }}</strong>. Si vous ne l'avez pas demande, ignorez cet email.
              </div>
              <p style="margin:20px 0 0;font-size:12px;line-height:1.6;color:#746b58;word-break:break-all;">Lien de secours: {{ .MagicLink }}</p>
            </td>
          </tr>
        </table>
      </td>
    </tr>
  </table>
</body>
</html>`
