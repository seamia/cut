// Copyright 2017-2025 Seamia Corporation. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"errors"
	"regexp"
	"strings"

	"github.com/seamia/libs/printer"
)

func funcReplace(input string, media printer.Printer, args []string) error {
	args, caseInsensitive := isCaseInsensitive(args)

	if len(args) != 2 {
		return errors.New("Requires 2 arguments: from to")
	}

	from := args[0]
	to := args[1]

	replace := func(s string) string {
		return strings.ReplaceAll(s, from, to)
	}

	if caseInsensitive {
		re := regexp.MustCompile(`(?i)` + from)
		replace = func(s string) string {
			return re.ReplaceAllString(s, to)
		}
	}

	lines := strings.Split(input, cr)
	for _, line := range lines {
		line = replace(line)
		media("%s\n", line)
	}
	return nil
}
