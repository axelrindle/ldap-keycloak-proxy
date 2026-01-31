package config

import (
	"log"

	"github.com/caarlos0/env/v11"
	"github.com/go-playground/validator/v10"
	_ "github.com/joho/godotenv/autoload"
)

func (c *Config) Load() {
	var err error

	err = env.ParseWithOptions(c, env.Options{
		DefaultValueTagName: "default",
		Prefix:              "LDAP_",
	})
	if err != nil {
		log.Fatal("failed to load config: ", err)
	}

	validate := validator.New(validator.WithRequiredStructEnabled())
	validate.RegisterValidation("ldapBaseDn", validatorLdapBaseDn)
	err = validate.Struct(c)
	if err != nil {
		log.Fatal("config validation failed: ", err)
	}
}
