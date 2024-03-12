package handlers

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"pokabook/issue-bot/internal/discord"
	"pokabook/issue-bot/pkg/utils"
)

func RevokeIssueHandle(s *discordgo.Session, e *discordgo.MessageReactionRemove) {

	_, err := utils.CheckChannel(s, e.ChannelID)
	if err != nil || e.Emoji.Name != "🙏" {
		return
	}

	userID, exists := discord.IssueMapping[e.ChannelID]
	if !exists || userID != e.UserID {
		return
	}
	delete(discord.IssueMapping, e.ChannelID)

	message := fmt.Sprintf("<@%s>가 이슈를 반환했습니다.\n", e.UserID)
	_, err = s.ChannelMessageSend(e.ChannelID, message)
	if err != nil {
		fmt.Println("멘션 메시지를 전송하는데 실패했습니다: ", err)
		return
	}

}
