package boundary

import (
	"context"
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"help_center_telegram_bot/internal/gate/postgres"
	"help_center_telegram_bot/internal/gate/redis"
	"help_center_telegram_bot/pkg/models"
	"log"
	"log/slog"
	"strconv"
)

type Gate struct {
	Log      *slog.Logger
	Bot      *tgbotapi.BotAPI
	Update   tgbotapi.Update
	Redis    redis.Handler
	Postgres postgres.Handler
}

func New(log *slog.Logger,
	bot *tgbotapi.BotAPI,
	update tgbotapi.Update,
	redis redis.Handler,
	postgres postgres.Handler) *Gate {
	return &Gate{
		Log:      log,
		Bot:      bot,
		Update:   update,
		Redis:    redis,
		Postgres: postgres,
	}
}

func (g *Gate) HandleStart() {
	msg := tgbotapi.NewMessage(g.Update.Message.Chat.ID, "Привет! Я бот службы поддержки. Используй /new для создания нового тикета.")
	_, err := g.Bot.Send(msg)
	if err != nil {
		g.Log.Error("Failed to send message:", msg, err)
		return
	}
}

func (g *Gate) HandleNewTicket(ctx context.Context) {
	msg := tgbotapi.NewMessage(g.Update.Message.Chat.ID, "Выберите подразделение:")
	msg.ReplyMarkup = tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("Support"),
			tgbotapi.NewKeyboardButton("IT"),
			tgbotapi.NewKeyboardButton("Billing"),
		),
	)
	_, err := g.Bot.Send(msg)
	if err != nil {
		g.Log.Error("Failed to send message:", msg, err)
		return
	}

	g.Log.Info("Сохраняем состояние в Redis")
	g.Redis.Set(ctx, g.Update.Message.Chat.ID, "state", "select_department")
}

func (g *Gate) HandleDepartmentSelection(ctx context.Context) {
	department := g.Update.Message.Text
	msg := tgbotapi.NewMessage(g.Update.Message.Chat.ID, "Введите заголовок (тему) тикета:")
	_, err := g.Bot.Send(msg)
	if err != nil {
		g.Log.Error("Failed to send message:", msg, err)
		return
	}

	g.Log.Info("Сохраняем состояние в Redis")
	g.Redis.Set(ctx, g.Update.Message.Chat.ID, "department", department)
	g.Redis.Set(ctx, g.Update.Message.Chat.ID, "state", "select_title")
}

func (g *Gate) HandleTitleInput(ctx context.Context) {
	title := g.Update.Message.Text
	msg := tgbotapi.NewMessage(g.Update.Message.Chat.ID, "Введите текст обращения:")
	_, err := g.Bot.Send(msg)
	if err != nil {
		g.Log.Error("Failed to send message:", msg, err)
		return
	}

	g.Log.Info("Сохраняем состояние в Redis")
	g.Redis.Set(ctx, g.Update.Message.Chat.ID, "title", title)
	g.Redis.Set(ctx, g.Update.Message.Chat.ID, "state", "select_description")
}

func (g *Gate) HandleDescriptionInput(ctx context.Context) {
	description := g.Update.Message.Text
	department := g.Redis.Get(ctx, g.Update.Message.Chat.ID, "department")
	title := g.Redis.Get(ctx, g.Update.Message.Chat.ID, "title")

	msg := tgbotapi.NewMessage(g.Update.Message.Chat.ID, "Результат:")
	msg.Text = "Подразделение: " + department + "\nЗаголовок: " + title + "\nТекст обращения: " + description
	msg.ReplyMarkup = tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("Submit"),
		),
	)
	_, err := g.Bot.Send(msg)
	if err != nil {
		g.Log.Error("Failed to send message:", msg, err)
		return
	}

	g.Redis.Set(ctx, g.Update.Message.Chat.ID, "description", description)
	g.Redis.Set(ctx, g.Update.Message.Chat.ID, "state", "submit")
}

func (g *Gate) HandleSubmit(ctx context.Context, adminChatId int64) {
	department := g.Redis.Get(ctx, g.Update.Message.Chat.ID, "department")
	title := g.Redis.Get(ctx, g.Update.Message.Chat.ID, "title")
	description := g.Redis.Get(ctx, g.Update.Message.Chat.ID, "description")

	g.Log.Info("Создаем тикет в PostgreSQL")
	ticket := models.Ticket{
		Department:  department,
		Title:       title,
		Description: description,
		ClientID:    g.Update.Message.Chat.ID,
	}

	err := g.Postgres.CreateTicket(ctx, ticket)
	if err != nil {
		log.Println("Error creating ticket:", err)
		msg := tgbotapi.NewMessage(g.Update.Message.Chat.ID, "Ошибка при создании тикета. Попробуйте позже.")
		_, err := g.Bot.Send(msg)
		if err != nil {
			g.Log.Error("Failed to send message:", msg, err)
			return
		}
		return
	}

	msg := tgbotapi.NewMessage(g.Update.Message.Chat.ID, "Тикет успешно создан!")
	_, err = g.Bot.Send(msg)
	if err != nil {
		g.Log.Error("Failed to send message:", msg, err)
		return
	}

	g.Log.Info("Отправляем сообщение админу")
	adminMsg := tgbotapi.NewMessage(adminChatId, "Новый тикет:\nПодразделение: "+department+"\nЗаголовок: "+title+"\nТекст обращения: "+description)
	adminMsg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Fast Reply", "fast_reply"),
		),
	)
	_, err = g.Bot.Send(adminMsg)
	if err != nil {
		g.Log.Error("Failed to send message:", adminMsg, err)
		return
	}
	g.Redis.Set(ctx, adminChatId, "client_chat_id", strconv.FormatInt(g.Update.Message.Chat.ID, 10))
}

func (g *Gate) HandleFastReply(ctx context.Context) {
	clientChatID := g.Redis.Get(ctx, g.Update.CallbackQuery.Message.Chat.ID, "client_chat_id")

	chatID, err := strconv.ParseInt(clientChatID, 10, 64)
	if err != nil {
		g.Log.Error("Failed to parse client chat ID", slog.String("error", err.Error()))
		return
	}

	g.Log.Info("Отправляем сообщение клиенту")
	msg := tgbotapi.NewMessage(chatID, "Hello world")
	_, err = g.Bot.Send(msg)
	if err != nil {
		g.Log.Error("Failed to send message:", msg, err)
		return
	}
}
