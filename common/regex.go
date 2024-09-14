package common

import "regexp"

var (
	RegexPositionalParameters = regexp.MustCompile(`=\?`)
	RegexNamedParameters      = regexp.MustCompile(`@\w+`)
	RegexEndpoint             = regexp.MustCompile(`^\/([a-zA-Z0-9-_]+\/?)*$`)
)
