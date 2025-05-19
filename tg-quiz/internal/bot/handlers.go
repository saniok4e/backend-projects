package bot

import (
	"fmt"
	"strconv"
	"tg-quiz/internal/game"
	"tg-quiz/internal/models"
	"tg-quiz/internal/utils"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// Handler представляет обработчик сообщений бота
type Handler struct {
	bot  *tgbotapi.BotAPI
	game *game.Game
}

// NewHandler создает новый обработчик
func NewHandler(bot *tgbotapi.BotAPI, game *game.Game) *Handler {
	return &Handler{
		bot:  bot,
		game: game,
	}
}

// HandleMessage обрабатывает входящие сообщения
func (h *Handler) HandleMessage(msg *tgbotapi.Message) {
	userID := msg.From.ID
	text := msg.Text

	// Проверяем, является ли сообщение ответом на вопрос с вводом текста
	if h.game.IsActive && h.game.CurrentQIndex >= 0 && h.game.CurrentQIndex < len(h.game.Questions) {
		if q := h.game.Questions[h.game.CurrentQIndex]; q.Type == models.InputQuestion {
			if user, exists := h.game.Users[userID]; exists {
				correctAnswer := q.Options[0]
				errors := utils.LevenshteinDistance(text, correctAnswer)

				var response string
				if errors <= 2 {
					h.game.AddScore(user.Team)
					response = fmt.Sprintf("Верно! 🎉\nПравильный ответ: %s\n%s", correctAnswer, utils.FormatScore(user.Team.Score))
				} else {
					response = fmt.Sprintf("Неверно 😢\nПравильный ответ: %s\n%s", correctAnswer, utils.FormatScore(user.Team.Score))
				}

				h.bot.Send(tgbotapi.NewMessage(user.ID, response))

				// Проверяем, был ли это последний вопрос
				if h.game.CurrentQIndex == len(h.game.Questions)-1 {
					h.game.EndGame(h.bot)
				}
				return
			}
		}
	}

	switch text {
	case "/start":
		if userID == h.game.AdminID {
			h.sendAdminPanel(msg.Chat.ID)
		} else if _, exists := h.game.Users[userID]; !exists {
			h.game.Users[userID] = &models.User{ID: userID}
			h.bot.Send(tgbotapi.NewMessage(msg.Chat.ID, "Привет! Пожалуйста, введи свое имя:"))
		} else {
			h.bot.Send(tgbotapi.NewMessage(msg.Chat.ID, "Ты уже зарегистрирован. Ожидай начала игры."))
		}
	case "/next":
		if userID != h.game.AdminID {
			h.bot.Send(tgbotapi.NewMessage(msg.Chat.ID, "У тебя нет прав для этой команды."))
			return
		}
		if !h.game.IsActive {
			h.game.IsActive = true
		}
		h.game.CurrentQIndex++
		if h.game.CurrentQIndex >= len(h.game.Questions) {
			h.game.EndGame(h.bot)
			return
		}
		h.game.SendQuestion(h.bot)
	case "/leaders":
		h.game.ShowLeaderboard(h.bot, msg.Chat.ID)
	default:
		if user, exists := h.game.Users[userID]; exists && user.Name == "" {
			user.Name = text
			h.game.SendTeamSelection(h.bot, user.ID)
		}
	}
}

// HandleCallback обрабатывает callback-запросы от кнопок
func (h *Handler) HandleCallback(cq *tgbotapi.CallbackQuery) {
	userID := cq.From.ID

	// Обработка выбора команды
	if user, exists := h.game.Users[userID]; exists && user.Team == nil {
		switch cq.Data {
		case "team_red":
			user.Team = h.game.RedTeam
			h.bot.Send(tgbotapi.NewMessage(user.ID, fmt.Sprintf("Привет, %s! Ты в команде 🔴 Красные. Ожидай начала игры.", user.Name)))
			return
		case "team_blue":
			user.Team = h.game.BlueTeam
			h.bot.Send(tgbotapi.NewMessage(user.ID, fmt.Sprintf("Привет, %s! Ты в команде 🔵 Синие. Ожидай начала игры.", user.Name)))
			return
		}
	}

	// Обработка админских кнопок
	if userID == h.game.AdminID {
		switch cq.Data {
		case "next_question":
			if !h.game.IsActive {
				h.game.IsActive = true
			}
			h.game.CurrentQIndex++
			if h.game.CurrentQIndex >= len(h.game.Questions) {
				h.game.EndGame(h.bot)
				return
			}
			h.game.SendQuestion(h.bot)
			h.sendAdminPanel(cq.Message.Chat.ID)
			return

		case "show_leaders":
			h.game.ShowLeaderboard(h.bot, cq.Message.Chat.ID)
			return

		case "confirm_reset":
			h.sendResetConfirmation(cq.Message.Chat.ID)
			return

		case "cancel_reset":
			h.bot.Send(tgbotapi.NewMessage(cq.Message.Chat.ID, "Сброс игры отменен"))
			h.sendAdminPanel(cq.Message.Chat.ID)
			return

		case "reset_game":
			h.game.ResetGame()
			h.bot.Send(tgbotapi.NewMessage(cq.Message.Chat.ID, "✅ Игра сброшена. Все очки обнулены."))
			h.sendAdminPanel(cq.Message.Chat.ID)
			return

		case "show_players":
			var playersList string
			for _, user := range h.game.Users {
				playersList += fmt.Sprintf("👤 %s - %s\n", user.Name, user.Team.Name)
			}
			if playersList == "" {
				playersList = "Пока нет участников"
			}
			h.bot.Send(tgbotapi.NewMessage(cq.Message.Chat.ID, "👥 Участники:\n\n"+playersList))
			return
		}
	}

	// Обработка ответов на вопросы
	user, exists := h.game.Users[userID]
	if !exists || user.Team == nil {
		return
	}

	selectedIdx, err := strconv.Atoi(cq.Data)
	if err != nil {
		return
	}

	correctIdx := h.game.Questions[h.game.CurrentQIndex].CorrectIdx
	var response string
	if selectedIdx == correctIdx {
		h.game.AddScore(user.Team)
		response = "Верно! 🎉\n" + utils.FormatScore(user.Team.Score)
	} else {
		response = "Неверно 😢\n" + utils.FormatScore(user.Team.Score)
	}

	h.bot.Send(tgbotapi.NewMessage(user.ID, response))

	// Удаляем клавиатуру
	edit := tgbotapi.NewEditMessageReplyMarkup(
		cq.Message.Chat.ID,
		cq.Message.MessageID,
		tgbotapi.InlineKeyboardMarkup{InlineKeyboard: [][]tgbotapi.InlineKeyboardButton{}},
	)
	if _, err := h.bot.Request(edit); err != nil {
		fmt.Printf("Ошибка при удалении клавиатуры: %v\n", err)
	}

	// Проверяем, был ли это последний вопрос
	if h.game.CurrentQIndex == len(h.game.Questions)-1 {
		h.game.EndGame(h.bot)
	}
}

// sendAdminPanel отправляет панель управления админу
func (h *Handler) sendAdminPanel(chatID int64) {
	msg := tgbotapi.NewMessage(chatID, "🎮 Панель управления игрой")

	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("▶️ Следующий вопрос", "next_question"),
			tgbotapi.NewInlineKeyboardButtonData("📊 Таблица лидеров", "show_leaders"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("🔄 Сбросить игру", "confirm_reset"),
			tgbotapi.NewInlineKeyboardButtonData("👥 Участники", "show_players"),
		),
	)

	msg.ReplyMarkup = keyboard
	h.bot.Send(msg)
}

// sendResetConfirmation отправляет подтверждение сброса игры
func (h *Handler) sendResetConfirmation(chatID int64) {
	msg := tgbotapi.NewMessage(chatID, "⚠️ Вы уверены, что хотите сбросить игру?\nВсе очки будут обнулены!")

	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("✅ Да, сбросить", "reset_game"),
			tgbotapi.NewInlineKeyboardButtonData("❌ Нет, отмена", "cancel_reset"),
		),
	)

	msg.ReplyMarkup = keyboard
	h.bot.Send(msg)
}
