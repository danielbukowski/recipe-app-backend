package argon

import (
	"github.com/alexedwards/argon2id"
)

type argonPasswordHasher struct {
	params argon2id.Params
}

func New(params argon2id.Params) *argonPasswordHasher {
	return &argonPasswordHasher{
		params: params,
	}
}

func (h *argonPasswordHasher) CreateHashFromPassword(password string) (string, error) {
	return argon2id.CreateHash(password, &h.params)
}
