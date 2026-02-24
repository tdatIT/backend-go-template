package security

import "golang.org/x/crypto/bcrypt"

const DefaultCost = bcrypt.DefaultCost

func HashPassword(password string, cost int) (string, error) {
	if cost <= 0 {
		cost = DefaultCost
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), cost)
	if err != nil {
		return "", err
	}

	return string(hash), nil
}

func ComparePassword(hash string, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
}
