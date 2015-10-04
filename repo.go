package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"sort"
	"sync"
)

func NewRepo(name string, accessToken string) *Repo {
	url := fmt.Sprintf("https://api.github.com/repos/%s/pulls?state=open&access_token=%s", name, accessToken)
	return &Repo{Name: name, URL: url}
}

type Repositories []*Repo

func (slice Repositories) Len() int {
	return len(slice)
}

func (slice Repositories) Less(i, j int) bool {
	return slice[i].Name < slice[j].Name
}

func (slice Repositories) Swap(i, j int) {
	slice[i], slice[j] = slice[j], slice[i]
}

type Repo struct {
	Name  string
	URL   string
	Pulls PullRequests
	Resp  *http.Response
}

func (r *Repo) Fetch(wg *sync.WaitGroup) {
	defer wg.Done()
	resp, err := http.Get(r.URL)
	if err != nil {
		fmt.Printf("Could not fetch PRs for %s: %s", r.Name, err.Error())
		return
	}
	r.Resp = resp

	defer r.Resp.Body.Close()
	body, err := ioutil.ReadAll(r.Resp.Body)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Could not read response: %s\n", err.Error())
		return
	}

	err = json.Unmarshal(body, &r.Pulls)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Could not unmarshal response: %s\n", err.Error())
		return
	}

	sort.Sort(r.Pulls)

	for _, pull := range r.Pulls {
		pull.Repo = r.Name
	}
}
