package sws

import (
	"fmt"
	"strings"
)

type ContentSecurityPolicy struct {
	Default []string
	Script  []string
	Image   []string
	Style   []string
	Connect []string
	Frame   []string
}

func (csp ContentSecurityPolicy) String() string {
	result := ""
	if len(csp.Default) > 0 {
		result += fmt.Sprintf("default-src %s;", strings.Join(csp.Default, " "))
	}

	if len(csp.Script) > 0 {
		result += fmt.Sprintf("script-src %s;", strings.Join(csp.Script, " "))
	}

	if len(csp.Style) > 0 {
		result += fmt.Sprintf("style-src %s;", strings.Join(csp.Style, " "))
	}

	if len(csp.Connect) > 0 {
		result += fmt.Sprintf("connect-src %s;", strings.Join(csp.Connect, " "))
	}

	if len(csp.Frame) > 0 {
		result += fmt.Sprintf("frame-src %s;", strings.Join(csp.Frame, " "))
	}

	if len(csp.Image) > 0 {
		result += fmt.Sprintf("img-src %s;", strings.Join(csp.Image, " "))
	}

	return result
}
