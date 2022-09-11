package bot

import (
	"fmt"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/hankpeeples/scrum-bot/utils"
)

// StandupInit initializes and runs the standup message thread creation loop
func StandupInit(s *discordgo.Session, m *discordgo.MessageCreate) {
	// Initialize timer duration
	duration := time.Second * 5
	// Hard coding channel IDs for simplicity
	// channelIDs := []string{"1016903999628259411", "1014940760568774666", "1016070677419270175"}
	channelIDs := []string{"1018326204161470526", "1018326221840449567", "1018326246364565645"}
	// number of channels
	numChannels := len(channelIDs)
	// Create message send ticker
	ticker := time.NewTicker(duration)

	// create initialized channels message
	var channels string
	for i, channel := range channelIDs {
		if i == len(channelIDs)-1 { // no trailing comma
			channels += fmt.Sprintf("<#%s>", channel)
			break
		}
		channels += fmt.Sprintf("<#%s>, ", channel)
	}

	// Send initialized confirmation in channel '!init' was used
	_, err := s.ChannelMessageSendEmbed(m.ChannelID, &discordgo.MessageEmbed{
		Title:       "Standup coordinator initialized",
		Description: fmt.Sprintf("Standup messages will be sent in %s: Monday - Friday, around 8am. \n\nA text file will be created to store responses should they be needed at a later date. Use `!getResponses <text channel name>` and the text file for your group will be uploaded to discord for your use.", channels),
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
			fullDate = time.Now().Local().Format(time.RubyDate)
			// Only want first part of date: `Fri Sep 02`
			date = fullDate[0:10]
			// Only day of the week
			day = date[0:3]

			day = "Mon"
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

					// send questions in thread
					_, err = s.ChannelMessageSendEmbed(thread.ID, &discordgo.MessageEmbed{
						Description: "Answer all questions in one single message and number each one.",
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
