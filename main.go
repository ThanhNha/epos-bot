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

func ReadConfig() (tgbotapi.Chat, error) {
	viper.SetConfigName("phab-bot")
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

func ConnectToPhab() (*gonduit.Conn, error) {
	phabricatorClient, err := gonduit.Dial(viper.GetString("phabricator.url"), &core.ClientOptions{
		APIToken: viper.GetString("phabricator.token"),
		Timeout:  time.Second * 20,
	})
	if err != nil {
		log.Fatalf("Error connecting to phabricator, %s", err)
	}
	return phabricatorClient, nil
}

func main() {

	// read all config

	chat, err := ReadConfig()

	//Connect to Phabricator
	if err != nil {
		log.Fatalf("Error getting chat with id %v, %s", viper.GetInt64("telegram.chat_id"), err)
	}

	log.Printf("Will send messages to %#v", chat.Title)

	// init phabricator bot telegram
	phab_bot.Init(telegramClient)
}
