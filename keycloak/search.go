package keycloak

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"go.uber.org/zap"
)

type keycloakUserRepresentation struct {
	Id            string              `json:"id"`
	Username      string              `json:"username"`
	FirstName     string              `json:"firstName"`
	LastName      string              `json:"lastName"`
	Email         string              `json:"email"`
	EmailVerified bool                `json:"emailVerified"`
	Attributes    map[string][]string `json:"attributes"`
}

func (k Keycloak) SearchUser(query string, attribute string) (*keycloakUserRepresentation, error) {
	clientToken, err := k.getClientToken()
	if err != nil {
		return nil, err
	}

	k.Logger.Debug("performing search request", zap.String("query", query))

	params := map[string]string{
		"max":                 "10",
		"briefRepresentation": "true",
		"emailVerified":       "true",
		"enabled":             "true",
		"exact":               "true",
		attribute:             query,
	}

	url := fmt.Sprintf("%s/admin/realms/%s/users?%s",
		k.Config.BaseUrl,
		k.Config.Realm,
		urlEncode(params))

	k.Logger.Debug("performing user query", zap.String("url", url))

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header = map[string][]string{
		"Authorization": {fmt.Sprintf("Bearer %s", clientToken)},
	}

	res, err := k.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("invalid Keycloak response: HTTP %s", res.Status)
	}

	buf := &bytes.Buffer{}
	_, err = io.Copy(buf, res.Body)
	if err != nil {
		return nil, err
	}

	k.Logger.Debug("user query response", zap.String("response", buf.String()))

	var result []keycloakUserRepresentation
	err = json.Unmarshal(buf.Bytes(), &result)
	if err != nil {
		return nil, err
	}

	if len(result) != 1 {
		return nil, nil
	}

	rep := result[0]

	return &rep, nil
}
