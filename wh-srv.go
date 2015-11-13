package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

func (m *merger) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// unmarshal payload
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		// something
		return
	}
	var payload whPayload
	err = json.Unmarshal(body, &payload)
	if err != nil {
		// something
		return
	}

	// check it's the right kind of payload + comment targets merge-bot
	// + comes from a whitelisted user (reviewer?)
	if payload.Action != "created" {
		// something
		return
	}
	if !strings.Contains(payload.Issue.HTMLURL, "/pulls/") {
		// something
		return
	}
	if _, present := m.reviewers[payload.Comment.User.Login]; !present {
		// something
		return
	}
	if payload.Comment.Body != fmt.Sprintf("@%s merge", m.ghUser) {
		// something
		return
	}

	otherLabels := []string{}
	reviewers := []string{}
	for _, l := range payload.Issue.Labels {
		if strings.HasPrefix(l.Name, m.reviewerLabelPrefix) {
			name := l.Name[len(m.reviewerLabelPrefix):]
			if _, present := m.reviewers[name]; present {
				reviewers = append(reviewers, name)
			}
		} else {
			otherLabels = append(otherLabels, l.Name)
		}
	}
	if len(reviewers) < m.numReviewersRequired {
		override := 0
		for _, l := range otherLabels {
			if limit, present := m.labelReviewOverride[l]; present {
				if limit > override {
					override = limit
				}
			}
			if override > 0 {
				if len(reviewers) < override {
					// something
					return
				}
			} else {
				// something
				return
			}
		}
	}

	// strip useful info and send down action channel
	m.mergeRequests <- &mergeRequest{
		reviewers: reviewers,
		prNum:     payload.Issue.Number,
	}
}
