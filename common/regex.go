package common

import "regexp"

var (
	RegexPositionalParameters          = regexp.MustCompile(`\?`)
	RegexNamedParameters               = regexp.MustCompile(`@\w+`)
	RegexStringInterpolationParameters = regexp.MustCompile(`\$param`)
	RegexEndpoint                      = regexp.MustCompile(`^\/([a-zA-Z0-9-_]+\/?)*$`)
	RegexCron                          = regexp.MustCompile(`(@(annually|yearly|monthly|weekly|daily|hourly|reboot))|(@every (\d+(ns|us|µs|ms|s|m|h))+)|((((\d+,)+\d+|(\d+(\/|-)\d+)|\d+|\*) ?){5,7})`)
	RegexJWT                           = regexp.MustCompile(`^([a-zA-Z0-9_=]+)\.([a-zA-Z0-9_=]+)\.([a-zA-Z0-9_\-\+\/=]*)`)
)
