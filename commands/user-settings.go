package commands

import (
	"context"
	"fmt"

	"github.com/Skye-31/WordleBot/models"
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
	"github.com/uptrace/bun"
)

func editUserSettings(db *bun.DB, event *events.ApplicationCommandInteractionEvent) {
	data := event.SlashCommandInteractionData()
	u := models.User{
		ID:                event.User().ID,
		CachedTag:         event.User().Tag(),
		Public:            false,
		DefaultWordLength: 5,
	}

	query := db.NewInsert().On("CONFLICT (id) DO UPDATE")
	columns := []string{"id", "tag"}
	if o, e := data.OptBool("public"); e {
		u.Public = o
		query.Set("public = ?", u.Public)
		columns = append(columns, "public")
	}
	if o, e := data.OptInt("default-word-size"); e {
		u.DefaultWordLength = uint8(o)
		query.Set("default_word_length = ?", u.DefaultWordLength)
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
			ID:        event.User().ID,
			CachedTag: event.User().Tag(),
		}
		if err := db.NewSelect().Model(&u).WherePK().Scan(context.TODO()); err != nil {
			u.Public = false
			u.DefaultWordLength = 5
		}
		m = &u
	}
	if err := event.CreateMessage(discord.MessageCreate{
		Embeds: []discord.Embed{
			discord.NewEmbedBuilder().
				SetAuthor(event.User().Tag(), "", event.User().EffectiveAvatarURL(discord.WithSize(128))).
				SetTitle("User Settings").
				AddField("Public", boolEmote(m.Public), true).
				AddField("Default Word Length", fmt.Sprintf("%d", m.DefaultWordLength), true).
				Build(),
		},
		Flags: discord.MessageFlagEphemeral,
	}); err != nil {
		event.Client().Logger().Error(err)
	}
}

func boolEmote(b bool) string {
	if b {
		return "✅"
	}
	return "❌"
}
