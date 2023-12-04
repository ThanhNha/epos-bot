package phab_bot

import (
	"epos-bot/pkg/util"
	"fmt"

	"github.com/uber/gonduit"
	"github.com/uber/gonduit/requests"
)

var phabricatorClient *gonduit.Conn

func GetListRevisionsOfWeek() ([]TableContent, error) {
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
		})

	}
	return RevisionList, nil
}
