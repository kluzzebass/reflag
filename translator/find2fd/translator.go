package find2fd

import (
	"regexp"
	"strings"

	"github.com/kluzzebass/reflag/translator"
)

func init() {
	translator.Register(&Translator{})
}

// Translator implements the find to fd flag translation
type Translator struct{}

func (t *Translator) Name() string        { return "find2fd" }
func (t *Translator) SourceTool() string  { return "find" }
func (t *Translator) TargetTool() string  { return "fd" }
func (t *Translator) IncludeInInit() bool { return true }

// Translate converts find arguments to fd arguments
func (t *Translator) Translate(args []string, mode string) []string {
	return translateFlags(args)
}

// Expressions that take a value
var expressionsWithValue = map[string]bool{
	"-name":     true,
	"-iname":    true,
	"-path":     true,
	"-ipath":    true,
	"-regex":    true,
	"-iregex":   true,
	"-type":     true,
	"-maxdepth": true,
	"-mindepth": true,
	"-size":     true,
	"-newer":    true,
	"-mtime":    true,
	"-atime":    true,
	"-ctime":    true,
	"-mmin":     true,
	"-amin":     true,
	"-cmin":     true,
	"-user":     true,
	"-group":    true,
	"-perm":     true,
}

// Expressions to ignore (no fd equivalent or default behavior)
var ignoredExpressions = map[string]bool{
	"-print":  true,
	"-print0": false, // handled specially
	"-a":      true,  // implicit AND
	"-and":    true,
	"-true":   true,
}

func translateFlags(args []string) []string {
	var fdArgs []string
	var pattern string
	var paths []string
	caseInsensitive := false
	skipNext := false

	// First pass: extract paths (arguments before first expression)
	i := 0
	for i < len(args) {
		arg := args[i]
		if strings.HasPrefix(arg, "-") || arg == "!" || arg == "(" || arg == ")" {
			break
		}
		// Skip "." as fd defaults to current directory
		if arg != "." {
			paths = append(paths, arg)
		}
		i++
	}

	// Second pass: process expressions
	for ; i < len(args); i++ {
		if skipNext {
			skipNext = false
			continue
		}

		arg := args[i]

		// Skip logical operators and grouping (fd doesn't support them the same way)
		if arg == "!" || arg == "-not" || arg == "(" || arg == ")" || arg == "-o" || arg == "-or" {
			continue
		}

		if ignoredExpressions[arg] {
			continue
		}

		// Handle expressions with values
		if expressionsWithValue[arg] && i+1 < len(args) {
			val := args[i+1]
			skipNext = true

			switch arg {
			case "-name":
				if pattern == "" {
					pattern = globToRegex(val)
				} else {
					// Multiple -name: fd doesn't support well, use glob
					fdArgs = append(fdArgs, "-g", val)
				}
			case "-iname":
				caseInsensitive = true
				if pattern == "" {
					pattern = globToRegex(val)
				} else {
					fdArgs = append(fdArgs, "-g", val)
				}
			case "-path":
				fdArgs = append(fdArgs, "-p", val)
			case "-ipath":
				fdArgs = append(fdArgs, "-i", "-p", val)
			case "-regex":
				if pattern == "" {
					pattern = val
				}
			case "-iregex":
				caseInsensitive = true
				if pattern == "" {
					pattern = val
				}
			case "-type":
				fdArgs = append(fdArgs, "-t", translateType(val))
			case "-maxdepth":
				fdArgs = append(fdArgs, "-d", val)
			case "-mindepth":
				fdArgs = append(fdArgs, "--min-depth", val)
			case "-size":
				fdArgs = append(fdArgs, "-S", val)
			case "-newer":
				fdArgs = append(fdArgs, "--newer", val)
			case "-mtime":
				fdArgs = append(fdArgs, translateMtime(val)...)
			case "-atime":
				fdArgs = append(fdArgs, translateAtime(val)...)
			case "-ctime":
				fdArgs = append(fdArgs, translateCtime(val)...)
			case "-mmin":
				fdArgs = append(fdArgs, translateMmin(val)...)
			case "-amin":
				fdArgs = append(fdArgs, translateAmin(val)...)
			case "-cmin":
				fdArgs = append(fdArgs, translateCmin(val)...)
			case "-user":
				fdArgs = append(fdArgs, "--owner", val)
			case "-group":
				fdArgs = append(fdArgs, "--owner", ":"+val)
			case "-perm":
				// fd doesn't have direct perm support, skip
			}
			continue
		}

		// Handle standalone expressions
		switch arg {
		case "-print0":
			fdArgs = append(fdArgs, "-0")
		case "-L", "-follow":
			fdArgs = append(fdArgs, "-L")
		case "-H":
			fdArgs = append(fdArgs, "-H")
		case "-P":
			// Default behavior, ignore
		case "-empty":
			fdArgs = append(fdArgs, "-t", "e")
		case "-executable":
			fdArgs = append(fdArgs, "-t", "x")
		case "-xdev", "-mount":
			fdArgs = append(fdArgs, "--one-file-system")
		case "-depth":
			// fd doesn't have depth-first, ignore
		case "-daystart":
			// fd doesn't support, ignore
		case "-delete":
			// Too dangerous to auto-translate
		case "-prune":
			// No direct equivalent
		case "-quit":
			fdArgs = append(fdArgs, "-1")
		}

		// Handle -exec and similar (warn but pass through somehow?)
		if arg == "-exec" || arg == "-execdir" || arg == "-ok" || arg == "-okdir" {
			// Skip until we find ; or +
			for i++; i < len(args); i++ {
				if args[i] == ";" || args[i] == "+" {
					break
				}
			}
		}
	}

	// Build final command - ensure we return empty slice not nil
	result := make([]string, 0)

	if caseInsensitive {
		result = append(result, "-i")
	}

	result = append(result, fdArgs...)

	// Add pattern if we have one, otherwise add match-all pattern if we have paths
	// fd syntax is: fd [PATTERN] [PATH]... - pattern must come before paths
	if pattern != "" {
		result = append(result, pattern)
	} else if len(paths) > 0 {
		// When searching a directory without a pattern, fd needs a match-all pattern
		result = append(result, ".")
	}

	// Add paths
	result = append(result, paths...)

	return result
}

// translateType converts find -type values to fd -t values
func translateType(t string) string {
	switch t {
	case "f":
		return "f"
	case "d":
		return "d"
	case "l":
		return "l"
	case "s":
		return "s" // socket
	case "p":
		return "p" // pipe
	case "b", "c":
		return "f" // block/char devices - approximate as file
	default:
		return t
	}
}

// globToRegex converts a shell glob pattern to a regex
func globToRegex(glob string) string {
	// Handle simple extension patterns specially
	if strings.HasPrefix(glob, "*.") && !strings.ContainsAny(glob[2:], "*?[]") {
		// Simple extension match: *.txt -> \.txt$
		ext := glob[1:] // .txt
		return regexp.QuoteMeta(ext) + "$"
	}

	// General glob to regex conversion
	var result strings.Builder
	for i := 0; i < len(glob); i++ {
		c := glob[i]
		switch c {
		case '*':
			if i+1 < len(glob) && glob[i+1] == '*' {
				// ** matches anything including /
				result.WriteString(".*")
				i++
			} else {
				// * matches anything except /
				result.WriteString("[^/]*")
			}
		case '?':
			result.WriteString("[^/]")
		case '.':
			result.WriteString("\\.")
		case '[':
			// Character class - pass through mostly
			result.WriteByte('[')
			i++
			if i < len(glob) && glob[i] == '!' {
				result.WriteByte('^')
				i++
			}
			for i < len(glob) && glob[i] != ']' {
				result.WriteByte(glob[i])
				i++
			}
			if i < len(glob) {
				result.WriteByte(']')
			}
		case '\\':
			if i+1 < len(glob) {
				result.WriteByte('\\')
				i++
				result.WriteByte(glob[i])
			}
		case '^', '$', '+', '{', '}', '|', '(':
			// Escape regex metacharacters
			result.WriteByte('\\')
			result.WriteByte(c)
		default:
			result.WriteByte(c)
		}
	}
	return result.String()
}

// Time translation helpers

func translateMtime(val string) []string {
	return translateTimeGeneric(val, "--changed-within", "--changed-before")
}

func translateAtime(val string) []string {
	// fd doesn't distinguish atime/mtime well, use changed
	return translateTimeGeneric(val, "--changed-within", "--changed-before")
}

func translateCtime(val string) []string {
	return translateTimeGeneric(val, "--changed-within", "--changed-before")
}

func translateTimeGeneric(val, withinFlag, beforeFlag string) []string {
	if strings.HasPrefix(val, "-") {
		// -N means within last N days
		days := val[1:]
		return []string{withinFlag, days + "d"}
	} else if strings.HasPrefix(val, "+") {
		// +N means more than N days ago
		days := val[1:]
		return []string{beforeFlag, days + "d"}
	} else {
		// Exact N days - approximate as within
		return []string{withinFlag, val + "d"}
	}
}

func translateMmin(val string) []string {
	return translateMinGeneric(val, "--changed-within", "--changed-before")
}

func translateAmin(val string) []string {
	return translateMinGeneric(val, "--changed-within", "--changed-before")
}

func translateCmin(val string) []string {
	return translateMinGeneric(val, "--changed-within", "--changed-before")
}

func translateMinGeneric(val, withinFlag, beforeFlag string) []string {
	if strings.HasPrefix(val, "-") {
		// -N means within last N minutes
		mins := val[1:]
		return []string{withinFlag, mins + "min"}
	} else if strings.HasPrefix(val, "+") {
		// +N means more than N minutes ago
		mins := val[1:]
		return []string{beforeFlag, mins + "min"}
	} else {
		return []string{withinFlag, val + "min"}
	}
}
