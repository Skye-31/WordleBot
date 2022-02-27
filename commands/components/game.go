package components

import (
	"github.com/DisgoOrg/disgo/core/events"
	"github.com/DisgoOrg/disgo/discord"
	"github.com/DisgoOrg/disgo/json"
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
