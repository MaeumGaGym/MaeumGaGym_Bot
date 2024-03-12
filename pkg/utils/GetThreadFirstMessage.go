package utils

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
)

func GetThreadFirstMessage(s *discordgo.Session, channelID string) (*discordgo.Message, error) {

	var lastMessageID string
	for {
		messages, err := s.ChannelMessages(channelID, 100, lastMessageID, "", "")
		if err != nil {
			return nil, err
		}

		if len(messages) == 0 {
			break
		}

		lastMessageID = messages[len(messages)-1].ID

		if len(messages) < 100 {
			return messages[len(messages)-1], nil
		}
	}

	return nil, fmt.Errorf("no first message found in thread")
}
