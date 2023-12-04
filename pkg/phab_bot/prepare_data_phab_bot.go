package phab_bot

import (
	"fmt"
	"strings"

	"github.com/spf13/viper"
	"github.com/uber/gonduit/requests"
)

func getRevisionInfo(phid string) (RevisionMessage, error) {

	var revisionItem RevisionMessage
	var Mention string
	var ListReviewer []string
	var AuthorText []string

	if !strings.Contains(phid, "DREV") {
		return revisionItem, nil
	}

	diff, err := phabricatorClient.DifferentialRevisionSearch(requests.DifferentialRevisionSearchRequest{
		Constraints: &requests.DifferentialRevisionSearchConstraints{
			PHIDs: []string{phid},
		},
		Attachments: &requests.DifferentialRevisionSearchAttachments{
			Reviewers: true,
			Projects:  true,
		},
	})

	if err != nil {
		return revisionItem, fmt.Errorf("error fetching feed, %s", err)
	}

	for _, item := range diff.Data {

		var StatusRevision string

		var IsCloseRevision bool

		StatusRevision = item.Fields.Status.Value

		IsCloseRevision = item.Fields.Status.Closed

		// don't need continue when is closed revision or don't need review
		if IsCloseRevision && StatusRevision == "changes-planned" {
			return revisionItem, nil
		}

		// get Reviewers
		var ReviewersIDs []string
		for _, id := range item.Attachments.Reviewers.Reviewers {
			ReviewersIDs = append(ReviewersIDs, id.ReviewerPHID)
		}

		// get Project Name
		var ProjectIDs []string
		for _, id := range item.Attachments.Projects.ProjectPHIDs {
			ProjectIDs = append(ProjectIDs, id)
		}

		ProjectList, err := getProjectInfo(ProjectIDs)

		if len(ReviewersIDs) != 0 {
			ListReviewer, err = getUserInfo(ReviewersIDs)
		}

		if (item.Fields.AuthorPHID) != "" {
			AuthorText, err = getUserInfo([]string{item.Fields.AuthorPHID})
		}

		if len(AuthorText) == 0 {
			AuthorText[0] = ""
		}

		User := strings.Join(ListReviewer, ", ")

		Mention = mapReviewerMentionTele(ListReviewer, AuthorText[0], StatusRevision)

		if err != nil {
			return revisionItem, fmt.Errorf("error fetching feed, %s", err)
		}

		revisionItem = RevisionMessage{
			Title:         item.Fields.Title,
			Author:        AuthorText[0],
			Status:        StatusRevision,
			IsClose:       IsCloseRevision,
			ListReviewers: User,
			ListProject:   ProjectList,
			MentionTele:   Mention,
		}

	}
	return revisionItem, nil

}
func getUserInfo(IDs []string) ([]string, error) {

	if len(IDs) == 0 {
		return nil, nil
	}

	users, err := phabricatorClient.PHIDLookup(requests.PHIDLookupRequest{
		Names: IDs,
	})

	if err != nil {
		return nil, err
	}

	var NameUsers []string
	for _, name := range users {

		NameUsers = append(NameUsers, name.Name)

	}

	return NameUsers, nil

}
func getProjectInfo(IDs []string) (string, error) {

	if len(IDs) == 0 {
		return "", nil
	}

	var ProjectNameText []string

	projects, err := phabricatorClient.PHIDLookup(requests.PHIDLookupRequest{
		Names: IDs,
	})

	for _, project := range projects {
		ProjectNameText = append(ProjectNameText, project.Name)
	}

	if err != nil {
		return "", err
	}
	Project := strings.Join(ProjectNameText, ", ")

	return Project, nil

}

func mapReviewerMentionTele(Reviewer []string, Author string, Status string) string {

	if len(Reviewer) == 0 || Status == "" {
		return ""
	}

	var mention string
	var reviewerConfigMap = viper.GetStringMapString("reviewers")

	switch Status {

	case "needs-review":
		// Get config map from phab_bot.json
		for _, name := range Reviewer {
			mentionTele := reviewerConfigMap[strings.ToLower(name)]

			if mentionTele != "" {
				mention += fmt.Sprintf("@%s ", mentionTele)
			}
		}
		break

	case "request-changes", "accepted":
		mentionTele := reviewerConfigMap[strings.ToLower(Author)]
		if mentionTele != "" {
			mention += fmt.Sprintf("@%s ", mentionTele)
		}
		break

	default:
		mention = ""
	}

	return mention
}

func PrepareMessage(data FeedItem) string {
	actionEmoji := ""
	for k, e := range EmojiForActions {
		if strings.Contains(data.Text, k) {
			actionEmoji = e
			break
		}
	}

	actionText := data.Text
	actionText = strings.ReplaceAll(actionText, data.Title, data.ShortTitle)

	authorText := data.Author

	if data.Type == "DREV" {
		authorText = data.AuthorRevision
	}

	text := fmt.Sprintf("%s <b>%s</b>\n%s %s\n\n", EmojiForTypes[data.Type], data.Title, actionEmoji, actionText)

	text += fmt.Sprintf("<b>Author</b>: %s\n\n", authorText)
	if data.Reviewers != "" && data.Type == "DREV" {
		text += fmt.Sprintf("<b>Reviewers</b>: %s\n", data.Reviewers)
	}

	if data.Projects != "" {
		text += fmt.Sprintf("<b>Projects Tag</b>: %s\n", data.Projects)
	}

	text += data.Mentions

	return text
}
