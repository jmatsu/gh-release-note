package internal

import "math"

func TakeN(s string, n float64) string {
	if s == "" {
		return ""
	}

	return s[0:int(math.Min(n, float64(len(s))))]
}
