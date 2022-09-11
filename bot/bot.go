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
			// Initialize timer duration
			duration := time.Hour * 24
			// Hard coding channel IDs for simplicity
			channelIDs := []string{"1016903999628259411", "1014940760568774666", "1016070677419270175"}
			// number of channels
			numChannels := len(channelIDs)
			// Create message send ticker
			ticker := time.NewTicker(duration)

			// Send initialized confirmation in channel '!init' was used
			_, err := s.ChannelMessageSendEmbed(m.ChannelID, &discordgo.MessageEmbed{
				Title:       "Standup coordinator initialized",
				Description: fmt.Sprintf("Standup messages will be sent in <#%s>, <#%s>, <#%s>:\nMonday - Friday, around 8am", channelIDs[0], channelIDs[1], channelIDs[2]),
				Color:       green,
			})
			if err != nil {
				utils.HandleEmbedFailure(s, m, err)
			}

			go func() {
				var fullDate, date, day string
				// Counter for looping through each text channelID
				i := 0
				// run every 24 hours
				_ = <-ticker.C
				for {
					if i >= numChannels {
						i = 0
						// Reset the ticker
						ticker.Reset(duration)
						_ = <-ticker.C
					}
					// Get current date and time
					fullDate = time.Now().UTC().Format(time.RubyDate)
					// Only want first part of date: `Fri Sep 02`
					date = fullDate[0:10]
					// Only day of the week
					day = date[0:3]

					// If it is saturday or sunday, no message
					if day == "Sat" || day == "Sun" {
						log.Infof("No standup today: %s", day)
						// 'i' and the ticker will reset so this message is only logged once
						i = numChannels
					} else {
						// Find each channels current state
						ch, err := s.State.Channel(channelIDs[i])
						if err != nil {
							log.Errorf("Error getting channel: %s", err)
							s.ChannelMessageSend(m.ChannelID, "Unable to find correct channels...")
							return
						}

						if !ch.IsThread() {
							msg, err := s.ChannelMessageSend(channelIDs[i], fmt.Sprintf("Standup Thread for `%s`", date))
							if err != nil {
								log.Errorf("Thread init msg: %v", err)
							}

							// Create thread. Thread might archive after 300min (5 hours).
							// Not sure what the archive duration actually does...
							thread, err := s.MessageThreadStart(channelIDs[i], msg.ID, "Standup meeting", 1440)
							if err != nil {
								log.Errorf("Thread start: %v", err)
							}

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
							log.Infof("Standup message sent [%d of %d]", i+1, numChannels)
						}
					}
					i++
				} // end of for loop
			}()
		}
	}
}
