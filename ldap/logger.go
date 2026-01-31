package ldap

import (
	"go.uber.org/zap"
)

type ServerLogger struct {
	logger *zap.SugaredLogger
}

func (s ServerLogger) Fatal(v ...interface{}) {
	s.logger.Fatal(v)
}
func (s ServerLogger) Fatalf(format string, v ...interface{}) {
	s.logger.Fatalf(format, v)
}
func (s ServerLogger) Fatalln(v ...interface{}) {
	s.logger.Fatalln(v)
}

func (s ServerLogger) Panic(v ...interface{}) {
	s.logger.Panic(v)
}
func (s ServerLogger) Panicf(format string, v ...interface{}) {
	s.logger.Panicf(format, v)
}
func (s ServerLogger) Panicln(v ...interface{}) {
	s.logger.Panicln(v)
}

func (s ServerLogger) Print(v ...interface{}) {
	s.logger.Debug(v)
}
func (s ServerLogger) Printf(format string, v ...interface{}) {
	s.logger.Debugf(format, v)
}
func (s ServerLogger) Println(v ...interface{}) {
	s.logger.Debugln(v)
}
