package ps2procs

import (
	"strings"

	"github.com/kluzzebass/reflag/translator"
)

func init() {
	translator.Register(&Translator{})
}

// Translator implements the ps to procs flag translation
type Translator struct{}

func (t *Translator) Name() string        { return "ps2procs" }
func (t *Translator) SourceTool() string  { return "ps" }
func (t *Translator) TargetTool() string  { return "procs" }
func (t *Translator) IncludeInInit() bool { return true }

// Translate converts ps arguments to procs arguments
func (t *Translator) Translate(args []string, mode string) []string {
	return translateFlags(args)
}

// Flags to ignore (procs shows all processes by default with good format)
var ignoredFlags = map[string]bool{
	"-e": true, // all processes (procs default)
	"-A": true, // all processes (procs default)
	"-f": true, // full format (procs default)
	"-l": true, // long format
	"-j": true, // job format
	"-v": true, // virtual memory format
	"-w": true, // wide output
	"-r": true, // running only (no equivalent)
	"-x": true, // BSD: include processes without tty
	"-a": true, // BSD: all with tty except session leaders
	"-d": true, // all except session leaders
	"-N": true, // negate selection
	"-T": true, // this terminal
	"-g": true, // session or group leaders
	"-s": true, // session leaders
	"-t": true, // by tty
	"-c": true, // command name only
	"-m": true, // threads
	"-L": true, // threads
}

// Column name mappings from ps to procs
var columnMap = map[string]string{
	"pid":      "pid",
	"ppid":     "ppid",
	"uid":      "uid",
	"user":     "user",
	"gid":      "gid",
	"group":    "group",
	"comm":     "command",
	"cmd":      "command",
	"command":  "command",
	"args":     "command",
	"%cpu":     "cpu",
	"pcpu":     "cpu",
	"cpu":      "cpu",
	"%mem":     "mem",
	"pmem":     "mem",
	"mem":      "mem",
	"rss":      "rss",
	"vsz":      "vsz",
	"vsize":    "vsz",
	"stat":     "state",
	"state":    "state",
	"tty":      "tty",
	"time":     "time",
	"etime":    "elapsed",
	"elapsed":  "elapsed",
	"nice":     "nice",
	"ni":       "nice",
	"pri":      "priority",
	"priority": "priority",
	"start":    "start_time",
	"stime":    "start_time",
	"lstart":   "start_time",
}

func translateFlags(args []string) []string {
	var procsArgs []string
	var searchTerms []string
	skipNext := false
	hasPagerFlag := false

	for i, arg := range args {
		if skipNext {
			skipNext = false
			continue
		}

		// Handle GNU long options
		if strings.HasPrefix(arg, "--") {
			if idx := strings.Index(arg, "="); idx != -1 {
				opt := arg[:idx]
				val := arg[idx+1:]

				switch opt {
				case "--sort":
					procsArgs = append(procsArgs, translateSort(val)...)
				case "--user", "--User":
					searchTerms = append(searchTerms, val)
				case "--pid":
					searchTerms = append(searchTerms, val)
				case "--pager":
					hasPagerFlag = true
					procsArgs = append(procsArgs, arg)
				}
				continue
			}

			switch arg {
			case "--forest":
				procsArgs = append(procsArgs, "--tree")
			case "--headers", "--no-headers":
				// Ignore
			case "--pager":
				hasPagerFlag = true
				procsArgs = append(procsArgs, arg)
				if i+1 < len(args) && !strings.HasPrefix(args[i+1], "-") {
					procsArgs = append(procsArgs, args[i+1])
					skipNext = true
				}
			default:
				if !ignoredFlags[arg] {
					procsArgs = append(procsArgs, arg)
				}
			}
			continue
		}

		// Handle UNIX-style options (with dash)
		if strings.HasPrefix(arg, "-") && len(arg) > 1 {
			flags := arg[1:]

			// Check for flags that take values
			for j, c := range flags {
				switch c {
				case 'u', 'U': // user
					remaining := flags[j+1:]
					var val string
					if len(remaining) > 0 {
						val = string(remaining)
					} else if i+1 < len(args) && !strings.HasPrefix(args[i+1], "-") {
						val = args[i+1]
						skipNext = true
					}
					if val != "" {
						searchTerms = append(searchTerms, val)
					}
					goto nextArg
				case 'p': // pid
					remaining := flags[j+1:]
					var val string
					if len(remaining) > 0 {
						val = string(remaining)
					} else if i+1 < len(args) && !strings.HasPrefix(args[i+1], "-") {
						val = args[i+1]
						skipNext = true
					}
					if val != "" {
						searchTerms = append(searchTerms, val)
					}
					goto nextArg
				case 'C': // command name
					remaining := flags[j+1:]
					var val string
					if len(remaining) > 0 {
						val = string(remaining)
					} else if i+1 < len(args) && !strings.HasPrefix(args[i+1], "-") {
						val = args[i+1]
						skipNext = true
					}
					if val != "" {
						searchTerms = append(searchTerms, val)
					}
					goto nextArg
				case 'o', 'O': // output format - skip value
					if j+1 >= len(flags) && i+1 < len(args) && !strings.HasPrefix(args[i+1], "-") {
						skipNext = true
					}
					goto nextArg
				case 'G', 'g': // group - skip value
					if j+1 >= len(flags) && i+1 < len(args) && !strings.HasPrefix(args[i+1], "-") {
						skipNext = true
					}
					goto nextArg
				case 't': // tty - skip value
					if j+1 >= len(flags) && i+1 < len(args) && !strings.HasPrefix(args[i+1], "-") {
						skipNext = true
					}
					goto nextArg
				case 'H': // tree view
					procsArgs = append(procsArgs, "--tree")
				case 'e', 'A', 'a', 'x', 'f', 'l', 'j', 'v', 'w', 'r', 'd', 'N', 'T', 's', 'c', 'm', 'L':
					// Ignored flags
				default:
					// Unknown flag, pass through
					procsArgs = append(procsArgs, "-"+string(c))
				}
			}
		nextArg:
			continue
		}

		// Handle BSD-style options (no dash) - like "aux", "ef"
		if len(arg) > 0 && !strings.HasPrefix(arg, "-") {
			// Check if it looks like BSD ps options (common patterns)
			if isBSDStyleOptions(arg) {
				for _, c := range arg {
					switch c {
					case 'f': // forest/tree (BSD)
						procsArgs = append(procsArgs, "--tree")
						// Most BSD flags can be ignored as procs shows all with good defaults
						// a, u, x, e, etc. are about process selection which procs handles
					}
				}
				continue
			}

			// Otherwise treat as a search term (could be PID or pattern)
			searchTerms = append(searchTerms, arg)
		}
	}

	// Add default --pager disable if user hasn't specified it
	if !hasPagerFlag {
		procsArgs = append([]string{"--pager", "disable"}, procsArgs...)
	}

	// Build result
	result := make([]string, 0, len(procsArgs)+len(searchTerms))
	result = append(result, procsArgs...)
	result = append(result, searchTerms...)
	return result
}

// isBSDStyleOptions checks if a string looks like BSD ps options
func isBSDStyleOptions(s string) bool {
	// BSD ps options are typically short combinations like "aux", "ef", "axjf"
	if len(s) > 5 || len(s) < 1 {
		return false
	}
	// Must contain at least one of the common BSD option chars
	hasCommonBSD := false
	commonBSD := "auxef"
	for _, c := range s {
		if strings.ContainsRune(commonBSD, c) {
			hasCommonBSD = true
			break
		}
	}
	if !hasCommonBSD {
		return false
	}
	// All chars must be valid BSD option chars
	validBSDChars := "adefghjlmnoprstuvwxAJLMNOPRSTUVWX"
	for _, c := range s {
		if !strings.ContainsRune(validBSDChars, c) {
			return false
		}
	}
	return true
}

// translateSort converts ps sort column to procs sort
func translateSort(col string) []string {
	// Remove leading +/- for direction
	desc := false
	if strings.HasPrefix(col, "-") {
		desc = true
		col = col[1:]
	} else if strings.HasPrefix(col, "+") {
		col = col[1:]
	}

	// Map column name
	procsCol := col
	if mapped, ok := columnMap[strings.ToLower(col)]; ok {
		procsCol = mapped
	}

	if desc {
		return []string{"--sortd", procsCol}
	}
	return []string{"--sorta", procsCol}
}
