package controllers

import (
	"net/http"
	"sync"
	"time"

	uuid "github.com/satori/go.uuid"
)

var sessions sync.Map

type session struct {
	username string
	expiry   time.Time
}

func (s session) isExpired() bool {
	return s.expiry.Before(time.Now())
}

//validSession returns username and true if the session is valid
func validSession(r *http.Request) (string, bool) {
	c, err := r.Cookie("session_token")
	if err == nil {
		if val, ok := sessions.Load(c.Value); ok {
			return val.(session).username, ok
		}
	}
	return "", false
}

func NewSessionToken(w http.ResponseWriter, username string) {
	sessionToken := uuid.NewV4()
	deleteSameUser(username)
	expiresAt := time.Now().Add(2 * time.Hour) //let's add Refresh
	sessions.Store(sessionToken.String(), session{username, expiresAt})
	http.SetCookie(w, &http.Cookie{
		Name:     "session_token",
		Value:    sessionToken.String(),
		Path:     "/",  //?
		HttpOnly: true, //?
		Expires:  expiresAt,
	})
}

func deleteSameUser(username string) {
	sessions.Range(func(key, value interface{}) bool {
		if username == value.(session).username {
			sessions.Delete(key)
		}
		return true
	})
}

func DeleteExpiredSessions() {
	for {
		sessions.Range(func(key, value interface{}) bool {
			if value.(session).isExpired() {
				sessions.Delete(key)
			}
			return true
		})
		time.Sleep(5 * time.Second)
	}
}
