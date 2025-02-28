package auth

import (
	"bytes"
	"encoding/json"
	"evolve/config"
	"evolve/util"

	"aidanwoods.dev/go-paseto"
)

func ValidateToken(token string) (map[string]any, error) {
	var logger = util.NewLogger()

	// logger.Info(token)

	parser := paseto.NewParserForValidNow()
	payLoad, err := parser.ParseV4Public(config.PublicKey, token, nil)
	if err != nil {
		logger.Error("failed to parse token")
		return nil, err
	}

	payLoadJSON := json.NewDecoder(bytes.NewReader(payLoad.ClaimsJSON()))
	var payLoadMap map[string]any
	if err = payLoadJSON.Decode(&payLoadMap); err != nil {
		logger.Error("failed to decode token payload")
		return nil, err
	}

	// logger.Info(fmt.Sprintf("Token payload: %v", payLoadMap))

	return payLoadMap, nil
}
