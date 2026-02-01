package main

import (
	_ "embed"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/axelrindle/ldap-keycloak-proxy/config"
	"github.com/axelrindle/ldap-keycloak-proxy/ldap"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

//go:embed banner.txt
var banner string

var (
	Version        = "dev"
	CommitHash     = "unknown"
	BuildTimestamp = "unknown"
)

var (
	showVersion bool
	healthcheck bool
)

func BuildVersion() string {
	return fmt.Sprintf("%s-%s (%s)", Version, CommitHash, BuildTimestamp)
}

func makeLogger(c *config.Config) (*zap.Logger, error) {
	if c.Environment == "production" {
		conf := zap.NewProductionConfig()
		conf.Level = c.ZapLoggerLevel()
		conf.DisableCaller = true
		return conf.Build()
	} else {
		conf := zap.NewDevelopmentConfig()
		conf.Level = c.ZapLoggerLevel()
		conf.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
		return conf.Build()
	}
}

func main() {
	println(banner)

	flag.BoolVar(&showVersion, "version", false, "show program version")
	flag.BoolVar(&healthcheck, "healthcheck", false, "run a healthcheck")
	flag.Parse()
	if showVersion {
		println(BuildVersion())
		return
	}

	config := &config.Config{}
	config.Load()

	logger, err := makeLogger(config)
	if err != nil {
		log.Fatal(err)
	}
	defer logger.Sync()

	if healthcheck {
		addressParts := strings.Split(config.Server.Address, ":")
		if len(addressParts) != 2 {
			logger.Fatal("invalid server address", zap.String("address", config.Server.Address))
		}
		port := addressParts[1]

		conn, err := net.Dial("tcp", fmt.Sprintf("127.0.0.1:%s", port))
		if err != nil {
			log.Fatal("healthcheck failed", zap.Error(err))
		}
		conn.Close()

		return
	}

	logger.Info(fmt.Sprintf("ldap base dn is %s", config.Ldap.BaseDn))

	server := &ldap.Server{
		Config: config,
		Logger: logger,
	}
	err = server.Init()
	if err != nil {
		logger.Fatal("server initialization failed", zap.Error(err))
	}

	go server.Boot()

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)
	<-signalChan

	logger.Info("Shutting down …")

	server.Shutdown()
}
