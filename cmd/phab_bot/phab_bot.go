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

	fmt.Println("feedItems")

	notifyTypes := phab_bot.GetNotifyTypes()

	notifyTypesMap := phab_bot.CreateNotifyTypesMap(notifyTypes)

	lastMsgTime := time.Now()
	feedItems := make([]phab_bot.FeedItem,0,1000)
	for {
		var flag bool
		currentTime := time.Now()
		hour := currentTime.Hour()

		util.ResetFlag(&flag, hour)

		CheckAndSendReport(&flag, hour, telegramClient)

		phab_bot.FetchFeedItems(feedItems)

		ProcessFeedItems(feedItems, notifyTypesMap, &lastMsgTime)

		pollInterval := viper.GetDuration("phabricator.poll_interval")

		time.Sleep(pollInterval)
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

func ProcessFeedItems(feedItems []phab_bot.FeedItem, notifyTypesMap map[string]bool, lastMsgTime *time.Time) {
	limit := 0

	for _, v := range feedItems {
		if !notifyTypesMap[v.Type] || v.TimeData.Before(*lastMsgTime) || v.TimeData.Equal(*lastMsgTime) || v.IsClose || v.Status == "changes-planned" {
			continue
		}

		text := phab_bot.PrepareMessage(v)
		topicIDRaw, ok := viper.GetStringMapString("telegram.topic_id_mapping")[strings.ToLower(v.Type)]
		if !ok {
			log.Printf("Topic ID not found for type: %s", v.Type)
			continue
		}

		topicID, _ := strconv.Atoi(topicIDRaw)

		err := SendMessageTele(viper.GetInt64("telegram.chat_id"), topicID, text, v)
		if err != nil {
			log.Printf("Error sending message: %v", err)
		}

		*lastMsgTime = v.TimeData
		limit++
		if limit > 10 {
			log.Printf("Limit reached, pausing for 5 seconds")
			time.Sleep(time.Second * 5)
			limit = 0
		}
	}
}

func CheckAndSendReport(flag *bool, hour int, telegramClient *tgbotapi.BotAPI) {
	if !*flag && (hour == 9 || hour == 18) {
		result, err := SendReportRevisions(telegramClient)
		if err != nil {
			log.Printf("Error sending report revisions: %v", err)
		} else {
			*flag = result
		}
	}
}

func SendMessageTele(ChatID int64, topicID int, message string, FeedItem phab_bot.FeedItem) error {
	msg := tgbotapi.NewThreadMessage(ChatID, topicID, message)
	msg.ParseMode = "HTML"
	msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonURL(FeedItem.ShortTitle, FeedItem.URL),
		),
	)

	_, err := telegramClient.Send(msg)

	return err
}
