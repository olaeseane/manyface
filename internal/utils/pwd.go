package utils

import "golang.org/x/crypto/argon2"

func HashIt(plain, salt string) []byte {
	hashed := argon2.IDKey([]byte(plain), []byte(salt), 1, 64*1024, 4, 32)
	return append([]byte(salt), hashed...)
}
