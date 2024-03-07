package authenticator

import (
	"errors"
	"fmt"
	"math/rand"
	"time"

	"github.com/udvarid/task-manager-golang/communicator"
	"github.com/udvarid/task-manager-golang/model"
	"github.com/udvarid/task-manager-golang/repository/sessionRepository"
)

var sessions = make(map[string]model.SessionWithTime)

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func IsValid(id string, session string) bool {
	sessionInMap, isPresent := sessions[id]
	if !isPresent {
		sessionInDb, err := sessionRepository.FindSession(id)
		if err == nil {
			sessions[id] = sessionInDb
			sessionInMap = sessionInDb
		} else {
			return false
		}
	}
	now := time.Now()
	diff := now.Sub(sessionInMap.SessDate)
	if diff.Minutes() > 30 {
		delete(sessions, id)
		sessionRepository.DeleteSession((id))
		return false
	}
	return sessionInMap.Session == session
}

func CheckIn(id string, session string) {
	sessionInMap, isPresent := sessions[id]
	if isPresent {
		sessionInMap.IsChecked = true
		sessions[id] = sessionInMap
		sessionRepository.AddSession(id, sessionInMap)
	} else {
		sessionInDb, err := sessionRepository.FindSession(id)
		if err == nil {
			sessionInDb.IsChecked = true
			sessions[id] = sessionInDb
			sessionRepository.AddSession(id, sessionInDb)
		}
	}
}

func GiveSession(id string) (string, error) {
	if len(id) == 0 {
		return "", errors.New("empty id")
	}
	sessionGenerated := randStringBytes(50)
	sess := model.SessionWithTime{
		Session:   sessionGenerated,
		SessDate:  time.Now(),
		IsChecked: false,
	}
	sessions[id] = sess
	sessionRepository.AddSession(id, sess)
	return sessionGenerated, nil
}

func Validate(activeConfiguration *model.Configuration, id string, session string) bool {
	isValidatedInTime := false
	if activeConfiguration.Environment == "local" {
		fmt.Println("Local environment, validation process skipped")
		isValidatedInTime = true
	} else {
		linkToSend := activeConfiguration.RemoteAddress + "checkin/" + id + "/" + session
		communicator.SendMessageWithLink(activeConfiguration, id, linkToSend)

		foundChecked := make(chan string)
		timer := time.NewTimer(60 * time.Second)
		go func() {
			for {
				time.Sleep(1 * time.Second)
				isCheckedAlready := isChecked(id, session)
				if isCheckedAlready {
					foundChecked <- "one"
				}

			}
		}()
		select {
		case <-foundChecked:
			fmt.Println("Id is validated")
			isValidatedInTime = true
		case <-timer.C:
			fmt.Println("Id is not validated in time")
		}
	}
	return isValidatedInTime
}

func isChecked(id string, session string) bool {
	sessionInMap, isPresent := sessions[id]
	if !isPresent {
		sessionInDb, err := sessionRepository.FindSession(id)
		if err == nil {
			sessions[id] = sessionInDb
			sessionInMap = sessionInDb
		} else {
			return false
		}

	}
	return sessionInMap.Session == session && sessionInMap.IsChecked
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
