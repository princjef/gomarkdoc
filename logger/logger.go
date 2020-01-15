package logger

import (
	"github.com/sirupsen/logrus"
	prefixed "github.com/x-cray/logrus-prefixed-formatter"
)

type (
	// Logger provides basic logging capabilities at different logging levels.
	Logger interface {
		Debug(a ...interface{})
		Debugf(format string, a ...interface{})
		Info(a ...interface{})
		Infof(format string, a ...interface{})
		Warn(a ...interface{})
		Warnf(format string, a ...interface{})
		Error(a ...interface{})
		Errorf(format string, a ...interface{})
	}

	// Level defines valid logging levels for a Logger.
	Level int

	// Option defines an option for configuring the logger.
	Option func(opts *options)

	// options defines options for configuring the logger
	options struct {
		fields map[string]interface{}
	}
)

// Valid logging levels
const (
	DebugLevel Level = iota + 1
	InfoLevel
	WarnLevel
	ErrorLevel
)

// New initializes a new Logger.
func New(level Level, opts ...Option) Logger {
	var options options
	for _, opt := range opts {
		opt(&options)
	}

	log := logrus.New()

	formatter := &prefixed.TextFormatter{
		DisableTimestamp: true,
	}
	formatter.SetColorScheme(&prefixed.ColorScheme{
		DebugLevelStyle: "cyan",
		PrefixStyle:     "black+h",
	})

	log.Formatter = formatter

	switch level {
	case DebugLevel:
		log.SetLevel(logrus.DebugLevel)
	case InfoLevel:
		log.SetLevel(logrus.InfoLevel)
	case WarnLevel:
		log.SetLevel(logrus.WarnLevel)
	case ErrorLevel:
		log.SetLevel(logrus.ErrorLevel)
	default:
		log.SetLevel(logrus.ErrorLevel)
	}

	if options.fields != nil {
		return log.WithFields(options.fields)
	}

	return log
}

// WithField sets the provided key/value pair for use on all logs.
func WithField(key string, value interface{}) Option {
	return func(opts *options) {
		if opts.fields == nil {
			opts.fields = make(map[string]interface{})
		}

		opts.fields[key] = value
	}
}
