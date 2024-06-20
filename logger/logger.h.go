package logger

import (
	"github.com/sirupsen/logrus"
	prefixed "github.com/x-cray/logrus-prefixed-formatter"
)

type LogConfig struct {
	Name  string            `yaml:"name"`
	Level string            `yaml:"level"`
	Files map[string]string `yaml:"files"`
	Trace bool              `yaml:"trace"`
}

var (
	stdFormatter = &prefixed.TextFormatter{
		FullTimestamp:   true,
		TimestampFormat: "2006-01-02 15:04:05",
		ForceFormatting: true,
		ForceColors:     true,
		DisableColors:   false,
	}
	fileFormatter = &prefixed.TextFormatter{
		FullTimestamp:   true,
		TimestampFormat: "2006-01-02 15:04:05",
		ForceFormatting: true,
		ForceColors:     false,
		DisableColors:   true,
	}
	logLevels = map[string]logrus.Level{
		"trace":   logrus.TraceLevel,
		"debug":   logrus.DebugLevel,
		"info":    logrus.InfoLevel,
		"warn":    logrus.WarnLevel,
		"warning": logrus.WarnLevel,
		"error":   logrus.ErrorLevel,
		"fatal":   logrus.FatalLevel,
		"panic":   logrus.PanicLevel,
	}
)
