package ldap

import (
	"time"

	"github.com/axelrindle/ldap-keycloak-proxy/config"
	"github.com/axelrindle/ldap-keycloak-proxy/keycloak"
	"github.com/vjeantet/ldapserver"
	"go.uber.org/zap"
)

type Server struct {
	Config *config.Config
	Logger *zap.Logger

	keycloak keycloak.Keycloak

	shouldRun bool
}

func (s *Server) Init() error {
	s.keycloak = keycloak.Keycloak{
		Config: s.Config.Keycloak,
		Logger: s.Logger,
	}

	err := s.keycloak.Init()
	if err != nil {
		return err
	}

	s.shouldRun = true

	return nil
}

func (s *Server) Boot() {
	if s.Config.Logging.TraceServer {
		ldapserver.Logger = ServerLogger{logger: s.Logger.Sugar()}
	} else {
		ldapserver.Logger = ldapserver.DiscardingLogger
	}

	mux := ldapserver.NewRouteMux()
	mux.Bind(s.handleBind)
	// mux.Extended(func(w ldapserver.ResponseWriter, m *ldapserver.Message) {
	// 	r := m.GetExtendedRequest()

	// 	res := ldapserver.NewExtendedResponse(ldapserver.LDAPResultSuccess)
	// 	res.SetResponseName(ldapserver.NoticeOfWhoAmI)
	// 	res.SeMatchedDN()
	// 	w.Write(res)
	// }).RequestName(ldapserver.NoticeOfWhoAmI)

	mux.Search(s.handleSearch).
		BaseDn(s.Config.Ldap.BaseDn)

	server := &ldapserver.Server{
		Handler:      mux,
		ReadTimeout:  time.Second * 5,
		WriteTimeout: time.Second * 5,
	}

	go server.ListenAndServe(s.Config.Server.Address)

	s.Logger.Info("LDAP server listening", zap.String("address", s.Config.Server.Address))
}

func (s *Server) Shutdown() {
	s.shouldRun = false
}
