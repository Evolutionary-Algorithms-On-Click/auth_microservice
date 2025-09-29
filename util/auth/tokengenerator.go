package auth

import (
	"encoding/json"
	"evolve/config"
	"evolve/util"
	"time"

	"aidanwoods.dev/go-paseto"
)

func Token(payload map[string]string) (string, error) {

	logger := util.SharedLogger
	userJson, err := json.Marshal(payload)
	if err != nil {
		logger.Error("failed to marshal user payload", err)
		return "", err
	}

	token := paseto.NewToken()
	token.SetIssuedAt(time.Now())
	token.SetNotBefore(time.Now())
	token.SetExpiration(time.Now().Add(48 * time.Hour))
	token.SetString("user", string(userJson))

	signed := token.V4Sign(config.PrivateKey, nil)
	return signed, nil
}
