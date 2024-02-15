package authenticator

import (
	"math/rand"
	"time"
)

var sessions = make(map[string]string)

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func IsValid(id string, session string) bool {
	sessionInMap, isPresent := sessions[id]
	return isPresent && sessionInMap == session
}

func GiveSession(id string) string {
	session := randStringBytes(50)
	sessions[id] = session
	return session
}

func randStringBytes(n int) string {
	rand.Seed(time.Now().UnixNano())
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}

// https://stackoverflow.com/questions/22892120/how-to-generate-a-random-string-of-a-fixed-length-in-go
