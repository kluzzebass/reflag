package more2moor

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
		// Basic flags
		{
			name:     "exit at EOF",
			input:    []string{"-e"},
			expected: []string{"--quit-if-one-screen"},
		},
		{
			name:     "version short",
			input:    []string{"-V"},
			expected: []string{"-version"},
		},

		// Ignored flags (no moor equivalent)
		{
			name:     "display help prompt",
			input:    []string{"-d"},
			expected: []string{},
		},
		{
			name:     "no pause at form feeds",
			input:    []string{"-l"},
			expected: []string{},
		},
		{
			name:     "count logical lines",
			input:    []string{"-f"},
			expected: []string{},
		},
		{
			name:     "clear screen",
			input:    []string{"-p"},
			expected: []string{},
		},
		{
			name:     "draw from top",
			input:    []string{"-c"},
			expected: []string{},
		},
		{
			name:     "squeeze blank lines",
			input:    []string{"-s"},
			expected: []string{},
		},
		{
			name:     "suppress underlining",
			input:    []string{"-u"},
			expected: []string{},
		},

		// Combined flags
		{
			name:     "combined flags with one translatable",
			input:    []string{"-dse"},
			expected: []string{"--quit-if-one-screen"},
		},
		{
			name:     "all ignored flags",
			input:    []string{"-dlfpcs"},
			expected: []string{},
		},

		// Long flags
		{
			name:     "version long",
			input:    []string{"--version"},
			expected: []string{"-version"},
		},
		{
			name:     "exit on eof long",
			input:    []string{"--exit-on-eof"},
			expected: []string{"--quit-if-one-screen"},
		},
		{
			name:     "no init long",
			input:    []string{"--no-init"},
			expected: []string{"--no-clear-on-exit"},
		},
		{
			name:     "help long ignored",
			input:    []string{"--help"},
			expected: []string{},
		},

		// Numeric line count (ignored)
		{
			name:     "numeric line count",
			input:    []string{"-10"},
			expected: []string{},
		},
		{
			name:     "numeric line count larger",
			input:    []string{"-50"},
			expected: []string{},
		},
		{
			name:     "lines option long",
			input:    []string{"--lines=20"},
			expected: []string{},
		},

		// -n with separate argument
		{
			name:     "lines with -n flag",
			input:    []string{"-n", "25"},
			expected: []string{},
		},
		{
			name:     "lines with -n flag and file",
			input:    []string{"-n", "25", "file.txt"},
			expected: []string{"file.txt"},
		},

		// Initial commands
		{
			name:     "start at line number",
			input:    []string{"+123"},
			expected: []string{"+123"},
		},
		{
			name:     "start at line with file",
			input:    []string{"+50", "file.txt"},
			expected: []string{"+50", "file.txt"},
		},
		{
			name:     "start at pattern ignored",
			input:    []string{"+/pattern", "file.txt"},
			expected: []string{"file.txt"},
		},

		// Files
		{
			name:     "single file",
			input:    []string{"file.txt"},
			expected: []string{"file.txt"},
		},
		{
			name:     "multiple files",
			input:    []string{"file1.txt", "file2.txt"},
			expected: []string{"file1.txt", "file2.txt"},
		},
		{
			name:     "flags and files",
			input:    []string{"-e", "file.txt"},
			expected: []string{"--quit-if-one-screen", "file.txt"},
		},
		{
			name:     "flags and multiple files",
			input:    []string{"-e", "file1.txt", "file2.txt"},
			expected: []string{"--quit-if-one-screen", "file1.txt", "file2.txt"},
		},

		// End of options marker
		{
			name:     "end of options marker",
			input:    []string{"-e", "--", "-file.txt"},
			expected: []string{"--quit-if-one-screen", "-file.txt"},
		},

		// Complex combinations
		{
			name:     "complex combination",
			input:    []string{"-dse", "+100", "file.txt"},
			expected: []string{"--quit-if-one-screen", "+100", "file.txt"},
		},
		{
			name:     "mixed short and long flags",
			input:    []string{"-e", "--no-init", "file.txt"},
			expected: []string{"--quit-if-one-screen", "--no-clear-on-exit", "file.txt"},
		},

		// Empty input
		{
			name:     "empty input",
			input:    []string{},
			expected: []string{},
		},

		// Only files
		{
			name:     "only files no flags",
			input:    []string{"file1.txt", "file2.txt", "file3.txt"},
			expected: []string{"file1.txt", "file2.txt", "file3.txt"},
		},

		// Unknown long flags pass through
		{
			name:     "unknown long flag passes through",
			input:    []string{"--unknown-flag", "file.txt"},
			expected: []string{"--unknown-flag", "file.txt"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := translateFlags(tt.input)
			// Handle nil vs empty slice comparison
			if len(result) == 0 && len(tt.expected) == 0 {
				return
			}
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("translateFlags(%v) = %v, want %v", tt.input, result, tt.expected)
			}
		})
	}
}

func TestTranslator(t *testing.T) {
	tr := &Translator{}

	// Test interface methods
	if tr.Name() != "more2moor" {
		t.Errorf("Name() = %s, want more2moor", tr.Name())
	}
	if tr.SourceTool() != "more" {
		t.Errorf("SourceTool() = %s, want more", tr.SourceTool())
	}
	if tr.TargetTool() != "moor" {
		t.Errorf("TargetTool() = %s, want moor", tr.TargetTool())
	}
	if tr.IncludeInInit() {
		t.Error("IncludeInInit() = true, want false")
	}

	// Test Translate method
	input := []string{"-e", "file.txt"}
	expected := []string{"--quit-if-one-screen", "file.txt"}
	result := tr.Translate(input, "")

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Translate(%v) = %v, want %v", input, result, expected)
	}
}
