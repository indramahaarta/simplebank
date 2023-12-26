package token

import (
	"fmt"
	"time"

	paseto "github.com/o1egl/paseto"
)

type PasetoMaker struct {
	paseto      *paseto.V2
	symetricKey []byte
}

// CreateToken implements Maker.
func (maker *PasetoMaker) CreateToken(username string, duration time.Duration) (string, error) {
	payload, err := NewPayload(username, duration)
	if err != nil {
		return "", err
	}

	return maker.paseto.Encrypt(maker.symetricKey, payload, nil)
}

// VerifyToken implements Maker.
func (maker *PasetoMaker) VerifyToken(token string) (*Payload, error) {
	payload := &Payload{}
	err := maker.paseto.Decrypt(token, maker.symetricKey, payload, nil)
	if err != nil {
		return nil, err
	}

	err = payload.Valid()
	if err != nil {
		return nil, err
	}

	return payload, nil
}

func NewPasetoMaker(symetricKey string) (Maker, error) {
	if len(symetricKey) < minSecretKeyzise {
		return nil, fmt.Errorf("token size must be at least %d character", minSecretKeyzise)
	}

	return &PasetoMaker{paseto: paseto.NewV2(), symetricKey: []byte(symetricKey)}, nil
}
