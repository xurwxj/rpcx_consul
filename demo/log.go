package demo

import (
	"fmt"
	"time"

	"github.com/rs/zerolog"
	"github.com/xurwxj/gtils/base"
	"github.com/xurwxj/viper"
)

// InitLog init log instance
func InitLog() {
	fmt.Println("initializing logger...")
	level := viper.GetString("log.level")
	if level == "" {
		level = "info"
	}
	switch level {
	case "info":
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	case "error":
		zerolog.SetGlobalLevel(zerolog.ErrorLevel)
	case "debug":
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	case "fatal":
		zerolog.SetGlobalLevel(zerolog.FatalLevel)
	case "panic":
		zerolog.SetGlobalLevel(zerolog.PanicLevel)
	case "trace":
		zerolog.SetGlobalLevel(zerolog.TraceLevel)
	case "warn":
		zerolog.SetGlobalLevel(zerolog.WarnLevel)
	}
	zerolog.TimeFieldFormat = time.RFC3339
	var c base.Config
	c.EncodeLogsAsJson = true
	switch viper.GetString("log.output") {
	case "", "file":
		c.ConsoleLoggingEnabled = false
		c.FileLoggingEnabled = true
	case "std":
		c.ConsoleLoggingEnabled = true
		c.FileLoggingEnabled = false
	case "fileStd":
		c.ConsoleLoggingEnabled = true
		c.FileLoggingEnabled = true
	}
	logPath := viper.GetString("log.path")
	if logPath == "" {
		logPath = "logs"
	}
	c.Directory = logPath
	logFile := viper.GetString("log.file")
	if logFile == "" {
		logFile = "service.log"
	}
	c.Filename = logFile
	logMax := viper.GetInt("log.max")
	if logMax == 0 {
		logMax = 10
	}
	c.MaxSize = logMax
	logMaxAge := viper.GetInt("log.maxAge")
	if logMaxAge == 0 {
		logMaxAge = 7
	}
	c.MaxAge = logMaxAge
	c.MaxBackups = viper.GetInt("log.maxBackups")
	c.LocalTime = viper.GetBool("log.localtime")
	Log = base.Configure(c)
	if Log != nil {
		fmt.Println("init log done")
		Log.Info().Msg("init log done")
	}
}
