package config

type LoggingConfig struct {
	Level       string `env:"LOG_LEVEL" default:"info" validate:"oneof=debug info warn error"`
	TraceServer bool   `env:"LOG_TRACE_SERVER" default:"false" validate:"boolean"`
}

type ServerConfig struct {
	Address string `env:"SERVER_ADDRESS" default:":1337"`
}

type KeycloakConfig struct {
	BaseUrl                   string `env:"KEYCLOAK_BASEURL,required"`
	SkipCertificateValidation bool   `env:"KEYCLOAK_SKIPCERTIFICATEVALIDATION" default:"false"`
	Realm                     string `env:"KEYCLOAK_REALM,required"`
	ClientID                  string `env:"KEYCLOAK_CLIENTID,required"`
	ClientSecret              string `env:"KEYCLOAK_CLIENTSECRET,required"`
}

type LdapConfig struct {
	BaseDn string `env:"BASE_DN,required" validate:"ldapBaseDn"`
}

type Config struct {
	Logging  LoggingConfig
	Server   ServerConfig
	Keycloak KeycloakConfig
	Ldap     LdapConfig

	Environment string `env:"ENVIRONMENT" default:"production" validate:"oneof=production development"`
}
