package ps2procs

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
		// Basic usage - most flags ignored since procs shows all by default
		{
			name:     "no args",
			input:    []string{},
			expected: []string{},
		},
		{
			name:     "aux BSD style",
			input:    []string{"aux"},
			expected: []string{},
		},
		{
			name:     "ef UNIX style",
			input:    []string{"-ef"},
			expected: []string{},
		},
		{
			name:     "all processes -e",
			input:    []string{"-e"},
			expected: []string{},
		},
		{
			name:     "all processes -A",
			input:    []string{"-A"},
			expected: []string{},
		},

		// Tree view
		{
			name:     "forest GNU",
			input:    []string{"--forest"},
			expected: []string{"--tree"},
		},
		{
			name:     "hierarchy -H",
			input:    []string{"-H"},
			expected: []string{"--tree"},
		},
		{
			name:     "forest BSD f",
			input:    []string{"axjf"},
			expected: []string{"--tree"},
		},

		// User filter
		{
			name:     "user filter -u",
			input:    []string{"-u", "root"},
			expected: []string{"root"},
		},
		{
			name:     "user filter -U",
			input:    []string{"-U", "www-data"},
			expected: []string{"www-data"},
		},
		{
			name:     "user filter attached",
			input:    []string{"-uroot"},
			expected: []string{"root"},
		},

		// PID filter
		{
			name:     "pid filter",
			input:    []string{"-p", "1234"},
			expected: []string{"1234"},
		},
		{
			name:     "pid filter attached",
			input:    []string{"-p1234"},
			expected: []string{"1234"},
		},

		// Command filter
		{
			name:     "command filter",
			input:    []string{"-C", "nginx"},
			expected: []string{"nginx"},
		},

		// Sort
		{
			name:     "sort ascending",
			input:    []string{"--sort=cpu"},
			expected: []string{"--sorta", "cpu"},
		},
		{
			name:     "sort descending",
			input:    []string{"--sort=-mem"},
			expected: []string{"--sortd", "mem"},
		},
		{
			name:     "sort with plus",
			input:    []string{"--sort=+pid"},
			expected: []string{"--sorta", "pid"},
		},
		{
			name:     "sort column mapping",
			input:    []string{"--sort=%cpu"},
			expected: []string{"--sorta", "cpu"},
		},

		// Combined
		{
			name:     "typical ps aux",
			input:    []string{"aux"},
			expected: []string{},
		},
		{
			name:     "ps with user",
			input:    []string{"-ef", "-u", "root"},
			expected: []string{"root"},
		},
		{
			name:     "ps tree with user",
			input:    []string{"--forest", "-u", "root"},
			expected: []string{"--tree", "root"},
		},

		// Search term (not BSD options)
		{
			name:     "search by name",
			input:    []string{"nginx"},
			expected: []string{"nginx"},
		},
		{
			name:     "search by pid",
			input:    []string{"1234"},
			expected: []string{"1234"},
		},

		// Ignored flags with values
		{
			name:     "output format ignored",
			input:    []string{"-o", "pid,comm"},
			expected: []string{},
		},
		{
			name:     "tty ignored",
			input:    []string{"-t", "pts/0"},
			expected: []string{},
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

func TestIsBSDStyleOptions(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"aux", true},
		{"ef", true},
		{"axjf", true},
		{"ax", true},
		{"u", true},
		{"nginx", false},   // no common BSD chars (auxef)
		{"1234", false},    // looks like a PID
		{"toolong", false}, // too many chars
		{"root", false},    // no common BSD chars
		{"ps", false},      // no common BSD chars
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := isBSDStyleOptions(tt.input)
			if result != tt.expected {
				t.Errorf("isBSDStyleOptions(%q) = %v, want %v", tt.input, result, tt.expected)
			}
		})
	}
}

func TestTranslatorInterface(t *testing.T) {
	tr := &Translator{}

	if tr.Name() != "ps2procs" {
		t.Errorf("Name() = %q, want %q", tr.Name(), "ps2procs")
	}
	if tr.SourceTool() != "ps" {
		t.Errorf("SourceTool() = %q, want %q", tr.SourceTool(), "ps")
	}
	if tr.TargetTool() != "procs" {
		t.Errorf("TargetTool() = %q, want %q", tr.TargetTool(), "procs")
	}
}
