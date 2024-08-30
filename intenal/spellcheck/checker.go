package spellcheck

type SpellChecker interface {
	CheckText(text string) ([]byte, error)
}
