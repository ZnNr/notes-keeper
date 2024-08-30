package spellcheck

import (
	"crypto/tls"
	"encoding/json"
	mistakes "github.com/ZnNr/notes-keeper.git/intenal/notes/notemodel"
	"io"

	"log"
	"log/slog"
	"net/http"

	"strings"
)

const op = "spellcheck.YaSpellCheck.CheckText"
const apiURL = "https://speller.yandex.net/services/spellservice.json/checkText?text="

type YaSpellCheck struct {
	logger *slog.Logger
}

func NewYaSpellChecker(logger *slog.Logger) *YaSpellCheck {
	return &YaSpellCheck{
		logger: logger,
	}
}

type InputSpell struct {
	Code int      `json:"code"`
	Pos  int      `json:"pos"`
	Row  int      `json:"row"`
	Col  int      `json:"col"`
	Len  int      `json:"len"`
	Word string   `json:"word"`
	S    []string `json:"s"`
}

func (y *YaSpellCheck) CheckText(text string) ([]byte, error) {

	y.logger = y.logger.With("op", op)
	var inputSpell []InputSpell

	modifiedText := strings.ReplaceAll(text, " ", "+")

	transport := &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
	}
	client := &http.Client{
		Transport: transport,
	}

	resp, err := client.Get(apiURL + modifiedText)
	if err != nil {
		y.logger.Error("cannot get response", slog.String("url", apiURL+modifiedText))
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		y.logger.Error("cannot read body", slog.String("url", apiURL+modifiedText))
		return nil, err
	}
	err = json.Unmarshal(body, &inputSpell)
	if err != nil {
		y.logger.Error("cannot unmarshal body", slog.String("url", apiURL+modifiedText))
		log.Fatalf("Error unmarshalling input JSON: %v", err)
	}
	var outputSpell []mistakes.Mistakes
	for _, item := range inputSpell {
		outputSpell = append(outputSpell, mistakes.Mistakes{
			OriginalWord: item.Word,
			CorrectWord:  item.S,
		})
	}

	mistakes, err := json.Marshal(outputSpell)
	if err != nil {
		y.logger.Error("cannot marshal body", slog.String("url", apiURL+modifiedText))
		log.Fatalf("Error marshalling output JSON: %v", err)
	}
	return mistakes, nil
}
