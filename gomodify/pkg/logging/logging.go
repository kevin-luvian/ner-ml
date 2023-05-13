package logging

import (
	"fmt"
	"runtime"
	"strings"

	"github.com/sirupsen/logrus"
	"github.com/ttys3/rotatefilehook"
)

var (
	DefaultCallerDepth = 2

	logger *logrus.Logger = logrus.New()
)

// Setup initialize the log instance
func Setup() {
	filePath := getLogFilePath()
	fileName := getLogFileName()
	logger = logrus.New()

	logger.SetFormatter(&logrus.TextFormatter{
		DisableTimestamp: false,
		FullTimestamp:    true,
		TimestampFormat:  "15:04:05",
		PadLevelText:     true,
		DisableSorting:   true,
	})

	rotateFileHook, err := rotatefilehook.NewRotateFileHook(rotatefilehook.RotateFileConfig{
		Filename:   fmt.Sprintf("%s/%s", filePath, fileName),
		MaxSize:    10, // the maximum size in megabytes
		MaxBackups: 5,  // the maximum number of old log files to retain
		MaxAge:     7,  // the maximum number of days to retain old log files
		LocalTime:  true,
		Level:      logrus.InfoLevel,
		Formatter:  &logrus.JSONFormatter{},
	})
	if err != nil {
		logger.Fatal(err)
	}

	logger.AddHook(rotateFileHook)
}

// Info output logs at info level
func Debugln(args ...interface{}) {
	addFields().Debugln(args...)
}

// Info output logs at info level
func Infoln(args ...interface{}) {
	addFields().Infoln(args...)
}

// Infof output logs at info level
func Infof(format string, args ...interface{}) {
	addFields().Infof(format, args...)
}

// Info output logs at info level
func Warnf(format string, args ...interface{}) {
	addFields().Warnf(format, args...)
}

// Info output logs at info level
func Warnln(args ...interface{}) {
	addFields().Warnln(args...)
}

// Error output logs at error level
func Errorln(args ...interface{}) {
	addFields().Errorln(args...)
}

// Error output logs at error level
func Errorf(format string, args ...interface{}) {
	addFields().Errorf(format, args...)
}

// Fatal output logs at fatal level
func Fatalln(args ...interface{}) {
	addFields().Logln(logrus.FatalLevel, args...)
	panic(1)
}

// setPrefix set the prefix of the log output
func addFields() *logrus.Entry {
	fields := map[string]interface{}{}

	_, file, line, ok := runtime.Caller(DefaultCallerDepth)
	if ok {
		baseIndex := strings.LastIndex(file, "/server") + len("/server")
		file = file[baseIndex:]

		fields = logrus.Fields{
			"caller": fmt.Sprintf("%s:%d", file, line),
		}
	}

	return logger.WithFields(fields)
}
