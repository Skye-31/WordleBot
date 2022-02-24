package commands

import (
	"github.com/DisgoOrg/disgo/core"
	"github.com/DisgoOrg/disgo/core/events"
	"github.com/DisgoOrg/disgo/discord"
	"github.com/Skye-31/WordleBot/types"
)

func Listener(_ *core.Bot, _ *types.WordsData) func(event *events.ApplicationCommandInteractionEvent) {
	return func(event *events.ApplicationCommandInteractionEvent) {
		data := event.SlashCommandInteractionData()
		n := commandName(data)
		switch n {
		default:
			_ = event.CreateMessage(discord.MessageCreate{Content: "Unknown command: " + n, Flags: discord.MessageFlagEphemeral})
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
