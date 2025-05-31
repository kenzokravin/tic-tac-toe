package rooms

import "github.com/google/uuid"

type Room struct {
	ID       string
	State    string
	Capacity int
	Full     bool
}

type Player struct {
	ID   uuid.UUID
	Name string
}
