package main

import (
	"fmt"
	"os"

	"github.com/davecgh/go-spew/spew"
	"github.com/nlopes/slack"
)

const tokenKey = "SLACKBOT_TOKEN"

func main() {
	token := os.Getenv(tokenKey)
	api := slack.New(token)
	//api.SetDebug(true)

	rtm := api.NewRTM()
	go rtm.ManageConnection()

Loop:
	for {
		select {
		case msg := <-rtm.IncomingEvents:
			spew.Dump(msg)
			fmt.Print("Event Received: ")
			switch ev := msg.Data.(type) {
			case *slack.HelloEvent:
				// Ignore hello

			case *slack.ConnectedEvent:
				spew.Dump(ev.Info)
				fmt.Println("Connection counter:", ev.ConnectionCount)
				rtm.SendMessage(rtm.NewOutgoingMessage("Slackbot reporting for duty!", ev.Info.Channels[0].ID))

			case *slack.MessageEvent:
				respond(rtm, ev)
				fmt.Printf("Message: %v\n", ev)

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

func respond(rtm *slack.RTM, msg *slack.MessageEvent) {
	omsg := rtm.NewOutgoingMessage("what?  stop bothering me!", msg.Channel)
	rtm.SendMessage(omsg)

}
