package util

import "unicode/utf8"

func TruncateString(s string, maxLength int) string {

	if utf8.RuneCountInString(s) > maxLength {
		runes := []rune(s)

		// can't append an ellipsis if maxLength < 3
		if maxLength < 3 {
			return string(runes[:maxLength])
		}

		return string(runes[:maxLength-3]) + "..."
	}

	return s
}
