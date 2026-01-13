package ping

import (
	"reflect"
	"testing"
)

func TestPreprocessor_Preprocess(t *testing.T) {
	p := &Preprocessor{}

	tests := []struct {
		name string
		args []string
		want []string
	}{
		{
			name: "HTTPS URL with path",
			args: []string{"https://vg.no/index.html"},
			want: []string{"vg.no"},
		},
		{
			name: "HTTP URL with query params",
			args: []string{"http://example.com/search?q=test"},
			want: []string{"example.com"},
		},
		{
			name: "with count flag",
			args: []string{"-c", "4", "https://google.com/"},
			want: []string{"-c", "4", "google.com"},
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
			name: "IP address",
			args: []string{"8.8.8.8"},
			want: []string{"8.8.8.8"},
		},
		{
			name: "URL without scheme",
			args: []string{"example.com/page.html"},
			want: []string{"example.com"},
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

func TestPreprocessor_ToolName(t *testing.T) {
	p := &Preprocessor{}
	if got := p.ToolName(); got != "ping" {
		t.Errorf("ToolName() = %v, want %v", got, "ping")
	}
}
