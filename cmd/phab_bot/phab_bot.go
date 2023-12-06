package phab_bot

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"epos-bot/pkg/phab_bot"
	"epos-bot/pkg/util"

	tgbotapi "github.com/musianisamuele/telegram-bot-api"
	"github.com/spf13/viper"
	"github.com/uber/gonduit"
)

var telegramClient *tgbotapi.BotAPI

var phabricatorClient *gonduit.Conn

func FeedActivity(telegramClient *tgbotapi.BotAPI) {

	notifyTypes := viper.GetStringSlice("telegram.notify_types")

	TopicForTypes := viper.GetStringMapString("telegram.topic_id_mapping")

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
		flag := false
		var CurrentTime = time.Now()
		hour := CurrentTime.Hour()
		minute := CurrentTime.Minute()
		//reset flag
		if hour != 9 && hour != 18 && flag {
			flag = false
		}

		// Check Condition to excute func send report
		if !flag && hour == 9 && minute == 20 || hour == 9 && minute == 25 {

			result, _ := SendReportRevisions(telegramClient)
			//sent, active flag to prevent send
			flag = result

		}

		feedItems, err := phab_bot.FetchFeed(viper.GetString("phabricator.url"), viper.GetString("phabricator.token"))

		if err != nil {
			log.Fatalf("Error fetching feed, %s", err)
		}

		var limit = 0

		for _, v := range feedItems {

			if !notifyTypesMap[v.Type] || v.TimeData.Before(lastMsgTime) || v.TimeData == lastMsgTime || v.IsClose || v.Status == "changes-planned" {
				continue
			}

			text := phab_bot.PrepareMessage(v)

			topicID, err := strconv.Atoi(TopicForTypes[strings.ToLower(v.Type)])

			msg := tgbotapi.NewThreadMessage(viper.GetInt64("telegram.chat_id"), topicID, text)
			msg.ParseMode = "HTML"
			msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(
				tgbotapi.NewInlineKeyboardRow(
					tgbotapi.NewInlineKeyboardButtonURL(v.ShortTitle, v.URL),
				),
			)

			_, err = telegramClient.Send(msg)

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

func SendReportRevisions(telegramClient *tgbotapi.BotAPI) (bool, error) {

	tableContent, _ := phab_bot.GetListRevisionsOfWeek(viper.GetString("phabricator.url"), viper.GetString("phabricator.token"))

	err := phab_bot.CreateHtmlFile(tableContent)

	if err != nil {
		fmt.Println("Error when get file:", err)
		return false, err
	}

	// Start get file to send telegram

	filePath := "./static/daily-report-revisions.html"
	file, err := os.Open(filePath)

	if err != nil {
		log.Fatal(err)
	}

	defer file.Close()

	fileInfo, err := file.Stat()
	if err != nil {
		log.Fatal(err)
	}

	filePart := tgbotapi.FileBytes{
		Name:  fileInfo.Name(),
		Bytes: util.GetFileBytes(file),
	}

	revisions := tgbotapi.DocumentConfig{
		BaseFile: tgbotapi.BaseFile{
			BaseChat: tgbotapi.BaseChat{ChatID: viper.GetInt64("telegram.chat_id"), MessageThreadID: 67, ReplyToMessageID: 0},
			File:     filePart,
		},
		Caption: "List revisions on active.",
	}

	_, err = telegramClient.Send(revisions)

	if err != nil {
		log.Printf("Sending message errir %v", err)
	}

	return true, nil

}
