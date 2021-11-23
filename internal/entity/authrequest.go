package entity

import (
	"time"
)

type AuthRequest struct {
	Application      string
	State            string
	ExpiresAt        time.Time
	DropOffUrl       string
	PkceCodeVerifier string
}
