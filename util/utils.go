package util

import (
	"crypto/rand"
	"math/big"
)

func DeduplicateSlice[T comparable](input []T) []T {
	unique := make([]T, 0)
	occurrenceMap := make(map[T]struct{})

	for _, val := range input {
		if _, ok := occurrenceMap[val]; !ok {
			occurrenceMap[val] = struct{}{}
			unique = append(unique, val)
		}
	}

	return unique
}

func GenerateRandomNumericCode(length int) (string, error) {
	var otpCharset = "0123456789"
	otp := ""
	for i := 0; i < length; i++ {
		num, err := rand.Int(rand.Reader, big.NewInt(int64(len(otpCharset))))
		if err != nil {
			return "", err
		}
		otp += string(otpCharset[num.Int64()])
	}

	return otp, nil
}
