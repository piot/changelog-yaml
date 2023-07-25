/*---------------------------------------------------------------------------------------------
 *  Copyright (c) Peter Bjorklund. All rights reserved.
 *  Licensed under the MIT License. See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package changelogyaml

import "fmt"

type CategoryType uint8

const (
	Added CategoryType = iota
	Changed
	Fixed
	Workaround
	Performance
	Tests
	Removed
	Improved
	Breaking
	Deprecated
	Refactored
	Experimental
	Docs
	Noted
)

type CategoryInfo struct {
	EmojiName string
	Name      string
}

func infoFromCategoryName(name CategoryType) CategoryInfo {
	lookup := map[CategoryType]CategoryInfo{
		Added:        {"star2", "added"},
		Changed:      {"hammer_and_wrench", "changed"},
		Fixed:        {"lady_beetle", "fixed"},
		Workaround:   {"see_no_evil", "workaround"},
		Performance:  {"zap", "performance"},
		Tests:        {"vertical_traffic_light", "test"},
		Removed:      {"fire", "removed"},
		Improved:     {"art", "improved"},
		Breaking:     {"triangular_flag_on_post", "breaking"},
		Deprecated:   {"spider_web", "deprecated"},
		Refactored:   {"recycle", "refactor"},
		Experimental: {"alembic", "experimental"},
		Docs:         {"book", "docs"},
		Noted:        {"beetle", "known issue"},
	}

	info, wasFound := lookup[name]
	if !wasFound {
		panic(fmt.Errorf("unknown '%v'", name))
	}

	return info
}
