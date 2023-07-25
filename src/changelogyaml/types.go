/*---------------------------------------------------------------------------------------------
 *  Copyright (c) Peter Bjorklund. All rights reserved.
 *  Licensed under the MIT License. See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package changelogyaml

import (
	"io"
	"log"

	"gopkg.in/yaml.v2"
)

const githubUrlPrefix = "https://github.com/"

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
