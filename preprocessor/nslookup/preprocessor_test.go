package nslookup

import (
	"reflect"
	"testing"
)

func TestPreprocessor_ToolName(t *testing.T) {
	p := &Preprocessor{}
	if got := p.ToolName(); got != "nslookup" {
		t.Errorf("ToolName() = %v, want %v", got, "nslookup")
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
			args: []string{"https://cloudflare.com/dns"},
			want: []string{"cloudflare.com"},
		},
		{
			name: "HTTP URL",
			args: []string{"http://example.com/page"},
			want: []string{"example.com"},
		},
		{
			name: "with server argument",
			args: []string{"https://example.com/", "8.8.8.8"},
			want: []string{"example.com", "8.8.8.8"},
		},
		{
			name: "plain hostname",
			args: []string{"example.com"},
			want: []string{"example.com"},
		},
		{
			name: "URL with port",
			args: []string{"https://example.com:8080/api"},
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
