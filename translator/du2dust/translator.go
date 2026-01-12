package du2dust

import (
	"strings"

	"github.com/kluzzebass/reflag/translator"
)

func init() {
	translator.Register(&Translator{})
}

// Translator implements the du to dust flag translation
type Translator struct{}

func (t *Translator) Name() string       { return "du2dust" }
func (t *Translator) SourceTool() string { return "du" }
func (t *Translator) TargetTool() string { return "dust" }

// Translate converts du arguments to dust arguments
func (t *Translator) Translate(args []string, mode string) []string {
	return translateFlags(args)
}

// Flags to ignore (dust handles automatically or no equivalent)
var ignoredFlags = map[string]bool{
	"-h":                true, // dust is human-readable by default
	"--human-readable":  true,
	"-c":                true, // dust shows total by default
	"--total":           true,
	"-P":                true, // don't follow symlinks (dust default)
	"--no-dereference":  true,
	"-l":                true, // count links
	"--count-links":     true,
	"-S":                true, // separate dirs
	"--separate-dirs":   true,
	"--time":            true, // no equivalent
	"--time-style":      true,
	"-0":                true, // null terminated
	"--null":            true,
	"-H":                true, // dereference args
	"--dereference-args": true,
	"-D":                true, // BSD/GNU dereference args
}

func translateFlags(args []string) []string {
	var dustArgs []string
	var paths []string
	skipNext := false

	for i, arg := range args {
		if skipNext {
			skipNext = false
			continue
		}

		// Handle --option=value format
		if strings.HasPrefix(arg, "--") {
			if idx := strings.Index(arg, "="); idx != -1 {
				opt := arg[:idx]
				val := arg[idx+1:]

				switch opt {
				case "--max-depth":
					dustArgs = append(dustArgs, "-d", val)
				case "--exclude":
					dustArgs = append(dustArgs, "-v", val)
				case "--threshold":
					dustArgs = append(dustArgs, "-z", val)
				case "--block-size":
					// Try to map common block sizes
					dustArgs = append(dustArgs, mapBlockSize(val)...)
				}
				continue
			}

			// Handle standalone long options
			switch arg {
			case "--summarize":
				dustArgs = append(dustArgs, "-d", "0")
			case "--all":
				dustArgs = append(dustArgs, "-F")
			case "--dereference":
				dustArgs = append(dustArgs, "-L")
			case "--one-file-system":
				dustArgs = append(dustArgs, "-x")
			case "--apparent-size":
				dustArgs = append(dustArgs, "-s")
			case "--si":
				dustArgs = append(dustArgs, "-o", "si")
			case "--bytes":
				dustArgs = append(dustArgs, "-o", "b")
			case "--inodes":
				dustArgs = append(dustArgs, "-f")
			default:
				if ignoredFlags[arg] {
					continue
				}
				// Pass through unknown long options
				dustArgs = append(dustArgs, arg)
			}
			continue
		}

		if strings.HasPrefix(arg, "-") && len(arg) > 1 {
			flags := arg[1:]
			for j, c := range flags {
				switch c {
				case 's': // summarize
					dustArgs = append(dustArgs, "-d", "0")
				case 'a': // all files
					dustArgs = append(dustArgs, "-F")
				case 'd': // max depth
					remaining := flags[j+1:]
					var val string
					if len(remaining) > 0 {
						val = string(remaining)
					} else if i+1 < len(args) && !strings.HasPrefix(args[i+1], "-") {
						val = args[i+1]
						skipNext = true
					}
					if val != "" {
						dustArgs = append(dustArgs, "-d", val)
					}
					goto nextArg
				case 'L': // follow symlinks
					dustArgs = append(dustArgs, "-L")
				case 'x': // one file system
					dustArgs = append(dustArgs, "-x")
				case 'b': // bytes (GNU)
					dustArgs = append(dustArgs, "-o", "b")
				case 'k': // kilobytes
					dustArgs = append(dustArgs, "-o", "kb")
				case 'm': // megabytes
					dustArgs = append(dustArgs, "-o", "mb")
				case 'g': // gigabytes (BSD)
					dustArgs = append(dustArgs, "-o", "gb")
				case 't': // threshold
					remaining := flags[j+1:]
					var val string
					if len(remaining) > 0 {
						val = string(remaining)
					} else if i+1 < len(args) && !strings.HasPrefix(args[i+1], "-") {
						val = args[i+1]
						skipNext = true
					}
					if val != "" {
						dustArgs = append(dustArgs, "-z", val)
					}
					goto nextArg
				case 'I': // BSD exclude pattern
					remaining := flags[j+1:]
					var val string
					if len(remaining) > 0 {
						val = string(remaining)
					} else if i+1 < len(args) && !strings.HasPrefix(args[i+1], "-") {
						val = args[i+1]
						skipNext = true
					}
					if val != "" {
						dustArgs = append(dustArgs, "-v", val)
					}
					goto nextArg
				case 'B': // block size
					remaining := flags[j+1:]
					var val string
					if len(remaining) > 0 {
						val = string(remaining)
					} else if i+1 < len(args) && !strings.HasPrefix(args[i+1], "-") {
						val = args[i+1]
						skipNext = true
					}
					if val != "" {
						dustArgs = append(dustArgs, mapBlockSize(val)...)
					}
					goto nextArg
				case 'X': // exclude from file - no equivalent, skip value
					if j+1 >= len(flags) && i+1 < len(args) {
						skipNext = true
					}
					goto nextArg
				case 'h', 'c', 'P', 'l', 'S', 'H', 'D', '0':
					// Ignored flags
				default:
					dustArgs = append(dustArgs, "-"+string(c))
				}
			}
		nextArg:
			continue
		}

		// Non-flag argument (path)
		paths = append(paths, arg)
	}

	// Build result
	result := make([]string, 0, len(dustArgs)+len(paths))
	result = append(result, dustArgs...)
	result = append(result, paths...)
	return result
}

// mapBlockSize converts du block size to dust output format
func mapBlockSize(size string) []string {
	size = strings.ToUpper(size)
	switch size {
	case "1", "1B":
		return []string{"-o", "b"}
	case "K", "KB", "1K", "1024":
		return []string{"-o", "kb"}
	case "M", "MB", "1M":
		return []string{"-o", "mb"}
	case "G", "GB", "1G":
		return []string{"-o", "gb"}
	default:
		// Can't translate arbitrary block sizes
		return nil
	}
}
