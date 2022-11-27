package bot

import (
	"fmt"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/hankpeeples/scrum-bot/utils"
)

// StandupInit initializes and runs the standup message thread creation loop
func StandupInit(s *discordgo.Session, channelIDs []*discordgo.Channel) {
	var diff int
	// Set 't' to hours. If dateStruct.Hour is == 8, the ticker needs a positive non-zero
	// duration. So this will change to ~ time.Second * 2.
	t := time.Hour

	dateStruct := utils.GetDate()
	// check current time is before or after 8am
	if dateStruct.Hour > 8 {
		// calculate time until 8am
		diff = (24 - dateStruct.Hour) + 8
		log.Infof("Waiting %d hour(s) before starting standup timer.", diff)
	} else if dateStruct.Hour < 8 {
		// calculate time until 8am
		diff = 8 - dateStruct.Hour
		log.Infof("Waiting %d hour(s) after starting standup timer.", diff)
	} else {
		// It is ~8am here! (Not regarding minutes)
		log.Info("It's ~8am! Sending standup messages in 2 seconds...")
		// set 't' to seconds to send messages now, will reset to 24 hr after
		t = time.Second
		diff = 2
	}

	// create message send ticker
	ticker := time.NewTicker(t * time.Duration(diff+3))

	// Counter for looping through each text channelID
	i := 0

	go func(i int) {

		// waiting initial time difference between bot start and 8 AM
		<-ticker.C

		// get channels again, making sure all will be used
		channelIDs = utils.GetStandupChannels(s)
		// number of channels
		numChannels := len(channelIDs)

		// re-acquire date
		dateStruct = utils.GetDate()

		// setting newTicker to 24 hours, if this is not done the standup messages will
		// be sent twice in a row on the first iteration
		newTicker := time.NewTicker(time.Hour * 24)

		for {
			// If it is saturday or sunday, no message
			if dateStruct.Day == "Sat" || dateStruct.Day == "Sun" {
				log.Infof("No standup today: %s", dateStruct.Day)
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
					msg, err := s.ChannelMessageSend(channelIDs[i].ID, fmt.Sprintf("Standup Thread for `%s`", dateStruct.FullDate))
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
			// see if all messages have been sent, if so...reset ticker
			if i >= numChannels {
				i = 0
				log.Info("Waiting 24hrs...")
				<-newTicker.C
				// get channels again, making sure all will be used
				channelIDs = utils.GetStandupChannels(s)
				// number of channels
				numChannels = len(channelIDs)
				// re-acquire date
				dateStruct = utils.GetDate()
			}
		} // end of for loop
	}(i)
}
