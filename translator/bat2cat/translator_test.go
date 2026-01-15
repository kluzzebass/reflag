package bat2cat

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
		// Basic usage - always adds plain mode flags
		{
			name:     "simple file",
			input:    []string{"file.txt"},
			expected: []string{"-p", "--paging=never", "--color=auto", "file.txt"},
		},
		{
			name:     "multiple files",
			input:    []string{"file1.txt", "file2.txt"},
			expected: []string{"-p", "--paging=never", "--color=auto", "file1.txt", "file2.txt"},
		},
		{
			name:     "stdin dash",
			input:    []string{"-"},
			expected: []string{"-p", "--paging=never", "--color=auto", "-"},
		},

		// Flags that map to cat equivalents
		{
			name:     "line numbers short",
			input:    []string{"-n", "file.txt"},
			expected: []string{"-p", "--paging=never", "--color=auto", "-n", "file.txt"},
		},
		{
			name:     "line numbers long",
			input:    []string{"--number", "file.txt"},
			expected: []string{"-p", "--paging=never", "--color=auto", "-n", "file.txt"},
		},
		{
			name:     "squeeze blank short",
			input:    []string{"-s", "file.txt"},
			expected: []string{"-p", "--paging=never", "--color=auto", "-s", "file.txt"},
		},
		{
			name:     "squeeze blank long",
			input:    []string{"--squeeze-blank", "file.txt"},
			expected: []string{"-p", "--paging=never", "--color=auto", "-s", "file.txt"},
		},
		{
			name:     "show all short",
			input:    []string{"-A", "file.txt"},
			expected: []string{"-p", "--paging=never", "--color=auto", "-A", "file.txt"},
		},
		{
			name:     "show all long",
			input:    []string{"--show-all", "file.txt"},
			expected: []string{"-p", "--paging=never", "--color=auto", "-A", "file.txt"},
		},
		{
			name:     "unbuffered short",
			input:    []string{"-u", "file.txt"},
			expected: []string{"-p", "--paging=never", "--color=auto", "-u", "file.txt"},
		},
		{
			name:     "unbuffered long",
			input:    []string{"--unbuffered", "file.txt"},
			expected: []string{"-p", "--paging=never", "--color=auto", "-u", "file.txt"},
		},

		// Combined flags
		{
			name:     "combined flags",
			input:    []string{"-ns", "file.txt"},
			expected: []string{"-p", "--paging=never", "--color=auto", "-n", "-s", "file.txt"},
		},
		{
			name:     "combined with number and squeeze",
			input:    []string{"-nsu", "file.txt"},
			expected: []string{"-p", "--paging=never", "--color=auto", "-n", "-s", "-u", "file.txt"},
		},

		// bat-specific flags that are ignored or overridden
		{
			name:     "plain flag ignored (already applied)",
			input:    []string{"-p", "file.txt"},
			expected: []string{"-p", "--paging=never", "--color=auto", "file.txt"},
		},
		{
			name:     "double plain ignored",
			input:    []string{"-pp", "file.txt"},
			expected: []string{"-p", "--paging=never", "--color=auto", "file.txt"},
		},
		{
			name:     "language ignored",
			input:    []string{"-l", "python", "file.py"},
			expected: []string{"-p", "--paging=never", "--color=auto", "file.py"},
		},
		{
			name:     "language long ignored",
			input:    []string{"--language=python", "file.py"},
			expected: []string{"-p", "--paging=never", "--color=auto", "file.py"},
		},
		{
			name:     "highlight ignored",
			input:    []string{"-H", "10:20", "file.txt"},
			expected: []string{"-p", "--paging=never", "--color=auto", "file.txt"},
		},
		{
			name:     "highlight long ignored",
			input:    []string{"--highlight-line=10:20", "file.txt"},
			expected: []string{"-p", "--paging=never", "--color=auto", "file.txt"},
		},
		{
			name:     "color ignored",
			input:    []string{"--color=always", "file.txt"},
			expected: []string{"-p", "--paging=never", "--color=auto", "file.txt"},
		},
		{
			name:     "theme ignored",
			input:    []string{"--theme=Monokai", "file.txt"},
			expected: []string{"-p", "--paging=never", "--color=auto", "file.txt"},
		},
		{
			name:     "style ignored",
			input:    []string{"--style=full", "file.txt"},
			expected: []string{"-p", "--paging=never", "--color=auto", "file.txt"},
		},
		{
			name:     "paging ignored",
			input:    []string{"--paging=always", "file.txt"},
			expected: []string{"-p", "--paging=never", "--color=auto", "file.txt"},
		},
		{
			name:     "line range ignored",
			input:    []string{"--line-range=10:20", "file.txt"},
			expected: []string{"-p", "--paging=never", "--color=auto", "file.txt"},
		},

		// Separator handling
		{
			name:     "separator with files",
			input:    []string{"-n", "--", "file.txt"},
			expected: []string{"-p", "--paging=never", "--color=auto", "-n", "--", "file.txt"},
		},
		{
			name:     "separator with dash file",
			input:    []string{"--", "-", "file.txt"},
			expected: []string{"-p", "--paging=never", "--color=auto", "--", "-", "file.txt"},
		},

		// Mixed valid and ignored flags
		{
			name:     "mixed flags",
			input:    []string{"-n", "--color=always", "-s", "--theme=dark", "file.txt"},
			expected: []string{"-p", "--paging=never", "--color=auto", "-n", "-s", "file.txt"},
		},
		{
			name:     "real world example 1",
			input:    []string{"--style=plain", "--paging=never", "file.txt"},
			expected: []string{"-p", "--paging=never", "--color=auto", "file.txt"},
		},
		{
			name:     "real world example 2",
			input:    []string{"-n", "--decorations=never", "--color=never", "file.txt"},
			expected: []string{"-p", "--paging=never", "--color=auto", "-n", "file.txt"},
		},

		// Force colorization ignored
		{
			name:     "force colorization",
			input:    []string{"-f", "file.txt"},
			expected: []string{"-p", "--paging=never", "--color=auto", "file.txt"},
		},

		// Diff mode ignored
		{
			name:     "diff mode",
			input:    []string{"-d", "file.txt"},
			expected: []string{"-p", "--paging=never", "--color=auto", "file.txt"},
		},

		// Multiple files with various flags
		{
			name:     "complex example",
			input:    []string{"-n", "-s", "file1.txt", "file2.txt", "file3.txt"},
			expected: []string{"-p", "--paging=never", "--color=auto", "-n", "-s", "file1.txt", "file2.txt", "file3.txt"},
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

func TestTranslatorInterface(t *testing.T) {
	tr := &Translator{}

	if tr.Name() != "cat2bat" {
		t.Errorf("Name() = %v, want cat2bat", tr.Name())
	}

	if tr.SourceTool() != "cat" {
		t.Errorf("SourceTool() = %v, want cat", tr.SourceTool())
	}

	if tr.TargetTool() != "bat" {
		t.Errorf("TargetTool() = %v, want bat", tr.TargetTool())
	}

	if tr.IncludeInInit() {
		t.Errorf("IncludeInInit() = true, want false")
	}
}

func TestTranslateMethod(t *testing.T) {
	tr := &Translator{}

	// Test that Translate method calls translateFlags correctly
	input := []string{"-n", "file.txt"}
	expected := []string{"-p", "--paging=never", "--color=auto", "-n", "file.txt"}

	result := tr.Translate(input, "")
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Translate(%v, '') = %v, want %v", input, result, expected)
	}
}
