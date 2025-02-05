package passwordhasher

import (
	"github.com/alexedwards/argon2id"
)

type argonPasswordHasher struct {
	params *argon2id.Params
}

func New(params *argon2id.Params) *argonPasswordHasher {
	return &argonPasswordHasher{
		params: params,
	}
}

func (h *argonPasswordHasher) CreateHashFromPassword(password string) (string, error) {
	return argon2id.CreateHash(password, h.params)
}

func (h *argonPasswordHasher) ComparePasswordAndHash(password, hash string) bool {
	ok, _ := argon2id.ComparePasswordAndHash(password, hash)
	return ok
}
