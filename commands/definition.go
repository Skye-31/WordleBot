package commands

import "github.com/DisgoOrg/disgo/discord"

var Definition = []discord.ApplicationCommandCreate{
	discord.SlashCommandCreate{
		Name:              "user",
		Description:       "Edit your settings",
		DefaultPermission: true,
		Options: []discord.ApplicationCommandOption{
			discord.ApplicationCommandOptionSubCommandGroup{
				Name:        "settings",
				Description: "View/Edit your settings",
				Options: []discord.ApplicationCommandOptionSubCommand{
					{
						Name:        "view",
						Description: "Shows you your current settings",
					},
					{
						Name:        "edit",
						Description: "Edits your settings",
					},
				},
			},
			discord.ApplicationCommandOptionSubCommand{
				Name:        "streak",
				Description: "View your current streak",
			},
			discord.ApplicationCommandOptionSubCommand{
				Name:        "stats",
				Description: "View your stats",
			},
		},
	},
}
