package urlparse

import (
	"net/url"
	"strings"
)

// ExtractHostname extracts the hostname from a URL string.
// If the input is not a valid URL, it returns the input unchanged.
// This handles various URL formats and edge cases:
// - Full URLs: https://example.com/path -> example.com
// - URLs with ports: http://example.com:8080 -> example.com
// - IPv4/IPv6 addresses: passed through unchanged
// - Plain hostnames: passed through unchanged
func ExtractHostname(input string) string {
	if input == "" {
		return input
	}

	// Try to parse as URL
	parsed, err := url.Parse(input)
	if err != nil {
		// Not a valid URL, return as-is
		return input
	}

	// If no scheme was detected, url.Parse might put everything in Path
	// Try adding a scheme and parsing again
	if parsed.Scheme == "" && parsed.Host == "" {
		// Check if it looks like a URL without scheme
		if strings.Contains(input, "/") || strings.Contains(input, "?") {
			parsed, err = url.Parse("http://" + input)
			if err != nil {
				return input
			}
		} else {
			// Plain hostname or IP, return as-is
			return input
		}
	}

	// Extract hostname (without port)
	hostname := parsed.Hostname()
	if hostname != "" {
		return hostname
	}

	// If we still don't have a hostname, return the original
	return input
}

// ProcessArgs processes a slice of arguments, extracting hostnames from URLs
// while preserving flags and other arguments unchanged.
// It identifies positional arguments (non-flags) and applies hostname extraction.
func ProcessArgs(args []string) []string {
	if len(args) == 0 {
		return args
	}

	result := make([]string, len(args))
	copy(result, args)

	for i, arg := range result {
		// Skip flags (arguments starting with -)
		if strings.HasPrefix(arg, "-") {
			continue
		}

		// Skip if it's clearly not a URL (no scheme indicators)
		if !looksLikeURL(arg) {
			continue
		}

		// Try to extract hostname
		extracted := ExtractHostname(arg)
		if extracted != arg {
			result[i] = extracted
		}
	}

	return result
}

// looksLikeURL returns true if the string might be a URL that needs processing
func looksLikeURL(s string) bool {
	// Check for common URL schemes
	if strings.HasPrefix(s, "http://") ||
		strings.HasPrefix(s, "https://") ||
		strings.HasPrefix(s, "ftp://") ||
		strings.HasPrefix(s, "ftps://") {
		return true
	}

	// Check if it contains URL-like components (path or query)
	// but avoid matching things like file paths
	if strings.Contains(s, "://") {
		return true
	}

	// Check for domain-like patterns with paths
	// e.g., example.com/path or www.example.com/page
	// But skip CIDR notation (e.g., 192.0.2.0/24)
	if strings.Contains(s, "/") && strings.Contains(s, ".") {
		parts := strings.Split(s, "/")
		if len(parts) > 1 && strings.Contains(parts[0], ".") {
			// Skip if it looks like CIDR notation (IP/prefix)
			// Check if the part after / is just digits (subnet prefix)
			if len(parts) == 2 {
				isDigits := true
				for _, c := range parts[1] {
					if c < '0' || c > '9' {
						isDigits = false
						break
					}
				}
				if isDigits {
					return false // It's CIDR notation, not a URL
				}
			}
			return true
		}
	}

	return false
}
