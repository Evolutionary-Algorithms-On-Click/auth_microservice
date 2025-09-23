package auth

import (
	"bytes"
	"encoding/json"
	"evolve/config"
	"evolve/util"
	"strings"

	"aidanwoods.dev/go-paseto"
)

func ValidateToken(token string) (map[string]string, error) {
	
	logger := util.Log_var
	// logger.Info(token)

	parser := paseto.NewParserForValidNow()
	payLoad, err := parser.ParseV4Public(config.PublicKey, token, nil)
	if err != nil {
		logger.Error("failed to parse token")
		return nil, err
	}

	payLoadJSON := json.NewDecoder(bytes.NewReader(payLoad.ClaimsJSON()))
	payLoadMap := make(map[string]string)
	if err = payLoadJSON.Decode(&payLoadMap); err != nil {
		logger.Error("failed to decode token payload")
		return nil, err
	}

	userJson := json.NewDecoder(strings.NewReader(payLoadMap["user"]))
	userMap := make(map[string]string)
	if err = userJson.Decode(&userMap); err != nil {
		logger.Error("failed to decode user json")
		return nil, err
	}

	// logger.Info(fmt.Sprintf("Token payload: %v", payLoadMap))
	return userMap, nil
}
