# reflag

![AI Slop Badge](ai-slop-badge.svg)

A tool that translates command-line flags between different CLI tools. Currently supports:

- `cat` → [bat](https://github.com/sharkdp/bat)
- `ls` → [eza](https://github.com/eza-community/eza)
- `grep` → [ripgrep](https://github.com/BurntSushi/ripgrep)
- `find` → [fd](https://github.com/sharkdp/fd)
- `df` → [duf](https://github.com/muesli/duf)
- `du` → [dust](https://github.com/bootandy/dust)
- `ps` → [procs](https://github.com/dalance/procs)
- `dig` → [doggo](https://github.com/mr-karan/doggo)
- `less` → [moor](https://github.com/walles/moor)

## Why?

You can't teach an old dog new tricks. After decades of muscle memory, your fingers type `ls -ltr` before your brain even registers the thought. Modern replacements like eza, fd, ripgrep, dust, and procs are genuinely better tools with nicer output, better defaults, and useful features—but they come with their own flag conventions that differ just enough to trip you up.

reflag bridges this gap. Instead of retraining years of muscle memory or giving up on better tools, you can keep typing the commands you know while getting the output you want. It's not about refusing to learn; it's about acknowledging that some habits are so deeply ingrained they're practically reflexes.

## Quick Start

### 1. Install reflag

Choose one of the following:

**Using Homebrew** (macOS):
```bash
brew tap kluzzebass/tap
brew install reflag
```

**From releases**:
```bash
# Download from https://github.com/kluzzebass/reflag/releases
# Make executable and move to your PATH
chmod +x reflag
sudo mv reflag /usr/local/bin/
```

**Using Go**:
```bash
go install github.com/kluzzebass/reflag@latest
```

**From source**:
```bash
git clone https://github.com/kluzzebass/reflag.git
cd reflag
make build
sudo mv reflag /usr/local/bin/
```

### 2. Install the modern tools

Install the tools you want to use. For example:

```bash
# macOS
brew install bat eza fd ripgrep dust procs doggo moor duf

# Linux (Ubuntu/Debian)
sudo apt install bat eza fd-find ripgrep dust procs

# Linux (Fedora)
sudo dnf install bat eza fd-find ripgrep dust procs

# Arch Linux
sudo pacman -S bat eza fd ripgrep dust procs
```

**Note:** `moor` may need to be installed separately on some platforms. See the [moor installation guide](https://github.com/walles/moor#installing).

### 3. Set up shell integration

Run this to enable automatic flag translation:

**bash** (`~/.bashrc`):
```bash
echo 'eval "$(reflag --init bash)"' >> ~/.bashrc && source ~/.bashrc
```

**zsh** (`~/.zshrc`):
```bash
echo 'eval "$(reflag --init bash)"' >> ~/.zshrc && source ~/.zshrc
```

**fish** (`~/.config/fish/config.fish`):
```fish
echo 'reflag --init fish | source' >> ~/.config/fish/config.fish && source ~/.config/fish/config.fish
```

### 4. Start using your familiar commands

```bash
ls -ltr             # Uses eza under the hood
grep -rni TODO      # Uses ripgrep
find . -name '*.go' # Uses fd
df -h               # Uses duf
du -h               # Uses dust
ps aux              # Uses procs
dig example.com MX  # Uses doggo
less -S file.txt    # Uses moor
```

That's it! Your muscle memory still works, but you get modern tool output.

**Note:** To bypass reflag and use the original command, use `command ls` or `/bin/ls`.

## Limitations

Flag translation is inherently imperfect. Here's what you should know:

**Not all flags have equivalents.** Some source tool flags simply don't map to anything in the target tool. reflag passes unrecognized flags through unchanged, which may cause errors or unexpected behavior in the target tool.

**Semantic differences exist.** Even when flags appear similar, subtle behavioral differences may exist. The translation aims for "close enough" rather than pixel-perfect compatibility.

**Flag interactions are complex.** Some flag combinations in the source tool may not translate cleanly when the target tool handles those interactions differently.

**Output format will differ.** While the information displayed should be similar, the exact formatting, colors, and layout will match the target tool, not the source. That's usually the point—eza's output is nicer than ls—but don't expect identical output.

**This is a convenience, not a compatibility layer.** reflag is meant to ease the transition to better tools, not to provide a perfect emulation of the source tool's behavior.

**Shell scripts are generally unaffected.** Scripts using `#!/bin/bash` run in non-interactive mode and don't source `~/.bashrc` or `~/.zshrc`, so they won't see the reflag functions. If you do encounter issues, you can bypass the functions with `command ls` or `/bin/ls`, or exclude specific translators from your init: `eval "$(reflag --init bash ls2eza grep2rg)"`.

## Installation

### Using Homebrew (macOS)

```bash
brew tap kluzzebass/tap
brew install reflag
```

### From releases

Download the latest binary from the [releases page](https://github.com/kluzzebass/reflag/releases).

### Using Go

```bash
go install github.com/kluzzebass/reflag@latest
```

### From source

```bash
git clone https://github.com/kluzzebass/reflag.git
cd reflag
make build
```

## Usage

### Explicit Mode

Specify the source and target tools explicitly:

```bash
$ reflag ls eza -la
eza -l -a

$ reflag ls eza -ltr
eza -l --sort=modified

$ reflag ls eza -lSh /tmp
eza -l --sort=size --reverse /tmp
```

### Shell Integration

Generate shell functions that wrap the source commands:

```bash
# Preview what would be generated
reflag --init bash

# Generate only specific translators (useful if conflicts exist)
reflag --init bash ls2eza grep2rg
```

**Recommended setup:** Add this to your shell config to automatically pick up new translators:

```bash
# ~/.bashrc or ~/.zshrc
eval "$(reflag --init bash)"

# Or for specific translators only:
eval "$(reflag --init bash ls2eza grep2rg)"
```

```fish
# ~/.config/fish/config.fish
reflag --init fish | source

# Or for specific translators only:
reflag --init fish ls2eza grep2rg | source
```

**Alternative:** Append the output once (won't auto-update with new translators):

```bash
reflag --init bash >> ~/.bashrc
reflag --init fish >> ~/.config/fish/config.fish
```

Or create functions manually:

```bash
# bash/zsh
ls() { eval "$(reflag ls eza "$@")"; }
```

```fish
# fish
function ls
    eval (reflag ls eza $argv)
end
```

### List Available Translators

```bash
$ reflag --list
df2duf: df -> duf
dig2doggo: dig -> doggo
du2dust: du -> dust
find2fd: find -> fd
grep2rg: grep -> rg
ls2eza: ls -> eza
ps2procs: ps -> procs
```

## ls2eza Translator

The ls2eza translator converts `ls` flags to `eza` equivalents.

### BSD vs GNU ls Compatibility

reflag supports both BSD ls (macOS, FreeBSD) and GNU ls (Linux) flag conventions. By default, it auto-detects based on your operating system:

- **macOS, FreeBSD, OpenBSD, NetBSD, DragonFly** -> BSD mode
- **Linux, Windows, others** -> GNU mode

Override with the `--mode` flag:

```bash
reflag --mode=bsd ls eza -T   # Force BSD mode
reflag --mode=gnu ls eza -T   # Force GNU mode
```

### Supported Flags

#### Display Format

| ls flag | eza equivalent | Description |
|---------|----------------|-------------|
| `-l` | `-l` | Long format |
| `-1` | `-1` | One entry per line |
| `-C` | `--grid` | Multi-column output |
| `-x` | `--across` | Sort grid across |
| `-m` | `--oneline` | Stream output |

#### Show/Hide Entries

| ls flag | eza equivalent | Description |
|---------|----------------|-------------|
| `-a` | `-a` | Show all (including . and ..) |
| `-A` | `-A` | Show hidden (except . and ..) |
| `-d` | `-d` | List directories themselves |
| `-R` | `--recurse` | Recurse into directories |

#### Sorting

| ls flag | eza equivalent | Description |
|---------|----------------|-------------|
| `-t` | `--sort=modified --reverse` | Sort by modification time |
| `-S` | `--sort=size --reverse` | Sort by size |
| `-c` | `--sort=changed --reverse` | Sort by change time |
| `-u` | `--sort=accessed --reverse` | Sort by access time |
| `-U` | `--sort=created --reverse` | Sort by creation time (BSD) |
| `-f` | `--sort=none -a` | Unsorted, show all |
| `-v` | `--sort=name` | Version/name sort |
| `-r` | `--reverse` | Reverse sort order |

#### Long Format Options

| ls flag | eza equivalent | Description |
|---------|----------------|-------------|
| `-i` | `--inode` | Show inode numbers |
| `-s` | `--blocksize` | Show allocated blocks |
| `-n` | `--numeric` | Numeric user/group IDs |
| `-o` | `-l --no-group` | Long format without group |
| `-g` | `-l --no-user` | Long format without owner |
| `-O` | `--flags` | Show file flags (BSD/macOS) |
| `-@` | `--extended` | Show extended attributes |
| `-h` | (default) | Human-readable sizes |

### Sort Order Handling

ls and eza have opposite default sort orders for time and size sorting. reflag automatically adds `--reverse` when needed so that the output matches ls behavior:

- `ls -lt` shows newest first -> `eza --sort=modified --reverse`
- `ls -ltr` shows oldest first -> `eza --sort=modified` (no reverse needed)
- `ls -lS` shows largest first -> `eza --sort=size --reverse`

### Conflicting Flags (BSD vs GNU)

| Flag | BSD ls | GNU ls |
|------|--------|--------|
| `-T` | Full timestamp display | Tab size (ignored) |
| `-X` | Don't cross filesystems (ignored) | Sort by extension |
| `-I` | Prevent auto -A (ignored) | Ignore pattern (`-I PATTERN`) |
| `-w` | Raw non-printable chars (ignored) | Output width (`-w COLS`) |
| `-D` | Date format (`-D FORMAT`) | Dired mode (ignored) |

## cat2bat Translator

The cat2bat translator converts `cat` commands to `bat` with flags that make bat behave like cat.

### Key Features

- **Plain output**: Always adds `-p` (plain style) to disable bat's decorations like line numbers, grid borders, and file headers
- **No paging**: Adds `--paging=never` to disable bat's automatic paging behavior
- **Smart colorization**: Uses `--color=auto` to enable syntax highlighting when output goes to a terminal, but disables it when piped or redirected

### Supported Flags

| cat flag | bat equivalent | Description |
|----------|----------------|-------------|
| `-n` | `-n` | Number all output lines |
| `--number` | `-n` | Number all output lines |
| `-s` | `-s` | Squeeze multiple blank lines into one |
| `--squeeze-blank` | `-s` | Squeeze multiple blank lines |
| `-A` | `-A` | Show non-printable characters |
| `--show-all` | `-A` | Show non-printable characters |
| `-u` | `-u` | Unbuffered output |
| `--unbuffered` | `-u` | Unbuffered output |

### Ignored Flags

All bat-specific flags are ignored since the goal is cat-like behavior:
- Syntax highlighting options (`-l`, `--language`, `--theme`)
- Line highlighting (`-H`, `--highlight-line`)
- Decorations (`--style`, `--decorations`)
- Git integration (`-d`, `--diff`)
- File metadata (`--file-name`)
- And other bat-specific features

### Examples

```bash
$ reflag cat bat README.md
bat -p --paging=never --color=auto README.md

$ reflag cat bat -n file.txt
bat -p --paging=never --color=auto -n file.txt

$ reflag cat bat -ns file1.txt file2.txt
bat -p --paging=never --color=auto -n -s file1.txt file2.txt

# Color output when viewing in terminal
$ reflag cat bat file.rs
# Shows syntax-highlighted Rust code

# No color when piping to other commands
$ reflag cat bat file.rs | grep "fn"
# Plain output suitable for grep
```

### Why Use This?

While bat is designed to be a cat replacement with syntax highlighting, sometimes you want:
- The simplicity of cat's output
- Syntax highlighting only when appropriate (terminal vs pipe)
- To maintain muscle memory for cat flags
- A gradual transition from cat to bat

## grep2rg Translator

The grep2rg translator converts `grep` flags to `rg` (ripgrep) equivalents.

### Key Differences

- **Recursion**: grep needs `-r` for recursive search; rg is recursive by default
- **Regex**: grep defaults to basic regex (`-G`); rg uses extended regex by default
- **Binary files**: grep searches binary files by default; rg skips them

### Supported Flags

#### Passthrough Flags (identical in both)

`-i`, `-v`, `-w`, `-x`, `-c`, `-l`, `-L`, `-n`, `-H`, `-h`, `-o`, `-q`, `-s`, `-F`, `-P`, `-a`, `-A`, `-B`, `-C`, `-m`, `-e`, `-f`

#### Translated Flags

| grep | rg | Notes |
|------|-----|-------|
| `-r`, `-R` | (default) | rg is recursive by default |
| `-E` | (default) | rg uses extended regex by default |
| `--include=GLOB` | `-g GLOB` | |
| `--exclude=GLOB` | `-g '!GLOB'` | |
| `--exclude-dir=DIR` | `-g '!DIR/'` | |
| `-Z`, `--null` | `-0` | Null-separated output |

### Examples

```bash
$ reflag grep rg -rni "TODO" .
rg -n -i TODO .

$ reflag grep rg --include='*.go' "func" src/
rg -g *.go func src/

$ reflag grep rg -A3 -B3 "error" file.txt
rg -A 3 -B 3 error file.txt
```

## find2fd Translator

The find2fd translator converts `find` expressions to `fd` syntax.

### Syntax Transformation

find and fd have fundamentally different command structures:
- **find**: `find [paths...] [expressions...]`
- **fd**: `fd [options...] [pattern] [paths...]`

The translator extracts the pattern from `-name`/`-iname` expressions and reorders arguments to match fd's expected format.

### Supported Expressions

| find | fd | Notes |
|------|-----|-------|
| `-name PATTERN` | `PATTERN` | Glob converted to regex |
| `-iname PATTERN` | `-i PATTERN` | Case insensitive |
| `-type f/d/l` | `-t f/d/l` | File type |
| `-maxdepth N` | `-d N` | |
| `-mindepth N` | `--min-depth N` | |
| `-path PATTERN` | `-p PATTERN` | Path pattern |
| `-regex PATTERN` | `PATTERN` | fd uses regex by default |
| `-size +/-N` | `-S +/-N` | |
| `-empty` | `-t e` | |
| `-executable` | `-t x` | |
| `-newer FILE` | `--newer FILE` | |
| `-mtime -N/+N` | `--changed-within/before Nd` | |
| `-print0` | `-0` | Null-separated output |
| `-L`, `-follow` | `-L` | Follow symlinks |
| `-xdev` | `--one-file-system` | |
| `-user USER` | `--owner USER` | |
| `-group GROUP` | `--owner :GROUP` | |

### Unsupported

- `-exec`, `-execdir`, `-ok` (too complex to translate safely)
- Logical operators `-o`, `-or`, `!`, `-not`, parentheses (fd handles these differently)
- `-print` (default behavior, ignored)

### Examples

```bash
$ reflag find fd . -name '*.go' -type f
fd -t f '\.go$'

$ reflag find fd /tmp -maxdepth 2 -name '*.txt'
fd -d 2 '\.txt$' /tmp

$ reflag find fd . -type d -name 'test*'
fd -t d 'test[^/]*'

$ reflag find fd . -mtime -7 -name '*.log'
fd --changed-within 7d '\.log$'
```

## df2duf Translator

The df2duf translator converts `df` (disk free) flags to `duf` (disk usage/free) equivalents.

### Key Differences

- **Output format**: df shows simple tabular output; duf provides colorized, modern output with usage bars
- **Human-readable**: duf uses human-readable sizes by default (df requires `-h`)
- **Filesystem focus**: Both tools show filesystem information rather than directory trees

### Supported Flags

#### Translated Flags

| df flag | duf equivalent | Description |
|---------|----------------|-------------|
| `-a`, `--all` | `-all` | Include pseudo/duplicate/inaccessible filesystems |
| `--inodes` | `-inodes` | Show inode information instead of block usage |
| `--exclude=PATTERN` | `-hide-mp PATTERN` | Exclude filesystems matching pattern |
| `-I PATTERN` | `-hide-mp PATTERN` | BSD: Exclude filesystems matching pattern |

#### Ignored Flags (duf defaults or not applicable)

| df flag | Reason |
|---------|--------|
| `-h`, `--human-readable` | duf is human-readable by default |
| `--si` | duf uses SI units by default |
| `-k`, `-m`, `-g` | duf shows optimal unit automatically |
| `-B`, `--block-size` | duf uses SI by default |
| `-l` | Count hard links (not applicable to filesystem stats) |
| `-x`, `--one-file-system` | duf shows all filesystems by default |
| `-t`, `--type` | Use duf's filtering options instead |
| `-T`, `--print-type` | duf shows filesystem type by default |
| `-i`, `-P`, `-H`, `-L` | Not applicable to filesystem display |
| `--total` | duf always shows summary |
| `--sync` | No equivalent |

### Examples

```bash
$ reflag df duf -h
duf

$ reflag df duf -a
duf -all

$ reflag df duf --inodes
duf -inodes

$ reflag df duf -h /
duf

$ reflag df duf --exclude=tmpfs
duf -hide-mp tmpfs
```

### Notes

- **Path arguments are ignored**: Unlike `df`, which can take filesystem or mount point arguments, `duf` displays all mounted filesystems by default. Path arguments in the `df` command are collected but not passed to `duf`.
- **Not included in shell init by default**: The translator sets `IncludeInInit() = false` because the behavioral differences between `df` and `duf` are significant enough that automatic substitution might cause confusion. To enable it, explicitly add it: `reflag --init bash +df2duf`

## dig2doggo Translator

The dig2doggo translator converts `dig` DNS query flags to `doggo` equivalents.

### Key Features

- **Query specification**: Translates domain names, query types (A, AAAA, MX, etc.), and classes (IN, CH, HS)
- **Nameserver selection**: Supports `@server` syntax with protocol prefixes (`@tcp://`, `@https://`, `@tls://`, `@quic://`)
- **Query flags**: Converts DNS flags like `+dnssec`, `+short`, `+recurse`, `+tcp`
- **EDNS options**: Supports `+nsid`, `+cookie`, `+padding`, `+ede`, `+subnet`

### Supported Flags

#### Short Flags

| dig | doggo | Description |
|-----|-------|-------------|
| `-4` | `-4` | IPv4 only |
| `-6` | `-6` | IPv6 only |
| `-q NAME` | `-q NAME` | Query name |
| `-t TYPE` | `-t TYPE` | Query type |
| `-c CLASS` | `-c CLASS` | Query class |
| `-x ADDR` | `-x` | Reverse lookup |
| `-m` | `--debug` | Debug mode |

#### Plus Options

| dig | doggo | Description |
|-----|-------|-------------|
| `+short` | `--short` | Short output |
| `+tcp` | `@tcp://` | Use TCP |
| `+dnssec` | `--do` | Request DNSSEC records |
| `+recurse` | `--rd` | Recursion desired |
| `+aa` | `--aa` | Authoritative answer |
| `+ad` | `--ad` | Authentic data |
| `+cd` | `--cd` | Checking disabled |
| `+nsid` | `--nsid` | Name server ID |
| `+cookie` | `--cookie` | DNS cookie |
| `+padding` | `--padding` | EDNS padding |
| `+ede` | `--ede` | Extended DNS errors |
| `+search` | `--search` | Use search list |
| `+timeout=N` | `--timeout Ns` | Query timeout |
| `+ndots=N` | `--ndots N` | Dots in name |
| `+subnet=ADDR` | `--ecs ADDR` | EDNS client subnet |

### Examples

```bash
$ reflag dig doggo example.com
doggo -q example.com

$ reflag dig doggo example.com MX
doggo -q example.com -t MX

$ reflag dig doggo @8.8.8.8 example.com A
doggo -q example.com -t A -n 8.8.8.8

$ reflag dig doggo +short example.com
doggo --short -q example.com

$ reflag dig doggo @1.1.1.1 +dnssec +short example.com
doggo --do --short -q example.com -n 1.1.1.1

$ reflag dig doggo @https://cloudflare-dns.com/dns-query example.com
doggo -q example.com -n @https://cloudflare-dns.com/dns-query
```

### Unsupported Features

- Batch mode (`-f FILE`)
- TSIG authentication (`-k`, `-y`)
- Zone transfers (AXFR/IXFR) - doggo has limited support
- Output formatting options (`+stats`, `+cmd`, `+comments`, etc.) - doggo has different output format
- Trace mode (`+trace`) - not available in doggo

## Adding New Translators

reflag is designed to be extensible. To add a new translator:

1. Create a new package under `translator/`
2. Implement the `translator.Translator` interface
3. Register it in `init()` using `translator.Register()`

See `translator/ls2eza/` for an example implementation.

## License

MIT
