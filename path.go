package sack

import (
	"fmt"
	"path/filepath"
	"strings"
	"unicode"
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
		c := str[0]
		cl := byte(unicode.ToLower(rune(c)))

		if c == '/' {
			str = strings.Replace(str, "/", "", 1)
		} else if (cl >= 'a' && cl <= 'z') && str[1] == ':' {
			str = strings.Replace(str, fmt.Sprintf("%c:\\", c), "", 1)
			str = strings.Replace(str, fmt.Sprintf("%c:/", c), "", 1)
		}
	}
	return str
}
