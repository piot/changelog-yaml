/*---------------------------------------------------------------------------------------------
 *  Copyright (c) Peter Bjorklund. All rights reserved.
 *  Licensed under the MIT License. See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package changelogyaml

import (
	"fmt"
	"io"
	"sort"
)

func WriteDocument(root *ChangelogYaml, outputFormatter Formatter, writer io.Writer) error {
	if _, err := fmt.Fprint(writer, outputFormatter.Heading(1, "Changelog")); err != nil {
		return err
	}

	for _, release := range root.Releases {
		completeReleaseLinkURL := fmt.Sprintf("%s%v/releases/tag/%v", githubUrlPrefix, root.Repo, release.Name)
		formattedReleaseLink := outputFormatter.Link(release.Name, completeReleaseLinkURL)

		releaseLink := fmt.Sprintf("%v (%v)", formattedReleaseLink, release.Date)
		releaseHeading := fmt.Sprintf("%s %v", outputFormatter.Emoji("bookmark"), releaseLink)
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
