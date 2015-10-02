package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/corylanou/questionbot"
	"github.com/nlopes/slack"
)

const tokenKey = "SLACKBOT_TOKEN"

var qbot *questionBot.Service

func main() {
	token := os.Getenv(tokenKey)
	api := slack.New(token)
	//api.SetDebug(true)

	rtm := api.NewRTM()
	go rtm.ManageConnection()

	// Load the questionBot service
	q, e := questionBot.NewService(questionBot.Config{DataPath: "./questionnaires.toml"})
	if e != nil {
		log.Fatal(e)
	}
	qbot = q

Loop:
	for {
		select {
		case msg := <-rtm.IncomingEvents:
			fmt.Print("Event Received: ")
			switch ev := msg.Data.(type) {
			case *slack.HelloEvent:
				// Ignore hello

			case *slack.ConnectedEvent:
				info := rtm.GetInfo()
				fmt.Println("Connection counter:", ev.ConnectionCount)
				rtm.SendMessage(rtm.NewOutgoingMessage(info.User.Name+" reporting for duty!", ev.Info.Channels[0].ID))

			case *slack.MessageEvent:
				info := rtm.GetInfo()
				prefix := fmt.Sprintf("<@%s>: ", info.User.ID)
				log.Println("prefix: ", prefix, " msg: ", ev.Text)
				if ev.User != info.User.ID && strings.HasPrefix(ev.Text, prefix) {
					respond(rtm, ev, prefix)
					fmt.Printf("Message: %v\n", ev)
				}

			case *slack.PresenceChangeEvent:
				fmt.Printf("Presence Change: %v\n", ev)

			case *slack.LatencyReport:
				fmt.Printf("Current latency: %v\n", ev.Value)

			case *slack.RTMError:
				fmt.Printf("Error: %s\n", ev.Error())

			case *slack.InvalidAuthEvent:
				fmt.Printf("Invalid credentials")
				break Loop

			default:

				// Ignore other events..
				// fmt.Printf("Unexpected: %v\n", msg.Data)
			}
		}
	}
}

func respond(rtm *slack.RTM, msg *slack.MessageEvent, prefix string) {
	var response string
	text := msg.Text
	text = strings.TrimPrefix(text, prefix)
	text = strings.TrimSpace(text)
	text = strings.ToLower(text)
	switch text {
	case "help":
		response = fmt.Sprintf("Available questionnaires: \n\n%s", qbot.Questionnaires.AvailableQuestionnaires())
	case "status":
		response = "status is.. YOU ARE AWESOME!"
	case "":
	default:
		response = fmt.Sprintf("I didn't understand your request. Available questionnaires: \n\n%s", qbot.Questionnaires.AvailableQuestionnaires())
	}
	omsg := rtm.NewOutgoingMessage(response, msg.Channel)
	rtm.SendMessage(omsg)
}
