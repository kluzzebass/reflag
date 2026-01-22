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
