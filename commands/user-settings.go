package commands

import (
	"context"
	"fmt"

	"github.com/DisgoOrg/disgo/core/events"
	"github.com/DisgoOrg/disgo/discord"
	"github.com/Skye-31/WordleBot/models"
	"github.com/uptrace/bun"
)

func editUserSettings(db *bun.DB, event *events.ApplicationCommandInteractionEvent) {
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

	if _, err := query.Model(&u).Column(columns...).Returning("*").Exec(context.TODO()); err != nil {
		_ = event.CreateMessage(discord.MessageCreate{Content: "Failed to edit user settings: " + err.Error()})
		return
	}
	viewUserSettings(db, &u, event)
}

func viewUserSettings(db *bun.DB, m *models.User, event *events.ApplicationCommandInteractionEvent) {
	if m == nil {
		u := models.User{
			ID:        event.User.ID,
			CachedTag: event.User.Tag(),
		}
		if err := db.NewSelect().Model(&u).WherePK().Scan(context.TODO()); err != nil {
			u.Public = false
			u.DefaultWordLength = 5
		}
		m = &u
	}
	_ = event.CreateMessage(discord.MessageCreate{
		Embeds: []discord.Embed{
			discord.NewEmbedBuilder().
				SetAuthor(event.User.Tag(), "", event.User.EffectiveAvatarURL(128)).
				SetTitle("User Settings").
				AddField("Public", boolEmote(m.Public), true).
				AddField("Default Word Length", fmt.Sprintf("%d", m.DefaultWordLength), true).
				Build(),
		},
		Flags: discord.MessageFlagEphemeral,
	})
}

func boolEmote(b bool) string {
	if b {
		return "✅"
	}
	return "❌"
}
