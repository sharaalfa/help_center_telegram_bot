package postgres

import (
	"database/sql"
	"errors"
	"github.com/DATA-DOG/go-sqlmock"
	_ "github.com/lib/pq" // Импортируем драйвер PostgreSQL
	"help_center_telegram_bot/pkg/logger"
	"log/slog"
	"testing"
)

var (
	sqlOpenWrapper = sql.Open
	initWrapper    = Init
)

func TestInit(t *testing.T) {

	db, mock, err := sqlmock.New(sqlmock.MonitorPingsOption(true))
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub db connection", err)
	}
	defer db.Close()

	log := logger.SetupLogger("local")

	originalSqlOpen := sqlOpenWrapper
	sqlOpenWrapper = func(driverName, dataSourceName string) (*sql.DB, error) {
		if dataSourceName == "invalid_url" {
			return nil, errors.New("failed to open db connection")
		}
		return db, nil
	}
	defer func() { sqlOpenWrapper = originalSqlOpen }()

	originalInitWrapper := initWrapper
	initWrapper = func(log slog.Logger, url string) (Handler, error) {
		db, err := sqlOpenWrapper("postgres", url)
		if err != nil {
			return Handler{}, err
		}
		if err := db.Ping(); err != nil {
			return Handler{}, err
		}
		return Handler{db}, nil
	}
	defer func() { initWrapper = originalInitWrapper }()

	tests := []struct {
		name            string
		log             slog.Logger
		url             string
		expectedHandler Handler
		expectedErr     error
	}{
		{
			name:            "success to open db",
			log:             *log,
			url:             "postgres://user:password@localhost:5432/dbname",
			expectedHandler: Handler{Db: db},
			expectedErr:     nil,
		},
		{
			name:            "failed to open db",
			log:             *log,
			url:             "invalid_url",
			expectedHandler: Handler{},
			expectedErr:     errors.New("failed to open db connection"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.url != "invalid_url" {
				mock.ExpectPing()
			}
			handler, err := initWrapper(tt.log, tt.url)
			if err != nil {
				if tt.expectedErr == nil || err.Error() != tt.expectedErr.Error() {
					t.Errorf("Init() error=%v, wantErr %v", err, tt.expectedErr)
					return
				}
			} else {
				if tt.expectedErr != nil {
					t.Errorf("Init() expected error=%v, but got none", tt.expectedErr)
					return
				}
			}

			if handler.Db != tt.expectedHandler.Db {
				t.Errorf("expected db to be %v, but got %v", tt.expectedHandler.Db, handler.Db)
			}
		})
	}

	// Проверяем, что все ожидания были выполнены
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}
