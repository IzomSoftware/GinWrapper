package file

import (
	"io"
	"strings"

	"github.com/gabriel-vasile/mimetype"
)

func GetFileExtension(filename string) string {
	str := strings.LastIndex(filename, ".")
	if str < 0 {
		return ""
	}
	return strings.ToLower(filename[str:])
}

func DetectMime(reader io.ReadSeeker) (string, error) {
	begin, err := reader.Seek(0, io.SeekCurrent)
	if err != nil {
		return "", err
	}

	defer reader.Seek(begin, io.SeekStart)

	if _, err := reader.Seek(0, io.SeekStart); err != nil {
		return "", err
	}

	mimetype, err := mimetype.DetectReader(reader)
	if err != nil {
		return "", err
	}

	return mimetype.String(), nil
}
