package glogger

//LogLevel represents the level of logging,
//which is categorized in the constants definition.
type LogLevel uint32

//String returns the string value of the provided level.
func (l LogLevel) String() string {
	switch l {
	case DebugLevel:
		return "Debug"
	case InfoLevel:
		return "Info"
	case ErrorLevel:
		return "Error"
	case WarnLevel:
		return "Warning"
	case TraceLevel:
		return "Trace"
	case PanicLevel:
		return "Panic"
	case FatalLevel:
		return "Fatal"
	default:
		return "unknown level"
	}
}

//Int returns the integer value of the provided log,
//this will be used in some cases like file rotate configurations.
func (l LogLevel) Int() int {
	return int(l)
}

// These are the different logging levels. You can set the logging level to log
// on your instance of logger, obtained when logger initializes.
const (
	// PanicLevel level, highest level of severity. Logs and then calls panic with the
	// message passed to Debug, Info, ...
	PanicLevel LogLevel = iota
	// FatalLevel level. Logs and then calls `logger.Exit(1)`. It will exit even if the
	// logging level is set to Panic.
	FatalLevel
	// ErrorLevel level. Logs. Used for errors that should definitely be noted.
	// Commonly used for hooks to send errors to an error tracking service.
	ErrorLevel
	// WarnLevel level. Non-critical entries that deserve eyes.
	WarnLevel
	// InfoLevel level. General operational entries about what's going on inside the
	// application.
	InfoLevel
	// DebugLevel level. Usually only enabled when debugging. Very verbose logging.
	DebugLevel
	// TraceLevel level. Designates finer-grained informational events than the Debug.
	TraceLevel
)

//LoggerConfig , represents the configuration of the logger object,
//which is needed in logger initialization.

type LoggerConfig struct {
	Filename       string            `json:"filename"`
	MaxSize        int               `json:"max_size"`
	MaxBackups     int               `json:"max_backups"`
	MaxAge         int               `json:"max_age"`
	Level          LogLevel          `json:"level"`
	PrettyPrint    bool              `json:"pretty_print"`
	STDOut         bool              `json:"std_out"`
	UseElastic     bool              `json:"use_elastic"`
	ElasticConfig  *ElasticConfig    `json:"elastic_config"`
	UseLogStash    bool              `json:"use_log_stash"`
	LogstashConfig *LogstashConfig   `json:"logstash_config"`
	UseOthers      bool              `json:"use_others"`
	NameFields     map[string]string `json:"name_fields"`
}

//ElasticConfig ...
type ElasticConfig struct {
	Address string `json:"address"`
	TimeOut string `json:"time_out"`
}

//LogstashConfig ...
type LogstashConfig struct {
	Address string `json:"address"`
}
