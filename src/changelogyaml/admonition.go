/*---------------------------------------------------------------------------------------------
 *  Copyright (c) Peter Bjorklund. All rights reserved.
 *  Licensed under the MIT License. See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package changelogyaml

import (
	"fmt"
	"regexp"
	"strings"
)

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
