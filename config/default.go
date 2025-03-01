package config

import "aidanwoods.dev/go-paseto"

const (
	HTTP_PORT = ":5000"
	GRPC_PORT = ":5001"
)

var (
	PrivateKey paseto.V4AsymmetricSecretKey
	PublicKey  paseto.V4AsymmetricPublicKey
)
