// Package auth collects structures and functions around the
// generation and processing of credentials.
package auth

import (
	"github.com/epinio/epinio/helpers/randstr"
	"golang.org/x/crypto/bcrypt"
)

// PasswordAuth wraps a set of password-based credentials
type PasswordAuth struct {
	Username string
	Password string
}

func (auth *PasswordAuth) Htpassword() (string, error) {
	hash, err := HashBcrypt(auth.Password)
	if err != nil {
		return "", err
	}
	return auth.Username + ":" + hash, nil
}

// HashBcrypt generates an Apache MD5 hash for a password.
// See https://github.com/foomo/htpasswd for the origin of this code.
// MIT licensed, as per `blob/master/LICENSE.txt`

func HashBcrypt(password string) (hash string, err error) {
	passwordBytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return
	}
	return string(passwordBytes), nil
}

func RandomPasswordAuth() (*PasswordAuth, error) {
	user, err := randstr.Hex16()
	if err != nil {
		return nil, err
	}

	password, err := randstr.Hex16()
	if err != nil {
		return nil, err
	}

	return &PasswordAuth{
		Username: user,
		Password: password,
	}, nil
}
