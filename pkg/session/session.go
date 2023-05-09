package session

import "time"

func GenerateSessionID() string {
	return ""
}

type Session struct {
	userID  int
	idLayer string
	lastUpd time.Time
}
