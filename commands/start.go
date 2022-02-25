package commands

import (
	"context"

	"github.com/DisgoOrg/disgo/core/events"
	"github.com/DisgoOrg/disgo/discord"
	"github.com/Skye-31/WordleBot/models"
	"github.com/Skye-31/WordleBot/types"
	"github.com/uptrace/bun"
)

func start(db *bun.DB, _ *types.WordsData, event *events.ApplicationCommandInteractionEvent) {
	//data := event.SlashCommandInteractionData()
	count, err := db.NewSelect().Model((*models.Game)(nil)).Where("id = ?", event.User.ID).Count(context.TODO())
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
					discord.NewPrimaryButton("Continue game", "start:continue").WithEmoji(discord.ComponentEmoji{Name: "➡"}),
					discord.NewDangerButton("Give up", "start:giveup").WithEmoji(discord.ComponentEmoji{Name: "🗑"}),
				)},
		})
		return
	}
	_ = event.CreateMessage(discord.MessageCreate{Content: "h"})
	//guess := data.Options.Bool("guess")
}
