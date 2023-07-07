/*---------------------------------------------------------------------------------------------
 *  Copyright (c) Peter Bjorklund. All rights reserved.
 *  Licensed under the MIT License. See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package main

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"regexp"
	"sort"
	"strconv"

	"gopkg.in/yaml.v2"
)

type Module struct {
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
}

type Release struct {
	Name    string
	Date    string
	Notice  string
	Modules map[string]Module `yaml:"modules"`
}

type ModuleDefinition struct {
	Repo        string
	Name        string
	Description string
}

type ChangelogYaml struct {
	Repo     string
	Releases []Release
	Modules  map[string]ModuleDefinition `yaml:"modules"`
}

func (c *ChangelogYaml) ReadYaml(filename io.Reader) *ChangelogYaml {
	yamlFile, err := ioutil.ReadAll(filename)
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
			usernameProfileLink := fmt.Sprintf("https://github.com/%v", usernameString)
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

func replacePullRequestLink(line string, moduleRepoURL string) (string, error) {
	re := regexp.MustCompile(`\#\d*`)
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
			pullRequestLink := fmt.Sprintf("https://%v/%v", moduleRepoURL, suffix)
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

func replaceCommitHashLink(line string, moduleRepoURL string) string {
	re := regexp.MustCompile(`\$[a-f\d]*`)
	allMatches := re.FindAllStringIndex(line, -1)

	lineToPrint := ""
	previousMatchPosition := 0

	if len(allMatches) > 0 {
		for _, match := range allMatches {
			commitHashString := line[match[0]+1 : match[1]]
			commitHashLink := fmt.Sprintf("https://%v/commit/%v", moduleRepoURL, commitHashString)
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

func moduleLines(moduleRepoUrl string, strings []string, name string, writer io.Writer) error {
	for _, line := range strings {
		categoryInfo := infoFromCategoryName(name)
		prefix := categoryInfo.Prefix

		if name == "breaking" {
			prefix += fmt.Sprintf("[%v]", categoryInfo.Name)
		}

		line, err := replacePullRequestLink(line, moduleRepoUrl)
		if err != nil {
			return err
		}

		line = replaceCommitHashLink(line, moduleRepoUrl)

		line = replaceAtProfileLink(line)

		fmt.Fprintf(writer, "* %v %v\n", prefix, line)
	}

	return nil
}

func moduleGroupedLines(moduleRepoUrl string, module *Module, writer io.Writer) {
	moduleLines(moduleRepoUrl, module.Breaking, "breaking", writer)
	moduleLines(moduleRepoUrl, module.Added, "added", writer)
	moduleLines(moduleRepoUrl, module.Fixed, "fixed", writer)
	moduleLines(moduleRepoUrl, module.Workaround, "workaround", writer)
	moduleLines(moduleRepoUrl, module.Changed, "changed", writer)
	moduleLines(moduleRepoUrl, module.Removed, "removed", writer)
	moduleLines(moduleRepoUrl, module.Improved, "improved", writer)
	moduleLines(moduleRepoUrl, module.Docs, "docs", writer)
	moduleLines(moduleRepoUrl, module.Tests, "tests", writer)
	moduleLines(moduleRepoUrl, module.Refactored, "refactored", writer)
	moduleLines(moduleRepoUrl, module.Performance, "performance", writer)
	moduleLines(moduleRepoUrl, module.Deprecated, "deprecated", writer)
	moduleLines(moduleRepoUrl, module.Experimental, "experimental", writer)
}

func writeToMarkdown(root *ChangelogYaml, writer io.Writer) error {
	fmt.Fprintln(writer, "# Changelog")

	for _, release := range root.Releases {
		releaseLink := fmt.Sprintf("[%v](https://%v/releases/tag/%v) (%v)", release.Name,
			root.Repo, release.Name, release.Date)
		fmt.Fprintf(writer, "\n## :bookmark: %v\n", releaseLink)

		if release.Notice != "" {
			notice := replaceAtProfileLink(release.Notice)
			fmt.Fprintf(writer, "\n%v\n", notice)
		}

		sortedModuleNames := make([]string, 0, len(release.Modules))
		for k := range release.Modules {
			sortedModuleNames = append(sortedModuleNames, k)
		}

		sort.Strings(sortedModuleNames)

		for _, moduleName := range sortedModuleNames {
			module := release.Modules[moduleName]

			info, found := root.Modules[moduleName]
			if !found {
				panic(fmt.Errorf("must have info for module '%s'", moduleName))
			}

			repoLink := fmt.Sprintf("https://%v", info.Repo)
			description := ""

			if info.Description != "" {
				description = fmt.Sprintf(" - %v", info.Description)
			}

			fmt.Fprintf(writer, "\n### [%v](%v)%v\n\n", moduleName, repoLink, description)
			moduleGroupedLines(info.Repo, &module, writer)
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
