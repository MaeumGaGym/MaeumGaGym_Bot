package handlers

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"pokabook/issue-bot/internal/discord"
	"pokabook/issue-bot/pkg/utils"
)

func CheckIssueFormatHandler(s *discordgo.Session, e *discordgo.ThreadCreate) {
	if e.ID == "" || e.LastMessageID != "" {
		return
	}

	if e.Type != discordgo.ChannelTypeGuildPublicThread || e.ParentID != discord.TargetForumChannelID {
		return
	}

	firstMessage, err := utils.GetThreadFirstMessage(s, e.ID)
	if err != nil {
		fmt.Println("메시지를 찾는데 실패했습니다: ", err)
		return
	}

	_, err = utils.ParseIssueContent(firstMessage.Content)
	if err != nil {
		_, err = s.ChannelMessageSend(e.ID, err.Error())
		if err != nil {
			fmt.Println("멘션 메시지를 전송하는데 실패했습니다: ", err)
			return
		}
		return
	}

}
