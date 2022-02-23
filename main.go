package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/DisgoOrg/disgo/core"
	"github.com/DisgoOrg/disgo/core/bot"
	"github.com/DisgoOrg/disgo/core/events"
	"github.com/DisgoOrg/disgo/discord"
	"github.com/DisgoOrg/disgo/httpserver"
	"github.com/DisgoOrg/disgo/info"
	"github.com/DisgoOrg/log"
	"github.com/DisgoOrg/snowflake"
)

var (
	token     = os.Getenv("disgo_token")
	publicKey = os.Getenv("disgo_public_key")
	guildID   = snowflake.GetSnowflakeEnv("disgo_guild_id")

	commands = []discord.ApplicationCommandCreate{
		discord.SlashCommandCreate{
			Name:              "user",
			Description:       "Edit your settings",
			DefaultPermission: true,
		},
	}
)

func main() {
	log.SetLevel(log.LevelInfo)
	log.Info("Starting Wordlebot on Disgo ", info.Version)

	disgo, err := bot.New(token,
		bot.WithHTTPServerOpts(
			httpserver.WithURL("/interactions"),
			httpserver.WithPort(":4567"),
			httpserver.WithPublicKey(publicKey),
		),
		bot.WithCacheOpts(core.WithCacheFlags(0), core.WithMemberCachePolicy(core.MemberCachePolicyNone), core.WithMessageCachePolicy(core.MessageCachePolicyNone)),
		bot.WithEventListeners(&events.ListenerAdapter{
			OnApplicationCommandInteraction: commandListener,
		}),
	)
	if err != nil {
		log.Fatal("error while building disgo instance: ", err)
		return
	}

	defer disgo.Close(context.TODO())

	if _, err = disgo.SetGuildCommands(guildID, commands); err != nil {
		log.Fatal("error while registering commands: ", err)
	}

	if err = disgo.StartHTTPServer(); err != nil {
		log.Fatal("error while starting http server: ", err)
	}

	log.Infof("WordleBot is now running. Press CTRL-C to exit.")
	s := make(chan os.Signal, 1)
	signal.Notify(s, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-s
}

func commandListener(event *events.ApplicationCommandInteractionEvent) {
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
