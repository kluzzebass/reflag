# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

ls2eza is a Go CLI tool that translates `ls` command-line flags to their `eza` equivalents. It outputs the translated eza command (does not execute it).

## Build Commands

```bash
go build -o ls2eza       # Build the binary
go run main.go -la       # Run directly without building
```

## Architecture

Single-file Go application (`main.go`) with these components:

1. **Flag mappings** - Static maps defining ls→eza translations:
   - `reverseNeeded` - Flags that need sort order correction (`t`, `S`, `c`, `u`, `U`)
   - `flagMap` - Short flag translations (30+ flags)
   - `longFlagMap` - Long option translations (`--all`, `--recursive`, etc.)
2. **`translateFlags()`** - Core logic that parses arguments, applies mappings, handles reverse-sort semantics, and deduplicates flags
3. **`main()`** - Entry point that outputs the shell-quoted eza command

### Key behavior: Reverse sort handling

ls and eza have opposite default sort orders for time (`-t`, `-c`, `-u`, `-U`) and size (`-S`). The tool uses XOR logic to determine when to add `--reverse`:
- `ls -lt` needs `--reverse` (ls shows newest first, eza shows oldest first)
- `ls -ltr` does NOT need `--reverse` (user explicitly wants oldest first)
- Same logic applies to `-c` (changed), `-u` (accessed), `-U` (created), `-S` (size)
