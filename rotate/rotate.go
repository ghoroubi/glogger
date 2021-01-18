package rotate

import (
	"github.com/sirupsen/logrus"
	"gopkg.in/natefinch/lumberjack.v2"
	"io"
)

// RotateFileConfig , represents the configuration of file to be rotated.
type RotateFileConfig struct {
	Filename   string
	MaxSize    int
	MaxBackups int
	MaxAge     int
	Level      logrus.Level
	Formatter  logrus.Formatter
}

// RotateFileHook , the hook which will be used by the logrus.
type RotateFileHook struct {
	Config    RotateFileConfig
	logWriter io.Writer
}

// NewRotateFileHook  , returns a new instance of logrus hook
func NewRotateFileHook(config RotateFileConfig) (logrus.Hook, error) {
	hook := RotateFileHook{
		Config: config,
	}
	hook.logWriter = &lumberjack.Logger{
		Filename:   config.Filename,
		MaxSize:    config.MaxSize,
		MaxBackups: config.MaxBackups,
		MaxAge:     config.MaxAge,
	}

	return &hook, nil
}

// Levels , returns levels of log.
func (hook *RotateFileHook) Levels() []logrus.Level {
	return logrus.AllLevels[:hook.Config.Level+1]
}

// Fire , implementation of Hook interface.Fire.
func (hook *RotateFileHook) Fire(entry *logrus.Entry) (err error) {
	b, err := hook.Config.Formatter.Format(entry)
	if err != nil {
		return err
	}
	_, err = hook.logWriter.Write(b)

	return err
}

