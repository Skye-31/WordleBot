package types

import (
	"encoding/json"
	"errors"
	"os"

	"github.com/DisgoOrg/log"
)

type WordsData struct {
	Four  []string `json:"4"`
	Five  []string `json:"5"`
	Six   []string `json:"6"`
	Seven []string `json:"7"`
	Eight []string `json:"8"`
}

func LoadWordsData(log log.Logger) (*WordsData, error) {
	log.Info("Loading words data...")
	file, err := os.Open("data.json")
	if os.IsNotExist(err) {
		return nil, errors.New("data.json not found")
	} else if err != nil {
		return nil, err
	}

	var words WordsData
	if err = json.NewDecoder(file).Decode(&words); err != nil {
		return nil, err
	}
	four, five, six, seven, eight := len(words.Four), len(words.Five), len(words.Six), len(words.Seven), len(words.Eight)
	log.Infof("Loaded %d words: (4: %d, 5: %d, 6: %d, 7: %d, 8: %d)\n", four+five+six+seven+eight, four, five, six, seven, eight)
	return &words, nil
}
