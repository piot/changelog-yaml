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
	"strconv"

	"gopkg.in/yaml.v2"
)

type Module struct {
	Name       string
	Improved   []string
	Changed    []string
	Added      []string
	Removed    []string
	Fixed      []string
	Workaround []string
	Deprecated []string
	Tests      []string
	Docs       []string
	Refactored []string
	Breaking   []string
}

type Release struct {
	Name    string
	Date    string
	Notice  string
	Modules []Module
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
		"added":       {":star2:", "added"},
		"changed":     {":hammer_and_wrench:", "changed"},
		"fixed":       {":lady_beetle:", "fixed"},
		"workaround":  {":see_no_evil:", "workaround"},
		"performance": {":zap:", "performance"},
		"tests":       {":vertical_traffic_light:", "test"},
		"removed":     {":fire:", "removed"},
		"improved":    {":art:", "improved"},
		"breaking":    {":rotating_light:", "breaking"},
		"deprecated":  {":rotating_light:", "deprecated"},
		"refactored":  {":recycle:", "refactor"},
		"docs":        {":book:", "docs"},
	}

	info, wasFound := lookup[name]
	if !wasFound {
		panic(fmt.Errorf("unknown '%v'", name))
	}

	return info
}

func moduleLines(moduleRepoUrl string, strings []string, name string, writer io.Writer) error {
	for _, line := range strings {
		categoryInfo := infoFromCategoryName(name)
		prefix := categoryInfo.Prefix
		if name == "breaking" {
			prefix += fmt.Sprintf("[%v]", categoryInfo.Name)
		}
		re := regexp.MustCompile(`\(\#\d*\)`)
		allMatches := re.FindAllStringIndex(line, -1)
		lineToPrint := line
		for _, match := range allMatches {
			matchString := line[match[0]+2 : match[1]-1]
			pullRequestId, err := strconv.Atoi(matchString)
			if err != nil {
				return err
			}
			suffix := fmt.Sprintf("pull/%v", pullRequestId)
			pullRequestLink := fmt.Sprintf("https://%v/%v", moduleRepoUrl, suffix)
			pullRequestCompleteLink := fmt.Sprintf("([#%v](%s))", pullRequestId, pullRequestLink)
			lineToPrint = line[:match[0]] + pullRequestCompleteLink + line[match[1]:]
		}
		fmt.Fprintf(writer, "* %v %v\n", prefix, lineToPrint)
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
	moduleLines(moduleRepoUrl, module.Improved, "improvements", writer)
	moduleLines(moduleRepoUrl, module.Docs, "docs", writer)
	moduleLines(moduleRepoUrl, module.Tests, "tests", writer)
	moduleLines(moduleRepoUrl, module.Refactored, "refactored", writer)
}

func writeToMarkdown(root *ChangelogYaml, writer io.Writer) {
	fmt.Fprintln(writer, "# Changelog")
	for _, release := range root.Releases {
		releaseLink := fmt.Sprintf("[%v](https://%v/releases/tag/%v) (%v)", release.Name, root.Repo, release.Name, release.Date)
		fmt.Fprintf(writer, "\n## :bookmark: %v\n", releaseLink)

		if release.Notice != "" {
			fmt.Fprintf(writer, "\n%v\n", release.Notice)
		}

		for _, module := range release.Modules {
			info, found := root.Modules[module.Name]
			if !found {
				panic(fmt.Errorf("must have info for module '%s'", module.Name))
			}
			repoLink := fmt.Sprintf("https://%v", info.Repo)
			description := ""
			if info.Description != "" {
				description = fmt.Sprintf(" - %v", info.Description)
			}
			fmt.Fprintf(writer, "\n### [%v](%v)%v\n\n", module.Name, repoLink, description)
			moduleGroupedLines(info.Repo, &module, writer)
		}
	}
}

func main() {
	var c ChangelogYaml

	reader := bufio.NewReader(os.Stdin)
	c.ReadYaml(reader)

	writeToMarkdown(&c, os.Stdout)
}
