package sack

import (
	"path/filepath"
	"strings"
)

func trimRelativePrefix(p string) string {
	p = strings.TrimPrefix(p, "../")
	p = strings.TrimPrefix(p, "./")
	return p
}

func CleanPath(namepath string, cleanDrive bool) string {
	str := filepath.ToSlash(filepath.Clean(namepath))
	str = trimRelativePrefix(str)

	if cleanDrive && len(str) > 0 {
		if str[0] == '/' {
			str = str[1:]
		} else if vol := filepath.VolumeName(namepath); vol != "" {
			str = strings.TrimPrefix(str, filepath.ToSlash(vol))
			str = strings.TrimPrefix(str, "/")
		}
	}
	return str
}
