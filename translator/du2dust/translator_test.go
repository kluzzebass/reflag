package du2dust

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
			name:     "no args",
			input:    []string{},
			expected: []string{},
		},
		{
			name:     "path only",
			input:    []string{"/tmp"},
			expected: []string{"/tmp"},
		},
		{
			name:     "multiple paths",
			input:    []string{"/tmp", "/var"},
			expected: []string{"/tmp", "/var"},
		},

		// Summarize
		{
			name:     "summarize short",
			input:    []string{"-s", "/tmp"},
			expected: []string{"-d", "0", "/tmp"},
		},
		{
			name:     "summarize long",
			input:    []string{"--summarize", "/tmp"},
			expected: []string{"-d", "0", "/tmp"},
		},

		// Depth
		{
			name:     "max depth short",
			input:    []string{"-d", "2", "/tmp"},
			expected: []string{"-d", "2", "/tmp"},
		},
		{
			name:     "max depth long",
			input:    []string{"--max-depth=3", "/tmp"},
			expected: []string{"-d", "3", "/tmp"},
		},
		{
			name:     "max depth attached",
			input:    []string{"-d2", "/tmp"},
			expected: []string{"-d", "2", "/tmp"},
		},

		// Human readable ignored (dust default)
		{
			name:     "human readable ignored",
			input:    []string{"-h", "/tmp"},
			expected: []string{"/tmp"},
		},
		{
			name:     "human readable long ignored",
			input:    []string{"--human-readable", "/tmp"},
			expected: []string{"/tmp"},
		},

		// All files
		{
			name:     "all files",
			input:    []string{"-a", "/tmp"},
			expected: []string{"-F", "/tmp"},
		},
		{
			name:     "all files long",
			input:    []string{"--all", "/tmp"},
			expected: []string{"-F", "/tmp"},
		},

		// Follow symlinks
		{
			name:     "follow symlinks",
			input:    []string{"-L", "/tmp"},
			expected: []string{"-L", "/tmp"},
		},
		{
			name:     "follow symlinks long",
			input:    []string{"--dereference", "/tmp"},
			expected: []string{"-L", "/tmp"},
		},

		// One file system
		{
			name:     "one file system",
			input:    []string{"-x", "/"},
			expected: []string{"-x", "/"},
		},
		{
			name:     "one file system long",
			input:    []string{"--one-file-system", "/"},
			expected: []string{"-x", "/"},
		},

		// Size formats
		{
			name:     "bytes",
			input:    []string{"-b", "/tmp"},
			expected: []string{"-o", "b", "/tmp"},
		},
		{
			name:     "kilobytes",
			input:    []string{"-k", "/tmp"},
			expected: []string{"-o", "kb", "/tmp"},
		},
		{
			name:     "megabytes",
			input:    []string{"-m", "/tmp"},
			expected: []string{"-o", "mb", "/tmp"},
		},
		{
			name:     "gigabytes BSD",
			input:    []string{"-g", "/tmp"},
			expected: []string{"-o", "gb", "/tmp"},
		},
		{
			name:     "si units",
			input:    []string{"--si", "/tmp"},
			expected: []string{"-o", "si", "/tmp"},
		},

		// Apparent size
		{
			name:     "apparent size",
			input:    []string{"--apparent-size", "/tmp"},
			expected: []string{"-s", "/tmp"},
		},

		// Exclude patterns
		{
			name:     "exclude GNU",
			input:    []string{"--exclude=*.log", "/tmp"},
			expected: []string{"-v", "*.log", "/tmp"},
		},
		{
			name:     "exclude BSD",
			input:    []string{"-I", "*.tmp", "/tmp"},
			expected: []string{"-v", "*.tmp", "/tmp"},
		},

		// Threshold
		{
			name:     "threshold short",
			input:    []string{"-t", "1M", "/tmp"},
			expected: []string{"-z", "1M", "/tmp"},
		},
		{
			name:     "threshold long",
			input:    []string{"--threshold=10M", "/tmp"},
			expected: []string{"-z", "10M", "/tmp"},
		},

		// Combined flags
		{
			name:     "combined sh",
			input:    []string{"-sh", "/tmp"},
			expected: []string{"-d", "0", "/tmp"},
		},
		{
			name:     "typical usage",
			input:    []string{"-shx", "/"},
			expected: []string{"-d", "0", "-x", "/"},
		},

		// Ignored flags
		{
			name:     "total ignored",
			input:    []string{"-c", "/tmp"},
			expected: []string{"/tmp"},
		},
		{
			name:     "count links ignored",
			input:    []string{"-l", "/tmp"},
			expected: []string{"/tmp"},
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

	if tr.Name() != "du2dust" {
		t.Errorf("Name() = %q, want %q", tr.Name(), "du2dust")
	}
	if tr.SourceTool() != "du" {
		t.Errorf("SourceTool() = %q, want %q", tr.SourceTool(), "du")
	}
	if tr.TargetTool() != "dust" {
		t.Errorf("TargetTool() = %q, want %q", tr.TargetTool(), "dust")
	}
}
