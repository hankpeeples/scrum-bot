// Package utils contains utility functions needed throughout the app
package utils

import (
	"fmt"
	"os"
	"time"

	"github.com/bwmarrin/discordgo"
	logger "github.com/withmandala/go-log"
)

const (
	// MainGuildID is the Guild ID of the group project server
	MainGuildID string = "938592071307116545"
	// TestGuildID is the Guild ID of my test server
	TestGuildID string = "706588663747969104"
	// MainParentCategoryID is the category ID of the group server 'stand ups'
	MainParentCategoryID string = "1019314319781019668"
	// TestParentCategoryID is the category ID of my server 'stand ups'
	TestParentCategoryID string = "1019361832877699153"
	// Prefix is the bot command prefix character
	Prefix string = "!"
)

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

// GetStandupChannels returns the standup channels we need to use
func GetStandupChannels(s *discordgo.Session) []string {
	// use 'stand ups' category as parentID to match stand-ups channels
	// TODO: Swap GuildID
	guild, err := s.GuildChannels(TestGuildID)
	if err != nil {
		log.Fatalf("Error getting guild channels: %v", err)
	}

	var standupChannels []string
	// loop through channels and pull out text channels
	for _, c := range guild {
		// Make sure channel is a text channel, and make sure its parent category ID matches
		// TODO: Swap ParentCategoryID
		if c.Type == discordgo.ChannelTypeGuildText && c.ParentID == TestParentCategoryID {
			standupChannels = append(standupChannels, c.ID)
		}
	}

	log.Infof("Found %d stand up channels.", len(standupChannels))
	return standupChannels
}
