package utils

import "strings"

func SanitiseMessage(b []byte) []byte {
	return []byte(strings.Trim(string(b), " \n\r"))
}
