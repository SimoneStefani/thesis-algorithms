package common

import (
	"crypto/sha256"
	"encoding/base64"
	"errors"
)

func HashTransaction(tr string) string {
	h := sha256.Sum256([]byte(tr))
	return base64.StdEncoding.EncodeToString(h[:])
}

func Includes(tr string, list []string) (int, error) {

	for i, transaction := range list {
		if tr == transaction {
			return i, nil
		}
	}
	return -1, errors.New("error: not in list")
}
