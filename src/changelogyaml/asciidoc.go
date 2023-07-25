/*---------------------------------------------------------------------------------------------
 *  Copyright (c) Peter Bjorklund. All rights reserved.
 *  Licensed under the MIT License. See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package changelogyaml

import (
	"fmt"
	"strings"
)

type AsciiDocFormatter struct {
}

func (a *AsciiDocFormatter) Heading(level int, header string) string {
	return strings.Repeat("=", level) + " " + header + "\n\n"
}

func (m *AsciiDocFormatter) BulletPoint(text string) string {
	return "* " + text + "\n"
}

func (m *AsciiDocFormatter) Emoji(name string) string {
	unicodeInt := emojiNameToUnicode(name)
	return fmt.Sprintf("&#x%X;", unicodeInt)
}

func (m *AsciiDocFormatter) Link(name string, link string) string {
	return fmt.Sprintf("%s[%s]", link, name)
}

func AdmonitionTypeToAsciiDocName(admonitionType AdmonitionType) string {
	switch admonitionType {
	case Note:
		return "NOTE"
	case Important:
		return "IMPORTANT"
	case Warning:
		return "WARNING"
	}
	// CAUTION and TIP is not supported yet

	panic(fmt.Errorf("unknown admonition"))
}

func (m *AsciiDocFormatter) Admonition(admonitionType AdmonitionType, text string) string {
	return fmt.Sprintf("%s: %s", AdmonitionTypeToAsciiDocName(admonitionType), text)
}
