package bot

import (
	"fmt"
	"io"
	"os"

	"github.com/bwmarrin/discordgo"
)

// GetResponses will upload the relevant text file for the given channel name
func GetResponses(s *discordgo.Session, m *discordgo.MessageCreate, command []string) {
	if len(command) < 2 {
		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("<@%s>, please provide your groups standup channel name...", m.Author.ID))
		return
	}

	// get file reader
	var r io.Reader
	var err error

	r, err = os.Open(command[1] + ".txt")
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("<@%s>, unable to find file! It may not be created yet...", m.Author.ID))
		log.Errorf("Error opening file: %v", err)
		return
	}

	// send file to text channel
	_, err = s.ChannelFileSend(m.ChannelID, command[1]+".txt", r)
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("<@%s>, unable to upload file!", m.Author.ID))
		log.Errorf("Error uploading file: %v", err)
		return
	}
}
