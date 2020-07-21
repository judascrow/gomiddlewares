package gomiddlewares

import (
	"os"
	"path"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

var (
	file *os.File
	e    error
)

func GoLogger() gin.HandlerFunc {

	timeFormat := "2006-01-02 15:04:05"

	writerSyncer := getLogWriter()
	encoder := getEncoder()

	core := zapcore.NewCore(encoder, writerSyncer, zapcore.DebugLevel)

	logger := zap.New(core, zap.AddCaller())

	return func(c *gin.Context) {
		start := time.Now()
		// some evil middlewares modify this values
		path := c.Request.URL.Path
		query := c.Request.URL.RawQuery

		c.Next()

		end := time.Now()
		latency := end.Sub(start)

		if len(c.Errors) > 0 {
			// Append error field if this is an erroneous request.
			for _, e := range c.Errors.Errors() {
				logger.Error(e)
			}
		} else {

			status := c.Writer.Status()

			loggerConfig := logger.Info
			switch {
			case status > 499:
				loggerConfig = logger.Fatal
			case status > 399:
				loggerConfig = logger.Error
			case status > 299:
				loggerConfig = logger.Warn
			}

			loggerConfig(path,
				zap.Int("status", status),
				zap.String("method", c.Request.Method),
				zap.String("path", path),
				zap.String("query", query),
				zap.String("ip", c.ClientIP()),
				zap.String("user-agent", c.Request.UserAgent()),
				zap.String("time", end.Format(timeFormat)),
				zap.Duration("latency", latency),
			)
		}
	}

}

func getEncoder() zapcore.Encoder {
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	return zapcore.NewJSONEncoder(encoderConfig)
}
func getLogWriter() zapcore.WriteSyncer {

	now := time.Now() //or time.Now().UTC()

	logFileName := now.Format("2006-01-02") + ".log"
	dirPath := path.Join(".", "logs")
	// Create directory if does not exist
	if _, err := os.Stat(dirPath); os.IsNotExist(err) {
		err = os.MkdirAll(dirPath, os.ModeDir)
		if err != nil {
			panic(err)
		}
	}

	file, e = os.OpenFile(path.Join(dirPath, logFileName), os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0664)
	if e != nil {
		panic(e)
	}

	lumberJackLogger := &lumberjack.Logger{
		Filename:   path.Join(dirPath, logFileName),
		MaxSize:    10,
		MaxBackups: 5,
		MaxAge:     30,
		Compress:   false,
	}
	return zapcore.AddSync(lumberJackLogger)
}

func CloseLogFile() {
	if err := file.Close(); err != nil {
		return
	}
}
