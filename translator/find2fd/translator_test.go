package find2fd

import (
	"reflect"
	"testing"
)

func TestTranslateFlags(t *testing.T) {
	tests := []struct {
		name     string
		input    []string
		expected []string
	}{
		// Basic usage
		{
			name:     "find current dir",
			input:    []string{"."},
			expected: []string{},
		},
		{
			name:     "find specific dir",
			input:    []string{"/tmp"},
			expected: []string{".", "/tmp"},
		},
		{
			name:     "find multiple dirs",
			input:    []string{"src", "lib"},
			expected: []string{".", "src", "lib"},
		},

		// -name patterns
		{
			name:     "name pattern simple",
			input:    []string{".", "-name", "*.txt"},
			expected: []string{"\\.txt$"},
		},
		{
			name:     "name pattern go files",
			input:    []string{".", "-name", "*.go"},
			expected: []string{"\\.go$"},
		},
		{
			name:     "name exact file",
			input:    []string{".", "-name", "Makefile"},
			expected: []string{"Makefile"},
		},
		{
			name:     "iname case insensitive",
			input:    []string{".", "-iname", "*.TXT"},
			expected: []string{"-i", "\\.TXT$"},
		},

		// -type
		{
			name:     "type file",
			input:    []string{".", "-type", "f"},
			expected: []string{"-t", "f"},
		},
		{
			name:     "type directory",
			input:    []string{".", "-type", "d"},
			expected: []string{"-t", "d"},
		},
		{
			name:     "type symlink",
			input:    []string{".", "-type", "l"},
			expected: []string{"-t", "l"},
		},

		// Depth
		{
			name:     "maxdepth",
			input:    []string{".", "-maxdepth", "2"},
			expected: []string{"-d", "2"},
		},
		{
			name:     "mindepth",
			input:    []string{".", "-mindepth", "1"},
			expected: []string{"--min-depth", "1"},
		},
		{
			name:     "both depths",
			input:    []string{".", "-mindepth", "1", "-maxdepth", "3"},
			expected: []string{"--min-depth", "1", "-d", "3"},
		},

		// Combined expressions
		{
			name:     "type and name",
			input:    []string{".", "-type", "f", "-name", "*.go"},
			expected: []string{"-t", "f", "\\.go$"},
		},
		{
			name:     "name and maxdepth",
			input:    []string{".", "-maxdepth", "2", "-name", "*.txt"},
			expected: []string{"-d", "2", "\\.txt$"},
		},
		{
			name:     "typical find usage",
			input:    []string{".", "-type", "f", "-name", "*.go", "-maxdepth", "3"},
			expected: []string{"-t", "f", "-d", "3", "\\.go$"},
		},

		// Path with expressions
		{
			name:     "src dir with type",
			input:    []string{"src", "-type", "f"},
			expected: []string{"-t", "f", ".", "src"},
		},

		// -print0
		{
			name:     "print0",
			input:    []string{".", "-name", "*.txt", "-print0"},
			expected: []string{"-0", "\\.txt$"},
		},

		// -print ignored
		{
			name:     "print ignored",
			input:    []string{".", "-name", "*.txt", "-print"},
			expected: []string{"\\.txt$"},
		},

		// Follow symlinks
		{
			name:     "follow symlinks L",
			input:    []string{"-L", ".", "-name", "*.txt"},
			expected: []string{"-L", "\\.txt$"},
		},
		{
			name:     "follow symlinks word",
			input:    []string{"-follow", ".", "-name", "*.txt"},
			expected: []string{"-L", "\\.txt$"},
		},

		// Empty and executable
		{
			name:     "empty",
			input:    []string{".", "-empty"},
			expected: []string{"-t", "e"},
		},
		{
			name:     "executable",
			input:    []string{".", "-executable"},
			expected: []string{"-t", "x"},
		},

		// Time expressions
		{
			name:     "mtime within",
			input:    []string{".", "-mtime", "-7"},
			expected: []string{"--changed-within", "7d"},
		},
		{
			name:     "mtime before",
			input:    []string{".", "-mtime", "+30"},
			expected: []string{"--changed-before", "30d"},
		},
		{
			name:     "mmin within",
			input:    []string{".", "-mmin", "-60"},
			expected: []string{"--changed-within", "60min"},
		},

		// Size
		{
			name:     "size",
			input:    []string{".", "-size", "+1M"},
			expected: []string{"-S", "+1M"},
		},

		// Newer than file
		{
			name:     "newer than file",
			input:    []string{".", "-newer", "reference.txt"},
			expected: []string{"--newer", "reference.txt"},
		},

		// User/group
		{
			name:     "user",
			input:    []string{".", "-user", "root"},
			expected: []string{"--owner", "root"},
		},
		{
			name:     "group",
			input:    []string{".", "-group", "wheel"},
			expected: []string{"--owner", ":wheel"},
		},

		// Logical operators ignored
		{
			name:     "and ignored",
			input:    []string{".", "-type", "f", "-a", "-name", "*.go"},
			expected: []string{"-t", "f", "\\.go$"},
		},
		{
			name:     "parens ignored",
			input:    []string{".", "(", "-name", "*.go", ")"},
			expected: []string{"\\.go$"},
		},

		// One file system
		{
			name:     "one file system",
			input:    []string{".", "-xdev"},
			expected: []string{"--one-file-system"},
		},

		// Regex
		{
			name:     "regex pattern",
			input:    []string{".", "-regex", ".*\\.go$"},
			expected: []string{".*\\.go$"},
		},
		{
			name:     "iregex pattern",
			input:    []string{".", "-iregex", ".*\\.GO$"},
			expected: []string{"-i", ".*\\.GO$"},
		},

		// -path
		{
			name:     "path pattern",
			input:    []string{".", "-path", "*/test/*"},
			expected: []string{"-p", "*/test/*"},
		},

		// Empty input
		{
			name:     "empty input",
			input:    []string{},
			expected: []string{},
		},

		// Quit/single result
		{
			name:     "quit",
			input:    []string{".", "-name", "*.go", "-quit"},
			expected: []string{"-1", "\\.go$"},
		},

		// Real-world find usage patterns
		// find ~/.local (list all files in directory)
		{
			name:     "find home subdirectory",
			input:    []string{"/Users/ove/.local"},
			expected: []string{".", "/Users/ove/.local"},
		},
		// find /var/log -name "*.log" -mtime -1 (recent logs)
		{
			name:     "find recent logs",
			input:    []string{"/var/log", "-name", "*.log", "-mtime", "-1"},
			expected: []string{"--changed-within", "1d", "\\.log$", "/var/log"},
		},
		// find . -type f -size +100M (large files)
		{
			name:     "find large files",
			input:    []string{".", "-type", "f", "-size", "+100M"},
			expected: []string{"-t", "f", "-S", "+100M"},
		},
		// find . -type f -name "*.log" -mtime +7 (old log files)
		{
			name:     "find old log files",
			input:    []string{".", "-type", "f", "-name", "*.log", "-mtime", "+7"},
			expected: []string{"-t", "f", "--changed-before", "7d", "\\.log$"},
		},
		// find /tmp /var/tmp -type f (multiple directories)
		{
			name:     "find in multiple directories",
			input:    []string{"/tmp", "/var/tmp", "-type", "f"},
			expected: []string{"-t", "f", ".", "/tmp", "/var/tmp"},
		},
		// find . -maxdepth 1 -type f (only current directory)
		{
			name:     "find files in current dir only",
			input:    []string{".", "-maxdepth", "1", "-type", "f"},
			expected: []string{"-d", "1", "-t", "f"},
		},
		// find ~ -name ".bashrc" (find dotfiles)
		{
			name:     "find dotfile in home",
			input:    []string{"/Users/ove", "-name", ".bashrc"},
			expected: []string{"\\.bashrc", "/Users/ove"},
		},
		// find . -type d -name "node_modules" (find directories by name)
		{
			name:     "find node_modules directories",
			input:    []string{".", "-type", "d", "-name", "node_modules"},
			expected: []string{"-t", "d", "node_modules"},
		},
		// find . -empty -type f (empty files)
		{
			name:     "find empty files",
			input:    []string{".", "-empty", "-type", "f"},
			expected: []string{"-t", "e", "-t", "f"},
		},
		// find . -amin -30 (accessed in last 30 minutes)
		{
			name:     "find recently accessed",
			input:    []string{".", "-amin", "-30"},
			expected: []string{"--changed-within", "30min"},
		},
		// find . -ctime -1 (changed in last day)
		{
			name:     "find recently changed ctime",
			input:    []string{".", "-ctime", "-1"},
			expected: []string{"--changed-within", "1d"},
		},
		// find . -type f -name "*.txt" -print0 | xargs -0 ... (null-separated)
		{
			name:     "find for xargs with null separator",
			input:    []string{".", "-type", "f", "-name", "*.txt", "-print0"},
			expected: []string{"-t", "f", "-0", "\\.txt$"},
		},
		// find /home -user root -type f (files owned by root)
		{
			name:     "find files owned by root",
			input:    []string{"/home", "-user", "root", "-type", "f"},
			expected: []string{"--owner", "root", "-t", "f", ".", "/home"},
		},
		// find . -L -type l (follow symlinks, find broken links)
		{
			name:     "find with follow symlinks",
			input:    []string{"-L", ".", "-type", "l"},
			expected: []string{"-L", "-t", "l"},
		},
		// find . -mindepth 2 -maxdepth 4 -type f (depth range)
		{
			name:     "find with depth range",
			input:    []string{".", "-mindepth", "2", "-maxdepth", "4", "-type", "f"},
			expected: []string{"--min-depth", "2", "-d", "4", "-t", "f"},
		},
		// find project/ -name "*.js" -type f (search in subdirectory)
		{
			name:     "find js files in project",
			input:    []string{"project/", "-name", "*.js", "-type", "f"},
			expected: []string{"-t", "f", "\\.js$", "project/"},
		},
		// find . -iname "readme*" (case insensitive glob)
		{
			name:     "find readme case insensitive",
			input:    []string{".", "-iname", "readme*"},
			expected: []string{"-i", "readme[^/]*"},
		},
		// find . -name "*.go" -newer go.mod (newer than reference)
		{
			name:     "find go files newer than go.mod",
			input:    []string{".", "-name", "*.go", "-newer", "go.mod"},
			expected: []string{"--newer", "go.mod", "\\.go$"},
		},
		// find /etc -type f -size +1k -size -100k (size range - partial support)
		{
			name:     "find config files by size",
			input:    []string{"/etc", "-type", "f", "-size", "+1k"},
			expected: []string{"-t", "f", "-S", "+1k", ".", "/etc"},
		},
		// find . -xdev -type f (single filesystem)
		{
			name:     "find on single filesystem",
			input:    []string{".", "-xdev", "-type", "f"},
			expected: []string{"--one-file-system", "-t", "f"},
		},
		// find . -mount -name "*.bak" (mount is alias for xdev)
		{
			name:     "find with mount option",
			input:    []string{".", "-mount", "-name", "*.bak"},
			expected: []string{"--one-file-system", "\\.bak$"},
		},
		// find . -path "*/test/*" -name "*.go" (path and name combined)
		{
			name:     "find test go files by path",
			input:    []string{".", "-path", "*/test/*", "-name", "*.go"},
			expected: []string{"-p", "*/test/*", "\\.go$"},
		},
		// find . -not -name "*.txt" (negation - limited support)
		{
			name:     "find with negation ignored",
			input:    []string{".", "-not", "-name", "*.txt"},
			expected: []string{"\\.txt$"},
		},
		// find . -name "*.tar.gz" (compound extension)
		{
			name:     "find tar.gz files",
			input:    []string{".", "-name", "*.tar.gz"},
			expected: []string{"\\.tar\\.gz$"},
		},
		// find . -H -name "*.sh" (H option)
		{
			name:     "find with H option",
			input:    []string{"-H", ".", "-name", "*.sh"},
			expected: []string{"-H", "\\.sh$"},
		},
		// find . -group staff -type f
		{
			name:     "find files by group",
			input:    []string{".", "-group", "staff", "-type", "f"},
			expected: []string{"--owner", ":staff", "-t", "f"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := translateFlags(tt.input)
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("translateFlags(%v) = %v, want %v", tt.input, result, tt.expected)
			}
		})
	}
}

func TestGlobToRegex(t *testing.T) {
	tests := []struct {
		glob     string
		expected string
	}{
		{"*.txt", "\\.txt$"},
		{"*.go", "\\.go$"},
		{"*.tar.gz", "\\.tar\\.gz$"},
		{"Makefile", "Makefile"},
		{"test*", "test[^/]*"},
		{"?oo", "[^/]oo"},
		{"file.txt", "file\\.txt"},
		{"[abc].txt", "[abc]\\.txt"},
		{"[!abc].txt", "[^abc]\\.txt"},
	}

	for _, tt := range tests {
		t.Run(tt.glob, func(t *testing.T) {
			result := globToRegex(tt.glob)
			if result != tt.expected {
				t.Errorf("globToRegex(%q) = %q, want %q", tt.glob, result, tt.expected)
			}
		})
	}
}

func TestTranslatorInterface(t *testing.T) {
	tr := &Translator{}

	if tr.Name() != "find2fd" {
		t.Errorf("Name() = %q, want %q", tr.Name(), "find2fd")
	}
	if tr.SourceTool() != "find" {
		t.Errorf("SourceTool() = %q, want %q", tr.SourceTool(), "find")
	}
	if tr.TargetTool() != "fd" {
		t.Errorf("TargetTool() = %q, want %q", tr.TargetTool(), "fd")
	}
}
