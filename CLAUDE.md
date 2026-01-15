# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

reflag is a Go CLI tool that translates command-line flags between traditional UNIX tools and their modern replacements. It supports 9 translators covering common tools like cat, ls, find, grep, du, ps, dig, less, and more.

## Build and Test Commands

```bash
make build               # Build with version info
make test                # Run all tests
make clean               # Remove binary
make build-all           # Cross-compile for all platforms
go run . ls eza -la      # Run directly without building
go test -run TestName    # Run specific test
```

## Releasing

Push a semver tag to trigger a GitHub Actions release:

```bash
git tag v1.0.0
git push origin v1.0.0
```

## Architecture

### Package Structure

```
reflag/
├── main.go                       # CLI entry point, shell integration
├── main_test.go                  # CLI tests
├── translator/
│   ├── translator.go             # Translator interface
│   ├── registry.go               # Global translator registry
│   ├── translator_test.go        # Registry tests
│   ├── cat2bat/                  # cat → bat (plain mode)
│   ├── ls2eza/                   # ls → eza (BSD/GNU modes)
│   ├── find2fd/                  # find → fd (glob→regex conversion)
│   ├── grep2rg/                  # grep → ripgrep
│   ├── du2dust/                  # du → dust
│   ├── ps2procs/                 # ps → procs
│   ├── dig2doggo/                # dig → doggo (optional)
│   ├── less2moor/                # less → moor (optional)
│   └── more2moor/                # more → moor (optional)
```

### Core Components

1. **Translator Interface** (`translator/translator.go`):
   - `Name()` - translator identifier (e.g., "ls2eza")
   - `SourceTool()` - source tool name (e.g., "ls")
   - `TargetTool()` - target tool name (e.g., "eza")
   - `Translate(args, mode)` - converts source args to target args
   - `IncludeInInit()` - returns true if included in `--init` by default

2. **Registry** (`translator/registry.go`):
   - `Register(t)` - register a translator
   - `Get(source, target)` - lookup by source/target
   - `GetByName(name)` - lookup by name
   - `List()` - list all registered translators
   - `MustGet()` - lookup with panic on failure
   - `PrintTable()` - formatted output of all translators

3. **CLI** (`main.go`):
   - Explicit mode: `reflag [--mode=MODE] <source> <target> [flags...]`
   - Built-in flags: `--list`, `--version`, `--license`, `--help`
   - Shell integration: `--init bash|zsh|fish [+translator...] [-translator...]`

### Available Translators

| Translator | Source | Target | IncludeInInit | Notes |
|------------|--------|--------|---------------|-------|
| cat2bat | cat | bat | true | Plain mode with auto colorization |
| ls2eza | ls | eza | true | BSD/GNU mode detection |
| find2fd | find | fd | true | Glob-to-regex conversion |
| grep2rg | grep | ripgrep | true | Pattern handling |
| du2dust | du | dust | true | Unit conversions |
| ps2procs | ps | procs | true | BSD/GNU styles |
| dig2doggo | dig | doggo | false | DNS query translation |
| less2moor | less | moor | false | Pager flags |
| more2moor | more | moor | false | Simple pager flags |

### ls2eza Translator

Located in `translator/ls2eza/`:

1. **Mode detection** - `getLSMode()` determines BSD vs GNU ls compatibility:
   - Auto-detects based on OS (darwin/freebsd → BSD, linux/others → GNU)
   - Override with `LS2EZA_MODE=bsd` or `LS2EZA_MODE=gnu`
   - Or use `--mode=bsd` or `--mode=gnu` on command line

2. **Flag mappings**:
   - `reverseNeeded` - flags that need sort order correction (`t`, `S`, `c`, `u`, `U`)
   - `flagMap` - short flag translations (30+ flags)
   - `longFlagMap` - long option translations
   - `longFlagPrefixes` - long options with =value

3. **Reverse sort handling** - XOR logic to match ls sort order:
   - `ls -lt` needs `--reverse` (ls shows newest first, eza shows oldest first)
   - `ls -ltr` does NOT need `--reverse` (user explicitly wants oldest first)

### BSD vs GNU Conflicts (ls2eza)

These flags have different meanings between BSD and GNU ls:
- `-T`: BSD=full time display, GNU=tab size (ignored)
- `-X`: BSD=don't cross filesystems (ignored), GNU=sort by extension
- `-I`: BSD=prevent auto -A (ignored), GNU=ignore pattern
- `-w`: BSD=raw non-printable (ignored), GNU=output width
- `-D`: BSD=date format, GNU=dired mode (ignored)

### Other Translators

- **cat2bat**: Converts cat commands to bat with plain mode flags (-p, --paging=never, --color=auto) to make bat behave like cat while preserving auto colorization
- **find2fd**: Converts find expressions to fd syntax, including glob-to-regex pattern conversion
- **grep2rg**: Translates grep flags to ripgrep, handles include/exclude patterns
- **du2dust**: Converts du flags including unit/block size mappings
- **ps2procs**: Handles both BSD-style (`ps aux`) and GNU-style (`ps -ef`) syntax
- **dig2doggo**: Translates DNS query flags including +options syntax
- **less2moor**: Converts less pager flags for the moor Rust-based pager
- **more2moor**: Converts more pager flags (simpler than less) for moor

## Adding a New Translator

1. Create package `translator/<name>/`
2. Implement `translator.Translator` interface
   - Set `IncludeInInit()` to `true` for core/commonly-used translators
   - Set `IncludeInInit()` to `false` for optional/experimental translators
3. Call `translator.Register()` in `init()`
4. Import in `main.go` with blank identifier: `_ "github.com/kluzzebass/reflag/translator/<name>"`
5. Add tests in `translator/<name>/translator_test.go`
6. Update `README.md` with new translator information

### Optional Translators

Translators returning `false` from `IncludeInInit()` are excluded from `./reflag --init` by default. This is useful for:
- Experimental or less commonly used translators
- Translators for niche tools
- New translators that need more testing

Use `+` and `-` prefixes to modify the default set:
```bash
reflag --init bash +dig2doggo          # Add dig2doggo to defaults
reflag --init zsh -ls2eza              # Remove ls2eza from defaults
reflag --init fish +dig2doggo -ps2procs  # Add and remove
```
