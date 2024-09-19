package postgres

import (
	"context"
	"help_center_telegram_bot/pkg/models"
	"time"
)

func (h *Handler) CreateTicket(ctx context.Context, ticket models.Ticket) error {
	ctxTimeout, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()
	query := `INSERT INTO tickets (department, title, description, client_id) VALUES ($1, $2, $3, $4)`
	_, err := h.Db.ExecContext(ctxTimeout, query, ticket.Department, ticket.Title, ticket.Description, ticket.ClientID)
	return err
}
