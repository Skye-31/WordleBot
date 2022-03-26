package commands

import (
	"github.com/disgoorg/disgo/discord"
)

var four, eight = 4, 8

var Definition = []discord.ApplicationCommandCreate{
	discord.SlashCommandCreate{
		CommandName:       "user",
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
						Options: []discord.ApplicationCommandOption{
							discord.ApplicationCommandOptionBool{
								Name:        "public",
								Description: "Whether to your public profile should show up to other users",
								Required:    false,
							},
							discord.ApplicationCommandOptionInt{
								Name:        "default-word-size",
								Description: "The default word size for your Wordle games",
								Required:    false,
								MinValue:    &four,
								MaxValue:    &eight,
							},
						},
					},
				},
			},
		},
	},
	//discord.SlashCommandCreate{
	//	Name:        "leaderboard",
	//	Description: "View the leaderboard",
	//	Options: []discord.ApplicationCommandOption{
	//		discord.ApplicationCommandOptionSubCommand{
	//			Name:        "streaks",
	//			Description: "View the leaderboard for streaks",
	//		},
	//		discord.ApplicationCommandOptionSubCommand{
	//			Name:        "total",
	//			Description: "View the overall leaderboard words",
	//		},
	//	},
	//	DefaultPermission: true,
	//},
	//discord.SlashCommandCreate{
	//	Name:              "streak",
	//	Description:       "View streak information",
	//	DefaultPermission: true,
	//	Options: []discord.ApplicationCommandOption{
	//		discord.ApplicationCommandOptionUser{
	//			Name:        "user",
	//			Description: "The user to view the streak of",
	//		},
	//	},
	//},
	discord.SlashCommandCreate{
		CommandName:       "stats",
		Description:       "View a user's stats",
		DefaultPermission: true,
		Options: []discord.ApplicationCommandOption{
			discord.ApplicationCommandOptionUser{
				Name:        "user",
				Description: "The user to view the streak of",
			},
		},
	},
	discord.SlashCommandCreate{
		CommandName:       "start",
		Description:       "Start a wordle game",
		DefaultPermission: true,
		Options: []discord.ApplicationCommandOption{
			discord.ApplicationCommandOptionString{
				Name:        "starting-word",
				Description: "The word to start the game with",
				Required:    true,
			},
			discord.ApplicationCommandOptionInt{
				Name:        "letters",
				Description: "The number of letters to use in the wordle. (Default: 5)",
				MinValue:    &four,
				MaxValue:    &eight,
			},
		},
	},
	discord.SlashCommandCreate{
		CommandName:       "github",
		Description:       "View the source code for this bot",
		DefaultPermission: true,
	},
}
