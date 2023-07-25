/*---------------------------------------------------------------------------------------------
 *  Copyright (c) Peter Bjorklund. All rights reserved.
 *  Licensed under the MIT License. See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package changelogyaml

import (
	"fmt"
	"regexp"
)

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
