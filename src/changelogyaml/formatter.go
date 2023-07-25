/*---------------------------------------------------------------------------------------------
 *  Copyright (c) Peter Bjorklund. All rights reserved.
 *  Licensed under the MIT License. See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package changelogyaml

type AdmonitionType uint8

const (
	Note AdmonitionType = iota
	Important
	Warning
)

type Formatter interface {
	Heading(level int, header string) string
	BulletPoint(text string) string
	Emoji(name string) string
	Link(name string, link string) string
	Admonition(admonitionType AdmonitionType, text string) string
}
