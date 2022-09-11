// Package bot will contain bot specific functions
package bot

import (
	"os"
	"os/signal"
	"strings"
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
	log.Info("Starting bot session...")
	// Create discord session
	dg, err := discordgo.New("Bot " + token)
	if err != nil {
		log.Fatalf("Error creating discord session: %v", err)
	}

	dg.AddHandler(ready)
	// Register messageCreate func as callback for message events
	dg.AddHandler(messageCreate)

	// Need information about guilds (which includes channels), messages
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

func ready(s *discordgo.Session, event *discordgo.Ready) {
	// Set status message
	err := s.UpdateGameStatus(0, "Standup Coordinator")
	if err != nil {
		log.Error("Status message was NOT updated...")
		return
	}
	log.Info("Bot status message updated successfully")
}

func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	// Hard coding channel IDs for simplicity
	// ! channelIDs := []string{"1016903999628259411", "1014940760568774666", "1016070677419270175"}
	channelIDs := []string{"1018326204161470526", "1018326221840449567", "1018326246364565645"}
	// Ignore all messages created by the bot itself
	if m.Author.ID == s.State.User.ID {
		return
	}

	// Only look for commands that begin with defined prefix character
	if strings.HasPrefix(m.Content, utils.Prefix) {
		// Log received commands
		log.Infof("[%s]: %s", m.Author, m.Content)

		// Grab command following prefix
		command := m.Content[1:]

		if command == "init" {
			StandupInit(s, m, channelIDs)
		} else if command == "getResponses" {
			GetResponses(s, m, command)
		}
	}

	// find channel (thread)
	thread, err := s.Channel(m.ChannelID)
	if err != nil {
		log.Errorf("Error getting thread: %v", err)
	}
	// only grab responses from threads
	if thread.IsThread() {
		// parse and save to file
		utils.SaveResponse(s, m, thread)
	}
}
