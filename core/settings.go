package core

import "github.com/pelletier/go-toml/v2"

type SettingsSecrets struct {
	PasswordSecret string `toml:"password_secret"`
}

var Secrets SettingsSecrets

func LoadSecrets(secretsToml []byte) {
	err := toml.Unmarshal(secretsToml, &Secrets)
	if err != nil {
		panic(err)
	}
}
