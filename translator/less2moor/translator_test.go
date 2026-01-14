package less2moor

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
		// Basic display flags
		{
			name:     "chop long lines",
			input:    []string{"-S"},
			expected: []string{"--wrap=false"},
		},
		{
			name:     "no line numbers",
			input:    []string{"-N"},
			expected: []string{"--no-linenumbers"},
		},
		{
			name:     "follow mode",
			input:    []string{"-F"},
			expected: []string{"--follow"},
		},
		{
			name:     "no clear on exit",
			input:    []string{"-X"},
			expected: []string{"--no-clear-on-exit"},
		},
		{
			name:     "quit if one screen",
			input:    []string{"-e"},
			expected: []string{"--quit-if-one-screen"},
		},
		{
			name:     "quit if one screen uppercase",
			input:    []string{"-E"},
			expected: []string{"--quit-if-one-screen"},
		},

		// Combined flags
		{
			name:     "combined flags",
			input:    []string{"-SX"},
			expected: []string{"--wrap=false", "--no-clear-on-exit"},
		},
		{
			name:     "combined with ignored flags",
			input:    []string{"-SXr"},
			expected: []string{"--wrap=false", "--no-clear-on-exit"},
		},

		// Long flags
		{
			name:     "long quit if one screen",
			input:    []string{"--quit-if-one-screen"},
			expected: []string{"--quit-if-one-screen"},
		},
		{
			name:     "long no init",
			input:    []string{"--no-init"},
			expected: []string{"--no-clear-on-exit"},
		},
		{
			name:     "long chop long lines",
			input:    []string{"--chop-long-lines"},
			expected: []string{"--wrap=false"},
		},

		// Tab size handling
		{
			name:     "tab size short flag attached",
			input:    []string{"-x4"},
			expected: []string{"-tab-size=4"},
		},
		{
			name:     "tab size long flag with equals",
			input:    []string{"--tabs=4"},
			expected: []string{"-tab-size=4"},
		},

		// Mouse support
		{
			name:     "mouse support",
			input:    []string{"--mouse"},
			expected: []string{"-mousemode=scroll"},
		},

		// Shift amount
		{
			name:     "shift amount short",
			input:    []string{"-#16"},
			expected: []string{"--shift"},
		},
		{
			name:     "shift amount long",
			input:    []string{"--shift=8"},
			expected: []string{"-shift=8"},
		},

		// Version and help
		{
			name:     "version short",
			input:    []string{"-V"},
			expected: []string{"-version"},
		},
		{
			name:     "version long",
			input:    []string{"--version"},
			expected: []string{"-version"},
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
			input:    []string{"-S", "file.txt"},
			expected: []string{"--wrap=false", "file.txt"},
		},
		{
			name:     "flags and multiple files",
			input:    []string{"-SX", "file1.txt", "file2.txt"},
			expected: []string{"--wrap=false", "--no-clear-on-exit", "file1.txt", "file2.txt"},
		},

		// End of options marker
		{
			name:     "end of options marker",
			input:    []string{"-S", "--", "-file.txt"},
			expected: []string{"--wrap=false", "-file.txt"},
		},

		// Ignored flags (should produce no output)
		{
			name:     "raw control chars",
			input:    []string{"-r"},
			expected: []string{},
		},
		{
			name:     "ANSI color support",
			input:    []string{"-R"},
			expected: []string{},
		},
		{
			name:     "quiet mode",
			input:    []string{"-q"},
			expected: []string{},
		},
		{
			name:     "completely quiet",
			input:    []string{"-Q"},
			expected: []string{},
		},
		{
			name:     "squeeze blank lines",
			input:    []string{"-s"},
			expected: []string{},
		},

		// Complex combinations
		{
			name:     "complex combination",
			input:    []string{"-SXF", "+100", "file.txt"},
			expected: []string{"--wrap=false", "--no-clear-on-exit", "--follow", "+100", "file.txt"},
		},
		{
			name:     "mixed short and long flags",
			input:    []string{"-S", "--quit-if-one-screen", "file.txt"},
			expected: []string{"--wrap=false", "--quit-if-one-screen", "file.txt"},
		},
		{
			name:     "all ignored flags",
			input:    []string{"-rRqs"},
			expected: []string{},
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
	if tr.Name() != "less2moor" {
		t.Errorf("Name() = %s, want less2moor", tr.Name())
	}
	if tr.SourceTool() != "less" {
		t.Errorf("SourceTool() = %s, want less", tr.SourceTool())
	}
	if tr.TargetTool() != "moor" {
		t.Errorf("TargetTool() = %s, want moor", tr.TargetTool())
	}
	if tr.IncludeInInit() {
		t.Error("IncludeInInit() = true, want false")
	}

	// Test Translate method
	input := []string{"-S", "file.txt"}
	expected := []string{"--wrap=false", "file.txt"}
	result := tr.Translate(input, "")

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Translate(%v) = %v, want %v", input, result, expected)
	}
}
