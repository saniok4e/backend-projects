package models

type User struct {
	ID    int64
	Name  string
	Score int
}

type QuestionType string

const (
	TextQuestion  QuestionType = "text"
	PhotoQuestion QuestionType = "photo"
	AudioQuestion QuestionType = "audio"
	InputQuestion QuestionType = "input"
)

type Question struct {
	Type       QuestionType
	ContentURL string   // Фото или аудио URL (пусто для текстовых)
	Text       string   // Сам вопрос / подпись
	Options    []string // Варианты ответа (для input это правильный ответ)
	CorrectIdx int      // Индекс правильного ответа
}
