package handlers

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"pokabook/issue-bot/internal/discord"
	"pokabook/issue-bot/pkg/asana"
	"pokabook/issue-bot/pkg/asana/api"
	"pokabook/issue-bot/pkg/utils"
)

func GrantIssueHandle(s *discordgo.Session, e *discordgo.MessageReactionAdd) {
	thread, err := utils.CheckChannel(s, e.ChannelID)

	if err != nil || e.Emoji.Name != "🙏" {
		return
	}

	if _, exists := discord.IssueMapping[e.MessageID]; exists {
		return
	}

	discord.IssueMapping[e.MessageID] = e.UserID

	pre, err := s.ChannelMessageSend(e.ChannelID, fmt.Sprintf("<@%s>에게 이슈를 할당중", e.UserID))
	if err != nil {
		fmt.Println("멘션 메시지를 전송하는데 실패했습니다: ", err)
		return
	}

	firstMessage, err := utils.GetThreadFirstMessage(s, thread.ID)
	if err != nil {
		fmt.Println("메시지를 찾는데 실패했습니다: ", err)
		return
	}

	contents, err := utils.ParseIssueContent(firstMessage.Content)
	if err != nil {
		_, err = s.ChannelMessageEdit(e.ChannelID, pre.ID, err.Error())
		if err != nil {
			fmt.Println("멘션 메시지를 전송하는데 실패했습니다: ", err)
			return
		}
		return
	}

	appliedTag := asana.TagsIdMap[thread.AppliedTags[0]]
	gid, url := api.CreateTask(thread.Name, thread.AppliedTags[0], e.UserID, contents.Details, appliedTag)

	message := fmt.Sprintf("<@%s>에게 [이슈(%s)](%s)가 할당되었습니다.\n", e.UserID, gid, url)
	_, err = s.ChannelMessageEdit(e.ChannelID, pre.ID, message)
	if err != nil {
		fmt.Println("멘션 메시지를 전송하는데 실패했습니다: ", err)
		return
	}
}
