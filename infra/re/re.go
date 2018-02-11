package re

import "regexp"

func Match(reg string, data string) []string {
	return regexp.MustCompile(reg).FindStringSubmatch(data)
}

