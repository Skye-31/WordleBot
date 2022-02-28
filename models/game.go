package models

import (
	"bytes"
	_ "embed"
	"fmt"
	"strings"

	"github.com/DisgoOrg/disgo/core"
	"github.com/DisgoOrg/disgo/discord"
	"github.com/DisgoOrg/snowflake"
	"github.com/golang/freetype/truetype"
	"golang.org/x/image/font"

	"github.com/fogleman/gg"
)

//go:embed static/arial.ttf
var fontBytes []byte

type Game struct {
	ID         snowflake.Snowflake `bun:"id,pk,nullzero"`
	Word       string              `bun:"word"`
	Guesses    Guesses             `bun:"guesses,array"`
	HasGivenUp bool                `bun:"-"`
}

func (g Game) MaxGuesses() int {
	return len(g.Word) + 1
}

func (g Game) IsCorrect() bool {
	last := g.Guesses[len(g.Guesses)-1]
	return last == g.Word
}

func (g Game) IsOver() bool {
	if g.IsCorrect() || g.HasGivenUp {
		return true
	}
	return len(g.Guesses) >= g.MaxGuesses()
}

// LetterStatus returns []0,1,2 for not guessed, guessed, and guessed in right position
func (g Game) LetterStatus(guess string) []int {
	s := make([]int, len(g.Word))
	letterCounts := make(map[string]int)

	for pos, char := range guess {
		letterCounts[string(char)]++
		if char == rune(g.Word[pos]) {
			s[pos] = 2
		} else {
			if strings.Count(g.Word, string(char)) >= letterCounts[string(char)] {
				s[pos] = 1
			}
		}
	}

	return s
}

type RenderReturnInfo struct {
	Embeds     []discord.Embed
	Components []discord.ContainerComponent
	Flags      discord.MessageFlags
}

func (g Game) Render(event *core.CreateInteraction) RenderReturnInfo {
	r := RenderReturnInfo{
		Embeds: []discord.Embed{
			discord.NewEmbedBuilder().
				SetAuthor(event.User.Tag(), "", event.User.EffectiveAvatarURL(128)).
				SetTitlef("Guess the word").
				SetColor(0x54f27c).
				SetImage("attachment://word-" + event.ID.String() + ".png").
				Build(),
		},
		Components: []discord.ContainerComponent{
			discord.NewActionRow(
				discord.NewPrimaryButton("Guess", "game:guess"),
				discord.NewDangerButton("Give up", "game:giveup"),
			),
		},
		Flags: 0,
	}
	if event.GuildID != nil {
		r.Flags = discord.MessageFlagEphemeral
	}
	if g.IsOver() {
		if g.IsCorrect() {
			r.Embeds[0].Title = "Correct!"
		} else {
			r.Embeds[0].Title = "Nice try!"
			r.Embeds[0].Color = 0xf25454
		}
		r.Embeds[0].Description = "The word was " + g.Word
		r.Components = []discord.ContainerComponent{
			discord.NewActionRow(discord.NewPrimaryButton("Get sharable image", discord.CustomID(fmt.Sprintf("game:share:%s:%s", strings.Join(g.Guesses, ","), g.Word)))),
		}
	}
	return r
}

func (g Game) RenderImage(shouldDrawLetters bool) (*bytes.Buffer, error) {
	width := len(g.Word)*50 + (len(g.Word)-1)*10
	height := g.MaxGuesses()*50 + (g.MaxGuesses()-1)*10
	dc := gg.NewContext(width, height)
	fontFace, err := loadFontFace(fontBytes, 30)
	if err != nil {
		return nil, err
	}
	dc.SetFontFace(fontFace)

	for i := 0; i < g.MaxGuesses(); i++ {
		guessStatus := make([]int, len(g.Word))
		letters := make([]string, len(g.Word))
		if i < len(g.Guesses) {
			guessStatus = g.LetterStatus(g.Guesses[i])
			letters = strings.Split(strings.ToUpper(g.Guesses[i]), "")
		}
		for j := range guessStatus {
			nextColour := getColourFromStatus(guessStatus[j])
			dc.Push()
			dc.SetHexColor(nextColour)
			dc.DrawRoundedRectangle(float64(j*60), float64(i*60), 50, 50, 5)
			dc.Fill()
			dc.Pop()
		}
		if shouldDrawLetters {
			for j := range letters {
				if letters[j] != "" {
					dc.Push()
					dc.SetHexColor("#fff")
					dc.DrawStringAnchored(letters[j], float64(j*60+25), float64(i*60+23), 0.5, 0.5)
					dc.Pop()
				}
			}
		}
	}
	var b bytes.Buffer
	err = dc.EncodePNG(&b)
	return &b, err
}

type Guesses []string

func getColourFromStatus(status int) string {
	switch status {
	case 1:
		return "#b59f3b"
	case 2:
		return "#538d4e"
	default:
		return "#3a3a3a"
	}
}

func loadFontFace(fontBytes []byte, points float64) (font.Face, error) {
	f, err := truetype.Parse(fontBytes)
	if err != nil {
		return nil, err
	}
	face := truetype.NewFace(f, &truetype.Options{
		Size: points,
		// Hinting: font.HintingFull,
	})
	return face, nil
}
