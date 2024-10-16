package common

import "regexp"

var (
	RegexPositionalParameters          = regexp.MustCompile(`=\?`)
	RegexNamedParameters               = regexp.MustCompile(`@\w+`)
	RegexStringInterpolationParameters = regexp.MustCompile(`\$param`)
	RegexEndpoint                      = regexp.MustCompile(`^\/([a-zA-Z0-9-_]+\/?)*$`)
)
