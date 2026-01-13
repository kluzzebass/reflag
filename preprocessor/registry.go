package preprocessor

import (
	"fmt"
	"io"
	"sort"
	"sync"
	"text/tabwriter"
)

var (
	registry = make(map[string]Preprocessor)
	mu       sync.RWMutex
)

// Register adds a preprocessor to the global registry
func Register(p Preprocessor) {
	mu.Lock()
	defer mu.Unlock()
	registry[p.ToolName()] = p
}

// Get returns a preprocessor for the given tool name
// Returns nil if no preprocessor is registered
func Get(toolName string) Preprocessor {
	mu.RLock()
	defer mu.RUnlock()
	return registry[toolName]
}

// List returns all registered preprocessor tool names
func List() []string {
	mu.RLock()
	defer mu.RUnlock()
	names := make([]string, 0, len(registry))
	for name := range registry {
		names = append(names, name)
	}
	return names
}

// PrintTable writes a formatted table of all preprocessors to the given writer
func PrintTable(w io.Writer) {
	mu.RLock()
	names := make([]string, 0, len(registry))
	for name := range registry {
		names = append(names, name)
	}
	mu.RUnlock()

	sort.Strings(names)

	tw := tabwriter.NewWriter(w, 0, 0, 2, ' ', 0)
	fmt.Fprintln(tw, "TOOL\tFEATURE")
	for _, name := range names {
		p := Get(name)
		if p != nil {
			fmt.Fprintf(tw, "%s\t%s\n", name, p.Description())
		}
	}
	tw.Flush()
}
