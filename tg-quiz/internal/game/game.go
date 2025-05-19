package game

import (
	"fmt"
	"sort"
	"strconv"
	"tg-quiz/internal/models"
	"tg-quiz/internal/utils"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// Game представляет состояние игры
type Game struct {
	Users         map[int64]*models.User
	Questions     []models.Question
	CurrentQIndex int
	IsActive      bool
	AdminID       int64
	RedTeam       *models.Team
	BlueTeam      *models.Team
}

// NewGame создает новый экземпляр игры
func NewGame(adminID int64) *Game {
	return &Game{
		Users:         make(map[int64]*models.User),
		CurrentQIndex: -1,
		IsActive:      false,
		AdminID:       adminID,
		RedTeam: &models.Team{
			Type: models.RedTeam,
			Name: "🔴 Красные",
		},
		BlueTeam: &models.Team{
			Type: models.BlueTeam,
			Name: "🔵 Синие",
		},
		Questions: []models.Question{
			{
				Type:       models.TextQuestion,
				Text:       "Столица России?",
				Options:    []string{"Москва", "Санкт-Петербург", "Казань", "Новосибирск"},
				CorrectIdx: 0,
			},
			{
				Type:       models.PhotoQuestion,
				ContentURL: "./storage/1424-1000x830.jpg",
				Text:       "Что изображено на фото?",
				Options:    []string{"Доминик Торетто", "Собака", "Хомяк", "Кролик"},
				CorrectIdx: 0,
			},
			{
				Type:       models.AudioQuestion,
				ContentURL: "./storage/music/audio_7.m4a",
				Options:    []string{"Шарик", "Иван", "Смешарик"},
				CorrectIdx: 0,
			},
			{
				Type:       models.InputQuestion,
				Text:       "Введите столицу Франции:",
				Options:    []string{"Париж"},
				CorrectIdx: 0,
			},
		},
	}
}

// SendTeamSelection отправляет сообщение с выбором команды
func (g *Game) SendTeamSelection(bot *tgbotapi.BotAPI, userID int64) {
	msg := tgbotapi.NewMessage(userID, "Выберите команду:")

	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("🔴 Красные", "team_red"),
			tgbotapi.NewInlineKeyboardButtonData("🔵 Синие", "team_blue"),
		),
	)

	msg.ReplyMarkup = keyboard
	bot.Send(msg)
}

// SendQuestion отправляет текущий вопрос всем пользователям
func (g *Game) SendQuestion(bot *tgbotapi.BotAPI) {
	q := g.Questions[g.CurrentQIndex]

	for _, user := range g.Users {
		var msg tgbotapi.Chattable

		switch q.Type {
		case models.TextQuestion:
			buttons := make([][]tgbotapi.InlineKeyboardButton, len(q.Options))
			for i, opt := range q.Options {
				buttons[i] = tgbotapi.NewInlineKeyboardRow(
					tgbotapi.NewInlineKeyboardButtonData(opt, strconv.Itoa(i)),
				)
			}
			textMsg := tgbotapi.NewMessage(user.ID, q.Text)
			textMsg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(buttons...)
			msg = textMsg

		case models.PhotoQuestion:
			photo := tgbotapi.NewPhoto(user.ID, tgbotapi.FilePath(q.ContentURL))
			photo.Caption = q.Text
			buttons := make([][]tgbotapi.InlineKeyboardButton, len(q.Options))
			for i, opt := range q.Options {
				buttons[i] = tgbotapi.NewInlineKeyboardRow(
					tgbotapi.NewInlineKeyboardButtonData(opt, strconv.Itoa(i)),
				)
			}
			photo.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(buttons...)
			msg = photo

		case models.AudioQuestion:
			audio := tgbotapi.NewAudio(user.ID, tgbotapi.FilePath(q.ContentURL))
			audio.Caption = q.Text
			buttons := make([][]tgbotapi.InlineKeyboardButton, len(q.Options))
			for i, opt := range q.Options {
				buttons[i] = tgbotapi.NewInlineKeyboardRow(
					tgbotapi.NewInlineKeyboardButtonData(opt, strconv.Itoa(i)),
				)
			}
			audio.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(buttons...)
			msg = audio

		case models.InputQuestion:
			textMsg := tgbotapi.NewMessage(user.ID, q.Text+"\n\nВведите ответ в чат:")
			msg = textMsg
		}

		if _, err := bot.Send(msg); err != nil {
			fmt.Printf("Ошибка при отправке вопроса: %v\n", err)
		}
	}
}

// ShowLeaderboard показывает таблицу лидеров
func (g *Game) ShowLeaderboard(bot *tgbotapi.BotAPI, chatID int64) {
	// Сортируем команды по счету
	teams := []*models.Team{g.RedTeam, g.BlueTeam}
	sort.Slice(teams, func(i, j int) bool {
		return teams[i].Score > teams[j].Score
	})

	var leaderboard string
	for i, team := range teams {
		medal := "🥉"
		if i == 0 {
			medal = "🥇"
		} else if i == 1 {
			medal = "🥈"
		}
		leaderboard += fmt.Sprintf("%s %s - %d %s\n",
			medal,
			team.Name,
			team.Score,
			utils.FormatScoreWord(team.Score))

		// Добавляем список игроков команды
		for _, user := range g.Users {
			if user.Team == team {
				leaderboard += fmt.Sprintf("   👤 %s\n", user.Name)
			}
		}
		leaderboard += "\n"
	}

	if leaderboard == "" {
		leaderboard = "Пока нет участников"
	}

	msg := tgbotapi.NewMessage(chatID, "🏆 Таблица лидеров:\n\n"+leaderboard)
	bot.Send(msg)
}

// EndGame завершает игру и показывает финальные результаты
func (g *Game) EndGame(bot *tgbotapi.BotAPI) {
	// Сортируем команды по счету
	teams := []*models.Team{g.RedTeam, g.BlueTeam}
	sort.Slice(teams, func(i, j int) bool {
		return teams[i].Score > teams[j].Score
	})

	var leaderboard string
	for i, team := range teams {
		medal := "🥉"
		if i == 0 {
			medal = "🥇"
		} else if i == 1 {
			medal = "🥈"
		}
		leaderboard += fmt.Sprintf("%s %s - %d %s\n",
			medal,
			team.Name,
			team.Score,
			utils.FormatScoreWord(team.Score))

		// Добавляем список игроков команды
		for _, user := range g.Users {
			if user.Team == team {
				leaderboard += fmt.Sprintf("   👤 %s\n", user.Name)
			}
		}
		leaderboard += "\n"
	}

	for _, user := range g.Users {
		msg := tgbotapi.NewMessage(user.ID, "🎮 Игра окончена!\n\n🏆 Финальная таблица лидеров:\n\n"+leaderboard)
		bot.Send(msg)
	}

	adminMsg := tgbotapi.NewMessage(g.AdminID, "🎮 Игра окончена!\n\n🏆 Финальная таблица лидеров:\n\n"+leaderboard)
	bot.Send(adminMsg)

	g.IsActive = false
	g.CurrentQIndex = -1
}

// ResetGame сбрасывает состояние игры
func (g *Game) ResetGame() {
	g.CurrentQIndex = -1
	g.IsActive = false
	g.RedTeam.Score = 0
	g.BlueTeam.Score = 0
}

// AddScore добавляет очки команде
func (g *Game) AddScore(team *models.Team) {
	team.Score++
}
