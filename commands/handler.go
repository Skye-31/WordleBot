package commands

import (
	"github.com/DisgoOrg/disgo/core"
	"github.com/DisgoOrg/disgo/core/events"
	"github.com/DisgoOrg/disgo/discord"
	"github.com/Skye-31/WordleBot/commands/components"
	"github.com/Skye-31/WordleBot/types"
	"github.com/uptrace/bun"
)

func Listener(_ *core.Bot, db *bun.DB, words *types.WordsData) func(event *events.ApplicationCommandInteractionEvent) {
	return func(event *events.ApplicationCommandInteractionEvent) {
		data := event.SlashCommandInteractionData()
		n := commandName(data)
		switch n {
		case "user/settings/edit":
			editUserSettings(db, event)
		case "user/settings/view":
			viewUserSettings(db, nil, event)
		case "start":
			start(db, words, event)
		default:
			_ = event.CreateMessage(discord.MessageCreate{Content: "Unknown command: " + n, Flags: discord.MessageFlagEphemeral})
		}
	}
}

func ComponentInteraction(_ *core.Bot, _ *bun.DB, _ *types.WordsData) func(event *events.ComponentInteractionEvent) {
	return func(event *events.ComponentInteractionEvent) {
		id := event.Data.ID()
		switch id {
		case "game:guess":
			components.Guess(event)
		default:
			_ = event.CreateMessage(discord.MessageCreate{Content: "Unknown component interaction: " + id.String(), Flags: discord.MessageFlagEphemeral})
		}
	}
}

func commandName(data core.SlashCommandInteractionData) string {
	name := data.Name()
	if data.SubCommandGroupName != nil {
		name += "/" + *data.SubCommandGroupName
	}
	if data.SubCommandName != nil {
		name += "/" + *data.SubCommandName
	}
	return name
}
