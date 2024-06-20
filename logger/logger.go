package logger

import (
	"fmt"
	"github.com/rifflock/lfshook"
	"github.com/sirupsen/logrus"
	"os"
	"path/filepath"
	"strings"
)

func New(config LogConfig) *logrus.Logger {
	logger := logrus.New()
	logLevel := logLevels[strings.ToLower(config.Level)]
	logFiles := make(map[string]string)
	logPaths := make(map[logrus.Level]string)

	for levelName, filename := range config.Files {
		logFiles[strings.ToLower(levelName)] = filename
	}
	for levelName, filename := range logFiles {
		if level, ok := logLevels[levelName]; ok {
			logPaths[level] = filename
			if err := os.MkdirAll(filepath.Dir(filename), os.ModePerm); err != nil {
				logrus.WithFields(
					logrus.Fields{"fileName": filename, "logLevel": levelName, "error": err.Error()},
				).Warning("can not open log file, log will not be saved")
				continue
			}
			if _, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY|os.O_APPEND, os.ModePerm); err != nil {
				logrus.WithFields(
					logrus.Fields{"fileName": filename, "logLevel": levelName, "error": err.Error()},
				).Warning("can not open log file, log will not be saved")
				continue
			}
		} else {
			logrus.WithFields(
				logrus.Fields{"fileName": filename, "logLevel": levelName, "error": "unknown log level"},
			).Warning("can not parse log level, log will not be saved")
		}
	}
	for _, level := range logrus.AllLevels {
		if _, ok := logPaths[level]; !ok {
			if len(config.Name) > 0 {
				logPaths[level] = fmt.Sprintf("%s.log", config.Name)
			} else {
				logPaths[level] = fmt.Sprintf("%s.log", level.String())
			}
		}
	}

	logger.SetLevel(logLevel)
	logger.SetOutput(os.Stdout)
	logger.SetFormatter(stdFormatter)
	logger.SetReportCaller(config.Trace)
	logger.AddHook(lfshook.NewHook(lfshook.PathMap(logPaths), fileFormatter))
	return logger
}
