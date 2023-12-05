package main

import (
	"epos-bot/cmd/phab_bot"
	"log"
	"time"

	tgbotapi "github.com/musianisamuele/telegram-bot-api"
	"github.com/spf13/viper"
	"github.com/uber/gonduit"
	"github.com/uber/gonduit/core"
)

var telegramClient *tgbotapi.BotAPI

var phabricatorClient *gonduit.Conn

func ReadConfig() (tgbotapi.Chat, error) {
	viper.SetConfigName("config")
	viper.SetConfigType("json")
	viper.AddConfigPath("./configs")
	err := viper.ReadInConfig()
	if err != nil {
		log.Fatalf("Error reading config file, %s", err)
	}

	telegramClient, err = tgbotapi.NewBotAPI(viper.GetString("telegram.token"))

	if err != nil {
		log.Fatalf("Error connecting to telegram, %s", err)
	}

	self, err := telegramClient.GetMe()

	if err != nil {
		log.Fatalf("Error getting telegram bot info, %s", err)
	}

	log.Printf("Authorized on account %s", self.UserName)

	chat, err := telegramClient.GetChat(tgbotapi.ChatInfoConfig{
		ChatConfig: tgbotapi.ChatConfig{
			ChatID: viper.GetInt64("telegram.chat_id"),
		},
	})

	return chat, err

}

func initPhabricatorClient() error {
	var err error
	phabricatorClient, err = gonduit.Dial(viper.GetString("phabricator.url"), &core.ClientOptions{
		APIToken: viper.GetString("phabricator.token"),
		Timeout:  time.Second * 20,
	})
	return err
}

func main() {

	// read all config

	chat, err := ReadConfig()

	//Connect to Phabricator
	if err != nil {
		log.Fatalf("Error getting chat with id %v, %s", viper.GetInt64("telegram.chat_id"), err)
	}

	log.Printf("Will send messages to %#v", chat.Title)

	// Initialize the phabricatorClient variable
	err = initPhabricatorClient()
	if err != nil {
		log.Fatalf("Error initializing phabricator client: %s", err)
	}

	// init phabricator bot telegram
	phab_bot.FeedActivity(telegramClient)

}
