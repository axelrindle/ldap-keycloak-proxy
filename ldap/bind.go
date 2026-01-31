package ldap

import (
	"github.com/vjeantet/ldapserver"
	"go.uber.org/zap"
)

func (s Server) handleBind(w ldapserver.ResponseWriter, m *ldapserver.Message) {
	r := m.GetBindRequest()

	username := string(r.Name())
	password := r.AuthenticationSimple().String()

	if username == "foobar" {
		s.Logger.Debug("ldap bind succeeded", zap.String("username", username), zap.String("password", password))
		res := ldapserver.NewBindResponse(ldapserver.LDAPResultSuccess)
		w.Write(res)
	} else {
		res := ldapserver.NewBindResponse(ldapserver.LDAPResultInvalidCredentials)
		res.SetDiagnosticMessage("invalid credentials")
		w.Write(res)
	}

	// TODO: Call Keycloak
}
