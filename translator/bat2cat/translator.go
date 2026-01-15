package bat2cat

import (
	"strings"

	"github.com/kluzzebass/reflag/translator"
)

func init() {
	translator.Register(&Translator{})
}

// Translator implements the cat to bat flag translation
type Translator struct{}

func (t *Translator) Name() string        { return "cat2bat" }
func (t *Translator) SourceTool() string  { return "cat" }
func (t *Translator) TargetTool() string  { return "bat" }
func (t *Translator) IncludeInInit() bool { return false }

// Translate converts cat arguments to bat arguments to make bat behave like cat
func (t *Translator) Translate(args []string, mode string) []string {
	return translateFlags(args)
}

// Map of bat short flags to cat equivalents
var flagMap = map[rune]string{
	'n': "-n", // --number → -n (line numbers)
	's': "-s", // --squeeze-blank → -s (squeeze blank lines)
	'u': "-u", // --unbuffered → -u (unbuffered, though bat ignores this)
	'A': "-A", // --show-all → approximates -A (show non-printable)
}

func translateFlags(args []string) []string {
	var result []string
	skipNext := false

	// To make bat behave like cat, we need to:
	// 1. Always add -p (plain style, no decorations)
	// 2. Always add --paging=never (disable pager)
	// 3. Allow default colorization with --color=auto
	result = append(result, "-p", "--paging=never", "--color=auto")

	for i, arg := range args {
		if skipNext {
			skipNext = false
			continue
		}

		// Handle -- separator (everything after is files)
		if arg == "--" {
			result = append(result, args[i:]...)
			break
		}

		// Handle long options
		if strings.HasPrefix(arg, "--") {
			handleLongFlag(arg, args, i, &result, &skipNext)
			continue
		}

		// Handle short options
		if strings.HasPrefix(arg, "-") && len(arg) > 1 && arg[1] != '-' {
			handleShortFlags(arg, args, i, &result, &skipNext)
			continue
		}

		// Regular file argument
		result = append(result, arg)
	}

	return result
}

func handleLongFlag(arg string, args []string, i int, result *[]string, skipNext *bool) {
	// Handle --option=value format
	if idx := strings.Index(arg, "="); idx != -1 {
		opt := arg[:idx]
		val := arg[idx+1:]

		switch opt {
		case "--number":
			*result = append(*result, "-n")
		case "--squeeze-blank":
			*result = append(*result, "-s")
		case "--show-all":
			*result = append(*result, "-A")
		case "--file-name":
			// cat doesn't have this, ignore
		case "--language", "--highlight-line", "--diff-context", "--tabs", "--wrap",
			"--terminal-width", "--color", "--italic-text", "--decorations", "--paging",
			"--pager", "--map-syntax", "--ignored-suffix", "--theme", "--theme-light",
			"--theme-dark", "--style", "--line-range", "--squeeze-limit", "--strip-ansi",
			"--nonprintable-notation", "--binary":
			// These are bat-specific features that cat doesn't have
			// They're overridden by our plain mode settings
		default:
			// Unknown option, might be a file starting with --
			*result = append(*result, arg)
		}
		// Suppress unused variable warning
		_ = val
		return
	}

	// Handle flags without =value
	switch arg {
	case "--number":
		*result = append(*result, "-n")
	case "--squeeze-blank":
		*result = append(*result, "-s")
	case "--show-all":
		*result = append(*result, "-A")
	case "--unbuffered":
		*result = append(*result, "-u")
	case "--plain", "--force-colorization", "--diff", "--list-themes", "--list-languages",
		"--chop-long-lines", "--diagnostic", "--acknowledgements", "--set-terminal-title",
		"--help", "--version":
		// These are bat-specific, ignore or they're already handled
	default:
		// Check if next arg is a value for this flag
		if i+1 < len(args) && !strings.HasPrefix(args[i+1], "-") {
			switch arg {
			case "--language", "--highlight-line", "--file-name", "--diff-context",
				"--tabs", "--wrap", "--terminal-width", "--color", "--italic-text",
				"--decorations", "--paging", "--pager", "--map-syntax", "--ignored-suffix",
				"--theme", "--theme-light", "--theme-dark", "--style", "--line-range",
				"--squeeze-limit", "--strip-ansi", "--nonprintable-notation", "--binary",
				"--completion":
				// These take values but are bat-specific, skip both flag and value
				*skipNext = true
			default:
				// Unknown flag with potential value, keep it (might be a file)
				*result = append(*result, arg)
			}
		} else {
			// Unknown flag without value, might be a file
			*result = append(*result, arg)
		}
	}
}

func handleShortFlags(arg string, args []string, i int, result *[]string, skipNext *bool) {
	flags := arg[1:] // Remove leading dash

	// Check for combined flags like -pp or -ns
	for j, flag := range flags {
		if mapped, ok := flagMap[flag]; ok {
			*result = append(*result, mapped)
		} else {
			// Flags without direct mapping or bat-specific flags
			switch flag {
			case 'p':
				// -p (plain) is already added by default, ignore
			case 'l', 'H', 'm':
				// These flags take values, skip the next argument
				// Only skip if this is the last flag in a combined set
				if j == len(flags)-1 && i+1 < len(args) {
					*skipNext = true
				}
			case 'd', 'f', 'L', 'r', 'S':
				// These are bat-specific flags that don't take values, ignore
			case 'V':
				// -V (version), ignore
			case 'h':
				// -h (help), ignore
			default:
				// Unknown single char flag
				// Could be a typo or actual flag, preserve it
				*result = append(*result, "-"+string(flag))
			}
		}
	}
}
