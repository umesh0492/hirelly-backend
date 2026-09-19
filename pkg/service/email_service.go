package service

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/umesh0492/go-app-kit/notifications"
	"hirelly-backend/pkg/config"
	"hirelly-backend/pkg/domain"
)

// EmailService defines the interface for transactional email operations.
type EmailService interface {
	SendMandateConfirmation(ctx context.Context, mandate *domain.Mandate) (string, error)
}

type appKitEmailService struct {
	sender     notifications.Sender
	adminEmail string
}

// NewEmailService constructs an EmailService using Abeta-dev/go-app-kit.
func NewEmailService(cfg *config.Config) EmailService {
	if cfg.BrevoAPIKey == "" {
		log.Println("⚠️ Warning: BREVO_API_KEY is not set. Outbound emails will be simulated.")
		return &appKitEmailService{adminEmail: cfg.AdminEmail}
	}

	sender, err := notifications.NewBrevoSender(notifications.BrevoConfig{
		APIKey:      cfg.BrevoAPIKey,
		SenderName:  cfg.SenderName,
		SenderEmail: cfg.SenderEmail,
	})
	if err != nil {
		log.Printf("⚠️ Warning: failed to initialize Brevo sender via go-app-kit: %v", err)
	}

	return &appKitEmailService{
		sender:     sender,
		adminEmail: cfg.AdminEmail,
	}
}

func (s *appKitEmailService) SendMandateConfirmation(ctx context.Context, m *domain.Mandate) (string, error) {
	if s.sender == nil {
		log.Printf("ℹ️ [Email Simulated] Mandate %s confirmation for %s", m.ID, m.Email)
		return "simulated_no_key", nil
	}

	htmlContent := fmt.Sprintf(`<!DOCTYPE html>
<html>
<head>
  <meta charset="utf-8">
  <style>
    body { font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, Helvetica, Arial, sans-serif; background-color: #06080e; color: #f8fafc; padding: 30px 15px; margin: 0; }
    .card { max-width: 600px; margin: 0 auto; background: #0b0f19; border: 1px solid #d4af37; border-radius: 16px; padding: 36px 30px; box-shadow: 0 20px 40px rgba(0,0,0,0.6); }
    .header { font-size: 24px; font-weight: 800; color: #d4af37; margin-bottom: 8px; letter-spacing: 0.05em; }
    .badge { background: rgba(212,175,55,0.15); color: #f5e6ad; padding: 4px 12px; border-radius: 6px; font-size: 11px; font-weight: 700; text-transform: uppercase; letter-spacing: 0.08em; border: 1px solid rgba(212,175,55,0.3); vertical-align: middle; }
    .greeting { font-size: 16px; color: #e2e8f0; margin-top: 20px; line-height: 1.6; }
    .field-box { background: rgba(255,255,255,0.03); border: 1px solid rgba(212,175,55,0.25); border-left: 4px solid #d4af37; padding: 18px 20px; margin: 24px 0; border-radius: 8px; font-size: 14px; line-height: 1.8; color: #cbd5e1; }
    .field-box strong { color: #f5e6ad; display: inline-block; min-width: 140px; }
    .field-box .value { color: #ffffff; }
    .cta-note { font-size: 14px; color: #94a3b8; line-height: 1.6; margin: 20px 0; }
    .footer { font-size: 12px; color: #64748b; border-top: 1px solid rgba(255,255,255,0.08); padding-top: 20px; margin-top: 30px; text-align: center; line-height: 1.6; }
    .footer a { color: #d4af37; text-decoration: none; }
  </style>
</head>
<body>
  <div class="card">
    <div class="header">
      HIRELLY <span class="badge">Executive Advisory</span>
    </div>
    <div class="greeting">
      <p>Hello <strong>%s</strong>,</p>
      <p>Thank you for initiating a confidential executive search mandate with <strong>Hirelly Executive Advisory</strong>.</p>
    </div>
    
    <div class="field-box">
      <div><strong>Mandate Type:</strong> <span class="value">%s</span></div>
      <div><strong>Client / Lead:</strong> <span class="value">%s</span></div>
      <div><strong>Work Email:</strong> <span class="value">%s</span></div>
      <div><strong>Organization:</strong> <span class="value">%s</span></div>
      <div><strong>Details:</strong> <span class="value">%s</span></div>
    </div>

    <p class="cta-note">
      Our industry practice leads across Bangalore, London, Dubai, and Singapore have been notified. A senior partner will review your requirements under strict non-disclosure protocol and connect with you directly within 24 business hours.
    </p>
    
    <div class="footer">
      © 2026 Hirelly International Private Limited — Global Executive Search Network<br>
      Direct Advisory Desk: <a href="mailto:connect@hirelly.in">connect@hirelly.in</a> · <a href="https://hirelly.in">hirelly.in</a>
    </div>
  </div>
</body>
</html>`, m.Name, m.Type, m.Name, m.Email, m.Org, m.Message)

	recipients := []string{m.Email}
	if s.adminEmail != "" && s.adminEmail != m.Email {
		recipients = append(recipients, s.adminEmail)
	}

	msg := notifications.Message{
		ID:         m.ID,
		Title:      fmt.Sprintf("⚡ Hirelly Executive Inquiry: %s (%s)", m.Type, m.Name),
		Body:       fmt.Sprintf("New executive mandate from %s (%s) for %s", m.Name, m.Email, m.Org),
		HTMLBody:   htmlContent,
		Priority:   notifications.PriorityHigh,
		Channels:   []notifications.Channel{notifications.ChannelEmail},
		Recipients: recipients,
		CreatedAt:  time.Now().UTC(),
	}

	if err := s.sender.Send(ctx, msg); err != nil {
		return "", fmt.Errorf("go-app-kit brevo dispatch error: %w", err)
	}

	return "dispatched_via_go_app_kit", nil
}
