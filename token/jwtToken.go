package token

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt"
)

const minJWTSecretSize = 32

type JWTToken struct {
	secretKey string
}

func NewJWTMaker(secretKey string) (Maker, error) {
	if len(secretKey) < minJWTSecretSize {
		return nil, fmt.Errorf("invalid Key : too short : %v", len(secretKey))
	}
	return &JWTToken{secretKey}, nil
}

func (maker *JWTToken) CreateToken(username, role string, duration time.Duration) (string, *Payload, error) {
	payload := &Payload{}

	payload, err := NewPayload(username, role, duration)
	if err != nil {
		return "", payload, err
	}

	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)
	token, err := jwtToken.SignedString([]byte(maker.secretKey))
	return token, payload, err
}

func (maker *JWTToken) VerifyToken(token string) (*Payload, error) {

	// To overcome the trivial forgery vulnerabilty
	keyFunc := func(token *jwt.Token) (interface{}, error) {
		_, ok := token.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, ErrInvalidToken
		}
		return []byte(maker.secretKey), nil
	}

	/* Could be 2 different errors : invalid or expired
	To differentiate the differences you need to look into the implementation of the library

	// The error from Parse if token is not valid
		type ValidationError struct {
			Inner  error  // stores the error returned by external dependencies, i.e.: KeyFunc
			Errors uint32 // bitfield.  see ValidationError... constants
			text   string // errors that do not have a valid error just have text
		}
	*/

	jwtToken, err := jwt.ParseWithClaims(token, &Payload{}, keyFunc)
	if err != nil {
		verr, ok := err.(*jwt.ValidationError)
		if ok && errors.Is(verr.Inner, ErrExpiredToken) {
			return nil, ErrExpiredToken
		}
		return nil, ErrInvalidToken
	}

	payload, ok := jwtToken.Claims.(*Payload)
	if !ok {
		return nil, ErrInvalidToken
	}

	return payload, nil
}
