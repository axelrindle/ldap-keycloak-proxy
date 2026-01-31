package keycloak

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/axelrindle/ldap-keycloak-proxy/config"
	"go.uber.org/zap"
)

type Keycloak struct {
	Config config.KeycloakConfig
	Logger *zap.Logger

	client http.Client
}

func (k *Keycloak) Init() error {
	k.client = http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: k.Config.SkipCertificateValidation,
			},
		},
	}

	return k.test()
}

type keycloakRealmResponse struct {
	Realm string `json:"realm"`
}

func (k Keycloak) test() error {
	url := fmt.Sprintf("%s/realms/%s", k.Config.BaseUrl, k.Config.Realm)
	res, err := k.client.Get(url)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("invalid status received from Keycloak: %d", res.StatusCode)
	}

	buf := &bytes.Buffer{}
	_, err = io.Copy(buf, res.Body)
	if err != nil {
		return err
	}

	payload := &keycloakRealmResponse{}
	err = json.Unmarshal(buf.Bytes(), &payload)
	if err != nil {
		return err
	}

	if payload.Realm != k.Config.Realm {
		return fmt.Errorf("returned realm %s does not match expected realm %s", payload.Realm, k.Config.Realm)
	}

	// initially verify client details
	_, err = k.getClientToken()
	if err != nil {
		return err
	}

	k.Logger.Info("verified Keycloak connection", zap.String("url", k.Config.BaseUrl), zap.String("realm", k.Config.Realm))

	return nil
}
