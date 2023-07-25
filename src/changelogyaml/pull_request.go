/*---------------------------------------------------------------------------------------------
 *  Copyright (c) Peter Bjorklund. All rights reserved.
 *  Licensed under the MIT License. See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package changelogyaml

import (
	"fmt"
	"regexp"
	"strconv"
)

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
