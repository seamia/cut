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
		"before":       {funcBefore, ""},
		"after":        {funcAfter, ""},
		"after.first":  {funcAfterFirst, "sep\t- return what is left after first occurrence of `sep`"},
		"after.last":   {funcAfterLast, "sep\t- return what is left after last occurrence of `sep`"},
		"before.first": {funcBeforeFirst, "sep\t- return what is found before first occurrence of `sep`"},
		"before.last":  {funcBeforeLast, "sep\t- return what is found before last occurrence of `sep`"},
		"first":        {funcFirst, "N\t- show only N first characters"},
		"last":         {funcLast, "N\t- show only N last characters"},
		"empty":        {funcEmpty, ""},
		"sort":         {funcSort, "\t- sort lines"},
		"space":        {funcSpace, "\t- trim spaces from head and tail"},
		"dedupe":       {funcDedupe, "\t- remove consequent duplicate lines"},
		"discard":      {funcDiscard, "pattern\t- discard the lines that match specified pattern"},
		"retain":       {funcRetain, "pattern\t- retain only the lines that match specified pattern"},
		"exact":        {funcExact, "what\t- exclude lines that match `what`"},
		"prepend":      {funcPrepend, "what\t- prepend `what` to original string"},
		"prior":        {funcPrior, ""},
		"header":       {funcHeader, "prefix\t- prepend header `prefix` to the whole input blob"},
		"footer":       {funcFooter, "suffix\t- append footer `suffix` to the whole input blob"},
		"append":       {funcAppend, "what\t- append `what` to original string"},
		"top":          {funcTop, "N\t- top N lines"},
		"bottom":       {funcBottom, "N\t- get bottom N lines"},
		"exclude":      {funcExclude, "file\t- exclude all lines containing any of the phrases found in file"},
		"include":      {funcInclude, "file\t- include only lines containing any of the phrases found in file"},
		"replace":      {funcReplace, "from to [/i]\t- find and replace"},
		"script":       {nil, ""}, // funcScript,
		"help":         {nil, ""},
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
