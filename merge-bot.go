package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os/exec"
	"path"
	"strings"

	"gopkg.in/yaml.v2"
)

func (m *merger) Run() error {
	go m.processRequests()
	return nil
}

func (m *merger) get(endpoint string, unmarshaler interface{}) error {
	return nil
}

func (m *merger) checkMergeable(mr *mergeRequest) error {
	masterRef := masterPayload{}
	err := m.get(path.Join(m.ghAPIBase, fmt.Sprintf(masterRefsAPIEndpoint, m.ghRepo)), &masterRef)
	if err != nil {
		return err
	}

	pr := prPayload{}
	err = m.get(path.Join(m.ghAPIBase, fmt.Sprintf(pullAPIEndoint, m.ghRepo, mr.prNum)), &pr)
	if err != nil {
		return err
	}

	if !pr.Mergeable {
		// something
		return fmt.Errorf("")
	}
	if pr.Base.Ref != "master" {
		// something
		return fmt.Errorf("")
	}
	// check up to date
	if pr.Base.SHA != masterRef.Object.SHA {
		// something
		return fmt.Errorf("")
	}

	statuses := statusesPayload{}
	err = m.get(path.Join(m.ghAPIBase, fmt.Sprintf(statusAPIEndpoint, m.ghRepo, pr.Base.SHA)), &statuses)
	if err != nil {
		return err
	}
	if statuses.State != "success" {
		// something
		return fmt.Errorf("")
	}

	return nil
}

func (m *merger) merge(mr *mergeRequest) error {
	cmd := exec.Command(m.mergeScript[0], m.mergeScript[1:]...)
	stderr := new(bytes.Buffer)
	cmd.Stderr = stderr
	cmd.Env = []string{
		fmt.Sprintf("MERGE_BRANCH=%s", ""),
		fmt.Sprintf("MERGE_USER=%s", ""),
	}
	err := cmd.Run()
	if err != nil {
		// something, mb comment depending on the exit code...
		return fmt.Errorf(stderr.String())
	}
	return nil
}

func (m *merger) processRequests() {
	for mr := range m.mergeRequests {
		err := m.checkMergeable(mr)
		if err != nil {
			// something
			continue
		}
		err = m.merge(mr)
		if err != nil {
			// something
			continue
		}
		// something something
	}
}

func main() {
	cont, err := ioutil.ReadFile("example-config.yml")
	if err != nil {
		fmt.Println(err)
		return
	}
	var c config
	err = yaml.Unmarshal(cont, &c)
	if err != nil {
		fmt.Println(err)
		return
	}

	m := newMerger(
		strings.Split(c.MergeScript, " "),
		c.NumReviewsRequired,
		c.ReviewLabelPrefix,
		c.Reviewers,
		c.NumReviewsOverrides,
		c.Github.User,
		c.Github.Token,
		c.Github.Repo,
		c.Github.APIBase,
		c.WebhookServerAddr,
	)

	m.Run()
}
