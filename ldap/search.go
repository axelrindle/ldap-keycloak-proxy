package ldap

import (
	"fmt"

	"github.com/lor00x/goldap/message"
	"github.com/vjeantet/ldapserver"
	"go.uber.org/zap"
)

func (s Server) handleSearch(w ldapserver.ResponseWriter, m *ldapserver.Message) {
	r := m.GetSearchRequest()

	queries := r.Attributes()

	if len(queries) != 1 {
		res := ldapserver.NewSearchResultDoneResponse(ldapserver.LDAPResultSizeLimitExceeded)
		w.Write(res)
	}

	for _, query := range queries {
		s.Logger.Debug("search request",
			zap.String("base dn", string(r.BaseObject())),
			zap.String("filter", r.FilterString()),
			zap.String("query", string(query)),
		)

		rep, err := s.keycloak.SearchUser(string(query), "username")
		if err != nil {
			s.Logger.Error("user query failed", zap.Error(err))
			break
		}

		if rep == nil {
			continue
		}

		e := ldapserver.NewSearchResultEntry(fmt.Sprintf("%s,cn=%s", r.BaseObject(), rep.Username))
		e.AddAttribute("cn", message.AttributeValue(rep.Username))
		e.AddAttribute("uid", message.AttributeValue(rep.Username))
		e.AddAttribute("givenName", message.AttributeValue(rep.FirstName))
		e.AddAttribute("sn", message.AttributeValue(rep.LastName))
		e.AddAttribute("mail", message.AttributeValue(rep.Email))
		e.AddAttribute("description", message.AttributeValue("Keycloak User"))
		e.AddAttribute("keycloakId", message.AttributeValue(rep.Id))
		w.Write(e)
	}

	res := ldapserver.NewSearchResultDoneResponse(ldapserver.LDAPResultSuccess)
	w.Write(res)
}
