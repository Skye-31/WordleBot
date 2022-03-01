package commands

import (
	"strings"

	"github.com/DisgoOrg/disgo/core"
	"github.com/DisgoOrg/disgo/core/events"
	"github.com/DisgoOrg/disgo/discord"
	"github.com/Skye-31/WordleBot/commands/components"
	"github.com/Skye-31/WordleBot/types"
	"github.com/uptrace/bun"
)

func Listener(db *bun.DB, words *types.WordsData) func(event *events.ApplicationCommandInteractionEvent) {
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
		case "stats":
			stats(db, event)
		default:
			_ = event.CreateMessage(discord.MessageCreate{Content: "Unknown command: " + n, Flags: discord.MessageFlagEphemeral})
		}
	}
}

func ComponentInteraction(db *bun.DB, _ *types.WordsData) func(event *events.ComponentInteractionEvent) {
	return func(event *events.ComponentInteractionEvent) {
		id := event.Data.ID()
		startID := strings.Split(id.String(), ":")[1]
		switch startID {
		case "guess":
			components.Guess(event)
		case "continue":
			components.Continue(db, event)
		case "giveup":
			components.GiveUp(db, event)
		case "share":
			components.Share(event)
		default:
			_ = event.CreateMessage(discord.MessageCreate{Content: "Unknown component interaction: " + id.String(), Flags: discord.MessageFlagEphemeral})
		}
	}
}

func ModalInteraction(db *bun.DB, wd *types.WordsData) func(event *events.ModalSubmitInteractionEvent) {
	return func(event *events.ModalSubmitInteractionEvent) {
		id := event.Data.CustomID
		switch id {
		case "game:guess:submit":
			components.GuessSubmit(db, wd, event)
		default:
			_ = event.CreateMessage(discord.MessageCreate{Content: "Unknown modal interaction: " + id.String(), Flags: discord.MessageFlagEphemeral})
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
