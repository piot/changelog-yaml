/*---------------------------------------------------------------------------------------------
 *  Copyright (c) Peter Bjorklund. All rights reserved.
 *  Licensed under the MIT License. See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package changelogyaml

import (
	"fmt"
	"io"
)

func convertTextLine(line string, repoShortUrl string, formatter Formatter) (string, error) {
	line, err := replacePullRequestLink(line, repoShortUrl, formatter)
	if err != nil {
		return "", err
	}

	line = replaceCommitHashLink(line, repoShortUrl, formatter)

	line = replaceAtProfileLink(line, formatter)

	return line, nil
}

func linesForRepo(repoShortUrl string, strings []string, categoryType CategoryType, formatter Formatter,
	writer io.Writer) error {
	for _, line := range strings {
		categoryInfo := infoFromCategoryName(categoryType)
		prefix := formatter.Emoji(categoryInfo.EmojiName)

		if categoryType == Breaking {
			prefix += fmt.Sprintf("[%v]", categoryInfo.Name)
		}

		line, err := convertTextLine(line, repoShortUrl, formatter)
		if err != nil {
			return err
		}

		completeLine := prefix + " " + line

		if _, err := fmt.Fprint(writer, formatter.BulletPoint(completeLine)); err != nil {
			return err
		}
	}

	return nil
}

type LineInfo struct {
	Category CategoryType
	Lines    []string
}

func textLinesForTheRepo(repoShortUrl string, repoChanges *Changes, formatter Formatter, writer io.Writer) error {
	lines := []LineInfo{
		{Unreleased, repoChanges.Unreleased},
		{Breaking, repoChanges.Breaking},
		{Added, repoChanges.Added},
		{Fixed, repoChanges.Fixed},
		{Workaround, repoChanges.Workaround},
		{Changed, repoChanges.Changed},
		{Removed, repoChanges.Removed},
		{Improved, repoChanges.Improved},
		{Docs, repoChanges.Docs},
		{Tests, repoChanges.Tests},
		{Refactored, repoChanges.Refactored},
		{Deprecated, repoChanges.Deprecated},
		{Experimental, repoChanges.Experimental},
		{Noted, repoChanges.Noted},
		{Performance, repoChanges.Performance},
		{Style, repoChanges.Style},
	}

	for _, line := range lines {
		if err := linesForRepo(repoShortUrl, line.Lines, line.Category, formatter, writer); err != nil {
			return err
		}
	}

	return nil
}
