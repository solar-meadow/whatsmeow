package pkg

import (
	"fmt"
	"os"
)

func CheckEnvExist(key string) (string, error) {
	if value := os.Getenv(key); value != "" {
		return value, nil
	}
	return "", fmt.Errorf("value for key %s not found", key)
}
