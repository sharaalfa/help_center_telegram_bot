package postgres

import (
	"context"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"

	"help_center_telegram_bot/pkg/models"
)

func TestCreateTicket(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	handler := &Handler{Db: db}

	ticket := models.Ticket{
		Department:  "Support",
		Title:       "Test Ticket",
		Description: "This is a test ticket",
		ClientID:    1,
	}

	mock.ExpectExec("INSERT INTO tickets").
		WithArgs(ticket.Department, ticket.Title, ticket.Description, ticket.ClientID).
		WillReturnResult(sqlmock.NewResult(1, 1)) // Возвращаем результат с ID 1 и 1 строкой затронутой

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	err = handler.CreateTicket(ctx, ticket)

	assert.NoError(t, err)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}
