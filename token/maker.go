package token

import "time"

// Maker : interface for creating tokens
type Maker interface {
	CreateToken(username string, role string, duration time.Duration) (string, *Payload, error) // we need the *Payload to get the token.ID
	VerifyToken(token string) (*Payload, error)
}
