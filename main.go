package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/DisgoOrg/disgo/core"
	"github.com/DisgoOrg/disgo/core/bot"
	"github.com/DisgoOrg/disgo/core/events"
	"github.com/DisgoOrg/disgo/httpserver"
	"github.com/DisgoOrg/disgo/info"
	"github.com/DisgoOrg/log"
	"github.com/Skye-31/WordleBot/commands"
	"github.com/Skye-31/WordleBot/types"
	"github.com/oasisprotocol/curve25519-voi/primitives/ed25519"
)

func main() {
	logger := log.New(log.Ltime | log.Lshortfile)
	logger.SetLevel(log.LevelInfo)
	logger.Info("Starting Wordlebot on Disgo ", info.Version)

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

	httpserver.Verify = func(publicKey httpserver.PublicKey, message, sig []byte) bool {
		return ed25519.Verify(publicKey, message, sig)
	}

	disgo, err := bot.New(config.Token,
		bot.WithLogger(logger),
		bot.WithHTTPServerOpts(
			httpserver.WithURL("/interactions"),
			httpserver.WithPort(":4567"),
			httpserver.WithPublicKey(config.PublicKey),
		),
		bot.WithCacheOpts(core.WithCacheFlags(core.CacheFlagsNone), core.WithMemberCachePolicy(core.MemberCachePolicyNone), core.WithMessageCachePolicy(core.MessageCachePolicyNone)),
	)
	if err != nil {
		logger.Fatal("error while building disgo instance: ", err)
	}
	disgo.AddEventListeners(&events.ListenerAdapter{
		OnApplicationCommandInteraction: commands.Listener(disgo, words),
	})

	defer disgo.Close(context.TODO())

	if config.DevMode {
		if config.DevGuildID == "" {
			logger.Fatal("DevMode is enabled but no DevGuildID is set")
			return
		}
		if _, err = disgo.SetGuildCommands(config.DevGuildID, commands.Definition); err != nil {
			logger.Fatal("error while registering commands: ", err)
			return
		}
	}

	if err = disgo.StartHTTPServer(); err != nil {
		logger.Fatal("error while starting http server: ", err)
		return
	}

	logger.Infof("WordleBot is now running. Press CTRL-C to exit.")
	s := make(chan os.Signal, 1)
	signal.Notify(s, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-s
}
