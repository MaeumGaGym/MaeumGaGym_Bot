package pkg

import (
	"bufio"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"strings"
)

var (
	targetForumChannelID = "1214467392189505566"

	IOS = "1214468969080037376"
	BAC = "1214469037669617677"
	FLU = "1214469454537297920"
	FRO = "1214469532098371645"
	DO  = "1214469699686105088"
)

type IssueSections struct {
	Details  string
	Labels   []string
	Priority string
}

func SendMajorMention(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID || m.ChannelID == "" {
		return
	}

	channel, err := s.State.Channel(m.ChannelID)
	if err != nil || channel.Type != discordgo.ChannelTypeGuildPublicThread || channel.ParentID != targetForumChannelID {
		return
	}

	messages, err := s.ChannelMessages(m.ChannelID, 100, "", "", "")
	if err != nil || len(messages) > 1 {
		return
	}

	var iOSRoleID = "1154312160575500298"
	var BACRoleID = "1153326597844238366"
	var FLURoleID = "1154312210705821756"
	var FRORoleID = "1154312086965465119"
	var DORoleID = "1154350370735280140"

	var mention string
	for _, tag := range channel.AppliedTags {
		switch tag {
		case IOS:
			mention += "<@&" + iOSRoleID + "> "
		case BAC:
			mention += "<@&" + BACRoleID + "> "
		case FLU:
			mention += "<@&" + FLURoleID + "> "
		case FRO:
			mention += "<@&" + FRORoleID + "> "
		case DO:
			mention += "<@&" + DORoleID + "> "
		}
	}

	if mention != "" {
		message := "관련 분야의 이슈가 있습니다: " + mention
		_, err := s.ChannelMessageSend(m.ChannelID, message)
		if err != nil {
			fmt.Println("멘션 메시지를 전송하는데 실패했습니다: ", err)
		}
	}

}

func SetIssueManager(s *discordgo.Session, m *discordgo.MessageReactionAdd) {

	if m.UserID == s.State.User.ID || m.ChannelID == "" {
		return
	}

	channel, err := s.State.Channel(m.ChannelID)
	if err != nil || channel.Type != discordgo.ChannelTypeGuildPublicThread || channel.ParentID != targetForumChannelID {
		return
	}

	if m.Emoji.Name != "🙏" {
		err = s.MessageReactionRemove(m.ChannelID, m.MessageID, m.Emoji.APIName(), m.UserID)
		if err != nil {
			fmt.Println("이모지를 삭제하는데 실패했습니다: ", err)
		}
		return
	}

	reactions, err := s.MessageReactions(m.ChannelID, m.MessageID, m.Emoji.APIName(), 100, "", "")
	if err != nil {
		fmt.Println("반응 목록을 가져오는데 실패했습니다: ", err)
		return
	}

	if len(reactions) < 2 {
		user, err := s.User(m.UserID)
		if err != nil {
			fmt.Println("유저 정보를 가져오는데 실패했습니다: ", err)
			return
		}

		message := fmt.Sprintf("<@%s>에게 이슈가 할당되었습니다.\n", user.ID)

		_, err = s.ChannelMessageSend(m.ChannelID, message)
		if err != nil {
			fmt.Println("멘션 메시지를 전송하는데 실패했습니다: ", err)
		}
	}
}

// 디스코드 메시지를 파싱
func parseIssueContent(content string) IssueSections {
	scanner := bufio.NewScanner(strings.NewReader(content))
	var section IssueSections
	currentSection := ""

	for scanner.Scan() {
		line := scanner.Text()
		trimmedLine := strings.TrimSpace(line)

		if strings.HasPrefix(trimmedLine, "### ") {
			switch trimmedLine {
			case "### 이슈 상세 내용":
				currentSection = "details"
			case "### 라벨":
				currentSection = "labels"
			case "### 우선순위":
				currentSection = "priority"
			default:
				currentSection = ""
			}
			continue
		}

		if strings.HasPrefix(trimmedLine, "-") {
			trimmedLine = strings.TrimPrefix(trimmedLine, "- ")
		}

		switch currentSection {
		case "details":
			section.Details += line + "\n"
		case "labels":
			if trimmedLine != "" {
				section.Labels = append(section.Labels, trimmedLine)
			}
		case "priority":
			if section.Priority == "" {
				section.Priority = trimmedLine
			}
		}
	}

	return section
}
