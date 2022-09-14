package bot

import (
	"fmt"
	"time"

	"github.com/bwmarrin/discordgo"
)

// StandupInit initializes and runs the standup message thread creation loop
func StandupInit(s *discordgo.Session, channelIDs []*discordgo.Channel) {
	// Initialize timer duration
	duration := time.Hour * 24
	// number of channels
	numChannels := len(channelIDs)
	// Create message send ticker
	ticker := time.NewTicker(duration)

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
				ch, err := s.State.Channel(channelIDs[i].ID)
				if err != nil {
					log.Errorf("Error getting channel: %s", err)
					s.ChannelMessageSend(channelIDs[i].ID, "Unable to find correct channels...")
					return
				}

				if !ch.IsThread() {
					msg, err := s.ChannelMessageSend(channelIDs[i].ID, fmt.Sprintf("Standup Thread for `%s`", date))
					if err != nil {
						log.Errorf("Thread init msg: %v", err)
					}

					// Create thread. Thread might archive after 300min (5 hours).
					// Not sure what the archive duration actually does...
					thread, err := s.MessageThreadStart(channelIDs[i].ID, msg.ID, "Standup meeting", 1440)
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
						log.Errorf("[CRITICAL] Unable to send standup embed [%d]: %v", i+1, err)
						break
					}
					log.Infof("Standup message sent [%d of %d]", i+1, numChannels)
				}
			}
			i++
		} // end of for loop
	}()
}
