package ray2sing

import (
	"net/url"
	"strings"
	"unicode"
	"unicode/utf8"

	T "github.com/sagernet/sing-box/option"
)

// HysteriaURLData holds the parsed data from a Hysteria URL.
type UrlSchema struct {
	Scheme   string
	Username string
	Password string
	Hostname string
	Port     uint16
	Name     string
	Params   map[string]string
}

func (u UrlSchema) GetServerOption() T.ServerOptions {
	return T.ServerOptions{
		Server:     u.Hostname,
		ServerPort: u.Port,
	}
}

func (u UrlSchema) GetRelayOptions() (*T.TurnRelayOptions, error) {
	return ParseTurnURL(u.Params["relay"])
}

func isPrintableUserInfo(s string) bool {
	if s == "" || !utf8.ValidString(s) {
		return false
	}
	for _, r := range s {
		if r == unicode.ReplacementChar || (!unicode.IsPrint(r) && !unicode.IsSpace(r)) {
			return false
		}
	}
	return true
}

// parseHysteria2 parses a given URL and returns a HysteriaURLData struct.
func ParseUrl(inputURL string, defaultPort uint16) (*UrlSchema, error) {
	parsedURL, err := url.Parse(inputURL)
	if err != nil {
		return nil, err
	}
	port := toUInt16(parsedURL.Port(), defaultPort)

	data := &UrlSchema{
		Scheme:   parsedURL.Scheme,
		Username: parsedURL.User.Username(),
		Password: getPassword(parsedURL),
		Hostname: parsedURL.Hostname(),
		Port:     port,
		Name:     parsedURL.Fragment,
		Params:   make(map[string]string),
	}

	// Shadowsocks-style userinfo may be base64(method:password). Never treat UUIDs as base64.
	if data.Username != "" && !isUUIDString(data.Username) {
		userInfo, err := decodeBase64IfNeeded(data.Username)
		if err == nil && userInfo != data.Username && isPrintableUserInfo(userInfo) {
			userDetails := strings.SplitN(userInfo, ":", 2)
			if len(userDetails) == 2 && userDetails[0] != "" && userDetails[1] != "" {
				data.Username = userDetails[0]
				data.Password = userDetails[1]
			}
		}
	}

	for key, values := range parsedURL.Query() {
		data.Params[strings.ReplaceAll(strings.ToLower(key), "_", "")] = strings.Join(values, ",")
	}

	return data, nil
}

func getPassword(u *url.URL) string {
	if password, ok := u.User.Password(); ok {
		return password
	}
	return ""
}
