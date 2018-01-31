package auth

import (
	"log"
	"time"
	"errors"
	"math/rand"
)

func init() {
    rand.Seed(time.Now().UnixNano())
}

func CheckPass(plain string, enc string) bool {
	if enc[:2] == "1$" {
		return checkBCrypt(plain, enc[2:])
	}
	return false
}

const letters = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func GenerateRandom(length int) []byte {
	b := make([]byte, length)

	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}

	return b
}

func EncryptPass(method int, plain string) (string, error) {
	if (method == 1) {
		enc, err := encryptBCrypt(plain)
		if err != nil {
			log.Panic(err)
		}
		return "1$"+string(enc), nil
	}
	return "", errors.New("Unsupported method")
}
