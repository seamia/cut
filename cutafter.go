// Copyright 2017-2025 Seamia Corporation. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"bufio"
	"io"
	"os"
	"strings"

	"github.com/seamia/libs/printer"
)

func main() {

	if len(os.Args) == 1 {
		printer.Stderr("Missing operation\n")
		usage()
		return
	}

	{
		// need to do it at runtime, due to "initialization loop" compile error
		prior := subFunc["script"]
		prior.Processor = funcScript
		subFunc["script"] = prior
	}

	operation, found := subFunc[strings.ToLower(os.Args[1])]
	if !found {
		printer.Stderr("Unknown operation: %v\n", os.Args[1])
		usage()
		return
	}

	if operation.Processor == nil {
		usage()
		return
	}

	if err := Run(operation.Processor, os.Args[2:]); err != nil {
		printer.Stderr("Error: %v\n", err)
		os.Exit(errorExitCode)
	}
	os.Exit(normalExitCode)
}

func Run(proc Processor, args []string) error {
	info, err := os.Stdin.Stat()
	if err != nil {
		printer.Stderr("There was an error accessing stdin (%v)\n", err)
		os.Exit(errorExitCode)
	}

	pipe := info.Mode()&os.ModeNamedPipe != 0
	char := info.Mode()&os.ModeCharDevice != 0

	if !pipe || char {
		printer.Stderr("The command is intended to work with pipes.\n")
		os.Exit(errorExitCode)
	}

	reader := bufio.NewReader(os.Stdin)
	var output []rune

	for {
		input, _, err := reader.ReadRune()
		if err != nil && err == io.EOF {
			break
		}
		output = append(output, input)
	}

	all := string(output)
	return proc(all, printer.Stdout, args)
}
