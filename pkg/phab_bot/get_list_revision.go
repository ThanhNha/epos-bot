package phab_bot

import (
	"epos-bot/pkg/util"
	"fmt"
	"log"
	"time"

	"github.com/uber/gonduit"
	"github.com/uber/gonduit/core"
	"github.com/uber/gonduit/requests"
)

var phabricatorClient *gonduit.Conn

func GetListRevisionsOfWeek(PhabricatorUrl string, PhabricatorToken string) ([]TableContent, error) {

	// Connect to Phabricator
	var err error

	phabricatorClient, err = gonduit.Dial(PhabricatorUrl, &core.ClientOptions{
		APIToken: PhabricatorToken,
		Timeout:  time.Second * 20,
	})

	if err != nil {
		log.Fatalf("Error connecting to phabricator, %s", err)
	}

	diff, err := phabricatorClient.DifferentialRevisionSearch(requests.DifferentialRevisionSearchRequest{
		QueryKey: "active",
		Attachments: &requests.DifferentialRevisionSearchAttachments{
			Projects: true,
		},
	})

	if err != nil {
		return nil, fmt.Errorf("error fetching feed, %s", err)
	}

	var phids []string

	for _, item := range diff.Data {
		phids = append(phids, item.ResponseObject.PHID, item.Fields.AuthorPHID)
	}
	phids = util.RemoveDuplicates(phids)

	lookedUpPhids, _ := phabricatorClient.PHIDLookup(requests.PHIDLookupRequest{
		Names: phids,
	})

	var RevisionList []TableContent

	for _, item := range diff.Data {
		RevisionList = append(RevisionList, TableContent{
			Name:   item.Fields.Title,
			Author: lookedUpPhids[item.Fields.AuthorPHID].Name,
			URL:    lookedUpPhids[item.ResponseObject.PHID].URI,
			Status: item.Fields.Status.Name,
		})

	}
	return RevisionList, nil
}
