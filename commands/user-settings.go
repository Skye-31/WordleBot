package commands

import (
	"context"

	"github.com/DisgoOrg/disgo/core"
	"github.com/DisgoOrg/disgo/core/events"
	"github.com/DisgoOrg/disgo/discord"
	"github.com/Skye-31/WordleBot/models"
	"github.com/Skye-31/WordleBot/types"
	"github.com/uptrace/bun"
)

func editUserSettings(_ *core.Bot, db *bun.DB, _ *types.WordsData, event *events.ApplicationCommandInteractionEvent) {
	options := event.SlashCommandInteractionData().Options
	u := models.User{
		ID:                event.User.ID,
		CachedTag:         event.User.Tag(),
		Public:            false,
		DefaultWordLength: 5,
	}

	query := db.NewInsert().On("CONFLICT (id) DO UPDATE")
	hasSetPublic := false
	hasSetDefaultWordLength := false
	if o := options.BoolOption("public"); o != nil {
		u.Public = o.Value
		query.Set("public = ?", u.Public)
		hasSetPublic = true
	}
	if o := options.IntOption("default-word-size"); o != nil {
		u.DefaultWordLength = uint8(o.Value)
		query.Set("default_word_length = ?", u.DefaultWordLength)
		hasSetDefaultWordLength = true
	}
	columns := []string{"id", "tag"}
	if hasSetPublic {
		columns = append(columns, "public")
	}
	if hasSetDefaultWordLength {
		columns = append(columns, "default_word_length")
	}

	if _, err := query.Model(&u).Column(columns...).Exec(context.TODO()); err != nil {
		_ = event.CreateMessage(discord.MessageCreate{Content: "Failed to edit user settings: " + err.Error()})
		return
	}
	_ = event.CreateMessage(discord.MessageCreate{Content: "User settings edited!", Flags: discord.MessageFlagEphemeral})
}
