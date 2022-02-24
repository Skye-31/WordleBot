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
		if data.CommandName == "say" {
			err := event.CreateMessage(discord.NewMessageCreateBuilder().
				SetContent(*data.Options.String("message")).
				Build(),
			)
			if err != nil {
				event.Bot().Logger.Error("error on sending response: ", err)
			}
		} else {
			_ = event.CreateMessage(discord.NewMessageCreateBuilder().
				SetContent("received command").
				Build(),
			)
		}
	}
}
