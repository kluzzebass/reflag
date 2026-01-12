# ls2eza

A simple tool that translates `ls` command flags to their [eza](https://github.com/eza-community/eza) equivalents.

## Installation

```bash
go install github.com/kluzzebass/ls2eza@latest
```

Or build from source:

```bash
git clone https://github.com/kluzzebass/ls2eza.git
cd ls2eza
go build -o ls2eza
```

## Usage

ls2eza takes ls-style arguments and outputs the equivalent eza command:

```bash
$ ls2eza -la
eza -l -a

$ ls2eza -ltr
eza -l --sort=modified

$ ls2eza -lSh /tmp
eza -l --sort=size --reverse /tmp
```

### Using with an alias

You can create a shell alias to automatically translate and execute:

```bash
alias ls='eval $(ls2eza "$@")'
```

Or for fish shell:

```fish
function ls
    eval (ls2eza $argv)
end
```

## Supported Flags

### Display Format

| ls flag | eza equivalent | Description |
|---------|----------------|-------------|
| `-l` | `-l` | Long format |
| `-1` | `-1` | One entry per line |
| `-C` | `--grid` | Multi-column output |
| `-x` | `--across` | Sort grid across |
| `-m` | `--oneline` | Stream output |

### Show/Hide Entries

| ls flag | eza equivalent | Description |
|---------|----------------|-------------|
| `-a` | `-a` | Show all (including . and ..) |
| `-A` | `-A` | Show hidden (except . and ..) |
| `-d` | `-d` | List directories themselves |
| `-R` | `--recurse` | Recurse into directories |

### Sorting

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

### Long Format Options

| ls flag | eza equivalent | Description |
|---------|----------------|-------------|
| `-i` | `--inode` | Show inode numbers |
| `-s` | `--blocksize` | Show allocated blocks |
| `-n` | `--numeric` | Numeric user/group IDs |
| `-o` | `-l --no-group` | Long format without group |
| `-g` | `-l --no-user` | Long format without owner |
| `-O` | `--flags` | Show file flags (BSD/macOS) |
| `-@` | `--extended` | Show extended attributes |
| `-T` | `--time-style=full-iso` | Full timestamp display |
| `-h` | (default) | Human-readable sizes |

### Indicators

| ls flag | eza equivalent | Description |
|---------|----------------|-------------|
| `-F` | `-F` | Append type indicators |
| `-p` | `--classify` | Append / to directories |

### Symlinks

| ls flag | eza equivalent | Description |
|---------|----------------|-------------|
| `-L` | `-X` | Dereference symlinks |
| `-H` | `-X` | Follow symlinks on command line |

### Other

| ls flag | eza equivalent | Description |
|---------|----------------|-------------|
| `-G` | (default) | Color output |
| `-k` | (ignored) | Block size handling |

### Long Options

| ls option | eza equivalent |
|-----------|----------------|
| `--all` | `-a` |
| `--almost-all` | `-A` |
| `--directory` | `-d` |
| `--recursive` | `--recurse` |
| `--human-readable` | (default) |
| `--inode` | `--inode` |
| `--numeric-uid-gid` | `--numeric` |
| `--classify` | `-F` |
| `--dereference` | `-X` |
| `--color=WHEN` | `--color=WHEN` |

### Sort Order Handling

ls and eza have opposite default sort orders for time and size sorting. ls2eza automatically adds `--reverse` when needed so that the output matches ls behavior:

- `ls -lt` shows newest first → `eza --sort=modified --reverse`
- `ls -ltr` shows oldest first → `eza --sort=modified` (no reverse needed)
- `ls -lS` shows largest first → `eza --sort=size --reverse`
- `ls -lc` shows newest changed first → `eza --sort=changed --reverse`

## License

MIT
