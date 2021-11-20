package common

import (
	"math/rand"
	"net/http"
	"time"

	"go.uber.org/zap"
	"golang.org/x/crypto/argon2"
)

// type CustomError struct {
// 	Error   error
// 	Message string
// 	Code    int
// }

func RandStringRunes(n int) string {
	letterRunes := []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	rand.Seed(time.Now().UnixNano())

	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

func HashIt(plain, salt string) []byte {
	hashed := argon2.IDKey([]byte(plain), []byte(salt), 1, 64*1024, 4, 32)
	return append([]byte(salt), hashed...)
}

func HandleError(w http.ResponseWriter, err error, status int, msg string, logger *zap.SugaredLogger) {
	var errMsg string
	if err != nil {
		errMsg = msg + " - " + err.Error()
	} else {
		errMsg = msg
	}
	http.Error(w, errMsg, status)
	logger.Error(errMsg)
}
