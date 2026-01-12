package main

import (
	"slices"
	"testing"

	"github.com/kluzzebass/reflag/translator"
	_ "github.com/kluzzebass/reflag/translator/ls2eza"
)

func TestShellQuote(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"simple", "simple"},
		{"with space", "'with space'"},
		{"with\ttab", "'with\ttab'"},
		{"with\nnewline", "'with\nnewline'"},
		{"with'quote", "'with'\"'\"'quote'"},
		{"with\"double", "'with\"double'"},
		{"with$dollar", "'with$dollar'"},
		{"with`backtick", "'with`backtick'"},
		{"with\\backslash", "'with\\backslash'"},
		{"with!exclaim", "'with!exclaim'"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := shellQuote(tt.input)
			if result != tt.expected {
				t.Errorf("shellQuote(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}

func TestDetectFromBinaryName(t *testing.T) {
	tests := []struct {
		name           string
		binaryName     string
		expectedSource string
		expectedTarget string
		expectedOk     bool
	}{
		{"ls2eza", "ls2eza", "ls", "eza", true},
		{"cat2bat", "cat2bat", "cat", "bat", true},
		{"grep2rg", "grep2rg", "grep", "rg", true},
		{"reflag", "reflag", "", "", false},
		{"no separator", "somecommand", "", "", false},
		{"empty source", "2eza", "", "", false},
		{"empty target", "ls2", "", "", false},
		{"multiple 2s", "ls2eza2foo", "ls", "eza2foo", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			source, target, ok := detectFromBinaryName(tt.binaryName)
			if source != tt.expectedSource || target != tt.expectedTarget || ok != tt.expectedOk {
				t.Errorf("detectFromBinaryName(%q) = (%q, %q, %v), want (%q, %q, %v)",
					tt.binaryName, source, target, ok,
					tt.expectedSource, tt.expectedTarget, tt.expectedOk)
			}
		})
	}
}

func TestParseInitArgs(t *testing.T) {
	tests := []struct {
		name            string
		args            []string
		expectedShell   string
		expectedFilters []string
	}{
		{
			name:            "no args defaults to bash",
			args:            []string{},
			expectedShell:   "bash",
			expectedFilters: nil,
		},
		{
			name:            "shell only",
			args:            []string{"fish"},
			expectedShell:   "fish",
			expectedFilters: nil,
		},
		{
			name:            "shell first then translators",
			args:            []string{"bash", "ls2eza", "grep2rg"},
			expectedShell:   "bash",
			expectedFilters: []string{"ls2eza", "grep2rg"},
		},
		{
			name:            "translators first then shell",
			args:            []string{"ls2eza", "bash", "grep2rg"},
			expectedShell:   "bash",
			expectedFilters: []string{"ls2eza", "grep2rg"},
		},
		{
			name:            "shell at end",
			args:            []string{"ls2eza", "grep2rg", "fish"},
			expectedShell:   "fish",
			expectedFilters: []string{"ls2eza", "grep2rg"},
		},
		{
			name:            "translators only defaults to bash",
			args:            []string{"ls2eza"},
			expectedShell:   "bash",
			expectedFilters: []string{"ls2eza"},
		},
		{
			name:            "zsh shell",
			args:            []string{"zsh", "find2fd"},
			expectedShell:   "zsh",
			expectedFilters: []string{"find2fd"},
		},
		{
			name:            "multiple shells takes last",
			args:            []string{"bash", "fish", "zsh"},
			expectedShell:   "zsh",
			expectedFilters: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			shell, filters := parseInitArgs(tt.args)
			if shell != tt.expectedShell {
				t.Errorf("parseInitArgs(%v) shell = %q, want %q", tt.args, shell, tt.expectedShell)
			}
			if !slices.Equal(filters, tt.expectedFilters) {
				t.Errorf("parseInitArgs(%v) filters = %v, want %v", tt.args, filters, tt.expectedFilters)
			}
		})
	}
}

func TestTranslatorRegistry(t *testing.T) {
	// ls2eza should be registered via init()
	tr := translator.Get("ls", "eza")
	if tr == nil {
		t.Fatal("ls2eza translator not registered")
	}

	if tr.Name() != "ls2eza" {
		t.Errorf("Name() = %q, want %q", tr.Name(), "ls2eza")
	}

	// Test GetByName
	tr2 := translator.GetByName("ls2eza")
	if tr2 == nil {
		t.Fatal("GetByName(ls2eza) returned nil")
	}
	if tr2.Name() != "ls2eza" {
		t.Errorf("GetByName(ls2eza).Name() = %q, want %q", tr2.Name(), "ls2eza")
	}

	// Test List
	names := translator.List()
	if !slices.Contains(names, "ls2eza") {
		t.Error("ls2eza not found in List()")
	}

	// Test Get for non-existent translator
	tr3 := translator.Get("foo", "bar")
	if tr3 != nil {
		t.Error("Get(foo, bar) should return nil")
	}
}
