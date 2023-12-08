package phab_bot

import (
	"fmt"
	"log"
	"sort"
	"time"

	"epos-bot/pkg/util"

	"github.com/spf13/viper"
	"github.com/uber/gonduit"
	"github.com/uber/gonduit/core"
	"github.com/uber/gonduit/requests"
	"github.com/uber/gonduit/responses"
)

// fetchFeed calls feed.query and then uses PHIDLookup to get the actual data for each feed item.
func FetchFeed(PhabricatorUrl string, PhabricatorToken string) ([]FeedItem, error) {
	var feed map[string]FeedQueryResponseItem
	var err error
	req := &FeedQueryRequest{
		After: "",
		View:  "text",
		Limit: 10,
	}

	phabricatorClient, err = gonduit.Dial(PhabricatorUrl, &core.ClientOptions{
		APIToken: PhabricatorToken,
		Timeout:  time.Second * 20,
	})
	if err != nil {
		log.Fatalf("Error connecting to phabricator, %s", err)
	}

	err = phabricatorClient.Call("feed.query", req, &feed)
	if err != nil {
		return nil, fmt.Errorf("error fetching feed, %s", err)
	}
	// transpose to a list and sort by epoch
	feedList := make([]FeedQueryResponseItem, len(feed))
	i := 0
	for _, v := range feed {

		feedList[i] = v
		i++
	}

	sort.Slice(feedList, func(i, j int) bool {
		return feedList[i].Epoch < feedList[j].Epoch
	})

	phids := make([]string, 0, len(feedList)*2)

	for _, v := range feedList {
		phids = append(phids, v.AuthorPHID, v.ObjectPHID)

	}

	lookedUpPhids, err := LookupID(phabricatorClient, phids)

	if err != nil {
		return nil, fmt.Errorf("error looking up phids, %s", err)
	}

	feedItems := make([]FeedItem, len(feed))

	for i, v := range feedList {

		diff, _ := getRevisionInfo(v.ObjectPHID)

		feedItems[i] = FeedItem{
			URL:            lookedUpPhids[v.ObjectPHID].URI,
			Title:          lookedUpPhids[v.ObjectPHID].FullName,
			ShortTitle:     lookedUpPhids[v.ObjectPHID].Name,
			Time:           time.Unix(int64(v.Epoch), 0).Format(time.RFC1123),
			Author:         lookedUpPhids[v.AuthorPHID].FullName,
			AuthorRevision: diff.Author,
			Status:         diff.Status,
			IsClose:        diff.IsClose,
			Type:           lookedUpPhids[v.ObjectPHID].Type,
			TypeName:       lookedUpPhids[v.ObjectPHID].TypeName,
			TimeData:       time.Unix(int64(v.Epoch), 0),
			Reviewers:      diff.ListReviewers,
			Projects:       diff.ListProject,
			Mentions:       diff.MentionTele,
			Text:           v.Text,
		}
	}

	return feedItems, nil

}

func FetchFeedItems() ([]FeedItem, error) {
	return FetchFeed(viper.GetString("phabricator.url"), viper.GetString("phabricator.token"))
}

func LookupID(phabricatorClient *gonduit.Conn, phids []string) (responses.PHIDLookupResponse, error) {

	phids = util.RemoveDuplicates(phids)

	lookedUpPhids, err := phabricatorClient.PHIDLookup(requests.PHIDLookupRequest{
		Names: phids,
	})

	if err != nil {
		return nil, fmt.Errorf("error looking up phids, %s", err)
	}

	return lookedUpPhids, nil
}
