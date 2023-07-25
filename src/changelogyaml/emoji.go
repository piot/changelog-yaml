/*---------------------------------------------------------------------------------------------
 *  Copyright (c) Peter Bjorklund. All rights reserved.
 *  Licensed under the MIT License. See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package changelogyaml

import "fmt"

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
