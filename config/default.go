package config

import "aidanwoods.dev/go-paseto"

const (
	PORT = ":5000"
)

var (
	PrivateKey paseto.V4AsymmetricSecretKey
	PublicKey  paseto.V4AsymmetricPublicKey
)
