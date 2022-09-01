// Package bot will contain bot specific functions
package bot

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/bwmarrin/discordgo"
	"github.com/hankpeeples/scrum-bot/utils"
)

var (
	log   = utils.NewLogger()
	red   = 0xf54248
	blue  = 0x42b9f5
	green = 0x28de4f
)

// Start will begin a new discord bot session
func Start(token string) {
	log.Info("Attempting to start bot session...")
	// Create discord session
	dg, err := discordgo.New("Bot " + token)
	if err != nil {
		log.Fatalf("Error creating discord session: %v", err)
	}

	// Register ready func as callback for ready events
	// dg.AddHandler(ready)
	// Register messageCreate func as callback for message events
	// dg.AddHandler(messageCreate)

	// Need information about guilds (which includes channels),
	// messages and voice states.
	dg.Identify.Intents = discordgo.IntentsGuilds | discordgo.IntentsGuildMessages

	// Open websocket and begin listening
	err = dg.Open()
	if err != nil {
		log.Fatalf("Error opening websocket: %v", err)
	}

	log.Info("Session open and listening ✅")
	// Wait here for ctrl-c or other termination signal
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	log.Warn("Session terminated!")
	// close discordgo session after kill signal is received
	dg.Close()
}
