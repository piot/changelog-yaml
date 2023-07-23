/*---------------------------------------------------------------------------------------------
 *  Copyright (c) Peter Bjorklund. All rights reserved.
 *  Licensed under the MIT License. See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"regexp"
	"sort"
	"strconv"

	"gopkg.in/yaml.v2"
)

type RepoChanges struct {
	Improved     []string
	Changed      []string
	Added        []string
	Removed      []string
	Fixed        []string
	Workaround   []string
	Deprecated   []string
	Tests        []string
	Docs         []string
	Refactored   []string
	Performance  []string
	Breaking     []string
	Experimental []string
	Noted        []string
}

type Release struct {
	Name   string
	Date   string
	Notice string
	Repos  map[string]RepoChanges `yaml:"repos"`
}

type RepoDefinition struct {
	Repo        string
	Name        string
	Description string
}

type ChangelogYaml struct {
	Repo     string
	Releases []Release
	Repos    map[string]RepoDefinition `yaml:"repos"`
}

func (c *ChangelogYaml) ReadYaml(filename io.Reader) *ChangelogYaml {
	yamlFile, err := io.ReadAll(filename)
	if err != nil {
		log.Fatalf("could not read file %v ", err)
	}

	err = yaml.Unmarshal(yamlFile, c)
	if err != nil {
		log.Fatalf("unmarshal error: %v", err)
	}

	return c
}

type CategoryInfo struct {
	Prefix string
	Name   string
}

const githubUrlPrefix = "https://github.com/"

func infoFromCategoryName(name string) CategoryInfo {
	lookup := map[string]CategoryInfo{
		"added":        {":star2:", "added"},
		"changed":      {":hammer_and_wrench:", "changed"},
		"fixed":        {":lady_beetle:", "fixed"},
		"workaround":   {":see_no_evil:", "workaround"},
		"performance":  {":zap:", "performance"},
		"tests":        {":vertical_traffic_light:", "test"},
		"removed":      {":fire:", "removed"},
		"improved":     {":art:", "improved"},
		"breaking":     {":triangular_flag_on_post:", "breaking"},
		"deprecated":   {":spider_web:", "deprecated"},
		"refactored":   {":recycle:", "refactor"},
		"experimental": {":alembic:", "experimental"},
		"docs":         {":book:", "docs"},
		"noted":        {":beetle:", "known issue"},
	}

	info, wasFound := lookup[name]
	if !wasFound {
		panic(fmt.Errorf("unknown '%v'", name))
	}

	return info
}

func replaceAtProfileLink(line string) string {
	re := regexp.MustCompile(`@[a-z\d-]*`)
	allMatches := re.FindAllStringIndex(line, -1)

	lineToPrint := ""
	previousMatchPosition := 0

	if len(allMatches) > 0 {
		for _, match := range allMatches {
			usernameString := line[match[0]+1 : match[1]]
			usernameProfileLink := fmt.Sprintf("%s%v", githubUrlPrefix, usernameString)
			usernameProfileLinkComplete := fmt.Sprintf("[@%v](%s)", usernameString, usernameProfileLink)
			lineToPrint += line[previousMatchPosition:match[0]] + usernameProfileLinkComplete
			previousMatchPosition = match[1]
		}

		lineToPrint += line[previousMatchPosition:]
	} else {
		lineToPrint = line
	}

	return lineToPrint
}

func replacePullRequestLink(line string, repoShortUrl string) (string, error) {
	re := regexp.MustCompile(`#\d*`)
	allMatches := re.FindAllStringIndex(line, -1)

	lineToPrint := ""
	previousMatchPosition := 0

	if len(allMatches) > 0 {
		for _, match := range allMatches {
			matchString := line[match[0]+1 : match[1]]

			pullRequestID, err := strconv.Atoi(matchString)
			if err != nil {
				return "", err
			}

			suffix := fmt.Sprintf("pull/%v", pullRequestID)
			pullRequestLink := fmt.Sprintf("%s%v/%v", githubUrlPrefix, repoShortUrl, suffix)
			pullRequestCompleteLink := fmt.Sprintf("[#%v](%s)", pullRequestID, pullRequestLink)
			lineToPrint += line[previousMatchPosition:match[0]] + pullRequestCompleteLink
			previousMatchPosition = match[1]
		}

		lineToPrint += line[previousMatchPosition:]
	} else {
		lineToPrint = line
	}

	return lineToPrint, nil
}

func replaceCommitHashLink(line string, repoShortUrl string) string {
	re := regexp.MustCompile(`\$[a-f\d]*`)
	allMatches := re.FindAllStringIndex(line, -1)

	lineToPrint := ""
	previousMatchPosition := 0

	if len(allMatches) > 0 {
		for _, match := range allMatches {
			commitHashString := line[match[0]+1 : match[1]]
			commitHashLink := fmt.Sprintf("%s%v/commit/%v", githubUrlPrefix, repoShortUrl, commitHashString)
			commitHashLinkComplete := fmt.Sprintf("[%v](%s)", commitHashString, commitHashLink)
			lineToPrint += line[previousMatchPosition:match[0]] + commitHashLinkComplete
			previousMatchPosition = match[1]
		}

		lineToPrint += line[previousMatchPosition:]
	} else {
		lineToPrint = line
	}

	return lineToPrint
}

func linesForRepo(repoShortUrl string, strings []string, name string, writer io.Writer) error {
	for _, line := range strings {
		categoryInfo := infoFromCategoryName(name)
		prefix := categoryInfo.Prefix

		if name == "breaking" {
			prefix += fmt.Sprintf("[%v]", categoryInfo.Name)
		}

		line, err := replacePullRequestLink(line, repoShortUrl)
		if err != nil {
			return err
		}

		line = replaceCommitHashLink(line, repoShortUrl)

		line = replaceAtProfileLink(line)

		if _, err := fmt.Fprintf(writer, "* %v %v\n", prefix, line); err != nil {
			return err
		}
	}

	return nil
}

type LineInfo struct {
	Name  string
	Lines []string
}

func textLinesForTheRepo(repoShortUrl string, repoChanges *RepoChanges, writer io.Writer) error {
	lines := []LineInfo{
		{"breaking", repoChanges.Breaking},
		{"added", repoChanges.Added},
		{"fixed", repoChanges.Fixed},
		{"workaround", repoChanges.Workaround},
		{"changed", repoChanges.Changed},
		{"removed", repoChanges.Removed},
		{"improved", repoChanges.Improved},
		{"docs", repoChanges.Docs},
		{"tests", repoChanges.Tests},
		{"refactored", repoChanges.Refactored},
		{"deprecated", repoChanges.Deprecated},
		{"experimental", repoChanges.Experimental},
		{"noted", repoChanges.Noted},
		{"performance", repoChanges.Performance},
	}

	for _, line := range lines {
		if err := linesForRepo(repoShortUrl, line.Lines, line.Name, writer); err != nil {
			return err
		}
	}

	return nil
}

func writeToMarkdown(root *ChangelogYaml, writer io.Writer) error {
	fmt.Fprintln(writer, "# Changelog")

	for _, release := range root.Releases {
		releaseLink := fmt.Sprintf("[%v](%s%v/releases/tag/%v) (%v)", release.Name,
			githubUrlPrefix, root.Repo, release.Name, release.Date)
		fmt.Fprintf(writer, "\n## :bookmark: %v\n", releaseLink)

		if release.Notice != "" {
			notice := replaceAtProfileLink(release.Notice)
			fmt.Fprintf(writer, "\n%v\n", notice)
		}

		sortedRepoNames := make([]string, 0, len(release.Repos))
		for k := range release.Repos {
			sortedRepoNames = append(sortedRepoNames, k)
		}

		sort.Strings(sortedRepoNames)

		for _, repoName := range sortedRepoNames {
			repoInfo := release.Repos[repoName]

			info, found := root.Repos[repoName]
			if !found {
				panic(fmt.Errorf("must have info for repoInfo '%s'", repoName))
			}

			repoLink := fmt.Sprintf("%s%v", githubUrlPrefix, info.Repo)
			description := ""

			if info.Description != "" {
				description = fmt.Sprintf(" - %v", info.Description)
			}

			if _, err := fmt.Fprintf(writer, "\n### [%v](%v)%v\n\n", repoName, repoLink, description); err != nil {
				return err
			}

			if err := textLinesForTheRepo(info.Repo, &repoInfo, writer); err != nil {
				return err
			}
		}
	}

	return nil
}

func main() {
	var c ChangelogYaml

	reader := bufio.NewReader(os.Stdin)
	c.ReadYaml(reader)

	err := writeToMarkdown(&c, os.Stdout)
	if err != nil {
		os.Exit(-2)
	}
}
