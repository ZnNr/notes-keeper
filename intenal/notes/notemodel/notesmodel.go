package notemodel

// Note Структура представляет собой модель заметки
type Note struct {
	Text     string
	Mistakes []Mistakes
	Id       int
	UserId   int
}

type Mistakes struct {
	OriginalWord string   `json:"original_word"`
	CorrectWord  []string `json:"correct_word"`
}
