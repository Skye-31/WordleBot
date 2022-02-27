package components

import (
	"context"
	"strings"

	"github.com/DisgoOrg/disgo/core/events"
	"github.com/DisgoOrg/disgo/discord"
	"github.com/DisgoOrg/disgo/json"
	"github.com/DisgoOrg/log"
	"github.com/Skye-31/WordleBot/models"
	"github.com/Skye-31/WordleBot/types"
	"github.com/uptrace/bun"
)

func Guess(event *events.ComponentInteractionEvent) {
	_ = event.CreateModal(discord.NewModalCreateBuilder().
		SetTitle("Enter your new guess").
		SetCustomID("game:guess:submit").
		AddActionRow(discord.TextInputComponent{
			CustomID:  "guess",
			Style:     discord.TextInputStyleShort,
			Label:     "Your guess",
			MinLength: *json.NewInt(4),
			MaxLength: 8,
			Required:  true,
		},
		).Build())
}

func GuessSubmit(db *bun.DB, wd *types.WordsData, event *events.ModalSubmitInteractionEvent) {
	c := event.Data.Components[0].Components()[0]
	textInput, ok := c.(discord.TextInputComponent)
	if !ok {
		_ = event.CreateMessage(discord.MessageCreate{Content: "Failed to process your interaction", Flags: 64})
		return
	}
	guess := strings.ToLower(textInput.Value)

	game := models.Game{
		ID: event.User.ID,
	}
	if err := db.NewSelect().Model(&game).WherePK().Scan(context.TODO(), &game); err != nil {
		_ = event.CreateMessage(discord.MessageCreate{Content: "Error fetching game information from database", Flags: 64})
		return
	}
	if len(game.Word) != len(guess) {
		_ = event.CreateMessage(discord.MessageCreate{Content: "Your guess must be the same length as the word", Flags: 64})
		return
	}
	words := wd.GetByLength(len(guess))
	if !words.Has(guess) {
		_ = event.CreateMessage(discord.MessageCreate{Content: "Your word is not in the word list.", Flags: 64})
		return
	}
	game.Guesses = append(game.Guesses, guess)

	r := game.Render(&event.CreateInteraction)
	b, err := game.RenderImage(true)
	if err != nil {
		_ = event.CreateMessage(discord.NewMessageCreateBuilder().SetContent("Error rendering image").SetFlags(64).Build())
		return
	}
	attachment := discord.NewFile("word-"+event.ID.String()+".png", b)
	if err := event.UpdateMessage(discord.NewMessageUpdateBuilder().
		SetEmbeds(r.Embeds...).
		SetContainerComponents(r.Components...).
		RetainAttachments().
		SetFiles(attachment).
		SetFlags(r.Flags).
		Build()); err != nil {
		log.Errorf("Error creating message: %s", err)
	}
	if game.IsOver() {
		if _, err = db.NewDelete().Model(&game).WherePK().Exec(context.TODO()); err != nil {
			log.Errorf("Error deleting game: %s", err)
		}
	} else {
		if _, err := db.NewUpdate().Model(&game).WherePK().Column("guesses").Exec(context.TODO()); err != nil {
			event.Bot().Logger.Errorf("Error updating game: %s", err)
		}
	}
}
