package types

import (
	_ "embed"
	"encoding/json"
	"math/rand"
	"time"

	"github.com/disgoorg/log"
)

//go:embed data.json
var wordBytes []byte

type WordsData struct {
	Four  WordLength `json:"4"`
	Five  WordLength `json:"5"`
	Six   WordLength `json:"6"`
	Seven WordLength `json:"7"`
	Eight WordLength `json:"8"`
}

func (w WordsData) GetByLength(i int) WordLength {
	switch i {
	case 4:
		return w.Four
	case 5:
		return w.Five
	case 6:
		return w.Six
	case 7:
		return w.Seven
	case 8:
		return w.Eight
	default:
		return WordLength{}
	}
}

type WordLength []string

func (w WordLength) Has(s string) bool {
	for _, v := range w {
		if v == s {
			return true
		}
	}
	return false
}

func (w WordLength) GetRandom() string {
	rand.Seed(time.Now().UnixNano())
	return w[rand.Intn(len(w))]
}

func LoadWordsData(log log.Logger) (*WordsData, error) {
	var words WordsData
	if err := json.Unmarshal(wordBytes, &words); err != nil {
		return nil, err
	}
	four, five, six, seven, eight := len(words.Four), len(words.Five), len(words.Six), len(words.Seven), len(words.Eight)
	log.Infof("Loaded %d words: (4: %d, 5: %d, 6: %d, 7: %d, 8: %d)\n", four+five+six+seven+eight, four, five, six, seven, eight)
	return &words, nil
}
