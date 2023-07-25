/*---------------------------------------------------------------------------------------------
 *  Copyright (c) Peter Bjorklund. All rights reserved.
 *  Licensed under the MIT License. See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package main

import (
	"bufio"
	"flag"
	"log"
	"os"

	"github.com/piot/changelog-yaml/changelogyaml"
)

func main() {
	var outputFormat = flag.String("format", "md", "output format: md or adoc")
	flag.Parse()

	var c changelogyaml.ChangelogYaml

	reader := bufio.NewReader(os.Stdin)
	c.ReadYaml(reader)

	var formatter changelogyaml.Formatter
	if *outputFormat == "adoc" || *outputFormat == "asciidoc" {
		formatter = &changelogyaml.AsciiDocFormatter{}
	} else {
		formatter = &changelogyaml.MarkdownFormatter{}
	}

	err := changelogyaml.WriteDocument(&c, formatter, os.Stdout)
	if err != nil {
		log.Println(err)
		os.Exit(-2)
	}
}
