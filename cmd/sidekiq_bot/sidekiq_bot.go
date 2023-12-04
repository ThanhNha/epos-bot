package sidekiq_bot

import (
	"context"
	"crypto/tls"
	"fmt"
	"log"
	"time"

	tgbotapi "github.com/musianisamuele/telegram-bot-api"
	"github.com/redis/go-redis/v9"
	"github.com/spf13/viper"
)

var ctx = context.Background()

var telegramClient *tgbotapi.BotAPI

// The function that send message to telegram
func send_message_to_telegram(queue_name_status string) {
	// Code for sending telegram the message

	telegramClient, _ = tgbotapi.NewBotAPI(viper.GetString("telegram.token"))

	// mapReviewerMentionTele()
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
	if err != nil {
		log.Printf("Getting chat error %v", err)
	}
	log.Printf("Will send messages to %#v", chat.Title)

	msg := tgbotapi.NewThreadMessage(viper.GetInt64("telegram.chat_id"), viper.GetInt("telegram.sidekiq_topic_id"), queue_name_status)
	msg.ParseMode = "HTML"

	_, err = telegramClient.Send(msg)
	if err != nil {
		log.Printf("Sending message error %v", err)
	}
}

func main() {

	viper.SetConfigName("config")
	viper.SetConfigType("json")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig()
	if err != nil {
		log.Fatalf("Error reading config file, %s", err)
	}

	redis_client := redis.NewClient(&redis.Options{
		TLSConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
		Addr:     viper.GetString("redis.host"),
		Password: viper.GetString("redis.password"), // no password set
		DB:       0,                                 // use default DB
	})

	// url := viper.GetString("redis.url")
	// opts, err := redis.ParseURL(url)
	// if err != nil {
	// 	panic(err)
	// }
	// redis_client := redis.NewClient(opts)

	queues, err := redis_client.SMembers(ctx, "queues").Result()
	if err != nil {
		fmt.Println(err)
		return
	}

	for {
		msg := ""
		for _, queue_name := range queues {
			queue_size, err := redis_client.LLen(ctx, fmt.Sprintf("queue:%s", queue_name)).Result()
			if err != nil {
				fmt.Println(err)
				continue
			}

			if queue_size >= 10 {
				msg += fmt.Sprintf("%s: %d\n", queue_name, int16(queue_size))
			}
		}
		if msg != "" {
			msg = fmt.Sprintf("<b>SIDEKIQ STATUS</b>\n%s", msg)
			send_message_to_telegram(msg)
		}
		msg = ""

		time.Sleep(viper.GetDuration("redis.poll_interval"))
	}
}
