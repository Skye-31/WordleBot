package commands

import (
	"context"
	"regexp"
	"strconv"
	"strings"

	"github.com/DisgoOrg/disgo/core/events"
	"github.com/DisgoOrg/disgo/discord"
	"github.com/Skye-31/WordleBot/models"
	"github.com/Skye-31/WordleBot/types"
	"github.com/uptrace/bun"
)

var regex = regexp.MustCompile(`^[a-z]{4,8}$`)

func start(db *bun.DB, wd *types.WordsData, event *events.ApplicationCommandInteractionEvent) {
	count, err := db.NewSelect().Model((*models.Game)(nil)).Where("id = ?", event.User.ID).Count(context.TODO())
	if err != nil {
		_ = event.CreateMessage(discord.MessageCreate{Content: "Error fetching game information from database", Flags: 64})
		return
	}
	if count != 0 {
		_ = event.CreateMessage(discord.MessageCreate{
			Content: "You are already in a game",
			Flags:   64,
			Components: []discord.ContainerComponent{
				discord.NewActionRow(
					discord.NewPrimaryButton("Continue game", "start:continue").WithEmoji(discord.ComponentEmoji{Name: "➡"}),
					discord.NewDangerButton("Give up", "start:giveup").WithEmoji(discord.ComponentEmoji{Name: "🗑"}),
				)},
		})
		return
	}
	data := event.SlashCommandInteractionData()
	guess := strings.ToLower(*data.Options.String("starting-word"))
	if !regex.Match([]byte(guess)) {
		_ = event.CreateMessage(discord.MessageCreate{Content: "Invalid starting word, please use a 4-8 character word", Flags: 64})
		return
	}
	u := models.User{
		ID:        event.User.ID,
		CachedTag: event.User.Tag(),
	}
	if err := db.NewSelect().Model(&u).WherePK().Scan(context.TODO()); err != nil {
		u.DefaultWordLength = 5
	}
	wordLength := int(u.DefaultWordLength)
	if data.Options.IntOption("letters") != nil {
		wordLength = *data.Options.Int("letters")
	}
	if wordLength != len(guess) {
		_ = event.CreateMessage(discord.MessageCreate{Content: "Your starting word does not match the length chosen (" + strconv.Itoa(wordLength) + ").", Flags: 64})
		return
	}
	words := wd.GetByLength(wordLength)
	if !words.Has(guess) {
		_ = event.CreateMessage(discord.MessageCreate{Content: "Your starting word is not in the word list.", Flags: 64})
		return
	}
	word := words.GetRandom()
	_ = event.CreateMessage(discord.MessageCreate{Content: word})
}
