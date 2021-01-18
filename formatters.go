package glogger

import (
	"github.com/sirupsen/logrus"
	"runtime"
	"time"
)

//Formatter , formats the lof output.
type Formatter logrus.Formatter

type FormatterConfig struct {
	// TimestampFormat sets the format used for marshaling timestamps.
	TimestampFormat string

	// DisableTimestamp allows disabling automatic timestamps in output
	DisableTimestamp bool

	// DisableHTMLEscape allows disabling html escaping in output
	DisableHTMLEscape bool

	// DataKey allows users to put all the log entry parameters into a nested dictionary at a given key.
	DataKey string

	// FieldMap allows users to customize the names of keys for default fields.
	// As an example:
	// formatter := &JSONFormatter{
	//   	FieldMap: FieldMap{
	// 		 FieldKeyTime:  "@timestamp",
	// 		 FieldKeyLevel: "@level",
	// 		 FieldKeyMsg:   "@message",
	// 		 FieldKeyFunc:  "@caller",
	//    },
	// }
	FieldMap logrus.FieldMap

	// CallerPrettyfier can be set by the user to modify the content
	// of the function and file keys in the json data when ReportCaller is
	// activated. If any of the returned value is the empty string the
	// corresponding key will be removed from json fields.
	CallerPrettyfier func(*runtime.Frame) (function string, file string)

	// PrettyPrint will indent all json logs
	PrettyPrint bool
}

//JSONFormatter returns the json formatter.
func JSONFormatter(config *FormatterConfig) Formatter {
	return &logrus.JSONFormatter{
		TimestampFormat:   config.TimestampFormat,
		DisableTimestamp:  config.DisableTimestamp,
		DisableHTMLEscape: config.DisableHTMLEscape,
		DataKey:           config.DataKey,
		FieldMap:          config.FieldMap,
		CallerPrettyfier:  config.CallerPrettyfier,
		PrettyPrint:       config.PrettyPrint,
	}
}

//DefaultJSONFormatter , returns a default JSON formatter,
//if the user doesn't like to use a default configuration.
func DefaultJSONFormatter() Formatter {
	return &logrus.JSONFormatter{
		FieldMap: logrus.FieldMap{
			logrus.FieldKeyTime: "@timestamp",
			logrus.FieldKeyMsg:  "message",
		},
		TimestampFormat: time.UnixDate,
	}
}
