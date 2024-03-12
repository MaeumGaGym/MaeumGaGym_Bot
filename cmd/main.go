package main

import (
	"github.com/bwmarrin/discordgo"
	"log"
	"os"
	"os/signal"
	"pokabook/issue-bot/internal/discord/handlers"
)

var (
	Token   = os.Getenv("DISCORD_TOKEN")
	Session *discordgo.Session
)

func init() {
	var err error
	Session, err = discordgo.New("Bot " + Token)
	if err != nil {
		log.Fatalf("클라이언트 생성 오류: %v", err)
	}

	err = Session.Open()
	if err != nil {
		log.Fatalf("세션 오픈 오류: %v", err)
	}

	log.Printf("%s (%s)에 로그인 됨", Session.State.User.String(), Session.State.User.ID)
}

func main() {

	Session.AddHandler(handlers.SendMajorMentionHandle)
	Session.AddHandler(handlers.GrantIssueHandle)
	Session.AddHandler(handlers.RevokeIssueHandle)
	Session.AddHandler(handlers.CheckIssueFormatHandler)

	defer Session.Close()

	stop := make(chan os.Signal)
	signal.Notify(stop, os.Interrupt)
	<-stop
	log.Println("봇 종료됨")
}
