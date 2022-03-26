package components

import (
	"context"
	"strconv"
	"strings"

	"github.com/Skye-31/WordleBot/models"
	"github.com/Skye-31/WordleBot/types"
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
	"github.com/disgoorg/log"
	"github.com/uptrace/bun"
)

func Guess(event *events.ComponentInteractionEvent) {
	minlength, maxlength := 4, 8
	if split := strings.Split(event.Data.CustomID().String(), ":"); len(split) > 2 {
		var err, err2 error
		minlength, err = strconv.Atoi(split[2])
		maxlength, err2 = strconv.Atoi(split[2])
		if err != nil {
			log.Error(err)
			minlength, maxlength = 4, 8
		}
		if err2 != nil {
			log.Error(err2)
			minlength, maxlength = 4, 8
		}
	}
	_ = event.CreateModal(discord.NewModalCreateBuilder().
		SetTitle("Enter your new guess").
		SetCustomID("game:guess:submit").
		AddActionRow(discord.TextInputComponent{
			CustomID:  "guess",
			Style:     discord.TextInputStyleShort,
			Label:     "Your guess",
			MinLength: &minlength,
			MaxLength: maxlength,
			Required:  true,
		},
		).Build())
}

func Continue(db *bun.DB, event *events.ComponentInteractionEvent) {
	game := models.Game{
		ID: event.User().ID,
	}
	if err := db.NewSelect().Model(&game).WherePK().Scan(context.TODO(), &game); err != nil {
		_ = event.CreateMessage(discord.MessageCreate{Content: "Error fetching game information from database", Flags: 64})
		return
	}
	r := game.Render(event.BaseInteraction)
	b, err := game.RenderImage(true)
	if err != nil {
		_ = event.CreateMessage(discord.NewMessageCreateBuilder().SetContent("Error rendering image").SetFlags(64).Build())
		return
	}
	attachment := discord.NewFile("word-"+event.ID().String()+".png", "Wordle Game"+
		"", b)
	if err := event.UpdateMessage(discord.NewMessageUpdateBuilder().
		SetEmbeds(r.Embeds...).
		SetContainerComponents(r.Components...).
		SetContent("").
		RetainAttachments().
		SetFiles(attachment).
		SetFlags(r.Flags).
		Build()); err != nil {
		log.Errorf("Error updating message: %s", err)
	}
}

func GiveUp(db *bun.DB, event *events.ComponentInteractionEvent) {
	game := models.Game{
		ID: event.User().ID,
	}
	if err := db.NewSelect().Model(&game).WherePK().Scan(context.TODO(), &game); err != nil {
		_ = event.CreateMessage(discord.MessageCreate{Content: "Error fetching game information from database", Flags: 64})
		return
	}
	game.HasGivenUp = true
	r := game.Render(event.BaseInteraction)
	b, err := game.RenderImage(true)
	if err != nil {
		_ = event.CreateMessage(discord.NewMessageCreateBuilder().SetContent("Error rendering image").SetFlags(64).Build())
		return
	}
	attachment := discord.NewFile("word-"+event.ID().String()+".png", "Wordle Result", b)
	if err := event.UpdateMessage(discord.NewMessageUpdateBuilder().
		SetEmbeds(r.Embeds...).
		SetContainerComponents(r.Components...).
		SetContent("").
		RetainAttachments().
		SetFiles(attachment).
		SetFlags(r.Flags).
		Build()); err != nil {
		log.Errorf("Error updating message: %s", err)
	}
	if _, err = db.NewDelete().Model(&game).WherePK().Exec(context.TODO()); err != nil {
		log.Errorf("Error deleting game: %s", err)
	}
	uStats, columnToUpdate := getUStats(game, event.BaseInteraction)
	_, err = db.NewInsert().Model(&uStats).On("CONFLICT (id) DO UPDATE").Set(columnToUpdate+" = array_append(user_stats."+columnToUpdate+", ?::SMALLINT)", 0).Exec(context.TODO())
	if err != nil {
		log.Errorf("Error updating user stats: %s", err)
	}
}

func Share(event *events.ComponentInteractionEvent) {
	id := event.Data.CustomID()
	split := strings.Split(id.String(), ":")
	guesses, word := strings.Split(split[2], ","), split[3]
	game := models.Game{
		ID:      event.User().ID,
		Guesses: guesses,
		Word:    word,
	}
	b, err := game.RenderImage(false)
	if err != nil {
		_ = event.CreateMessage(discord.NewMessageCreateBuilder().SetContent("Error rendering image").SetFlags(64).Build())
		return
	}
	attachment := discord.NewFile("word-"+event.ID().String()+".png", "Wordle Share", b)
	if err := event.CreateMessage(discord.NewMessageCreateBuilder().
		SetContent("Your sharable wordle **" + event.User().Tag() + "**!").
		SetFlags(64).
		SetFiles(attachment).
		Build()); err != nil {
		log.Errorf("Error updating message: %s", err)
	}
}

func GuessSubmit(db *bun.DB, wd *types.WordsData, event *events.ModalSubmitInteractionEvent) {
	guess := strings.ToLower(event.Data.Text("guess"))

	game := models.Game{
		ID: event.User().ID,
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

	r := game.Render(event.BaseInteraction)
	b, err := game.RenderImage(true)
	if err != nil {
		_ = event.CreateMessage(discord.NewMessageCreateBuilder().SetContent("Error rendering image").SetFlags(64).Build())
		return
	}
	attachment := discord.NewFile("word-"+event.ID().String()+".png", "Wordle Game", b)
	if err := event.UpdateMessage(discord.NewMessageUpdateBuilder().
		SetEmbeds(r.Embeds...).
		SetContainerComponents(r.Components...).
		RetainAttachments().
		SetFiles(attachment).
		SetFlags(r.Flags).
		Build()); err != nil {
		log.Errorf("Error updating message: %s", err)
	}
	if game.IsOver() {
		if _, err = db.NewDelete().Model(&game).WherePK().Exec(context.TODO()); err != nil {
			log.Errorf("Error deleting game: %s", err)
		}
		uStats, columnToUpdate := getUStats(game, event.BaseInteraction)
		_, err := db.NewInsert().Model(&uStats).On("CONFLICT (id) DO UPDATE").Set(columnToUpdate+" = array_append(user_stats."+columnToUpdate+", ?::SMALLINT)", len(game.Guesses)).Exec(context.TODO())
		if err != nil {
			log.Errorf("Error updating user stats: %s", err)
		}
	} else {
		if _, err := db.NewUpdate().Model(&game).WherePK().Column("guesses").Exec(context.TODO()); err != nil {
			event.Client().Logger().Errorf("Error updating game: %s", err)
		}
	}
}

func getUStats(game models.Game, event discord.BaseInteraction) (models.UserStats, string) {
	ustats := models.UserStats{
		ID: event.User().ID,
	}
	columnToUpdate := ""
	n := []int{len(game.Guesses)}
	if game.HasGivenUp {
		n[0] = 0
	}
	switch len(game.Word) {
	case 4:
		columnToUpdate = "four"
		ustats.Four = n
	case 5:
		columnToUpdate = "five"
		ustats.Five = n
	case 6:
		columnToUpdate = "six"
		ustats.Six = n
	case 7:
		columnToUpdate = "seven"
		ustats.Seven = n
	case 8:
		columnToUpdate = "eight"
		ustats.Eight = n
	}
	return ustats, columnToUpdate
}
