package phab_bot

import (
	"time"

	"github.com/uber/gonduit/requests"
)

// FeedQueryRequest is the request struct for feed.query.
type FeedQueryRequest struct {
	After            string `json:"after,omitempty"`
	View             string `json:"view,omitempty"`
	Limit            uint64 `json:"limit"`
	requests.Request        // Includes __conduit__ field needed for authentication.
}

type FeedQueryResponseItem struct {
	Class            string `json:"class"`
	Epoch            int    `json:"epoch"`
	AuthorPHID       string `json:"authorPHID"`
	ChronologicalKey string `json:"chronologicalKey"`
	ObjectPHID       string `json:"objectPHID"`
	Text             string `json:"text"`
}

type FeedItem struct {
	URL              string
	Title            string
	Time             string
	Author           string
	AuthorRevision   string
	Type             string
	TypeName         string
	Reviewers        string
	Subscribers      string
	Mentions         string
	Projects         string
	Status           string
	IsClose          bool
	ChronologicalKey string
	Text             string `json:"text"`
	ShortTitle       string
	TimeData         time.Time
}

type RevisionMessage struct {
	Title           string
	Author          string
	Status          string
	IsClose         bool
	ListReviewers   string
	ListSubscribers string
	ListProject     string
	MentionTele     string
}

type RevisionQueryResponseItem struct {
	AuthorPHID string `json:"authorPHID"`
	ObjectPHID string `json:"objectPHID"`
}

type Revision struct {
	Id      string
	Name    string
	Author  string
	Status  string
	URL     string
	Project string
}

type ConfigMapMention struct {
	Key   string
	Value string
}

var EmojiForTypes = map[string]string{
	"TASK": "🎯",
	"WIKI": "📖",
	//"USER": "",
	//"CEVT": "📅",
	//"PROJ": "ℹ️",
	"DREV": "🚀",
}

var TopicForTypes = map[string]int{
	"WIKI": 244,
	"TASK": 67,
	"DREV": 2,
}

var EmojiForActions = map[string]string{

	"created":  "\U0001F4A1",
	"added":    "\U0001F4AC",
	"lowered":  "\U0001F53B",
	"raised":   "\U0001F53A",
	"awarded":  "\U0001F3C6",
	"triaged":  "\U0001F4AD",
	"updated":  "\U0001F449",
	"changed":  "\U0000270F\U0000FE0F ",
	"claimed":  "\U0001F44C",
	"set":      "\U0000270F\U0000FE0F ",
	"reopened": "\U0001F504",
	"closed":   "\U0001F510",
	"renamed":  "\U0001F449",
	"edited":   "\U0001F4DD",
	"Review":   "\U0001F4DD",
	"Revision": "\U0001F504",
	"Planned":  "\U0001F4A1",
	"Accepted": "\U00002705",
}
