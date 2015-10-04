package main

import (
	"fmt"
)

type PullRequests []*Pull

func (slice PullRequests) Len() int {
	return len(slice)
}

func (slice PullRequests) Less(i, j int) bool {
	return slice[i].Title < slice[j].Title
}

func (slice PullRequests) Swap(i, j int) {
	slice[i], slice[j] = slice[j], slice[i]
}

type Pull struct {
	Repo    string `json:"-"`
	Id      int    `json:"id"`
	HtmlUrl string `json:"html_url"`
	Title   string `json:"title"`
}

func (p Pull) String(selection int) string {
	return fmt.Sprintf("[%d] (%s) %s", selection, p.Repo, p.Title)
}
