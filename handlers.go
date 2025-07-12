// Copyright 2017-2025 Seamia Corporation. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"errors"
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/seamia/libs/printer"
)

func funcPrior(all string, media printer.Printer, args []string) error {
	return funcAfterEx(all, media, args, false)
}

func funcAfter(all string, media printer.Printer, args []string) error {
	return funcAfterEx(all, media, args, true)
}

func funcAfterEx(all string, media printer.Printer, args []string, removePrefix bool) error {
	after := ""
	indexFunc := strings.Index

	if len(args) > 0 {
		after = args[0]
		args = args[1:]
	} else {
		return errors.New("missing required argument")
	}

	if len(args) > 0 {
		if args[0] == "/i" {
			indexFunc = IndexCaseInsensitive
		} else {
			return fmt.Errorf("unknown argument (%v)", args[0])
		}
	}

	prefix := len(after)
	if !removePrefix {
		prefix = 0
	}

	for _, line := range strings.Split(all, cr) {
		if index := indexFunc(line, after); index >= 0 {
			remains := line[index+prefix:]
			media("%s\n", remains)
		}
	}
	return nil
}

func funcBefore(all string, media printer.Printer, args []string) error {
	before := ""
	indexFunc := strings.Index

	if len(args) > 0 {
		before = args[0]
		args = args[1:]
	} else {
		return errors.New("missing required argument")
	}

	if len(args) > 0 {
		if args[0] == "/i" {
			indexFunc = IndexCaseInsensitive
		} else {
			return fmt.Errorf("unknown argument (%v)", args[0])
		}
	}

	for _, line := range strings.Split(all, cr) {
		if index := indexFunc(line, before); index >= 0 {
			remains := line[:index]
			media("%s\n", remains)
		}
	}
	return nil
}

func funcEmpty(all string, media printer.Printer, args []string) error {
	trimFunc := func(i string) string { return i }

	if len(args) > 0 {
		if args[0] == "/i" {
			trimFunc = func(i string) string {
				return strings.Trim(i, " \t\r")
			}
		} else {
			return fmt.Errorf("unknown argument (%v)", args[0])
		}
	}

	for _, line := range strings.Split(all, cr) {
		if len(trimFunc(line)) > 0 {
			media("%s\n", line)
		}
	}
	return nil
}

func funcFirst(all string, media printer.Printer, args []string) error {

	if cutOff, err := getInteger(args); err == nil {
		for _, line := range strings.Split(all, cr) {
			if len(line) > cutOff {
				line = line[cutOff:]
			} else {
				line = ""
			}
			media("%s\n", line)

		}
		return nil
	} else {
		return err
	}
}

func funcLast(all string, media printer.Printer, args []string) error {

	if cutOff, err := getInteger(args); err == nil {
		for _, line := range strings.Split(all, cr) {
			if len(line) > cutOff {
				line = line[:cutOff]
			}
			media("%s\n", line)

		}
		return nil
	} else {
		return err
	}
}

func funcTop(all string, media printer.Printer, args []string) error {

	if cutOff, err := getInteger(args); err == nil {
		lines := strings.Split(all, cr)
		if len(lines) > cutOff {
			lines = lines[:cutOff]
		}
		for _, line := range lines {
			media("%s\n", line)
		}
		return nil
	} else {
		return err
	}
}

func funcBottom(all string, media printer.Printer, args []string) error {

	if cutOff, err := getInteger(args); err == nil {
		lines := strings.Split(all, cr)
		if len(lines) > cutOff {
			lines = lines[len(lines)-cutOff:]
		}
		for _, line := range lines {
			media("%s\n", line)
		}
		return nil
	} else {
		return err
	}
}

func funcSpace(all string, media printer.Printer, args []string) error {
	for _, line := range strings.Split(all, cr) {
		line = strings.Trim(line, " \t\r")
		media("%s\n", line)

	}
	return nil
}

func funcSort(all string, media printer.Printer, args []string) error {
	lines := strings.Split(all, cr)
	sortFunc := func(i, j int) bool {
		return lines[i] < lines[j]
	}

	if len(args) > 0 {
		if args[0] == "/i" {
			sortFunc = func(i, j int) bool {
				return lines[i] > lines[j]
			}
		} else {
			return fmt.Errorf("unknown argument (%v)", args[0])
		}
	}

	sort.Slice(lines, sortFunc)

	for _, line := range lines {
		media("%s\n", line)
	}
	return nil
}

func funcDedupe(all string, media printer.Printer, args []string) error {
	lines := strings.Split(all, cr)
	previous := ""
	for _, line := range lines {
		if previous != line {
			media("%s\n", line)
			previous = line
		}
	}
	return nil
}

func constructFindFunc(pattern string) Find {

	var find Find

	search := strings.TrimSuffix(pattern, star)
	search = strings.TrimPrefix(search, star)

	if strings.HasPrefix(pattern, star) && strings.HasSuffix(pattern, star) {
		find = func(s string) bool {
			return strings.Contains(s, search)
		}
	} else if strings.HasPrefix(pattern, star) {
		find = func(s string) bool {
			return strings.HasSuffix(s, search)
		}
	} else if strings.HasSuffix(pattern, star) {
		find = func(s string) bool {
			return strings.HasPrefix(s, search)
		}
	} else if len(pattern) != 0 {
		find = func(s string) bool {
			return s == search
		}
	} else {
		find = func(s string) bool {
			return len(strings.TrimSpace(s)) == 0
		}
	}

	return find
}

// discard lines that match specified pattern
func funcDiscard(all string, media printer.Printer, args []string) error {
	pattern := ""
	if len(args) > 0 {
		pattern = args[0]
	}

	find := constructFindFunc(pattern)
	lines := strings.Split(all, cr)
	for _, line := range lines {
		line = strings.TrimSuffix(line, "\r")
		if !find(line) {
			media("%s\n", line)
		}
	}
	return nil
}

// use only lines that match specified pattern
func funcRetain(all string, media printer.Printer, args []string) error {
	pattern := ""
	if len(args) > 0 {
		pattern = args[0]
	}

	find := constructFindFunc(pattern)
	lines := strings.Split(all, cr)
	for _, line := range lines {
		line = strings.TrimSuffix(line, "\r")
		if find(line) {
			media("%s\n", line)
		}
	}
	return nil
}

// "cut exact" - exclude all the line that are equal to specified text
func funcExact(all string, media printer.Printer, args []string) error {
	if len(args) != 1 {
		return errors.New("missing required one argument")
	}
	exact := args[0]
	
	for _, line := range strings.Split(all, cr) {
		trim := strings.Trim(line, " \t\r")
		if trim != exact {
			media("%s\n", line)
		}
	}
	return nil
}

func funcPrepend(all string, media printer.Printer, args []string) error {
	if len(args) != 1 {
		return errors.New("too few/many arguments")
	}

	prefix := args[0]
	lines := strings.Split(all, cr)
	for _, line := range lines {
		media("%s%s\n", prefix, line)
	}
	return nil
}

func funcAppend(all string, media printer.Printer, args []string) error {
	if len(args) != 1 {
		return errors.New("too few/many arguments")
	}

	suffix := args[0]
	lines := strings.Split(all, cr)
	for _, line := range lines {
		media("%s%s\n", line, suffix)
	}
	return nil
}

func funcScript(all string, media printer.Printer, args []string) error {
	if len(args) != 1 {
		return errors.New("too few/many arguments")
	}

	script := args[0]
	raw, err := os.ReadFile(script)
	if err != nil {
		return fmt.Errorf("failed to open script (%s), due to: %v", script, err)
	}
	for i, line := range strings.Split(string(raw), cr) {
		line = strings.Trim(line, " \t\r\n")
		if strings.HasPrefix(line, "#") {
			// it is a comment - disregard
			continue
		}

		params := strings.Split(line, " ")
		command := params[0]
		params = params[1:]

		operation, found := subFunc[strings.ToLower(command)]
		if !found || operation.Processor == nil {
			return fmt.Errorf("unknown command (%s) in script (%s) at line %v", command, script, i)
		}

		panic("NIY")
	}

	/*
		lines := strings.Split(all, cr)
		for _, line := range lines {
			media("%s%s\n", line, suffix)
		}
	*/
	return nil
}

type Differentiator func(s string) string

func getDifferentiator(separator string, before, first bool) Differentiator {
	return func(s string) string {
		parts := strings.Split(s, separator)
		if first {
			if before {
				// before.first
				return parts[0]
			} else {
				// after.first
				return strings.Join(parts[1:], separator)
			}
		} else {
			// last
			if before {
				// before.last
				if len(parts) < 2 {
					return ""
				} else {
					return strings.Join(parts[:len(parts)-1], separator)
				}
			} else {
				// after.last
				return parts[len(parts)-1]
			}
		}
	}
}

/*
func getPart(line, separator string, before, first bool) string {
	parts := strings.Split(line, separator)
	if first {
		if before {
			// before.first
			return parts[0]
		} else {
			// after.first
			return strings.Join(parts[1:], separator)
		}
	} else {
		// last
		if before {
			// before.last
			if len(parts) < 2 {
				return ""
			} else {
				return strings.Join(parts[:len(parts)-1], separator)
			}
		} else {
			// after.last
			return parts[len(parts)-1]
		}
	}
}
*/

func funcBeforeOrAfterFirstOrLast(all string, media printer.Printer, args []string, before, first bool) error {

	separator := ""
	if len(args) == 1 {
		separator = args[0]
	} else {
		return errors.New("missing required argument")
	}

	if len(separator) == 0 {
		return errors.New("empty separator")
	}

	differentiator := getDifferentiator(separator, before, first)

	for _, line := range strings.Split(all, cr) {
		if modified := differentiator(line); len(modified) > 0 {
			media("%s\n", modified)
		}
	}
	return nil
}

func funcAfterFirst(all string, media printer.Printer, args []string) error {
	return funcBeforeOrAfterFirstOrLast(all, media, args, false, true)
}

func funcAfterLast(all string, media printer.Printer, args []string) error {
	return funcBeforeOrAfterFirstOrLast(all, media, args, false, false)
}

func funcBeforeFirst(all string, media printer.Printer, args []string) error {
	return funcBeforeOrAfterFirstOrLast(all, media, args, true, true)
}

func funcBeforeLast(all string, media printer.Printer, args []string) error {
	return funcBeforeOrAfterFirstOrLast(all, media, args, true, false)
}

func funcHeader(all string, media printer.Printer, args []string) error {
	if len(args) < 1 {
		return errors.New("missing required argument")
	}

	header := strings.Join(args, cr)
	media("%s%s%s", header, cr, all)
	return nil
}

func funcFooter(all string, media printer.Printer, args []string) error {
	if len(args) < 1 {
		return errors.New("missing required argument")
	}

	footer := strings.Join(args, cr)
	media("%s%s%s", all, cr, footer)
	return nil
}
