package app

import (
	"context"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	_ "github.com/lib/pq"
	"help_center_telegram_bot/internal/boundary"
	"help_center_telegram_bot/internal/config"
	"help_center_telegram_bot/internal/gate/postgres"
	"help_center_telegram_bot/internal/gate/redis"
	"help_center_telegram_bot/pkg/logger"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
)

func New() {
	ctx := context.Background()

	loadConfig, err := config.LoadConfig()
	if err != nil {
		fmt.Printf("Error loading loadConfig: %v \n", err)
	}

	log := logger.SetupLogger(loadConfig.Env)

	bot, err := tgbotapi.NewBotAPI(loadConfig.Conf.TelegramToken)
	if err != nil {
		log.Error("Error connecting to Telegram", slog.String("error", err.Error()))
	}

	log.Info("Authorized on account", slog.String("username", bot.Self.UserName))

	postgresHandler, err := postgres.Init(*log, loadConfig.Conf.PostgresUrl)
	if err != nil {
		log.Error("Error connecting to PostgreSQL", slog.String("error", err.Error()))
	}

	redisHandler, client, err := redis.Init(ctx, log, loadConfig.Conf.RedisUrl)
	if err != nil {
		log.Error("Error connecting to Redis", slog.String("error", err.Error()))
	}

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := bot.GetUpdatesChan(u)

	if err != nil {
		log.Error("Error getting updates channel", slog.String("error", err.Error()))
	}

	go func() {
		for update := range updates {
			gate := boundary.Gate{
				Log:      log,
				Bot:      bot,
				Update:   update,
				Redis:    redisHandler,
				Postgres: postgresHandler,
			}
			if update.Message != nil {
				switch update.Message.Text {
				case "/start":
					gate.HandleStart()
				case "/new":
					gate.HandleNewTicket(ctx)
				case "Support", "IT", "Billing":
					gate.HandleDepartmentSelection(ctx)
				case "Submit":
					state := gate.Redis.Get(ctx, gate.Update.Message.Chat.ID, "department")
					switch state {
					case "Support":
						gate.HandleSubmit(ctx, loadConfig.Conf.SupportAdminChatID)
					case "IT":
						gate.HandleSubmit(ctx, loadConfig.Conf.ITAdminChatID)
					case "Billing":
						gate.HandleSubmit(ctx, loadConfig.Conf.BillingAdminChatID)

					}

				default:
					state := gate.Redis.Get(ctx, gate.Update.Message.Chat.ID, "state")
					switch state {
					case "select_title":
						gate.HandleTitleInput(ctx)
					case "select_description":
						gate.HandleDescriptionInput(ctx)
					}
				}
			} else if update.CallbackQuery != nil {
				switch update.CallbackQuery.Data {
				case "fast_reply":
					gate.HandleFastReply(ctx)
				}
			}
		}
	}()

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)
	<-signalChan

	log.Info("Shutting down...")
	bot.StopReceivingUpdates()
	err = postgresHandler.Db.Close()
	if err != nil {
		log.Error("Failed to Close PostgresSQL connection", slog.String("error", err.Error()))
		return
	}
	err = client.Close()
	if err != nil {
		log.Error("Failed to Close Redis connection", slog.String("error", err.Error()))
		return
	}

}
