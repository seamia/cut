// Copyright 2017-2025 Seamia Corporation. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/seamia/libs/printer"
)

func funcExclude(input string, media printer.Printer, args []string) error {
	return funcIncludeExclude(input, media, args, false)
}

func funcInclude(input string, media printer.Printer, args []string) error {
	return funcIncludeExclude(input, media, args, true)
}

func funcIncludeExclude(input string, media printer.Printer, args []string, include bool) error {

	fileName := ""
	caseInsensetive := false

	if len(args) == 1 {
		fileName = args[0]
	} else if len(args) == 2 {
		if args[0] == "/i" {
			fileName = args[1]
			caseInsensetive = true
		} else {
			fileName = args[0]
		}
	} else {
		return fmt.Errorf("missing required argument or too many arguments (%v)", args)
	}

	raw, err := os.ReadFile(fileName)
	if err != nil {
		return fmt.Errorf("failed to read file [%v], due to: %v", fileName, err)
	}

	strContains := strings.Contains
	if caseInsensetive {
		strContains = func(s, substr string) bool {
			return strings.Contains(strings.ToLower(s), substr)
		}
	}

	phrases := []string{}
	for _, line := range strings.Split(string(raw), "\n") {
		line = strings.Trim(line, "\r\n")
		if len(line) > 0 {
			if caseInsensetive {
				phrases = append(phrases, strings.ToLower(line))
			} else {
				phrases = append(phrases, line)
			}
		}
	}

	contains := func(line string) bool {
		for _, phrase := range phrases {
			if strContains(line, phrase) {
				return true
			}
		}
		return false
	}

	for _, line := range strings.Split(input, cr) {
		if contains(line) {
			if include {
				media("%s"+cr, line)
			}
		} else {
			if !include {
				media("%s"+cr, line)
			}
		}
	}
	return nil
}
