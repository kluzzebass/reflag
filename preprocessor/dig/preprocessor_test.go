package dig

import (
	"reflect"
	"testing"
)

func TestPreprocessor_ToolName(t *testing.T) {
	p := &Preprocessor{}
	if got := p.ToolName(); got != "dig" {
		t.Errorf("ToolName() = %v, want %v", got, "dig")
	}
}

func TestPreprocessor_Description(t *testing.T) {
	p := &Preprocessor{}
	desc := p.Description()
	if desc == "" {
		t.Error("Description() returned empty string")
	}
}

func TestPreprocessor_Preprocess(t *testing.T) {
	p := &Preprocessor{}

	tests := []struct {
		name string
		args []string
		want []string
	}{
		{
			name: "HTTPS URL with path",
			args: []string{"https://example.com/page.html"},
			want: []string{"example.com"},
		},
		{
			name: "HTTP URL with query params",
			args: []string{"http://cloudflare.com/dns?q=test"},
			want: []string{"cloudflare.com"},
		},
		{
			name: "with query type",
			args: []string{"https://example.com/page", "MX"},
			want: []string{"example.com", "MX"},
		},
		{
			name: "with flags and URL",
			args: []string{"+short", "https://google.com/"},
			want: []string{"+short", "google.com"},
		},
		{
			name: "plain hostname",
			args: []string{"example.com"},
			want: []string{"example.com"},
		},
		{
			name: "IP address",
			args: []string{"8.8.8.8"},
			want: []string{"8.8.8.8"},
		},
		{
			name: "URL with port",
			args: []string{"https://example.com:8443/api"},
			want: []string{"example.com"},
		},
		{
			name: "empty args",
			args: []string{},
			want: []string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := p.Preprocess(tt.args)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Preprocess(%v) = %v, want %v", tt.args, got, tt.want)
			}
		})
	}
}
