package commands

import (
	"github.com/DisgoOrg/disgo/discord"
	"github.com/DisgoOrg/disgo/json"
)

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
		},
	},
	discord.SlashCommandCreate{
		Name:        "streak",
		Description: "View streak information",
		Options: []discord.ApplicationCommandOption{
			discord.ApplicationCommandOptionSubCommand{
				Name:        "top",
				Description: "View the leaderboard of streaks",
			},
			discord.ApplicationCommandOptionSubCommand{
				Name:        "user",
				Description: "View a user's streak",
				Options: []discord.ApplicationCommandOption{
					discord.ApplicationCommandOptionUser{
						Name:        "user",
						Description: "The user to view the streak of",
					},
				},
			},
		},
	},
	discord.SlashCommandCreate{
		Name:        "stats",
		Description: "View a user's stats",
		Options: []discord.ApplicationCommandOption{
			discord.ApplicationCommandOptionSubCommand{
				Name:        "user",
				Description: "View a user's streak",
				Options: []discord.ApplicationCommandOption{
					discord.ApplicationCommandOptionUser{
						Name:        "user",
						Description: "The user to view the streak of",
					},
				},
			},
		},
	},
	discord.SlashCommandCreate{
		Name:        "start",
		Description: "Start a wordle game",
		Options: []discord.ApplicationCommandOption{
			discord.ApplicationCommandOptionInt{
				Name:        "letters",
				Description: "The number of letters to use in the wordle. (Default: 5)",
				MinValue:    json.NewInt(4),
				MaxValue:    json.NewInt(8),
			},
		},
	},
}
