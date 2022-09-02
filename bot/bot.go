// Package bot will contain bot specific functions
package bot

import (
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

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
			var channelID string = m.ChannelID
			// Create message send timeout
			timeout := time.NewTicker(time.Hour * 24)

			_, err := s.ChannelMessageSendEmbed(channelID, &discordgo.MessageEmbed{
				Title:       "Standup coordinator initialized",
				Description: fmt.Sprintf("Standup messages will be sent in <#%s>,\nMonday - Friday at 8am.", channelID),
				Color:       green,
			})
			if err != nil {
				utils.HandleEmbedFailure(s, m, err)
			}

			go func(s *discordgo.Session, channelID string) {
				for {
					// run function every 24 hours
					_ = <-timeout.C
					log.Info("Standup message sent...")
					// Get current date/day
					day := time.Now().UTC().Format(time.RubyDate)

					// If it is saturday or sunday, no message
					if day[0:3] == "Sat" || day[0:3] == "Sun" {
						log.Warnf("No standup today: %s", day[0:3])
						timeout.Reset(time.Hour * 24)
					}

					ch, err := s.State.Channel(channelID)
					if err != nil {
						log.Error("Error getting channel: %s", err)
					}

					if !ch.IsThread() {
						msg, _ := s.ChannelMessageSend(channelID, fmt.Sprintf("Standup Thread for `%s`", day[0:10]))
						// Create thread
						thread, err := s.MessageThreadStart(channelID, msg.ID, "Standup meeting", 60)

						_, err = s.ChannelMessageSendEmbed(thread.ID, &discordgo.MessageEmbed{
							Description: "Answer each question in this thread...",
							Fields: []*discordgo.MessageEmbedField{
								{
									Name:  "1. What did you work on last working day?",
									Value: "__",
								},
								{
									Name:  "2. What are you going to work on today?",
									Value: "__",
								},
								{
									Name:  "3. Are there any blocks to your workflow?",
									Value: "__",
								},
							},
							Color: blue,
						})
						if err != nil {
							utils.HandleEmbedFailure(s, m, err)
						}
					}

				}
			}(s, channelID)
		}
	}
}
