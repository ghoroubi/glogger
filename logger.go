package glogger

import (
	"context"
	"errors"
	logrustash "github.com/bshuster-repo/logrus-logstash-hook"
	"github.com/mattn/go-colorable"
	"github.com/olivere/elastic/v7"
	"github.com/sirupsen/logrus"
	"glogger/rotate"
	"gopkg.in/sohlich/elogrus.v7"
	"net"
	"time"
)

// NewLogger ...
// Initializes and returns the logger instance
func NewLogger(config *LoggerConfig) *logrus.Logger {
	var (
		// logger instance
		appLogger = logrus.StandardLogger()
	)

	// File rotator init
	// Setting up a file as logger with custom attributes which comes from config file
	rotateFileHook, err := rotate.NewRotateFileHook(rotate.RotateFileConfig{
		Filename:   config.Filename,
		MaxSize:    config.MaxAge,
		MaxBackups: config.MaxBackups,
		MaxAge:     config.MaxAge,
		Level:      logrus.Level(config.Level),
		Formatter: &logrus.JSONFormatter{
			DisableHTMLEscape: true,
			PrettyPrint:       config.PrettyPrint,
		},
	})
	if err != nil {
		logrus.Fatalf("Failed to initialize file rotate hook: %v", err)
	}

	// RotateFileHook for logrus
	// Which streams the log data into the provided file hook
	appLogger.AddHook(rotateFileHook)

	// Set level for main logger instance
	appLogger.SetLevel(logrus.Level(config.Level))

	// Set report caller which reports the caller function of logger operation
	appLogger.SetReportCaller(false)

	// Print to StdOut only in debug mode
	if config.STDOut {
		appLogger.SetOutput(colorable.NewColorableStdout())
	}

	// Stream to elastic and logstash
	// only in production mode
	if config.UseLogStash {
		// Adding logstash hook which streams the log data to the logstash engine
		logstashHook := getLogstashHook("tcp",
			config.LogstashConfig.Address,
			config.NameFields["service_name"])
		if logstashHook != nil {
			appLogger.AddHook(logstashHook)
		}

		// Adding elastic hook which streams the log data to the elastic dataset
		elasticHook, err := getElasticHook(config.ElasticConfig.Address,
			config.NameFields["service_name"], config.ElasticConfig.TimeOut)
		if err != nil {
			logrus.Warningf("Failed to initialize elastic hook: %v \n", err)
			return appLogger
		}

		// Adding elastic hook
		appLogger.AddHook(elasticHook)
	}

	return appLogger
}

// Logstash hook for logrus
func getLogstashHook(network string, addr string, logo string) logrus.Hook {
	if network == "" {
		network = "tcp"
	}
	n, err := net.Dial("tcp", addr)
	hook := logrustash.New(n, &logrus.JSONFormatter{})
	if err != nil {
		return nil
	}
	return hook
}

// ElasticSearch hook for logrus
func getElasticHook(addr string, svcName string, timeOut string) (logrus.Hook, error) {

	// Validate elastic address
	if addr == "" {
		return nil, errors.New("empty address")
	}

	// New elastic client
	d, err := time.ParseDuration(timeOut)
	if err != nil {
		return nil, err
	}

	// The timeout of elastic connection is too long,
	// Thus a context with timeout is required which
	// Gets the amount from config file
	ctx, cancelFunc := context.WithTimeout(context.Background(), d)
	// cancel context as soon as deadline occurs
	defer cancelFunc()

	// Creating elastic client with above context
	client, err := elastic.DialContext(ctx, elastic.SetURL(addr))
	if err != nil {
		return nil, err
	}

	// New elastic async hook
	hook, err := elogrus.NewAsyncElasticHook(client, "localhost", logrus.DebugLevel, svcName)
	if err != nil {
		return nil, err
	}

	return hook, nil
}
