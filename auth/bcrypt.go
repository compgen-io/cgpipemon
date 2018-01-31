package auth


import (
	"golang.org/x/crypto/bcrypt"
)

func checkBCrypt(plain string, enc string) bool {
    return bcrypt.CompareHashAndPassword([]byte(enc), []byte(plain)) == nil
}

func encryptBCrypt(plain string) (string, error) {
	enc, err := bcrypt.GenerateFromPassword([]byte(plain), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	return string(enc), nil
}