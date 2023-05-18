/*---------------------------------------------------------------------------------------------
 *  Copyright (c) Peter Bjorklund. All rights reserved.
 *  Licensed under the MIT License. See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"

	"gopkg.in/yaml.v2"
)

/*
New features/added/feat
Changes/Changed
Removed
Bugs fixed/bug fixes/fix/fixed
Security

Docs
Refactor
Performance
Tests
Breaking
Deprecate
*/

/*
Not included:
Styles
Chores
ci
revert
builds
*/

type Module struct {
	Name       string
	Improved   []string
	Changed    []string
	Added      []string
	Removed    []string
	Fixed      []string
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

func (c *ChangelogYaml) ReadYaml(filename string) *ChangelogYaml {
	yamlFile, err := ioutil.ReadFile(filename)
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

func moduleLines(strings []string, name string, writer io.Writer) {
	for _, line := range strings {
		categoryInfo := infoFromCategoryName(name)
		prefix := categoryInfo.Prefix
		if name == "breaking" {
			prefix += fmt.Sprintf("[%v]", categoryInfo.Name)
		}
		fmt.Fprintf(writer, "* %v %v\n", prefix, line)
	}
}

func moduleGroupedLines(module *Module, writer io.Writer) {
	moduleLines(module.Breaking, "breaking", writer)
	moduleLines(module.Added, "added", writer)
	moduleLines(module.Fixed, "fixed", writer)
	moduleLines(module.Changed, "changed", writer)
	moduleLines(module.Removed, "removed", writer)
	moduleLines(module.Improved, "improvements", writer)
	moduleLines(module.Docs, "docs", writer)
	moduleLines(module.Tests, "tests", writer)
	moduleLines(module.Refactored, "refactored", writer)
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
			fmt.Fprintf(writer, "\n### [%v](%v) - %v\n\n", module.Name, repoLink, info.Description)
			moduleGroupedLines(&module, writer)
		}
	}
}

func main() {
	var c ChangelogYaml

	filename := os.Args[1]
	c.ReadYaml(filename)

	writeToMarkdown(&c, os.Stdout)
}
