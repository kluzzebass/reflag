package ping6

import (
	"reflect"
	"testing"
)

func TestPreprocessor_ToolName(t *testing.T) {
	p := &Preprocessor{}
	if got := p.ToolName(); got != "ping6" {
		t.Errorf("ToolName() = %v, want %v", got, "ping6")
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
			args: []string{"https://ipv6.google.com/test"},
			want: []string{"ipv6.google.com"},
		},
		{
			name: "HTTP URL",
			args: []string{"http://example.com/page"},
			want: []string{"example.com"},
		},
		{
			name: "with count flag",
			args: []string{"-c", "3", "https://cloudflare.com/"},
			want: []string{"-c", "3", "cloudflare.com"},
		},
		{
			name: "with multiple flags",
			args: []string{"-c", "10", "-i", "0.5", "https://example.com:8080/api"},
			want: []string{"-c", "10", "-i", "0.5", "example.com"},
		},
		{
			name: "plain hostname",
			args: []string{"example.com"},
			want: []string{"example.com"},
		},
		{
			name: "IPv6 address",
			args: []string{"2001:4860:4860::8888"},
			want: []string{"2001:4860:4860::8888"},
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
