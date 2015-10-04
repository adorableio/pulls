package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/mitchellh/go-homedir"
	"gopkg.in/yaml.v2"
)

type RepoConfig []string

func LoadRepoConfig(path string) RepoConfig {
	path, err := filepath.Abs(path)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Could not find current working directory")
		return make(RepoConfig, 0)
	}

	path, err = findRepoConfig(path)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Could not find repo config file anywhere in current folder hierarchy")
		return make(RepoConfig, 0)
	}

	file, err := ioutil.ReadFile(path)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Could not open config file %s: %s", path, err.Error())
		os.Exit(1)
	}

	config := make(RepoConfig, 0)
	err = yaml.Unmarshal(file, &config)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Could not deserialize config file %s", err.Error())
		os.Exit(1)
	}

	return config
}

func findRepoConfig(path string) (string, error) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		// path doesn't exist
		subpath, filename := filepath.Split(path)

		if subpath == "/" {
			return "", fmt.Errorf("Could not find repo config file")
		}

		subpathParts := strings.Split(subpath, "/")
		subpathParts = subpathParts[0 : len(subpathParts)-2]
		subpathWithFilename := append(subpathParts, filename)
		newpath := "/" + filepath.Join(subpathWithFilename...)
		return findRepoConfig(newpath)
	}

	return path, nil
}

type GithubConfig map[string]string

func LoadGithubConfig() GithubConfig {
	home, err := homedir.Dir()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Could not find home directory: %s", err.Error())
		os.Exit(1)
	}

	filepath := path.Join(home, ".pulls.github.yml")
	file, err := ioutil.ReadFile(filepath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Could not open config file %s: %s", filepath, err.Error())
		os.Exit(1)
	}

	config := make(GithubConfig)
	err = yaml.Unmarshal(file, &config)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Could not deserialize config file %s", err.Error())
		os.Exit(1)
	}

	return config
}
