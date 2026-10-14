package file

import "strings"

func GetFileExtension(filename string) string {
	str := strings.LastIndex(filename, ".")
	if str < 0 {
		return ""
	}
	return strings.ToLower(filename[str:])
}
