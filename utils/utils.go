// Package utils contains utility functions needed throughout the app
package utils

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/bwmarrin/discordgo"
	logger "github.com/withmandala/go-log"
)

const (
	// mainGID is the Guild ID of the group project server
	mainGID string = "938592071307116545"
	// testGID is the Guild ID of my test server
	testGID string = "706588663747969104"
	// GuildID is the Guild ID the bot is currently using (main or test)
	GuildID string = mainGID
	// mainParentCID is the category ID of the group server 'stand ups'
	mainParentCID string = "1019314319781019668"
	// testParentCID is the category ID of my server 'stand ups'
	testParentCID string = "1019361832877699153"
	parentCID     string = mainParentCID
	// Prefix is the bot command prefix character
	Prefix string = "!"
)

// DateStruct holds separated time/date information
type DateStruct struct {
	// FullDate : `Wed Sep 14`
	FullDate string
	// Day : `Thu`
	Day string
	// Hour is the hour of day as an integer
	Hour int
}

var log = NewLogger()

// NewLogger returns a new instance of a logger
func NewLogger() *logger.Logger {
	return logger.New(os.Stdout).WithColor()
}

// HandleEmbedFailure delivers a message to the channel that something
// went wrong when trying to send an embed
func HandleEmbedFailure(s *discordgo.Session, m *discordgo.MessageCreate, err error) {
	s.ChannelMessageSend(m.ChannelID, "Something broke... Couldn't send embedded message.")
	log.Error("Embed error: ", err)
}

// SaveResponse takes a thread response and saves it to a text file
func SaveResponse(s *discordgo.Session, m *discordgo.MessageCreate, thread *discordgo.Channel) {
	// find thread parent channel
	parent, err := s.Channel(thread.ParentID)
	if err != nil {
		log.Errorf("Error getting parent: %v", err)
	}
	// get message author details
	author, err := s.GuildMember(m.GuildID, m.Author.ID)
	if err != nil {
		log.Errorf("Error finding author: %v", err)
	}

	// create response file with append perms
	f, err := os.OpenFile(parent.Name+".txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Errorf("Error creating file: %v", err)
	}

	// close file on function exit
	defer f.Close()

	// write new response to file
	fmt.Fprintln(f, fmt.Sprintf("Author: %s | Date: %s", author.Nick, time.Now().Local().Format(time.RubyDate)))
	fmt.Fprintln(f, m.Content)
	fmt.Fprintln(f, "----------------------------------------------------------------------------")

	log.Infof("Response: [%s] -> [%s/%s]", author.Nick, parent.Name, thread.Name)
}

func getAllGuildChannels(s *discordgo.Session) []*discordgo.Channel {
	// use 'stand ups' category as parentID to match stand-ups channels
	guild, err := s.GuildChannels(GuildID)
	if err != nil {
		log.Fatalf("Error getting guild channels: %v", err)
	}

	return guild
}

// GetStandupChannels returns the standup channels we need to use
func GetStandupChannels(s *discordgo.Session) []*discordgo.Channel {
	channels := getAllGuildChannels(s)

	var standupChannels []*discordgo.Channel
	// loop through channels and pull out text channels
	for _, c := range channels {
		// Make sure channel is a text channel, and make sure its parent category ID matches
		if c.Type == discordgo.ChannelTypeGuildText && c.ParentID == parentCID {
			standupChannels = append(standupChannels, c)
		}
	}

	log.Infof("Found %d stand up channels.", len(standupChannels))
	return standupChannels
}

// FindGeneral returns the ID of a channel named 'general'
func FindGeneral(s *discordgo.Session) string {
	channels := getAllGuildChannels(s)

	// loop through channels and pull out text channels
	for _, c := range channels {
		// Make sure channel is a text channel, and make sure its name is 'general'
		if c.Type == discordgo.ChannelTypeGuildText && c.Name == "general" {
			log.Infof("Found '%s' channel.", c.Name)
			return c.ID
		}
	}

	return ""
}

// CreateChannelsPrintString makes a formatted string of all channels that will be sent messages
func CreateChannelsPrintString(channelIDs []*discordgo.Channel) string {
	var channels string

	for i, channel := range channelIDs {
		if i == len(channelIDs)-1 {
			channels += fmt.Sprintf("<#%s>.", channel.ID)
			break
		}
		channels += fmt.Sprintf("<#%s>, ", channel.ID)
	}

	return channels
}

// GetDate returns the current time, date, and day
func GetDate() *DateStruct {
	// Get current date and time
	fullDate := time.Now().Local().Format(time.RubyDate)

	// Convert hour section of time to int for comparison
	hourInt, err := strconv.Atoi(fullDate[11:13])
	if err != nil {
		log.Errorf("Error converting hour to int: %v", err)
	}

	d := &DateStruct{
		FullDate: fullDate[0:10],
		Day:      fullDate[0:3],
		Hour:     hourInt,
	}

	return d
}

// SendHeartbeat keeps the bot from timing out while waiting for
// standup ticker duration to expire
func SendHeartbeat(s *discordgo.Session) {
	duration := time.Hour * 1
	log.Info("Sending heartbeat every hour to keep connection alive...")
	heartbeatTicker := time.NewTicker(duration)

	defer heartbeatTicker.Stop()

	for {
		<-heartbeatTicker.C
		heartbeatTicker.Reset(duration)
	}
}
