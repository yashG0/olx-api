package utils

import (
	"math/rand/v2"
)

const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func GetRequestId() string {
	result := make([]byte, 10)
	for i := range result {
		result[i] = charset[rand.IntN(len(charset))]
	}

	return string(result)
}
