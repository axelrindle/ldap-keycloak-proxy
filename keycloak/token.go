package keycloak

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
)

type keycloakTokenResponse struct {
	Token     string `json:"access_token"`
	ExpiresIn int64  `json:"expires_in"`

	Error            string `json:"error"`
	ErrorDescription string `json:"error_description"`
}

func (k Keycloak) getClientToken() (string, error) {
	k.Logger.Debug("retrieving new access token")

	params := map[string]string{
		"grant_type":    "client_credentials",
		"client_id":     k.Config.ClientID,
		"client_secret": k.Config.ClientSecret,
	}

	url := fmt.Sprintf("%s/realms/%s/protocol/openid-connect/token",
		k.Config.BaseUrl,
		k.Config.Realm)

	payload := new(bytes.Buffer)
	_, err := payload.WriteString(urlEncode(params))
	if err != nil {
		return "", err
	}

	res, err := k.client.Post(url, "application/x-www-form-urlencoded", payload)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	buf := &bytes.Buffer{}
	_, err = io.Copy(buf, res.Body)
	if err != nil {
		return "", err
	}

	result := &keycloakTokenResponse{}
	err = json.Unmarshal(buf.Bytes(), &result)
	if err != nil {
		return "", err
	}

	if result.Error != "" {
		return "", errors.New(result.ErrorDescription)
	}

	// TODO: use expiry for caching
	return result.Token, nil
}
