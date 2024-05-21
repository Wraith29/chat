package log

import "fmt"

func getFormat(lvl, msg string) string {
	return fmt.Sprintf("%s: %s\n", lvl, msg)
}

func Info(msg string, args ...any) {
	fmt.Printf(getFormat("INF", msg), args...)
}

func Err(msg string, args ...any) {
	fmt.Printf(getFormat("ERR", msg), args...)
}
