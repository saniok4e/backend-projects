package models

// QuestionType определяет тип вопроса
type QuestionType string

const (
	TextQuestion  QuestionType = "text"
	PhotoQuestion QuestionType = "photo"
	AudioQuestion QuestionType = "audio"
	InputQuestion QuestionType = "input"
)

// TeamType определяет тип команды
type TeamType string

const (
	RedTeam  TeamType = "red"
	BlueTeam TeamType = "blue"
)

// Team представляет команду
type Team struct {
	Type  TeamType
	Name  string
	Score int
}

// Question представляет структуру вопроса
type Question struct {
	Type       QuestionType
	Text       string
	ContentURL string
	Options    []string
	CorrectIdx int
}

// User представляет структуру пользователя
type User struct {
	ID   int64
	Name string
	Team *Team
}
