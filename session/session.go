package session

import (
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"

	"github.com/PatrikOlin/gordle/db"
	g "github.com/PatrikOlin/gordle/guess"
	"github.com/PatrikOlin/gordle/rules"
	st "github.com/PatrikOlin/gordle/stats"
	"github.com/PatrikOlin/gordle/word"
)

type Session struct {
	ID           uuid.UUID `json:"id" db:"id"`
	Word         string    `json:"word,omitempty" db:"word"`
	Status       string    `json:"status" db:"status"`
	Guesses      []g.Guess `json:"guesses" db:"-"`
	NumOfGuesses int       `json:"numberOfGuesses" db:"number_of_guesses"`
	CreatedAt    int       `json:"createdAt" db:"created_at"`
}

type FinishedSession struct {
	Session
	Stats st.Stats `json:"stats"`
}

func Create(userToken string) Session {
	s := Session{
		ID:           uuid.New(),
		Word:         word.New(),
		Status:       "unsolved",
		Guesses:      make([]g.Guess, 0, 6),
		NumOfGuesses: 0,
		CreatedAt:    int(time.Now().Unix()),
	}

	persistSession(s, userToken)

	return s
}

func Get(userToken string) (Session, error) {
	var session Session
	stmt := `
		SELECT s.id, status, word, number_of_guesses, created_at FROM sessions s
		JOIN user_game_sessions usg on s.id = usg.game_id
		JOIN user_sessions us on us.token = usg.user_token
		WHERE us.token = $1
		ORDER BY s.created_at DESC LIMIT 1`

	// stmt := "SELECT * FROM sessions s INNER JOIN user_game_sessions ugs ON ugs.token = WHERE id = $1"
	err := db.DBClient.Get(&session, stmt, userToken)

	if err != nil {
		fmt.Println(err)
		return session, err
	}

	fmt.Println("session ", session.ID, session.Word)

	session.GetGuesses()
	return session, nil
}

func persistSession(s Session, userToken string) {
	stmt := "INSERT INTO sessions (id, word, status, number_of_guesses, created_at) VALUES ($1, $2, $3, $4, $5)"

	tx, err := db.DBClient.Begin()
	if err != nil {
		log.Fatalln(err)
	}

	_, err = tx.Exec(stmt, s.ID, s.Word, s.Status, s.NumOfGuesses, s.CreatedAt)
	if err != nil {
		log.Fatalln(err)
	}

	stmt2 := "INSERT INTO user_game_sessions (user_token, game_id) VALUES ($1, $2)"

	_, err = tx.Exec(stmt2, userToken, s.ID)
	if err != nil {
		log.Fatalln(err)
	}

	err = tx.Commit()
	if err != nil {
		log.Fatalln(err)
	}
}

func (s *Session) Update() {
	stmt := "UPDATE sessions SET status=:status, number_of_guesses=:number_of_guesses WHERE id=:id"

	_, err := db.DBClient.NamedExec(stmt, s)
	if err != nil {
		log.Fatalln(err)
	}

}

func (s *Session) GetGuesses() {
	stmt := "SELECT word, word_state FROM guesses WHERE session_id = ?"
	db.DBClient.Select(&s.Guesses, stmt, s.ID)
}

func (s *Session) SetWordVisibility() {
	if s.Status == "unsolved" && s.NumOfGuesses != rules.Get().MaxGuesses {
		s.Word = ""
	}
}

func (s *Session) IsAlive() bool {
	if s.NumOfGuesses >= rules.Get().MaxGuesses {
		return false
	}
	if s.Status == "solved" {
		return false
	}

	return true
}

func (s Session) GetStats(userToken string) FinishedSession {
	sessions := []Session{}
	stmt := `
		SELECT s.id, status, word, number_of_guesses, created_at FROM sessions s
		JOIN user_game_sessions usg on s.id = usg.game_id
		JOIN user_sessions us on us.token = usg.user_token
		WHERE us.token = $1
		ORDER BY s.created_at ASC`

	err := db.DBClient.Select(&sessions, stmt, userToken)
	if err != nil {
		log.Fatalln(err)
	}

	stats := st.Stats{}
	ms := 0
	cs := 0

	for _, session := range sessions {
		if session.Status == "solved" {
			cs += 1
			if cs > ms {
				ms = cs
			}
			stats.WinDistribution.Add(session.NumOfGuesses)
		} else {
			if cs > ms {
				ms = cs
			}
			cs = 0

			stats.WinDistribution.Unsolved += 1
		}
	}

	stats.CurrentStreak = cs
	stats.MaxStreak = ms
	stats.CalculateWinPercentage()
	fmt.Println(stats.WinPercentage)

	fs := FinishedSession{
		Session: s,
		Stats:   stats,
	}

	return fs
}
