package entity

import (
	"time"
)

type AuthRequest struct {
	State            string
	ExpiresAt        time.Time
	DropOffUrl       string
	PkceCodeVerifier string
}
