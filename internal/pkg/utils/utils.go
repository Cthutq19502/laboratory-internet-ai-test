package utils

import (
	"encoding/json"
	"github.com/google/uuid"
	"regexp"
)

func GenerateUUID() string {
	id := uuid.New()
	return id.String()
}

func ParseJson[T any](response []byte) (*T, error) {

	var j T
	err := json.Unmarshal(response, &j)
	if err != nil {
		return nil, err
	}

	return &j, nil
}

func IsValidEmail(email string) bool {
	if email == "" {
		return false
	}

	const emailRegex = `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	re := regexp.MustCompile(emailRegex)
	return re.MatchString(email)
}
