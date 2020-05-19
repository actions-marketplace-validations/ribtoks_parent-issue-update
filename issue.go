package main

import (
	"errors"
	"fmt"

	"github.com/google/go-github/v31/github"
)

type IssueStatus int

const (
	StatusOpened IssueStatus = iota
	StatusClosed
)

var (
	errIssueNotFound = errors.New("issue not found")
)

type Issue struct {
	ID       int
	Title    string
	Body     string
	Status   IssueStatus
	Children []*Issue
}

func (i *Issue) ToMap() map[int]*Issue {
	issueMap := make(map[int]*Issue)
	issueMap[i.ID] = i

	for _, ci := range i.Children {
		ciMap := ci.ToMap()
		for k, v := range ciMap {
			issueMap[k] = v
		}
	}

	return issueMap
}

func (i *Issue) FormatTitle(spaces int) string {
	status := " "
	if i.Status == StatusClosed {
		status = "x"
	}

	prefix := make([]rune, spaces)
	for i := range prefix {
		prefix[i] = ' '
	}

	return fmt.Sprintf("%s- [%s] %s #%v", string(prefix), status, i.Title, i.ID)
}

func NewIssue(i *github.Issue) *Issue {
	issue := &Issue{
		ID:     i.GetNumber(),
		Title:  i.GetTitle(),
		Body:   i.GetBody(),
		Status: StatusOpened,
	}

	if i.GetState() == "closed" {
		issue.Status = StatusClosed
	}
	return issue
}