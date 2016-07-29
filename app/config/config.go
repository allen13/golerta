package config

import (
  "log"
  "github.com/BurntSushi/toml"
  "github.com/allen13/golerta/app/auth/ldap"
)

type GolertaConfig struct {
  Golerta golerta
  Ldap ldap.LDAPAuthProvider
}

type golerta struct {
  SigningKey string `toml:"signing_key"`
  AuthProvider string `toml:"auth_provider"`
}

func BuildConfig(configFile string)(config GolertaConfig){
  _, err := toml.DecodeFile(configFile, &config)

  if err != nil {
    log.Fatal("config file error: " + err.Error())
  }

  setDefaultConfigs(&config)
  return
}

func setDefaultConfigs(config* GolertaConfig){
  if config.Golerta.AuthProvider == ""{
    config.Golerta.AuthProvider = "ldap"
  }
}