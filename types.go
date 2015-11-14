package main

import (
	"net/http"
)

const (
	pullAPIEndoint        = "repos/%s/pulls/%d"
	masterRefsAPIEndpoint = "repos/%s/git/refs/heads/master"
	statusAPIEndpoint     = "repos/%s/commits/%s/status"
)

type config struct {
	ReviewLabelPrefix string   `yaml:"reviewLabelPrefix"`
	Reviewers         []string `yaml:"reviewers"`

	NumReviewsRequired  int            `yaml:"numReviewsRequired"`
	NumReviewsOverrides map[string]int `yaml:"numReviewsOverrides,flow"`

	MergeScript string `yaml:"mergeScript"`

	Github struct {
		User    string `yaml:"user"`
		Token   string `yaml:"token"`
		Repo    string `yaml:"repo"`
		APIBase string `yaml:"apiBase"`
	} `yaml:"github"`

	WebhookServerAddr string `yaml:"webhookServerAddr"`
}

// this will probably need to actually contain a bunch more info
// to pass through to the merge script...
type mergeRequest struct {
	reviewers []string
	prNum     int
	head      string
}

type merger struct {
	mergeRequests chan *mergeRequest

	mergeScript []string

	numReviewersRequired int
	reviewerLabelPrefix  string
	reviewers            map[string]struct{}
	labelReviewOverride  map[string]int // if matching label is present this num. reviewers required instead

	ghUser    string
	ghToken   string
	ghRepo    string
	ghAPIBase string

	whSrvAddr string
	whSecret  []byte

	client *http.Client
}

func newMerger(mergeScript []string, numReviewers int, reviewPrefix string, reviewers []string, overrides map[string]int, ghUser, ghToken, ghRepo, ghAPIBase string, whSrvAddr string, whSecret []byte) *merger {
	reviewerMap := make(map[string]struct{})
	for _, r := range reviewers {
		reviewerMap[r] = struct{}{}
	}
	return &merger{
		client:               new(http.Client),
		mergeScript:          mergeScript,
		numReviewersRequired: numReviewers,
		reviewerLabelPrefix:  reviewPrefix,
		reviewers:            reviewerMap,
		labelReviewOverride:  overrides,
		ghUser:               ghUser,
		ghToken:              ghToken,
		ghAPIBase:            ghAPIBase,
		whSrvAddr:            whSrvAddr,
		whSecret:             whSecret,
	}
}

// GitHub API types (many are minimal since we don't require all of the
// information provided by the API)

type whPayload struct {
	Action string `json:"action"`
	Issue  struct {
		Number  int    `json:"number"`
		APIURL  string `json:"url"`
		HTMLURL string `json:"html_url"`
		State   string `json:"state"`
		Labels  []struct {
			Name string `json:"name"`
		} `json:"labels"`
	} `json:"issue"`
	Comment struct {
		Body string `json:"body"`
		User struct {
			Login string `json:"login"`
		} `json:"user"`
	} `json:"comment"`
}

type masterPayload struct {
	Object struct {
		SHA string `json:"sha"`
	} `json:"object"`
}

type prPayload struct {
	Mergeable bool `json:"mergeable"`
	Head      struct {
		SHA string `json:"sha"`
	} `json:"head"`
	Base struct {
		Ref string `json:"ref"`
		SHA string `json:"sha"`
	} `json:"base"`
}

type statusesPayload struct {
	State string `json:"state"`
}
