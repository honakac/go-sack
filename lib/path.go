package lib

import (
	"fmt"
	"path/filepath"
	"strings"
)

func trimRelativePrefix(p string) string {
	p = strings.TrimPrefix(p, "../")
	p = strings.TrimPrefix(p, `..\`)
	p = strings.TrimPrefix(p, "./")
	p = strings.TrimPrefix(p, `.\`)
	return p
}

func CleanPath(namepath string, cleanDrive bool) string {
	str := filepath.ToSlash(filepath.Clean(namepath))
	str = trimRelativePrefix(str)

	if cleanDrive {
		if str[0] == '/' {
			str = strings.Replace(str, "/", "", 1)
		} else if (str[0] >= 'A' || str[0] <= 'Z' || str[0] >= 'a' || str[0] <= 'z') && str[1] == ':' {
			str = strings.Replace(str, fmt.Sprintf("%c:\\", str[0]), "", 1)
			str = strings.Replace(str, fmt.Sprintf("%c:/", str[0]), "", 1)
		}
	}
	return str
}
