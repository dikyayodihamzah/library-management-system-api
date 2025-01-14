package lib

import "golang.org/x/crypto/bcrypt"

// PasswordEncrypt Password Encrypt
func PasswordEncrypt(plain string) string {
	if hashed, err := bcrypt.GenerateFromPassword([]byte(plain), bcrypt.DefaultCost); nil == err {
		return string(hashed)
	}

	return ""
}

// PasswordCompare Password Compare
func PasswordCompare(encrypted, password string) bool {
	return nil == bcrypt.CompareHashAndPassword([]byte(encrypted), []byte(password))
}
