package main

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/hex"
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

	if r.Method != "POST" {
		fmt.Printf("Invalid request method: %s\n", r.Method)
		return
	}

	// Get signature
	githubSignature := r.Header.Get("X-Hub-Signature")
	if githubSignature == "" {
		fmt.Println("No signature on request")
		return
	}

	// Verify signature
	mac := hmac.New(sha1.New, m.whSecret)
	mac.Write(body)
	expectedMAC := mac.Sum(nil)
	if len(githubSignature) <= 5 {
		fmt.Println("Invalid signature on request, no actual signature")
		return
	}
	sigBytes, err := hex.DecodeString(githubSignature[5:])
	if err != nil {
		fmt.Printf("Invalid signature on request, %s", err)
		return
	}
	if match := hmac.Equal(sigBytes, expectedMAC); !match {
		fmt.Printf("Invalid signature on request, provided: %x, expected: sha1=%x", githubSignature, expectedMAC)
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
