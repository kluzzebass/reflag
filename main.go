package main

import (
	"fmt"
	"os"
	"strings"
)

// Flags that need --reverse in eza to match ls default behavior
// (ls shows newest/largest first, eza shows oldest/smallest first)
var reverseNeeded = map[rune]bool{
	't': true, // time sort: ls=newest first, eza=oldest first
	'S': true, // size sort: ls=largest first, eza=smallest first
}

// Simple 1:1 flag mappings
var flagMap = map[rune][]string{
	'l': {"-l"},
	'a': {"-a"},
	'A': {"-A"},
	'h': {}, // eza uses human-readable by default
	't': {"--sort=modified"},
	'S': {"--sort=size"},
	'R': {"--recurse"},
	'1': {"-1"},
	'd': {"-d"},
	'F': {"-F"},
	'G': {}, // color is default in eza
	'i': {"--inode"},
	's': {"--blocksize"},
	'n': {"--numeric"},
	'o': {"-l", "--no-group"},
	'g': {"-l", "--no-user"},
	'p': {"--classify"},
	'c': {"--sort=changed"},
	'u': {"--sort=accessed"},
	'x': {"--across"},
	'C': {"--grid"},
	'T': {"--tree"},
}

func translateFlags(args []string) []string {
	var ezaArgs []string
	var paths []string
	userReverse := false
	needsReverse := false

	for _, arg := range args {
		if strings.HasPrefix(arg, "--") {
			// Long option - pass through
			if arg == "--reverse" {
				userReverse = true
			} else {
				ezaArgs = append(ezaArgs, arg)
			}
		} else if strings.HasPrefix(arg, "-") && len(arg) > 1 {
			// Short options - translate each character
			for _, c := range arg[1:] {
				if c == 'r' {
					userReverse = true
					continue
				}
				if reverseNeeded[c] {
					needsReverse = true
				}
				if mapped, ok := flagMap[c]; ok {
					ezaArgs = append(ezaArgs, mapped...)
				} else {
					// Unknown flag - try passing it through
					ezaArgs = append(ezaArgs, "-"+string(c))
				}
			}
		} else {
			// Not a flag - it's a path
			paths = append(paths, arg)
		}
	}

	// XOR logic: reverse if exactly one of (needsReverse, userReverse) is true
	// - ls -lt → need reverse to get newest first (needsReverse=true, userReverse=false) → add --reverse
	// - ls -ltr → user wants oldest first (needsReverse=true, userReverse=true) → don't add --reverse
	// - ls -lr → user wants reverse alpha (needsReverse=false, userReverse=true) → add --reverse
	if needsReverse != userReverse {
		ezaArgs = append(ezaArgs, "--reverse")
	}

	// Deduplicate flags
	seen := make(map[string]bool)
	var deduped []string
	for _, f := range ezaArgs {
		if !seen[f] {
			seen[f] = true
			deduped = append(deduped, f)
		}
	}

	return append(deduped, paths...)
}

func shellQuote(s string) string {
	if strings.ContainsAny(s, " \t\n\"'\\$`!") {
		return "'" + strings.ReplaceAll(s, "'", "'\"'\"'") + "'"
	}
	return s
}

func main() {
	args := os.Args[1:]
	ezaArgs := translateFlags(args)

	// Build and print the command
	parts := make([]string, len(ezaArgs)+1)
	parts[0] = "eza"
	for i, arg := range ezaArgs {
		parts[i+1] = shellQuote(arg)
	}
	fmt.Println(strings.Join(parts, " "))
}
