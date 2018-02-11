package str

import "strconv"

func ToString(args ...interface{}) string {
	result := ""
	for _, arg := range args {
		switch val := arg.(type) {
		case int:
			result += strconv.Itoa(val)
		case string:
			result += val
		}
	}
	return result
}