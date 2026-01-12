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
	'c': true, // change time sort: ls=newest first, eza=oldest first
	'u': true, // access time sort: ls=newest first, eza=oldest first
	'U': true, // creation time sort (BSD): ls=newest first, eza=oldest first
}

// Simple 1:1 flag mappings
var flagMap = map[rune][]string{
	// Display format
	'l': {"-l"},                // long format
	'1': {"-1"},                // one entry per line
	'C': {"--grid"},            // multi-column output (default in terminal)
	'x': {"--across"},          // sort grid across rather than down
	'm': {"--oneline"},         // stream output (eza doesn't have comma-separated, use oneline)

	// Show/hide entries
	'a': {"-a"},                // show all including . and ..
	'A': {"-A"},                // show hidden but not . and ..
	'd': {"-d"},                // list directories themselves, not contents
	'R': {"--recurse"},         // recurse into directories

	// Sorting
	't': {"--sort=modified"},   // sort by modification time
	'S': {"--sort=size"},       // sort by size
	'c': {"--sort=changed"},    // sort by change time
	'u': {"--sort=accessed"},   // sort by access time
	'U': {"--sort=created"},    // sort by creation time (BSD)
	'f': {"--sort=none", "-a"}, // unsorted, show all
	'v': {"--sort=name"},       // natural version sort (approximate)

	// File size display
	'h': {},                    // human-readable (default in eza)
	'k': {},                    // 1024-byte blocks (eza handles differently)
	's': {"--blocksize"},       // show allocated blocks

	// Indicators and classification
	'F': {"-F"},                // append file type indicators (*/=>@|)
	'p': {"--classify"},        // append / to directories

	// Long format options
	'i': {"--inode"},           // show inode numbers
	'n': {"--numeric"},         // numeric user/group IDs
	'o': {"-l", "--no-group"},  // long format without group (BSD)
	'g': {"-l", "--no-user"},   // long format without owner (GNU style)
	'O': {"--flags"},           // show file flags (BSD/macOS)
	'e': {},                    // show ACL (no eza equivalent)
	'T': {"--time-style=full-iso"}, // full time info (not tree! ls -T shows full timestamp)
	'@': {"--extended"},        // show extended attributes

	// Symlink handling
	'L': {"-X"},                // dereference symlinks
	'H': {"-X"},                // follow symlinks on command line
	'P': {},                    // don't follow symlinks (default)

	// Color
	'G': {},                    // color output (default in eza)

	// Misc
	'q': {},                    // replace non-printable with ? (no eza equivalent)
	'w': {},                    // raw non-printable (no eza equivalent)
	'b': {},                    // C-style escapes (no eza equivalent)
	'B': {},                    // octal escapes (no eza equivalent)
}

// Long option mappings (ls long options to eza equivalents)
var longFlagMap = map[string][]string{
	"--all":             {"-a"},
	"--almost-all":      {"-A"},
	"--directory":       {"-d"},
	"--recursive":       {"--recurse"},
	"--human-readable":  {}, // default in eza
	"--inode":           {"--inode"},
	"--numeric-uid-gid": {"--numeric"},
	"--classify":        {"-F"},
	"--file-type":       {"--classify"},
	"--dereference":     {"-X"},
	"--no-group":        {"--no-group"},
}

func translateFlags(args []string) []string {
	var ezaArgs []string
	var paths []string
	userReverse := false
	needsReverse := false

	for _, arg := range args {
		if strings.HasPrefix(arg, "--") {
			// Long option handling
			if arg == "--reverse" {
				userReverse = true
			} else if strings.HasPrefix(arg, "--color") {
				// Pass through --color options (--color, --color=auto, etc.)
				ezaArgs = append(ezaArgs, arg)
			} else if strings.HasPrefix(arg, "--sort=") {
				// Pass through sort options
				ezaArgs = append(ezaArgs, arg)
			} else if strings.HasPrefix(arg, "--time=") {
				// Pass through time options
				ezaArgs = append(ezaArgs, arg)
			} else if mapped, ok := longFlagMap[arg]; ok {
				ezaArgs = append(ezaArgs, mapped...)
			} else {
				// Unknown long option - pass through
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
