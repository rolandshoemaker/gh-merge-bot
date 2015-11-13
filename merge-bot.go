package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"path"
	"strings"
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
	}
}

func main() {

}
