package whois

import (
	"reflect"
	"testing"
)

func TestPreprocessor_ToolName(t *testing.T) {
	p := &Preprocessor{}
	if got := p.ToolName(); got != "whois" {
		t.Errorf("ToolName() = %v, want %v", got, "whois")
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
			args: []string{"https://github.com/kluzzebass/reflag"},
			want: []string{"github.com"},
		},
		{
			name: "HTTP URL with query params",
			args: []string{"http://example.com/page?q=test"},
			want: []string{"example.com"},
		},
		{
			name: "with host flag",
			args: []string{"-h", "whois.arin.net", "https://cloudflare.com/"},
			want: []string{"-h", "whois.arin.net", "cloudflare.com"},
		},
		{
			name: "plain domain",
			args: []string{"example.com"},
			want: []string{"example.com"},
		},
		{
			name: "subdomain",
			args: []string{"www.example.com"},
			want: []string{"www.example.com"},
		},
		{
			name: "URL with port",
			args: []string{"https://example.com:8080/api"},
			want: []string{"example.com"},
		},
		{
			name: "URL without scheme",
			args: []string{"example.com/page.html"},
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
