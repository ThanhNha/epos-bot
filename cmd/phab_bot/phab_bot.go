package phab_bot

import (
	"log"
	"time"

	"epos-bot/pkg/phab_bot"

	tgbotapi "github.com/musianisamuele/telegram-bot-api"
	"github.com/spf13/viper"
)

var telegramClient *tgbotapi.BotAPI

func Init(telegramClient *tgbotapi.BotAPI) {
	// // //send flie to tele

	// filePath := "static/example.html"
	// file, err := os.Open(filePath)

	// if err != nil {
	// 	log.Fatal(err)
	// }

	// defer file.Close()

	// fileInfo, err := file.Stat()
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// filePart := tgbotapi.FileBytes{
	// 	Name:  fileInfo.Name(),
	// 	Bytes: util.GetFileBytes(file),
	// }

	// revisions := tgbotapi.DocumentConfig{
	// 	BaseFile: tgbotapi.BaseFile{
	// 		BaseChat: tgbotapi.BaseChat{ChatID: viper.GetInt64("telegram.chat_id"), MessageThreadID: 67, ReplyToMessageID: 0},
	// 		File:     filePart,
	// 	},
	// 	Caption: "List revisions is active.",
	// }

	// _, err = telegramClient.Send(revisions)

	notifyTypes := viper.GetStringSlice("telegram.notify_types")

	if len(notifyTypes) == 0 {
		notifyTypes = []string{"TASK", "DREV", "WIKI"}
		log.Printf("No notify types specified, defaulting to %v", notifyTypes)
	}

	notifyTypesMap := make(map[string]bool)
	for _, v := range notifyTypes {
		notifyTypesMap[v] = true
	}

	var lastMsgTime = time.Now()

	for {
		feedItems, err := phab_bot.FetchFeed(viper.GetString("phabricator.url"), viper.GetString("phabricator.token"))

		// phab_bot.CreateHtmlFile()

		if err != nil {
			log.Fatalf("Error fetching feed, %s", err)
		}
		// log.Printf("Fetched feed, got %v items", len(feedItems))

		var limit = 0

		for _, v := range feedItems {

			if !notifyTypesMap[v.Type] || v.TimeData.Before(lastMsgTime) || v.TimeData == lastMsgTime || v.IsClose || v.Status == "changes-planned" {
				continue
			}

			text := phab_bot.PrepareMessage(v)

			msg := tgbotapi.NewThreadMessage(viper.GetInt64("telegram.chat_id"), phab_bot.TopicForTypes[v.Type], text)
			msg.ParseMode = "HTML"
			msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(
				tgbotapi.NewInlineKeyboardRow(
					tgbotapi.NewInlineKeyboardButtonURL(v.ShortTitle, v.URL),
				),
			)

			_, err := telegramClient.Send(msg)

			if err != nil {
				log.Printf("Sending message errir %v", err)
			}

			lastMsgTime = v.TimeData
			limit++
			if limit > 10 {
				log.Printf("Limit reached, stopping for 5 seconds")
				time.Sleep(time.Second * 5)
				limit = 0
			}

		}

		time.Sleep(viper.GetDuration("phabricator.poll_interval"))
	}
}
