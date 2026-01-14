package less2moor

import (
	"strings"

	"github.com/kluzzebass/reflag/translator"
)

func init() {
	translator.Register(&Translator{})
}

// Translator implements the less to moor flag translation
type Translator struct{}

func (t *Translator) Name() string        { return "less2moor" }
func (t *Translator) SourceTool() string  { return "less" }
func (t *Translator) TargetTool() string  { return "moor" }
func (t *Translator) IncludeInInit() bool { return false }

// Translate converts less arguments to moor arguments
func (t *Translator) Translate(args []string, mode string) []string {
	return translateFlags(args)
}

// Simple 1:1 flag mappings from less to moor
var flagMap = map[rune][]string{
	// Display options
	'S': {"--wrap=false"}, // -S: chop long lines (moor wraps by default)
	'N': {"--no-linenumbers"},
	'F': {"--follow"}, // -F: follow mode (like tail -f)

	// Quit behavior
	'e': {"--quit-if-one-screen"}, // -e: quit at EOF on second attempt
	'E': {"--quit-if-one-screen"}, // -E: quit at EOF immediately
	'f': {},                       // -f: force open non-regular files (moor handles this)
	'X': {"--no-clear-on-exit"},   // -X: don't clear screen on exit
	'K': {},                       // -K: exit on Ctrl-C (default in moor)

	// Search options
	'i': {}, // -i: case-insensitive search (moor handles interactively)
	'I': {}, // -I: never case-insensitive (moor handles interactively)
	'g': {}, // -g: highlight only current match
	'G': {}, // -G: no highlighting (moor handles highlighting differently)
	'W': {}, // -W: highlight first unread line after forward
	'w': {}, // -w: highlight first unread line after page
	's': {}, // -s: squeeze blank lines (no direct moor equivalent)
	'r': {}, // -r: raw control chars (moor handles this automatically)
	'R': {}, // -R: ANSI color support (moor supports this by default)
	'q': {}, // -q: quiet (no bell)
	'Q': {}, // -Q: completely quiet

	// Tab handling
	// -x is handled specially as it takes an argument

	// Line numbers
	'n': {}, // -n: suppress line numbers (moor doesn't show by default)
	'J': {}, // -J: status column (not in moor)
	'j': {}, // -j: target line (not in moor)

	// Display size
	'z': {}, // -z: window size (not in moor)
	'y': {}, // -y: max forward scroll (not in moor)

	// Scrolling
	'c': {}, // -c: repaint from top (moor handles automatically)
	'C': {}, // -C: clear before repaint (moor handles automatically)
	'd': {}, // -d: suppress error for dumb terminal (not needed)
	'u': {}, // -u: backspace/return as printable (moor handles automatically)
	'U': {}, // -U: treat backspaces as control chars (moor handles automatically)

	// Misc
	'h': {},           // -h: max backward scroll (not in moor)
	'p': {},           // -p: start at pattern (would need special handling)
	't': {},           // -t: tag (not in moor)
	'T': {},           // -T: tags file (not in moor)
	'k': {},           // -k: lesskey file (not in moor)
	'o': {},           // -o: log file (not in moor)
	'O': {},           // -O: log file, overwrite (not in moor)
	'V': {"-version"}, // -V: version
	'?': {},           // -?: help (moor has --help)
	'm': {},           // -m: medium prompt (not in moor)
	'M': {},           // -M: long prompt (not in moor)
	'P': {},           // -P: prompt (not in moor)
	'a': {},           // -a: search after EOF (not in moor)
	'A': {},           // -A: no search after EOF (not in moor)
	'b': {},           // -b: buffer size (not in moor)
	'B': {},           // -B: auto buffer (not in moor)
	'D': {},           // -D: color descriptor (moor uses --style)
	'#': {"--shift"},  // -#: horizontal scroll amount
	'~': {},           // -~: blank lines after EOF (moor handles differently)
	'L': {},           // -L: ignore LESSOPEN (not in moor)
	'v': {},           // -v: use vi (not in moor)
}

// Long option mappings from less to moor
var longFlagMap = map[string][]string{
	"--quit-if-one-screen":  {"--quit-if-one-screen"},
	"--no-init":             {"--no-clear-on-exit"},
	"--chop-long-lines":     {"--wrap=false"},
	"--RAW-CONTROL-CHARS":   {}, // moor handles ANSI by default
	"--raw-control-chars":   {}, // moor handles control chars automatically
	"--squeeze-blank-lines": {}, // no direct equivalent
	"--follow-name":         {"--follow"},
	"--SILENT":              {}, // no bell in moor anyway
	"--silent":              {}, // no bell in moor anyway
	"--QUIET":               {}, // no bell in moor anyway
	"--quiet":               {}, // no bell in moor anyway
	"--version":             {"-version"},
	"--help":                {},            // moor has --help
	"--tabs":                {"-tab-size"}, // handled specially with argument
	"--mouse":               {"-mousemode=scroll"},
	"--MOUSE":               {"-mousemode=scroll"},
	"--no-keypad":           {}, // not relevant
	"--use-color":           {}, // moor uses color by default
	"--tilde":               {}, // moor handles EOF display differently
	"--shift":               {"-shift"},
	"--hilite-unread":       {}, // no equivalent
	"--HILITE-UNREAD":       {}, // no equivalent
	"--window":              {}, // no equivalent
	"--max-forw-scroll":     {}, // no equivalent
	"--line-num-width":      {}, // no equivalent
	"--status-col-width":    {}, // no equivalent
	"--underline-special":   {}, // moor handles automatically
	"--UNDERLINE-SPECIAL":   {}, // moor handles automatically
}

// Long flag prefixes that take values with = syntax
var longFlagPrefixes = []string{
	"--tabs=",
	"--tag=",
	"--tag-file=",
	"--quotes=",
	"--shift=",
	"--wheel-lines=",
	"--window=",
	"--max-forw-scroll=",
	"--line-num-width=",
	"--status-col-width=",
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
			// Other + commands like +/pattern aren't supported in moor
			continue
		}

		// Check long flag prefixes with values
		if inOptions && strings.HasPrefix(arg, "--") {
			handled := false
			for _, prefix := range longFlagPrefixes {
				if strings.HasPrefix(arg, prefix) {
					value := strings.TrimPrefix(arg, prefix)
					handled = true

					switch prefix {
					case "--tabs=":
						result = append(result, "-tab-size="+value)
					case "--shift=":
						result = append(result, "-shift="+value)
					case "--wheel-lines=":
						// moor doesn't have this option
					default:
						// Other options don't have moor equivalents
					}
					break
				}
			}
			if handled {
				continue
			}

			// Check for exact long flag matches
			if mapped, ok := longFlagMap[arg]; ok {
				result = append(result, mapped...)
				continue
			}

			// Unknown long flag - pass through (moor might handle it)
			result = append(result, arg)
			continue
		}

		// Short flags
		if inOptions && strings.HasPrefix(arg, "-") && len(arg) > 1 && arg[1] != '-' {
			// Check if it's a flag with an argument attached (like -x8)
			firstRune := rune(arg[1])
			if firstRune == 'x' && len(arg) > 2 {
				// -x with tab size
				tabSize := arg[2:]
				result = append(result, "-tab-size="+tabSize)
				continue
			}

			// Check for flags that take separate arguments
			if firstRune == 't' || firstRune == 'T' || firstRune == 'p' ||
				firstRune == 'P' || firstRune == 'o' || firstRune == 'O' ||
				firstRune == 'k' || firstRune == 'D' {
				// These flags take arguments in less but aren't supported in moor
				// Skip the flag and its argument
				if len(arg) == 2 && i+1 < len(args) {
					i++ // skip next arg
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
