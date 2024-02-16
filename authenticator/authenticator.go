package authenticator

import (
	"math/rand"
	"time"
)

type SessionWithTime struct {
	session   string
	sessDate  time.Time
	isChecked bool
}

var sessions = make(map[string]SessionWithTime)

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func IsValid(id string, session string) bool {
	sessionInMap, isPresent := sessions[id]
	if !isPresent {
		return false
	}
	now := time.Now()
	diff := now.Sub(sessionInMap.sessDate)
	if diff.Minutes() > 30 {
		delete(sessions, id)
		return false
	}
	return sessionInMap.session == session
}

func IsChecked(id string, session string) bool {
	sessionInMap, isPresent := sessions[id]
	if !isPresent {
		return false
	}
	return sessionInMap.session == session && sessionInMap.isChecked
}

func CheckIn(id string, session string) {
	sessionInMap, isPresent := sessions[id]
	if isPresent {
		sessionInMap.isChecked = true
		sessions[id] = sessionInMap
	}
}

func GiveSession(id string) string {
	sessionGenerated := randStringBytes(50)
	sess := SessionWithTime{
		session:   sessionGenerated,
		sessDate:  time.Now(),
		isChecked: false,
	}
	sessions[id] = sess
	return sessionGenerated
}

func randStringBytes(n int) string {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	rng := rand.New(r)
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rng.Intn(len(letterBytes))]
	}
	return string(b)
}
