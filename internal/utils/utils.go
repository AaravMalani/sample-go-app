package utils

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"

	"github.com/pkg/errors"
)

const (
	UtilsGenRandom        = "utils.GenerateRandomString"
	RandomGenerationError = "Unable to generate random bytes of length %d in %s"
)

func GenerateRandomString(length int) (string, error) {
	bytes := make([]byte, length)
	_, err := rand.Read(bytes)
	if err != nil {
		return "", errors.Wrap(err, fmt.Sprintf(RandomGenerationError, length, UtilsGenRandom))
	}
	return hex.EncodeToString(bytes), nil
}
