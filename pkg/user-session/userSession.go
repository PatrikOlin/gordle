package usersession

import (
	"log"

	"github.com/PatrikOlin/gordle/pkg/db"
	"github.com/google/uuid"
)

type UserSession struct {
	Token         uuid.UUID `json:"user_session_token" db:"token"`
	FinishedDaily bool      `json:"finished_daily" db:"finished_daily"`
}

func Create() UserSession {
	us := UserSession{
		Token:         uuid.New(),
		FinishedDaily: false,
	}

	persistUserSession(us)
	return us
}

func GetUserSession(userToken string) (UserSession, error) {
	var us UserSession
	stmt := "SELECT token, finished_daily FROM user_session WHERE token = $1"

	err := db.DBClient.Get(&us, stmt, userToken)
	if err != nil {
		return us, err
	}

	return us, nil
}

func (us *UserSession) SetDailyFinished() {
	stmt := "UPDATE user_session us SET finished_daily = TRUE WHERE token = $1"

	_, err := db.DBClient.Exec(stmt, us.Token)
	if err != nil {
		log.Fatalln(err)
	}
}

func persistUserSession(us UserSession) {
	stmt := "INSERT INTO user_session (token, finished_daily) VALUES ($1, $2)"

	_, err := db.DBClient.Exec(stmt, us.Token, us.FinishedDaily)
	if err != nil {
		log.Fatalln(err)
	}
}
