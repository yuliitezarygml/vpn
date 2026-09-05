package ray2sing

import (
	"encoding/base64"
	"fmt"
	"regexp"
	"strings"
	"unicode/utf8"
)

var uuidRegex = regexp.MustCompile(`(?i)^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)

func isUUIDString(s string) bool {
	return uuidRegex.MatchString(strings.TrimSpace(s))
}

func looksLikeBase64(s string) bool {
	if len(s) < 4 {
		return false
	}
	valid := 0
	for _, r := range s {
		switch {
		case r >= 'A' && r <= 'Z',
			r >= 'a' && r <= 'z',
			r >= '0' && r <= '9',
			r == '+', r == '/', r == '-', r == '_', r == '=':
			valid++
		default:
			return false
		}
	}
	return valid >= len(s)-1 // tolerate 1–2 bad chars
}

// decodeBase64FaultTolerant tries many variants and returns the first successful decode.
// If none succeed it returns an error containing the debug attempts.
func decodeBase64FaultTolerant(raw string) (string, error) {
	raw = strings.TrimSpace(raw)

	// UUIDs are hex+hyphens; URLEncoding can "successfully" decode them into garbage.
	if isUUIDString(raw) {
		return raw, fmt.Errorf("not base64: uuid")
	}

	padded := raw
	if m := len(padded) % 4; m != 0 {
		padded += strings.Repeat("=", 4-m)
	}

	data, err := base64.StdEncoding.DecodeString(padded)
	if err == nil && utf8.Valid(data) {
		return string(data), nil
	}

	data, err = base64.URLEncoding.DecodeString(padded)
	if err != nil {
		return raw, err
	}
	if !utf8.Valid(data) {
		return raw, fmt.Errorf("not base64: binary result")
	}
	return string(data), nil
}
