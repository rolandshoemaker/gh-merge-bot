package main

import (
	"net/http"
)

const (
	pullAPIEndoint        = "repos/%s/pulls/%d"
	masterRefsAPIEndpoint = "repos/%s/git/refs/heads/master"
	statusAPIEndpoint     = "repos/%s/commits/%s/status"
)

// this will probably need to actually contain a bunch more info
// to pass through to the merge script...
type mergeRequest struct {
	reviewers []string
	prNum     int
	head      string
	client    *http.Client
}

type merger struct {
	mergeRequests chan *mergeRequest

	mergeScript string

	numReviewersRequired int
	reviewerLabelPrefix  string
	reviewers            map[string]struct{}
	labelReviewOverride  map[string]int // if matching label is present this num. reviewers required instead

	ghUser  string
	ghToken string
	ghRepo  string

	whSrvAddr string
	ghAPIBase string
}

func newMerger() *merger {
	return &merger{}
}

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
