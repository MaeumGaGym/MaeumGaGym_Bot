package handlers

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"pokabook/issue-bot/internal/discord"
)

var TagsToRoleIDs = map[string]string{
	"1214468969080037376": "1154312160575500298", //iOS
	"1214469037669617677": "1153326597844238366", //Backend
	"1214469454537297920": "1154312210705821756", //Flutter
	"1214469532098371645": "1154312086965465119", //Frontend
	"1214469699686105088": "1154350370735280140", //DevOps
}

func SendMajorMentionHandle(s *discordgo.Session, e *discordgo.ThreadCreate) {
	if e.ID == "" || e.LastMessageID != "" {
		return
	}

	if e.Type != discordgo.ChannelTypeGuildPublicThread || e.ParentID != discord.TargetForumChannelID {
		return
	}

	var mention string
	for _, tag := range e.AppliedTags {
		if roleID, exists := TagsToRoleIDs[tag]; exists {
			mention += "<@&" + roleID + "> "
		}
	}

	if len(mention) > 0 {
		_, err := s.ChannelMessageSend(e.ID, "관련 분야의 이슈가 있습니다: "+mention)
		if err != nil {
			fmt.Printf("멘션 메시지를 전송하는데 실패했습니다: %+v\n", err)
		}
	}
}
