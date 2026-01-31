package config

import (
	"regexp"

	"github.com/go-playground/validator/v10"
)

var regex = regexp.MustCompile(`^(dc=[a-zA-Z0-9\-]+)(,?dc=[a-zA-Z0-9\-]+)*$`)

func validatorLdapBaseDn(fl validator.FieldLevel) bool {
	return regex.MatchString(fl.Field().String())
}
