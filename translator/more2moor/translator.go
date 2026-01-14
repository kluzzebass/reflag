package more2moor

import (
	"strings"

	"github.com/kluzzebass/reflag/translator"
)

func init() {
	translator.Register(&Translator{})
}

// Translator implements the more to moor flag translation
type Translator struct{}

func (t *Translator) Name() string        { return "more2moor" }
func (t *Translator) SourceTool() string  { return "more" }
func (t *Translator) TargetTool() string  { return "moor" }
func (t *Translator) IncludeInInit() bool { return false }

// Translate converts more arguments to moor arguments
func (t *Translator) Translate(args []string, mode string) []string {
	return translateFlags(args)
}

// Simple 1:1 flag mappings from more to moor
// more has fewer flags than less, so this is a simpler translator
var flagMap = map[rune][]string{
	// Display options
	'd': {}, // -d: display help prompt on invalid key (moor handles interactively)
	'l': {}, // -l: do not pause at form feeds (no moor equivalent)
	'f': {}, // -f: count logical lines (no moor equivalent)
	'p': {}, // -p: clear screen before display (moor handles automatically)
	'c': {}, // -c: draw from top of screen (moor handles automatically)
	's': {}, // -s: squeeze blank lines (no moor equivalent)
	'u': {}, // -u: suppress underlining (moor handles automatically)

	// Exit behavior (GNU more)
	'e': {"--quit-if-one-screen"}, // -e: exit at end of file

	// Version (some implementations)
	'V': {"-version"}, // -V: version
}

// Long option mappings from more to moor
var longFlagMap = map[string][]string{
	"--help":       {}, // moor has --help
	"--version":    {"-version"},
	"--exit-on-eof": {"--quit-if-one-screen"},
	"--no-init":    {"--no-clear-on-exit"},
	"--plain":      {}, // -p: suppress underlining (moor handles automatically)
	"--squeeze":    {}, // -s: squeeze blank lines (no moor equivalent)
	"--print-over": {}, // -p: clear and display (moor handles automatically)
	"--clean-print": {}, // -c: draw from top (moor handles automatically)
}

func translateFlags(args []string) []string {
	var result []string
	var files []string
	var initialCommand string
	inOptions := true

	for i := 0; i < len(args); i++ {
		arg := args[i]

		// Handle end of options marker
		if arg == "--" {
			inOptions = false
			continue
		}

		// Handle + commands (initial commands)
		if inOptions && strings.HasPrefix(arg, "+") {
			// moor supports +linenum for jumping to a line
			if len(arg) > 1 && arg[1] >= '0' && arg[1] <= '9' {
				// Extract line number
				initialCommand = arg
			}
			// +/pattern isn't supported in moor
			continue
		}

		// Check long flags
		if inOptions && strings.HasPrefix(arg, "--") {
			if mapped, ok := longFlagMap[arg]; ok {
				result = append(result, mapped...)
				continue
			}

			// Handle --lines=N (GNU more)
			if strings.HasPrefix(arg, "--lines=") {
				// moor doesn't have a lines option, ignore
				continue
			}

			// Unknown long flag - pass through (moor might handle it)
			result = append(result, arg)
			continue
		}

		// Short flags
		if inOptions && strings.HasPrefix(arg, "-") && len(arg) > 1 && arg[1] != '-' {
			// Check if it's a numeric argument like -10 (number of lines)
			if arg[1] >= '0' && arg[1] <= '9' {
				// -num sets screen size, no moor equivalent
				continue
			}

			// Check for -n flag with separate argument (number of lines)
			if arg == "-n" {
				if i+1 < len(args) {
					i++ // skip next arg (the number)
				}
				continue
			}

			// Process bundled short flags
			for j := 1; j < len(arg); j++ {
				flag := rune(arg[j])
				if mapped, ok := flagMap[flag]; ok {
					result = append(result, mapped...)
				}
				// Unknown flags are silently ignored
			}
			continue
		}

		// Everything else is a file
		files = append(files, arg)
	}

	// Add initial command if present (like +123 for line number)
	if initialCommand != "" {
		result = append(result, initialCommand)
	}

	// Add files at the end
	result = append(result, files...)

	return result
}
