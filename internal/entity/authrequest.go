package entity

import (
	"net/url"
	"time"
)

type AuthRequest struct {
	State            string
	ExpiresAt        time.Time
	DropoffUrl       url.URL
	PkceCodeVerifier string
}
