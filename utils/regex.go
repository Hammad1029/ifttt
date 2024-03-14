package utils

var Regex = struct {
	Endpoint string
	UUID     string
}{
	Endpoint: "^(?=.*?-).+$",
	UUID:     "[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}",
}
