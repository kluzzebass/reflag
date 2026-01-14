package main

import (
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/kluzzebass/reflag/translator"
	_ "github.com/kluzzebass/reflag/translator/dig2doggo" // Register dig2doggo translator
	_ "github.com/kluzzebass/reflag/translator/du2dust"   // Register du2dust translator
	_ "github.com/kluzzebass/reflag/translator/find2fd"   // Register find2fd translator
	_ "github.com/kluzzebass/reflag/translator/grep2rg"   // Register grep2rg translator
	_ "github.com/kluzzebass/reflag/translator/less2moor" // Register less2moor translator
	_ "github.com/kluzzebass/reflag/translator/ls2eza"    // Register ls2eza translator
	_ "github.com/kluzzebass/reflag/translator/ps2procs"  // Register ps2procs translator
)

// Version information - set via ldflags at build time
var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

func shellQuote(s string) string {
	if strings.ContainsAny(s, " \t\n\"'\\$`!") {
		return "'" + strings.ReplaceAll(s, "'", "'\"'\"'") + "'"
	}
	return s
}

func printVersion(name string) {
	fmt.Printf("%s %s\n", name, version)
	if commit != "none" {
		fmt.Printf("  commit: %s\n", commit)
	}
	if date != "unknown" {
		fmt.Printf("  built:  %s\n", date)
	}
	fmt.Println("\nCopyright (c) 2026 Jan Fredrik Leversund")
	fmt.Println("Licensed under the MIT License")
}

func printLicense() {
	licenseText := `MIT License

Copyright (c) 2026 Jan Fredrik Leversund

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.`
	fmt.Println(licenseText)
}

// parseInitArgs parses --init arguments, returning shell type and translator names
// Shell can appear anywhere in args; defaults to "bash" if not specified
func parseInitArgs(args []string) (shell string, filterNames []string) {
	shell = "bash"
	for _, arg := range args {
		switch arg {
		case "bash", "zsh", "fish":
			shell = arg
		default:
			filterNames = append(filterNames, arg)
		}
	}
	return
}

func printInit(shell string, filterNames []string) {
	var names []string
	if len(filterNames) > 0 {
		// Use specified translators
		for _, name := range filterNames {
			if translator.GetByName(name) != nil {
				names = append(names, name)
			} else {
				fmt.Fprintf(os.Stderr, "warning: unknown translator %q\n", name)
			}
		}
	} else {
		// Use all translators that should be included in init by default
		for _, name := range translator.List() {
			t := translator.GetByName(name)
			if t != nil && t.IncludeInInit() {
				names = append(names, name)
			}
		}
	}
	sort.Strings(names)

	switch shell {
	case "fish":
		fmt.Println("# reflag shell init - add to your ~/.config/fish/config.fish")
		fmt.Println()
		for _, name := range names {
			t := translator.GetByName(name)
			fmt.Printf("functions -e %s 2>/dev/null\n", t.SourceTool())
			fmt.Printf("function %s\n", t.SourceTool())
			fmt.Printf("    eval (reflag %s %s $argv)\n", t.SourceTool(), t.TargetTool())
			fmt.Println("end")
			fmt.Println()
		}
	default: // bash, zsh
		fmt.Println("# reflag shell init - add to your ~/.bashrc or ~/.zshrc")
		fmt.Println()
		for _, name := range names {
			t := translator.GetByName(name)
			fmt.Printf("unalias %s 2>/dev/null\n", t.SourceTool())
			fmt.Printf("%s() {\n", t.SourceTool())
			fmt.Printf("    eval \"$(reflag %s %s \"$@\")\"\n", t.SourceTool(), t.TargetTool())
			fmt.Println("}")
			fmt.Println()
		}
	}
}

func printUsage() {
	fmt.Println("reflag - translate command-line flags between tools")
	fmt.Println()
	fmt.Println("Quick setup:")
	fmt.Println("  echo 'eval \"$(reflag --init)\"' >> ~/.bashrc")
	fmt.Println("  echo 'eval \"$(reflag --init)\"' >> ~/.zshrc")
	fmt.Println("  echo 'reflag --init fish | source' >> ~/.config/fish/config.fish")
	fmt.Println()
	fmt.Println("Usage:")
	fmt.Println("  reflag [--mode=MODE] <source> <target> [flags...]")
	fmt.Println("  reflag --list")
	fmt.Println("  reflag --init [bash|zsh|fish] [translator...]")
	fmt.Println("  reflag --version")
	fmt.Println("  reflag --license")
	fmt.Println()
	fmt.Println("Options:")
	fmt.Println("  --mode=MODE    Set dialect mode (e.g., bsd or gnu for ls2eza)")
	fmt.Println("                 Auto-detects from OS if not specified")
	fmt.Println()
	fmt.Println("Available translators:")
	translator.PrintTable(os.Stdout)
}

func runTranslator(t translator.Translator, args []string, mode string) {
	// Handle version flag
	for _, arg := range args {
		if arg == "-V" || arg == "--version" {
			printVersion(t.Name())
			return
		}
	}

	translatedArgs := t.Translate(args, mode)

	// Build and print the command
	parts := make([]string, len(translatedArgs)+1)
	parts[0] = t.TargetTool()
	for i, arg := range translatedArgs {
		parts[i+1] = shellQuote(arg)
	}
	fmt.Println(strings.Join(parts, " "))
}

func main() {
	args := os.Args[1:]

	// Handle reflag's own flags
	if len(args) == 0 {
		printUsage()
		return
	}

	switch args[0] {
	case "--version", "-V":
		printVersion("reflag")
		return
	case "--license":
		printLicense()
		return
	case "--list", "-l":
		translator.PrintTable(os.Stdout)
		return
	case "--help", "-h":
		printUsage()
		return
	case "--init":
		shell, filterNames := parseInitArgs(args[1:])
		printInit(shell, filterNames)
		return
	}

	// Parse --mode flag if present
	mode := ""
	if strings.HasPrefix(args[0], "--mode=") {
		mode = strings.TrimPrefix(args[0], "--mode=")
		args = args[1:]
	} else if args[0] == "--mode" && len(args) > 1 {
		mode = args[1]
		args = args[2:]
	}

	// Explicit mode: reflag [--mode=MODE] <source> <target> [flags...]
	if len(args) < 2 {
		fmt.Fprintln(os.Stderr, "error: expected <source> <target> arguments")
		fmt.Fprintln(os.Stderr, "usage: reflag [--mode=MODE] <source> <target> [flags...]")
		os.Exit(1)
	}

	source, target := args[0], args[1]
	t := translator.Get(source, target)
	if t == nil {
		fmt.Fprintf(os.Stderr, "error: no translator registered for %s to %s\n", source, target)
		fmt.Fprintln(os.Stderr, "use 'reflag --list' to see available translators")
		os.Exit(1)
	}

	runTranslator(t, args[2:], mode)
}
