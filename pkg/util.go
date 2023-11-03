package pkg

import (
	"fmt"
	"os"
	"regexp"
	"strings"
)

func ExtractPhoneNumber(message string) (string, bool) {

	message = strings.ToLower(message)
	if strings.Contains(message, "код") {
		regex := `\+(\d{11})`
		re := regexp.MustCompile(regex)
		matches := re.FindStringSubmatch(message)
		if len(matches) > 1 {
			return matches[1], true
		}
	}
	return "", false
}

func ExtractNumber(message string) (string, bool) {
	message = strings.ToLower(message)
	if strings.Contains(message, "код") {
		regex := `(\+7|8)(\d{10})`
		re := regexp.MustCompile(regex)
		matches := re.FindStringSubmatch(message)
		if len(matches) > 2 {
			if matches[1] == "8" {
				// Преобразовать номер, начинающийся с 8, в формат +7
				phoneNumber := "+7" + matches[2]
				return phoneNumber, true
			}
			phoneNumber := matches[1] + matches[2]

			// Проверка, что номер содержит только цифры
			for _, char := range phoneNumber {
				if char != '+' && (char < '0' || char > '9') {
					return "", false
				}
			}

			return phoneNumber, true
		}
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
