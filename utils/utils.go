// Package utils contains utility functions needed throughout the app
package utils

import (
	"os"

	"github.com/bwmarrin/discordgo"
	logger "github.com/withmandala/go-log"
)

// Prefix is the bot command character prefix
const Prefix string = "!"

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

	log.Infof("Response: [%s] -> [%s/%s]", author.Nick, parent.Name, thread.Name)
}
