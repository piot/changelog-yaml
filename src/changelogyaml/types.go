/*---------------------------------------------------------------------------------------------
 *  Copyright (c) Peter Bjorklund. All rights reserved.
 *  Licensed under the MIT License. See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package changelogyaml

import (
	"io"
	"log"

	"gopkg.in/yaml.v3"
)

const githubUrlPrefix = "https://github.com/"

type Changes struct {
	// Keep a changelog https://keepachangelog.com/en/1.1.0/

	// Added denotes new features or functionalities introduced in the software.
	Added []string

	// Changed indicates changes to existing features or functionalities.
	Changed []string

	// Deprecated signifies functionalities that are no longer recommended and will be removed in future versions.
	Deprecated []string

	// Removed lists functionalities or features that have been removed from the software. Should have been set as Deprecated in version prior to being removed. Is implicitly breaking changes.
	Removed []string

	// Fixed enumerates fixes for bugs or issues in the software.
	Fixed []string

	// Security includes changes related to security enhancements or fixes.
	Security []string

	// ---------------- Others ---------------

	// Improved lists improvements made to existing functionalities without adding new features.
	Improved []string

	// Workaround provides workarounds or temporary solutions for known issues or limitations.
	Workaround []string

	// Tests includes changes or additions to testing procedures or test cases.
	Tests []string

	// Docs lists changes or additions to documentation, such as README files or inline code comments.
	Docs []string

	// Refactored denotes changes made to improve code structure or organization without changing external behavior.
	Refactored []string

	// Performance includes changes aimed at improving the performance of the software.
	Performance []string

	// Breaking denotes changes that may break backward compatibility with previous versions. Changed, but breaks the API compatibilty.
	Breaking []string

	// Experimental lists experimental features or functionalities that are not yet stable or fully supported and might be removed with short or no notice in future versions.
	Experimental []string

	// Noted provides a place to note any other significant changes not covered by the above categories.
	Noted []string

	// Style denotes changes related to coding style, formatting, or other stylistic aspects.
	Style []string

	// Unreleased contains a list of changes that are planned but not yet released in any version.
	// These changes typically represent work that is in progress or pending release in a future version.
	// Once a version is released, the changes listed in Unreleased are moved to the appropriate category (e.g., Added, Changed, Fixed, etc.).
	Unreleased []string
}

type Section struct {
	Order   int
	Notice  string
	Changes Changes
}

type Release struct {
	Name     string
	Date     string
	Notice   string
	Repos    map[string]Changes `yaml:"repos"`
	Sections map[string]Section `yaml:"sections"`
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
