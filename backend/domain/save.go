package domain

import (
	"time"

	"github.com/google/uuid"
)

type Save struct {
	ID        uuid.UUID
	PlayerID  uuid.UUID
	SceneID   string
	Honor     int
	Rage      int
	Karma     int
	CreatedAt time.Time
}
