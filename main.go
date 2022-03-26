package main

import (
	"context"
	"flag"
	"os"
	"os/signal"
	"syscall"

	"github.com/Skye-31/WordleBot/commands"
	"github.com/Skye-31/WordleBot/types"
	"github.com/disgoorg/disgo"
	"github.com/disgoorg/disgo/bot"
	"github.com/disgoorg/disgo/cache"
	"github.com/disgoorg/disgo/events"
	"github.com/disgoorg/disgo/httpserver"
	"github.com/disgoorg/log"
	"github.com/oasisprotocol/curve25519-voi/primitives/ed25519"
)

var (
	syncCommands = flag.Bool("sync-commands", false, "Whether the bot should sync commands")
	syncDB       = flag.Bool("sync-db", false, "Whether the bot should sync the database")
)

func main() {
	flag.Parse()
	logger := log.New(log.Ltime | log.Lshortfile)
	logger.SetLevel(log.LevelInfo)
	logger.Info("Starting Wordlebot on Disgo ", disgo.Version)
	logger.Infof("Syncing commands: %t, db: %t\n", *syncCommands, *syncDB)

	config, err := types.LoadConfig(logger)
	if err != nil {
		logger.Error(err)
		return
	}
	logger.SetLevel(config.LogLevel)

	words, err := types.LoadWordsData(logger)
	if err != nil {
		logger.Error("Failed to load words data: ", err)
		return
	}
	_ = words

	db, err := types.SetUpDatabase(config, logger, *syncDB)
	if err != nil {
		logger.Fatal("error while setting up database: ", err)
		return
	}

	httpserver.Verify = func(publicKey httpserver.PublicKey, message, sig []byte) bool {
		return ed25519.Verify(publicKey, message, sig)
	}
	client, err := disgo.New(config.Token,
		bot.WithLogger(logger),
		bot.WithHTTPServerConfigOpts(
			httpserver.WithURL("/interactions"),
			httpserver.WithAddress(":4567"),
			httpserver.WithPublicKey(config.PublicKey),
		),
		bot.WithCacheConfigOpts(cache.WithCacheFlags(cache.FlagsNone), cache.WithMemberCachePolicy(cache.MemberCachePolicyNone), cache.WithMessageCachePolicy(cache.MessageCachePolicyNone)),
		bot.WithEventListeners(&events.ListenerAdapter{
			OnApplicationCommandInteraction: commands.Listener(db, words),
			OnComponentInteraction:          commands.ComponentInteraction(db, words),
			OnModalSubmit:                   commands.ModalInteraction(db, words),
		}),
	)
	if err != nil {
		logger.Fatal("error while building disgo instance: ", err)
	}

	defer client.Close(context.TODO())
	if *syncCommands {
		if config.DevMode {
			if config.DevGuildID == "" {
				logger.Fatal("DevMode is enabled but no DevGuildID is set")
				return
			}
			if _, err = client.Rest().Applications().SetGuildCommands(client.ApplicationID(), config.DevGuildID, commands.Definition); err != nil {
				logger.Fatal("error while registering commands: ", err)
				return
			}
		} else {
			if _, err = client.Rest().Applications().SetGlobalCommands(client.ApplicationID(), commands.Definition); err != nil {
				logger.Fatal("error while registering commands: ", err)
				return
			}
		}
	}

	if err = client.StartHTTPServer(); err != nil {
		logger.Fatal("error while starting http server: ", err)
		return
	}

	logger.Infof("WordleBot is now running. Press CTRL-C to exit.")
	s := make(chan os.Signal, 1)
	signal.Notify(s, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-s
}
