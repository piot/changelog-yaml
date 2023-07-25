/*---------------------------------------------------------------------------------------------
 *  Copyright (c) Peter Bjorklund. All rights reserved.
 *  Licensed under the MIT License. See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package changelogyaml

import (
	"fmt"
	"strings"
)

type MarkdownFormatter struct {
}

func (m *MarkdownFormatter) Heading(level int, header string) string {
	return strings.Repeat("#", level) + " " + header + "\n\n"
}

func (m *MarkdownFormatter) BulletPoint(text string) string {
	return "* " + text + "\n"
}

func (m *MarkdownFormatter) Emoji(name string) string {
	return ":" + name + ":"
}

func (m *MarkdownFormatter) Link(name string, link string) string {
	return fmt.Sprintf("[%s](%s)", name, link)
}

func AdmonitionTypeToGithubName(admonitionType AdmonitionType) string {
	switch admonitionType {
	case Note:
		return "NOTE"
	case Important:
		return "IMPORTANT"
	case Warning:
		return "WARNING"
	}

	panic(fmt.Errorf("unknown admonition"))
}

func (m *MarkdownFormatter) Admonition(admonitionType AdmonitionType, text string) string {
	return fmt.Sprintf("> [!%s]\\\n> %s", AdmonitionTypeToGithubName(admonitionType), text)
}
