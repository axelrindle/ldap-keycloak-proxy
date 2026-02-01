package ldap

import (
	"fmt"
	"regexp"

	"github.com/lor00x/goldap/message"
	"github.com/vjeantet/ldapserver"
	"go.uber.org/zap"
)

var uidRegex = regexp.MustCompile(`uid=(?P<uid>[[:word:]]+)\)`)

func (s Server) handleSearch(w ldapserver.ResponseWriter, m *ldapserver.Message) {
	r := m.GetSearchRequest()

	matches := uidRegex.FindStringSubmatch(r.FilterString())
	if matches == nil {
		res := ldapserver.NewSearchResultDoneResponse(ldapserver.LDAPResultOperationsError)
		w.Write(res)
		return
	}
	uid := matches[uidRegex.SubexpIndex("uid")]

	s.Logger.Debug("search request",
		zap.String("base dn", string(r.BaseObject())),
		zap.String("filter", r.FilterString()),
		zap.String("query", string(uid)),
	)

	rep, err := s.keycloak.SearchUser(string(uid), "username")
	if err != nil {
		s.Logger.Error("user query failed", zap.Error(err))
		res := ldapserver.NewSearchResultDoneResponse(ldapserver.LDAPResultOther)
		w.Write(res)
		return
	}

	if rep == nil {
		res := ldapserver.NewSearchResultDoneResponse(ldapserver.LDAPResultNoSuchObject)
		w.Write(res)
		return
	}

	e := ldapserver.NewSearchResultEntry(fmt.Sprintf("%s,cn=%s", r.BaseObject(), rep.Username))
	e.AddAttribute("cn", message.AttributeValue(rep.Username))
	e.AddAttribute("uid", message.AttributeValue(rep.Username))
	e.AddAttribute("givenName", message.AttributeValue(rep.FirstName))
	e.AddAttribute("sn", message.AttributeValue(rep.LastName))
	e.AddAttribute("displayName", message.AttributeValue(fmt.Sprintf("%s %s", rep.FirstName, rep.LastName)))
	e.AddAttribute("mail", message.AttributeValue(rep.Email))
	e.AddAttribute("description", message.AttributeValue("Keycloak User"))
	e.AddAttribute("keycloakId", message.AttributeValue(rep.Id))
	w.Write(e)

	res := ldapserver.NewSearchResultDoneResponse(ldapserver.LDAPResultSuccess)
	w.Write(res)
}
