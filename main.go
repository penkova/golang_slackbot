package main

import (
	"encoding/json"
	"fmt"
	"github.com/nlopes/slack"
	"io/ioutil"
	"log"
	"strings"
)

var (
	BotKey        Token
	SlackClient   *slack.Client
	connectedUser *slack.UserDetails
	welcPref      = []string{"hi", "Hi", "hello", "Hello", "Howdy", "Wassup", "Hey", "Привет", "Здравствуйте"}
)

type Token struct {
	Token string `json:"token"`
}

//reading token for my bot of token.json
func init() {
	file, err := ioutil.ReadFile("./token.json")
	if err != nil {
		log.Fatal("File doesn't exist")
	}
	if err := json.Unmarshal(file, &BotKey); err != nil {
		log.Fatal("Cannot parse token.json")
	}
}

func main() {
	// New создает Slack клиента из предоставленного токена и опций.
	SlackClient = slack.New(BotKey.Token)
	//NewRTM возвращает RTM, который обеспечивает полностью управляемое соединение с протоколом Real-Time Messaging от Slack на основе веб-приложений
	rtm := SlackClient.NewRTM()

	//ManageConnection мы подключаемся к Slack RTM API и будет обрабатывать все входящие и исходящие события.
	go rtm.ManageConnection()

	Run(rtm)
}

func Run(rtm *slack.RTM) int {
	for {
		select {
		case msg := <-rtm.IncomingEvents:
			switch ev := msg.Data.(type) {

			//ConnectedEvent используется, когда мы подключаемся к Slack
			case *slack.ConnectedEvent:
				connectedUser = ev.Info.User
				log.Printf("[INFO] Connected: user_id=%s name=%s \n", connectedUser.ID, connectedUser.Name)

			//HelloEvent представляет событие hello
			case *slack.HelloEvent:
				log.Print("Hello Event")

			//MessageEvent представляет Slack Message (используется как тип события для входящего сообщения)
			case *slack.MessageEvent:
				handleMessageEvent(ev, rtm)

			//InvalidAuthEvent используется в случае, если мы не можем даже аутентифицироваться с помощью API
			case *slack.InvalidAuthEvent:
				log.Print("Invalid credentials")
				return 1
			}
		}
	}
}

func handleMessageEvent(ev *slack.MessageEvent, rtm *slack.RTM) {
	channelInfo, err := SlackClient.GetChannelInfo(ev.Channel)
	if err != nil {
		//log.Fatalln(err)
	}

	fmt.Println("chanel info", channelInfo)
	fmt.Println("chanel info 1", ev.Channel)

	// поменять название функции и написать проверку в каком канале нахожусь
	reactionBotForClient(ev, rtm)

}

func reactionBotForClient(ev *slack.MessageEvent, rtm *slack.RTM) {
	botId := connectedUser.ID
	receivedText := ev.Text
	checkPrefBot := strings.HasPrefix(receivedText, "<@"+botId+">")

	clientInfo, err := SlackClient.GetUserInfo(ev.User)
	if err != nil {
		log.Fatalln(err)
	}

	if checkPrefBot == true {
		for _, val := range welcPref {
			checkRecivedTxt := strings.Contains(receivedText, val)
			if checkRecivedTxt == true {
				log.Printf("Message: %v\n", ev)
				rtm.SendMessage(rtm.NewOutgoingMessage("Hello, "+"<@"+clientInfo.ID+">!", ev.Channel))
			} else {
				continue
			}
		}
		if receivedText != ""{
			rtm.SendMessage(rtm.NewOutgoingMessage("I'm sorry. I don't know it. I`m your friend. I can say \"Hello\" for you.", ev.Channel))
		}
	}
}
