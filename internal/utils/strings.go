package utils

import "unicode"

// MixCase чередует верхний и нижний регистр букв, включая кириллицу.
func MixCase(text string) string {
	result := []rune(text)
	upper := true
	for i, letter := range result {
		if !unicode.IsUpper(letter) && !unicode.IsLower(letter) {
			continue
		}
		if upper {
			result[i] = unicode.ToUpper(letter)
		} else {
			result[i] = unicode.ToLower(letter)
		}
		upper = !upper
	}
	return string(result)
}
