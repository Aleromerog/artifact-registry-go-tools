package npmrc

import (
	"fmt"
	"regexp"
	"strings"
)

var registryRegex = regexp.MustCompile(`(@[a-zA-Z0-9-*~][a-zA-Z0-9-*._~]*:)?registry=https:(\/\/[a-zA-Z0-9-]+[-]npm[.]pkg[.]dev\/.*\/)`)
var authTokenRegex = regexp.MustCompile(`(\/\/[a-zA-Z0-9-]+[-]npm[.]pkg[.]dev\/.*\/):_authToken=(.*)`)
var passwordRegex = regexp.MustCompile(`(\/\/[a-zA-Z0-9-]+[-]npm[.]pkg[.]dev\/.*\/):_password=(.*)`)

var configType = struct {
	Default   string
	Registry  string
	AuthToken string
	Password  string
}{
	Default:   "Default",
	Registry:  "Registry",
	AuthToken: "AuthToken",
	Password:  "Password",
}

type Config = struct {
	Type     string
	Scope    string
	Registry string
	Token    string
	Password string
}

func parseConfig(text string) (string, Config) {
	m := registryRegex.FindStringSubmatch(text)
	if m != nil {
		s := ""
		if m[1] != "" {
			s = strings.TrimSuffix(m[1], ":")
		}
		return m[0], Config{
			Type:     configType.Registry,
			Scope:    s,
			Registry: m[2],
		}
	}
	m = authTokenRegex.FindStringSubmatch(text)
	if m != nil {
		return m[0], Config{
			Type:     configType.AuthToken,
			Registry: m[1],
			Token:    m[2],
		}
	}
	m = passwordRegex.FindStringSubmatch(text)
	if m != nil {
		return fmt.Sprintf("%s:_password=%s", m[1], m[2]), Config{
			Type:     configType.Password,
			Registry: m[1],
			Password: m[2],
		}
	}
	return text, Config{
		Type: configType.Default,
	}
}
