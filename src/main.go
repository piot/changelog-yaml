/*---------------------------------------------------------------------------------------------
 *  Copyright (c) Peter Bjorklund. All rights reserved.
 *  Licensed under the MIT License. See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"regexp"
	"sort"
	"strconv"
	"strings"

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
		"added":        {"star2", "added"},
		"changed":      {"hammer_and_wrench", "changed"},
		"fixed":        {"lady_beetle", "fixed"},
		"workaround":   {"see_no_evil", "workaround"},
		"performance":  {"zap", "performance"},
		"tests":        {"vertical_traffic_light", "test"},
		"removed":      {"fire", "removed"},
		"improved":     {"art", "improved"},
		"breaking":     {"triangular_flag_on_post", "breaking"},
		"deprecated":   {"spider_web", "deprecated"},
		"refactored":   {"recycle", "refactor"},
		"experimental": {"alembic", "experimental"},
		"docs":         {"book", "docs"},
		"noted":        {"beetle", "known issue"},
	}

	info, wasFound := lookup[name]
	if !wasFound {
		panic(fmt.Errorf("unknown '%v'", name))
	}

	return info
}

func stringToAdmonitionType(name string) AdmonitionType {
	switch name {
	case "WARNING":
		return Warning
	case "NOTE":
		return Note
	case "IMPORTANT":
		return Important
	}

	panic(fmt.Errorf("unknown admonition: '%s'", name))
}

func replaceAdmonition(line string, formatter Formatter) string {
	re := regexp.MustCompile(`(WARNING|TIP|NOTE|IMPORTANT|CAUTION):\s.*`)
	allMatches := re.FindAllStringIndex(line, -1)

	lineToPrint := ""
	previousMatchPosition := 0

	if len(allMatches) > 0 {
		for _, match := range allMatches {
			matchString := line[match[0]:match[1]]

			parts := strings.Split(matchString, ":")

			lineToPrint += line[previousMatchPosition:match[0]] + formatter.Admonition(stringToAdmonitionType(parts[0]),
				parts[1][1:])
			previousMatchPosition = match[1]
		}

		lineToPrint += line[previousMatchPosition:]
	} else {
		lineToPrint = line
	}

	return lineToPrint
}

func replaceAtProfileLink(line string, formatter Formatter) string {
	re := regexp.MustCompile(`@[a-z\d-]*`)
	allMatches := re.FindAllStringIndex(line, -1)

	lineToPrint := ""
	previousMatchPosition := 0

	if len(allMatches) > 0 {
		for _, match := range allMatches {
			usernameString := line[match[0]+1 : match[1]]
			usernameProfileLink := fmt.Sprintf("%s%v", githubUrlPrefix, usernameString)
			usernameProfileLinkComplete := formatter.Link("@"+usernameString, usernameProfileLink)
			lineToPrint += line[previousMatchPosition:match[0]] + usernameProfileLinkComplete
			previousMatchPosition = match[1]
		}

		lineToPrint += line[previousMatchPosition:]
	} else {
		lineToPrint = line
	}

	return lineToPrint
}

func replacePullRequestLink(line string, repoShortUrl string, formatter Formatter) (string, error) {
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
			pullRequestCompleteLink := formatter.Link(fmt.Sprintf("#%v", pullRequestID), pullRequestLink)
			lineToPrint += line[previousMatchPosition:match[0]] + pullRequestCompleteLink
			previousMatchPosition = match[1]
		}

		lineToPrint += line[previousMatchPosition:]
	} else {
		lineToPrint = line
	}

	return lineToPrint, nil
}

func replaceCommitHashLink(line string, repoShortUrl string, formatter Formatter) string {
	re := regexp.MustCompile(`\$[a-f\d]*`)
	allMatches := re.FindAllStringIndex(line, -1)

	lineToPrint := ""
	previousMatchPosition := 0

	if len(allMatches) > 0 {
		for _, match := range allMatches {
			commitHashString := line[match[0]+1 : match[1]]
			commitHashLink := fmt.Sprintf("%s%v/commit/%v", githubUrlPrefix, repoShortUrl, commitHashString)
			commitHashLinkComplete := formatter.Link(commitHashString, commitHashLink)
			lineToPrint += line[previousMatchPosition:match[0]] + commitHashLinkComplete
			previousMatchPosition = match[1]
		}

		lineToPrint += line[previousMatchPosition:]
	} else {
		lineToPrint = line
	}

	return lineToPrint
}

func linesForRepo(repoShortUrl string, strings []string, name string, formatter Formatter, writer io.Writer) error {
	for _, line := range strings {
		categoryInfo := infoFromCategoryName(name)
		iconName := categoryInfo.Prefix

		prefix := formatter.Icon(iconName)

		if name == "breaking" {
			prefix += fmt.Sprintf("[%v]", categoryInfo.Name)
		}

		line, err := replacePullRequestLink(line, repoShortUrl, formatter)
		if err != nil {
			return err
		}

		line = replaceCommitHashLink(line, repoShortUrl, formatter)

		line = replaceAtProfileLink(line, formatter)

		completeLine := prefix + " " + line

		if _, err := fmt.Fprint(writer, formatter.BulletPoint(completeLine)); err != nil {
			return err
		}
	}

	return nil
}

type LineInfo struct {
	Name  string
	Lines []string
}

func textLinesForTheRepo(repoShortUrl string, repoChanges *RepoChanges, formatter Formatter, writer io.Writer) error {
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
		if err := linesForRepo(repoShortUrl, line.Lines, line.Name, formatter, writer); err != nil {
			return err
		}
	}

	return nil
}

type AdmonitionType uint8

const (
	Note AdmonitionType = iota
	Important
	Warning
)

type Formatter interface {
	Heading(level int, header string) string
	BulletPoint(text string) string
	Icon(name string) string
	Link(name string, link string) string
	Admonition(admonitionType AdmonitionType, text string) string
}

func emojiNameToUnicode(name string) int {
	lookup := map[string]int{
		"bookmark":                0x1F516,
		"triangular_flag_on_post": 0x1f6a9,
		"star2":                   0x1f31f,
		"hammer_and_wrench":       0x1f6e0,
		"lady_beetle":             0x1f41e,
		"see_no_evil":             0x1f648,
		"zap":                     0x26a1,
		"vertical_traffic_light":  0x1f6a6,
		"fire":                    0x1f525,
		"art":                     0x1f3a8,
		"spider_web":              0x1f578,
		"recycle":                 0x267b,
		"alembic":                 0x2697,
		"book":                    0x1F4D6,
		"noted":                   0x1FAB2,
	}

	replacement, wasFound := lookup[name]
	if !wasFound {
		panic(fmt.Errorf("can not replace %s", name))
	}

	return replacement
}

type MarkdownFormatter struct {
}

func (m *MarkdownFormatter) Heading(level int, header string) string {
	return strings.Repeat("#", level) + " " + header + "\n\n"
}

func (m *MarkdownFormatter) BulletPoint(text string) string {
	return "* " + text + "\n"
}

func (m *MarkdownFormatter) Icon(name string) string {
	return ":" + name + ":"
}

func (m *MarkdownFormatter) Link(name string, link string) string {
	return fmt.Sprintf("[%s](%s)", name, link)
}

func AdmonitionTypeToGithubName(admonitionType AdmonitionType) string {
	switch admonitionType {
	case Note:
		return "NOTE"
	case Important:
		return "IMPORTANT"
	case Warning:
		return "WARNING"
	}

	panic(fmt.Errorf("unknown admonition"))
}

func (m *MarkdownFormatter) Admonition(admonitionType AdmonitionType, text string) string {
	return fmt.Sprintf("> [!%s]\\\n> %s", AdmonitionTypeToGithubName(admonitionType), text)
}

type AsciiDocFormatter struct {
}

func (a *AsciiDocFormatter) Heading(level int, header string) string {
	return strings.Repeat("=", level) + " " + header + "\n\n"
}

func (m *AsciiDocFormatter) BulletPoint(text string) string {
	return "* " + text + "\n"
}

func (m *AsciiDocFormatter) Icon(name string) string {
	unicodeInt := emojiNameToUnicode(name)
	return fmt.Sprintf("&#x%X;", unicodeInt)
}

func (m *AsciiDocFormatter) Link(name string, link string) string {
	return fmt.Sprintf("%s[%s]", link, name)
}

func AdmonitionTypeToAsciiDocName(admonitionType AdmonitionType) string {
	switch admonitionType {
	case Note:
		return "NOTE"
	case Important:
		return "IMPORTANT"
	case Warning:
		return "WARNING"
	}
	// CAUTION and TIP is not supported yet

	panic(fmt.Errorf("unknown admonition"))
}

func (m *AsciiDocFormatter) Admonition(admonitionType AdmonitionType, text string) string {
	return fmt.Sprintf("%s: %s", AdmonitionTypeToAsciiDocName(admonitionType), text)
}

func writeToMarkdown(root *ChangelogYaml, outputFormatter Formatter, writer io.Writer) error {
	if _, err := fmt.Fprint(writer, outputFormatter.Heading(1, "Changelog")); err != nil {
		return err
	}

	for _, release := range root.Releases {
		completeReleaseLinkURL := fmt.Sprintf("%s%v/releases/tag/%v", githubUrlPrefix, root.Repo, release.Name)
		formattedReleaseLink := outputFormatter.Link(release.Name, completeReleaseLinkURL)

		releaseLink := fmt.Sprintf("%v (%v)", formattedReleaseLink, release.Date)
		releaseHeading := fmt.Sprintf("%s %v", outputFormatter.Icon("bookmark"), releaseLink)
		if _, err := fmt.Fprint(writer, outputFormatter.Heading(2, releaseHeading)); err != nil {
			return err
		}

		if release.Notice != "" {
			notice := replaceAdmonition(release.Notice, outputFormatter)
			notice = replaceAtProfileLink(notice, outputFormatter)
			fmt.Fprintf(writer, "%v\n\n", notice)
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

			repoURL := fmt.Sprintf("%s%v", githubUrlPrefix, info.Repo)
			description := ""

			if info.Description != "" {
				description = fmt.Sprintf(" - %v", info.Description)
			}

			formattedRepoLink := outputFormatter.Link(repoName, repoURL)

			completeLine := fmt.Sprintf("%s%v", formattedRepoLink, description)

			if _, err := fmt.Fprint(writer, outputFormatter.Heading(3, completeLine)); err != nil {
				return err
			}

			if err := textLinesForTheRepo(info.Repo, &repoInfo, outputFormatter, writer); err != nil {
				return err
			}

			fmt.Fprintf(writer, "\n")
		}
	}

	return nil
}

func main() {
	var outputFormat = flag.String("format", "markdown", "output format: md or adoc")
	flag.Parse()

	var c ChangelogYaml

	reader := bufio.NewReader(os.Stdin)
	c.ReadYaml(reader)

	var formatter Formatter
	if *outputFormat == "adoc" || *outputFormat == "asciidoc" {
		formatter = &AsciiDocFormatter{}
	} else {
		formatter = &MarkdownFormatter{}
	}

	err := writeToMarkdown(&c, formatter, os.Stdout)
	if err != nil {
		os.Exit(-2)
	}
}
