package urlparse

import (
	"reflect"
	"testing"
)

func TestExtractHostname(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{
			name:  "full HTTPS URL with path",
			input: "https://vg.no/index.html",
			want:  "vg.no",
		},
		{
			name:  "HTTP URL with port and path",
			input: "http://example.com:8080/path/to/page",
			want:  "example.com",
		},
		{
			name:  "HTTPS URL with query params",
			input: "https://google.com/search?q=test",
			want:  "google.com",
		},
		{
			name:  "plain hostname",
			input: "example.com",
			want:  "example.com",
		},
		{
			name:  "plain hostname with subdomain",
			input: "www.example.com",
			want:  "www.example.com",
		},
		{
			name:  "IPv4 address",
			input: "192.168.1.1",
			want:  "192.168.1.1",
		},
		{
			name:  "URL without scheme but with path",
			input: "example.com/page.html",
			want:  "example.com",
		},
		{
			name:  "FTP URL",
			input: "ftp://ftp.example.com/file.txt",
			want:  "ftp.example.com",
		},
		{
			name:  "URL with authentication",
			input: "https://user:pass@example.com/path",
			want:  "example.com",
		},
		{
			name:  "empty string",
			input: "",
			want:  "",
		},
		{
			name:  "localhost",
			input: "localhost",
			want:  "localhost",
		},
		{
			name:  "localhost with port",
			input: "http://localhost:8080",
			want:  "localhost",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ExtractHostname(tt.input)
			if got != tt.want {
				t.Errorf("ExtractHostname(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}

func TestProcessArgs(t *testing.T) {
	tests := []struct {
		name string
		args []string
		want []string
	}{
		{
			name: "ping with HTTPS URL",
			args: []string{"https://vg.no/index.html"},
			want: []string{"vg.no"},
		},
		{
			name: "ping with flags and URL",
			args: []string{"-c", "4", "https://example.com/page"},
			want: []string{"-c", "4", "example.com"},
		},
		{
			name: "multiple arguments with URL",
			args: []string{"-t", "10", "http://google.com/search?q=test"},
			want: []string{"-t", "10", "google.com"},
		},
		{
			name: "plain hostname unchanged",
			args: []string{"example.com"},
			want: []string{"example.com"},
		},
		{
			name: "flags only",
			args: []string{"-v", "-c", "5"},
			want: []string{"-v", "-c", "5"},
		},
		{
			name: "empty args",
			args: []string{},
			want: []string{},
		},
		{
			name: "mixed flags and hostnames",
			args: []string{"-c", "3", "example.com", "-t", "5"},
			want: []string{"-c", "3", "example.com", "-t", "5"},
		},
		{
			name: "URL with port",
			args: []string{"https://example.com:8443/api"},
			want: []string{"example.com"},
		},
		{
			name: "CIDR notation unchanged",
			args: []string{"+subnet=192.0.2.0/24", "example.com"},
			want: []string{"+subnet=192.0.2.0/24", "example.com"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ProcessArgs(tt.args)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ProcessArgs(%v) = %v, want %v", tt.args, got, tt.want)
			}
		})
	}
}

func TestLooksLikeURL(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  bool
	}{
		{
			name:  "HTTPS URL",
			input: "https://example.com",
			want:  true,
		},
		{
			name:  "HTTP URL",
			input: "http://example.com/path",
			want:  true,
		},
		{
			name:  "FTP URL",
			input: "ftp://ftp.example.com",
			want:  true,
		},
		{
			name:  "custom scheme",
			input: "custom://host",
			want:  true,
		},
		{
			name:  "domain with path",
			input: "example.com/page",
			want:  true,
		},
		{
			name:  "plain hostname",
			input: "example.com",
			want:  false,
		},
		{
			name:  "plain word",
			input: "localhost",
			want:  false,
		},
		{
			name:  "IP address",
			input: "192.168.1.1",
			want:  false,
		},
		{
			name:  "file path",
			input: "/usr/local/bin",
			want:  false,
		},
		{
			name:  "CIDR notation IPv4",
			input: "192.0.2.0/24",
			want:  false,
		},
		{
			name:  "CIDR notation IPv6",
			input: "2001:db8::/32",
			want:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := looksLikeURL(tt.input)
			if got != tt.want {
				t.Errorf("looksLikeURL(%q) = %v, want %v", tt.input, got, tt.want)
			}
		})
	}
}
