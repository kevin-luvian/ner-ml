package logging

import (
	"fmt"
	"time"

	"github.com/kevin-luvian/gomodify/pkg/setting"
)

// getLogFilePath get the log file save path
func getLogFilePath() string {
	return setting.App.LogSavePath
}

// getLogFileName get the save name of the log file
func getLogFileName() string {
	s := setting.App

	return fmt.Sprintf("%s%s.%s",
		s.LogSaveName,
		time.Now().Format(s.TimeFormat),
		s.LogFileExt,
	)
}
