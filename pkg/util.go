package pkg

import (
	"fmt"
	"os"
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

func WriteToFile(filename, text string) error {
	file, err := os.Create(filename)
	if err != nil {
		fmt.Println(err)
		return err
	}
	defer file.Close()

	_, err = file.WriteString(text)
	if err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}
