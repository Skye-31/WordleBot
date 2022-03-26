package commands

import (
	"context"
	"regexp"
	"strconv"
	"strings"

	"github.com/Skye-31/WordleBot/models"
	"github.com/Skye-31/WordleBot/types"
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
	"github.com/disgoorg/log"
	"github.com/uptrace/bun"
)

var regex = regexp.MustCompile(`^[a-z]{4,8}$`)

func start(db *bun.DB, wd *types.WordsData, event *events.ApplicationCommandInteractionEvent) {
	count, err := db.NewSelect().Model((*models.Game)(nil)).Where("id = ?", event.User().ID).Count(context.TODO())
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
					discord.NewPrimaryButton("Continue game", "game:continue"),
					discord.NewDangerButton("Give up", "game:giveup"),
				)},
		})
		return
	}
	data := event.SlashCommandInteractionData()
	guess := strings.ToLower(data.String("starting-word"))
	if !regex.Match([]byte(guess)) {
		_ = event.CreateMessage(discord.MessageCreate{Content: "Invalid starting word, please use a 4-8 character word", Flags: 64})
		return
	}
	u := models.User{
		ID:        event.User().ID,
		CachedTag: event.User().Tag(),
	}
	if err := db.NewSelect().Model(&u).WherePK().Scan(context.TODO()); err != nil {
		u.DefaultWordLength = 5
	}
	wordLength := int(u.DefaultWordLength)
	if data.Int("letters") != 0 {
		wordLength = data.Int("letters")
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
	game := models.Game{
		ID:      event.User().ID,
		Word:    word,
		Guesses: []string{guess},
	}
	_, err = db.NewInsert().Model(&game).Exec(context.TODO())
	if err != nil {
		_ = event.CreateMessage(discord.MessageCreate{Content: "Error creating game", Flags: 64})
		return
	}

	r := game.Render(event.BaseInteraction)
	b, err := game.RenderImage(true)
	if err != nil {
		_ = event.CreateMessage(discord.NewMessageCreateBuilder().SetContent("Error rendering image").SetFlags(64).Build())
		return
	}
	attachment := discord.NewFile("word-"+event.ID().String()+".png", "Wordle Game", b)
	if err := event.CreateMessage(discord.NewMessageCreateBuilder().
		SetEmbeds(r.Embeds...).
		SetContainerComponents(r.Components...).
		SetFiles(attachment).
		SetFlags(r.Flags).
		Build()); err != nil {
		log.Errorf("Error creating message: %s", err)
	}
}
