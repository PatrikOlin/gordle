package session

import (
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"

	dw "github.com/PatrikOlin/gordle/pkg/daily"
	"github.com/PatrikOlin/gordle/pkg/db"
	g "github.com/PatrikOlin/gordle/pkg/guess"
	"github.com/PatrikOlin/gordle/pkg/rules"
	st "github.com/PatrikOlin/gordle/pkg/stats"
	us "github.com/PatrikOlin/gordle/pkg/user-session"
	w "github.com/PatrikOlin/gordle/pkg/word"
)

type Session struct {
	ID           uuid.UUID `json:"id" db:"session_id"`
	Word         string    `json:"word,omitempty" db:"word"`
	Status       string    `json:"status" db:"status"`
	Guesses      []g.Guess `json:"guesses" db:"-"`
	NumOfGuesses int       `json:"numberOfGuesses" db:"number_of_guesses"`
	CreatedAt    int       `json:"createdAt" db:"created_at"`
	FinishedAt   int       `json:"finishedAt,omitempty" db:"finished_at"`
	IsDaily      bool      `json:"isDaily" db:"daily"`
}

type FinishedSession struct {
	Session
	Stats st.Stats `json:"stats"`
}

func Create(userSession us.UserSession) Session {
	var word string
	var isDaily bool
	if userSession.FinishedDaily {
		word = w.New()
		isDaily = false
	} else {
		word = dw.GetDailyWord()
		isDaily = true
	}

	s := Session{
		ID:           uuid.New(),
		Word:         word,
		Status:       "unsolved",
		Guesses:      make([]g.Guess, 0, 6),
		NumOfGuesses: 0,
		CreatedAt:    int(time.Now().Unix()),
		FinishedAt:   0,
		IsDaily:      isDaily,
	}

	persistSession(s, userSession.Token.String())

	return s
}

func Get(userSession us.UserSession) (Session, error) {
	var session Session
	stmt := `
		SELECT session_id, status, s.word, number_of_guesses, created_at, finished_at FROM game_session s
		JOIN user_game_session usg on s.session_id = usg.game_id
		JOIN user_session us on us.token = usg.user_token
		WHERE us.token = $1
		ORDER BY s.created_at DESC LIMIT 1`

	// stmt := "SELECT * FROM sessions s INNER JOIN user_game_sessions ugs ON ugs.token = WHERE id = $1"
	err := db.DBClient.Get(&session, stmt, userSession.Token)

	if err != nil {
		fmt.Println("Session.Get", err)
		return session, err
	}

	fmt.Println("session ", session.ID, session.Word)
	session.GetGuesses()
	return session, nil
}

func persistSession(s Session, userToken string) {
	stmt := "INSERT INTO game_session (session_id, word, status, number_of_guesses, created_at, finished_at, daily) VALUES ($1, $2, $3, $4, $5, $6, $7)"

	tx, err := db.DBClient.Begin()
	if err != nil {
		log.Fatalln(err)
	}

	_, err = tx.Exec(stmt, s.ID, s.Word, s.Status, s.NumOfGuesses, s.CreatedAt, s.FinishedAt, s.IsDaily)
	if err != nil {
		log.Fatalln(err)
	}

	stmt2 := "INSERT INTO user_game_session (user_token, game_id) VALUES ($1, $2)"

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
	stmt := "UPDATE game_session SET status=:status, number_of_guesses=:number_of_guesses WHERE session_id=:session_id"

	_, err := db.DBClient.NamedExec(stmt, s)
	if err != nil {
		log.Fatalln(err)
	}

}

func (s *Session) GetGuesses() {
	stmt := "SELECT word, word_state FROM guess WHERE session_id = $1"
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
		SELECT session_id, status, word, number_of_guesses, created_at FROM game_session s
		JOIN user_game_session usg on s.session_id = usg.game_id
		JOIN user_session us on us.token = usg.user_token
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

	fs := FinishedSession{
		Session: s,
		Stats:   stats,
	}

	return fs
}
