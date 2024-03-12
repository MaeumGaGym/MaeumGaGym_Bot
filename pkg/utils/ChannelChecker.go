package utils

import (
	"errors"
	"github.com/bwmarrin/discordgo"
)

func CheckChannel(s *discordgo.Session, channelId string) (*discordgo.Channel, error) {
	thread, err := s.State.Channel(channelId)

	if thread.Type != discordgo.ChannelTypeGuildPublicThread {
		return nil, errors.New("channel is not public thread")
	}

	return thread, err
}
