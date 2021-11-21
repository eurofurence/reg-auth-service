package entity

import (
	"time"
)

type AuthRequest struct {
	State            string
	ExpiresAt        time.Time
	DropoffUrl       string
	PkceCodeVerifier string
}
