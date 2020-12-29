package helper

import "strings"

// SanitizeHost takes a host value and adds root '.' suffix
func SanitizeHost(host string) string {
	return strings.TrimRight(host, ".") + "."
}
