package userSession

import (
	"log"

	"github.com/PatrikOlin/gordle/db"
	"github.com/google/uuid"
)

type UserSession struct {
	Token uuid.UUID `json:"user_session_token" db:"token"`
}

func Create() UserSession {
	us := UserSession{
		Token: uuid.New(),
	}

	persistUserSession(us)
	return us
}

func persistUserSession(us UserSession) {
	stmt := "INSERT INTO user_sessions (token) VALUES (?)"

	_, err := db.DBClient.Exec(stmt, us.Token)
	if err != nil {
		log.Fatalln(err)
	}
}
