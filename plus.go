// Copyright 2017-2025 Seamia Corporation. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"errors"
	"fmt"
	"sort"
	"strconv"
	"strings"

	"github.com/seamia/libs/printer"
)

const (
	errorExitCode  = 123
	normalExitCode = 0
	cr             = "\n"
	star           = "*"
)

type (
	Processor func(input string, output printer.Printer, args []string) error
	Find      func(s string) bool

	pair struct {
		Processor
		string
	}
)

var (
	subFunc = map[string]pair{
		"before":  {funcBefore, ""},
		"after":   {funcAfter, ""},
		"first":   {funcFirst, ""},
		"last":    {funcLast, ""},
		"empty":   {funcEmpty, ""},
		"sort":    {funcSort, "\t- sort lines"},
		"space":   {funcSpace, ""},
		"dedupe":  {funcDedupe, "\t- remove consequent duplicate lines"},
		"discard": {funcDiscard, ""},
		"retain":  {funcRetain, ""},
		"prepend": {funcPrepend, ""},
		"prior":   {funcPrior, ""},
		"append":  {funcAppend, ""},
		"top":     {funcTop, "N\t- top N lines"},
		"bottom":  {funcBottom, "N\t- get bottom N lines"},
		"exclude": {funcExclude, "file\t- exclude all lines containing any of the phrases found in file"},
		"include": {funcInclude, "file\t- include only lines containing any of the phrases found in file"},
		"replace": {funcReplace, "from to [/i]\t- find and replace"},
		"script":  {nil, ""}, // funcScript,
		"help":    {nil, ""},
	}
)

func usage() {
	printer.Stderr("Usage:\n")

	list := []string{}
	for key, value := range subFunc {
		if value.Processor != nil {
			list = append(list, key)
		}
	}
	sort.Strings(list)
	for _, name := range list {
		comment := subFunc[name].string
		printer.Stderr("\tcut %s %s\n", name, comment)
	}
	printer.Stderr("see https://github.com/seamia/cut for details\n")
}

func getInteger(args []string) (int, error) {
	if len(args) > 0 {
		if i, err := strconv.Atoi(args[0]); err != nil {
			return 0, fmt.Errorf("failed to process (%s) as int, due to: %v", args[0], err)
		} else if i < 0 {
			return 0, fmt.Errorf("cut off value (%v) cannot be negative", i)
		} else {
			return i, nil
		}
	} else {
		return 0, errors.New("missing required argument")
	}
}

func IndexCaseInsensitive(s, substr string) int {
	return strings.Index(strings.ToLower(s), strings.ToLower(substr))
}
