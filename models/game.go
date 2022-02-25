package models

import "github.com/DisgoOrg/snowflake"

type Game struct {
	ID      snowflake.Snowflake `bun:"id,pk,nullzero"`
	Word    string              `bun:"word"`
	Guesses []Guess             `bun:"guesses,array"`
}

type Guess string

func (g Game) MaxGuesses() int {
	return len(g.Word) + 1
}

func (g Game) IsOver() bool {
	return len(g.Guesses) >= g.MaxGuesses()
}
