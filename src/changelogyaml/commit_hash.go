/*---------------------------------------------------------------------------------------------
 *  Copyright (c) Peter Bjorklund. All rights reserved.
 *  Licensed under the MIT License. See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package changelogyaml

import (
	"fmt"
	"regexp"
)

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
