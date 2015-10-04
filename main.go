package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"sort"
	"strconv"
	"strings"
	"sync"
)

func main() {
	var (
		list   bool
		config string
	)

	flag.BoolVar(&list, "l", false, "print list of pull requests and exit")
	flag.StringVar(&config, "c", "repos.yml", "path to config file")
	flag.Parse()

	repoConfig := LoadRepoConfig(config)
	githubConfig := LoadGithubConfig()

	pulls := make(PullRequests, 0)
	repos := make(Repositories, 0)

	wg := new(sync.WaitGroup)
	for _, name := range repoConfig {
		r := NewRepo(name, githubConfig["accessToken"])
		repos = append(repos, r)

		wg.Add(1)
		go r.Fetch(wg)
	}
	wg.Wait()

	sort.Sort(repos)
	for _, repo := range repos {
		pulls = append(pulls, repo.Pulls...)
	}

	if len(pulls) == 0 {
		fmt.Println("You have no open pull requests; have a great day!")
		os.Exit(0)
	}

	for i, pull := range pulls {
		if i > 0 && pulls[i].Repo != pulls[i-1].Repo {
			fmt.Println("")
		}
		fmt.Println(pull.String(i + 1))
	}

	if !list {
		goInteractive(pulls)
	}

	os.Exit(0)
}

func goInteractive(pulls PullRequests) {
	selected := getSelection(len(pulls))

	if selected == -1 {
		os.Exit(1)
	}

	if selected < 1 || selected > int64(len(pulls)) {
		fmt.Println("Invalid selection.")
		os.Exit(1)
	}

	selectedPull := pulls[selected-1]

	open := fmt.Sprintf("open %s", selectedPull.HtmlUrl)
	err := exec.Command("/bin/sh", "-c", open).Run()
	if err != nil {
		fmt.Printf("Could not open url: %s", err.Error())
		os.Exit(1)
	}
}

func getSelection(numPulls int) int64 {
	reader := bufio.NewReader(os.Stdin)
	fmt.Printf("\nSelect [1-%d] or [e]xit: ", numPulls)
	text, _ := reader.ReadString('\n')

	input := strings.Trim(text, " \n")

	if input == "e" || input == "" {
		return -1
	}

	selected, err := strconv.ParseInt(input, 10, 0)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Could not parse input: %s", err.Error())
		os.Exit(1)
	}

	return selected
}
