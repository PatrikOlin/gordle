package models

import "github.com/gofrs/uuid"

type Session struct {
	ID     uuid.UUID
	Word   string
	Status string
}
