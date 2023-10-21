package main

import (
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
)

func init() {
	// loads values from .env into the system
	if err := godotenv.Load(); err != nil {
		fmt.Print("No .env file found")
		os.Exit(1)
	}
}

func caseContains(s, substr string) bool {
	s, substr = strings.ToUpper(s), strings.ToUpper(substr)
	return strings.Contains(s, substr)
}

func main() {
	discordToken, exists := os.LookupEnv("DISCORD_TOKEN")
	if !exists {
		fmt.Print("No DISCORD_TOKEN set in .env file")
		os.Exit(1)
	}

	dg, err := discordgo.New("Bot " + discordToken)
	if err != nil {
		fmt.Println("Error creating Discord session,", err)
		return
	}

	dg.AddHandler(messageCreate)

	dg.Identify.Intents = discordgo.IntentsGuildMessages

	err = dg.Open()
	if err != nil {
		fmt.Println("Error opening connection,", err)
		return
	}

	fmt.Println("Bot now running")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc

	dg.Close()
}

func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}

	if caseContains(m.Content, "https://twitter.com") || caseContains(m.Content, "https://www.twitter.com") || caseContains(m.Content, "https://x.com") {
		s.ChannelTyping(m.ChannelID)

		link := ""

		if caseContains(m.Content, "https://twitter.com") {
			link = "https://twitter.com"
		} else if caseContains(m.Content, "https://www.twitter.com") {
			link = "https://www.twitter.com"
		} else {
			link = "https://x.com"
		}

		editedMsg := strings.Replace(m.Content, link, "https://fxtwitter.com", 1)

		_, err := s.ChannelMessageSendReply(m.ChannelID, editedMsg, m.Reference())
		if err != nil {
			fmt.Println("Error", err)
		}
		//s.ChannelMessageDelete(m.ChannelID, m.Message.ID)
	}
}
