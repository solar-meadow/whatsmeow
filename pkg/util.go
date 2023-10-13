package pkg

import (
	"regexp"
)

func ExtractPhoneNumber(message string) (string, bool) {
	regex := `код \+(\d{11})`
	re := regexp.MustCompile(regex)
	matches := re.FindStringSubmatch(message)
	if len(matches) > 1 {
		return matches[1], true
	}
	return "", false
}
