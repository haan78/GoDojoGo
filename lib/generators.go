package lib

import (
	"crypto/rand"
	"fmt"
	"math/big"
)

const (
	lowerChars  = "abcdefghijklmnopqrstuvwxyz"
	upperChars  = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	numberChars = "0123456789"
	allChars    = lowerChars + upperChars + numberChars
)

func GeneratePassword(passwordLength int) (string, error) {
	password := make([]byte, passwordLength)

	// Ensure at least one uppercase
	upperIndex, err := randomInt(len(upperChars))
	if err != nil {
		return "", err
	}
	password[0] = upperChars[upperIndex]

	// Ensure at least one number
	numberIndex, err := randomInt(len(numberChars))
	if err != nil {
		return "", err
	}
	password[1] = numberChars[numberIndex]

	// Fill remaining characters randomly
	for i := 2; i < passwordLength; i++ {
		idx, err := randomInt(len(allChars))
		if err != nil {
			return "", err
		}
		password[i] = allChars[idx]
	}

	// Shuffle the password so required chars aren't predictable
	err = shuffle(password)
	if err != nil {
		return "", err
	}

	return string(password), nil
}

func randomInt(max int) (int, error) {
	n, err := rand.Int(rand.Reader, big.NewInt(int64(max)))
	if err != nil {
		return 0, err
	}
	return int(n.Int64()), nil
}

func shuffle(data []byte) error {
	for i := len(data) - 1; i > 0; i-- {
		j, err := randomInt(i + 1)
		if err != nil {
			return err
		}
		data[i], data[j] = data[j], data[i]
	}
	return nil
}

func GenerateRandomCode() (string, error) {
	max := big.NewInt(1000000) // 0 to 999999
	n, err := rand.Int(rand.Reader, max)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%06d", n.Int64()), nil
}
