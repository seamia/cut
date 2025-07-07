package main

import (
	"slices"
)

func found(args []string, value string) ([]string, bool) {
	if index := slices.Index(args, value); index >= 0 {
		args = append(args[:index], args[index+1:]...)
		return args, true
	}

	return args, false
}

func isCaseInsensitive(args []string) ([]string, bool) {
	return found(args, "/i")
}
