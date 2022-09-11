package bot

import (
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"
)

// GetResponses will upload the relevant text file for the given channel name
func GetResponses(s *discordgo.Session, m *discordgo.MessageCreate, command string) {
	// get args
	args := strings.Split(command, " ")
	if len(args) < 2 {
		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("<@%s>, please provide your groups standup channel name...", m.Author.ID))
		return
	}
}
