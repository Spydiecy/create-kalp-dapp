package kalpsdk

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
)

const (
	defaultLogFormat       = "[%lvl%]: %time% - %msg%" // Default log format will output [INFO]: 2006-01-02T15:04:05Z07:00 - Log message
	defaultTimestampFormat = time.RFC3339
	defaultLogLevel        = logrus.DebugLevel
)

var defaultLogOutput = os.Stdout
var channelName string
var isLogLevelSet bool
var chaincodeLogger = &ChaincodeLogger{
	Logger: &logrus.Logger{
		Hooks:    make(logrus.LevelHooks),
		ExitFunc: os.Exit,
	},
	StackTrace: true,
}

// Formatter implements the logrus.Formatter interface.
type Formatter struct {
	TimestampFormat string // Timestamp format
	LogFormat       string // Available standard keys: time, msg, lvl. Custom fields should be wrapped inside %, e.g., %time% %msg%
}

// Format builds the log message.
func (f *Formatter) Format(entry *logrus.Entry) ([]byte, error) {
	output := f.LogFormat
	if output == "" {
		output = defaultLogFormat
	}

	timestampFormat := f.TimestampFormat
	if timestampFormat == "" {
		timestampFormat = defaultTimestampFormat
	}

	output = strings.Replace(output, "%time%", entry.Time.Format(timestampFormat), 1)
	output = strings.Replace(output, "%msg%", entry.Message, 1)
	output = strings.Replace(output, "%lvl%", getColorByLogLevel(entry.Level), 1)

	for k, val := range entry.Data {
		switch v := val.(type) {
		case string:
			output = strings.Replace(output, "%"+k+"%", v, 1)
		case int:
			s := strconv.Itoa(v)
			output = strings.Replace(output, "%"+k+"%", s, 1)
		case bool:
			s := strconv.FormatBool(v)
			output = strings.Replace(output, "%"+k+"%", s, 1)
		}
	}
	return []byte(output), nil
}

// getColorByLogLevel returns the ANSI escape code for the log level color.
func getColorByLogLevel(level logrus.Level) string {
	switch level {
	case logrus.DebugLevel:
		return "\x1b[36m" + strings.ToUpper(level.String()) + "\x1b[0m" // Cyan
	case logrus.InfoLevel:
		return "\x1b[32m" + strings.ToUpper(level.String()) + "\x1b[0m" // Green
	case logrus.WarnLevel:
		return "\x1b[33m" + strings.ToUpper(level.String()) + "\x1b[0m" // Yellow
	case logrus.ErrorLevel, logrus.FatalLevel, logrus.PanicLevel:
		return "\x1b[31m" + strings.ToUpper(level.String()) + "\x1b[0m" // Red
	default:
		return ""
	}
}

// setupChaincodeLogging sets up the chaincode logger with default configurations.
func setupChaincodeLogging() {
	if chaincodeLogger.Logger.Formatter == nil {
		formatter := &Formatter{
			TimestampFormat: "2006-01-02 15:04:05.000 MST",
			LogFormat:       "\x1b[33m[KALP-SDK] %time%\x1b[0m %lvl% - %msg%\n",
		}
		chaincodeLogger.Logger.SetFormatter(formatter)
	}

	if !isLogLevelSet {
		chaincodeLogger.Logger.SetLevel(defaultLogLevel)
	}

	if chaincodeLogger.Logger.Out == nil {
		chaincodeLogger.Logger.SetOutput(defaultLogOutput)
	}
}

// ChaincodeLogger is an abstraction of a logging object for use by chaincodes.
type ChaincodeLogger struct {
	Logger     *logrus.Logger
	StackTrace bool
}

// NewLogger returns the logger instance for ChaincodeLogger.
func NewLogger() *ChaincodeLogger {
	return chaincodeLogger
}

// SetChaincodeLogLevel sets the log level for the chaincode logger.
func (chLogger *ChaincodeLogger) SetChaincodeLogLevel(level string) {
	l, err := logrus.ParseLevel(level)
	if err != nil {
		chLogger.Errorf("Invalid log level '%s'. It must be one of: debug, info, warning, error, fatal, panic", level)
		return
	}
	isLogLevelSet = true
	chLogger.Logger.SetLevel(l)
}

// SetChaincodeOutput sets the output for the chaincode logger.
func (chLogger *ChaincodeLogger) SetChaincodeOutput(output io.Writer) {
	chLogger.Logger.SetOutput(output)
}

// SetChaincodeFormatter sets the formatter for the chaincode logger.
func (chLogger *ChaincodeLogger) SetChaincodeFormatter(formatter *Formatter) {
	chLogger.Logger.SetFormatter(formatter)
}

// DisableStackTrace disables the stack trace in the log message.
func (chLogger *ChaincodeLogger) DisableStackTrace() {
	chLogger.StackTrace = false
}

// getCallerInfo returns the caller information in the format: [channelName] [filename:line] functionName
func getCallerInfo() string {
	pc, file, line, _ := runtime.Caller(2)
	funcPtr := runtime.FuncForPC(pc)
	functionName := "<unknown>"
	if funcPtr != nil {
		functionName = filepath.Base(funcPtr.Name())
	}
	return fmt.Sprintf("[%s] [%s:%d] %s ", channelName, filepath.Base(file), line, functionName)
}

// Log functions for logging messages at different levels with caller information

// Trace logs a message at the Trace level.
func (c *ChaincodeLogger) Trace(args ...interface{}) {
	if c.StackTrace {
		args = append([]interface{}{getCallerInfo()}, args...)
	}
	c.Logger.Trace(args...)
}

// Debug logs a message at the Debug level.
func (c *ChaincodeLogger) Debug(args ...interface{}) {
	if c.StackTrace {
		args = append([]interface{}{getCallerInfo()}, args...)
	}
	c.Logger.Debug(args...)
}

// Info logs a message at the Info level.
func (c *ChaincodeLogger) Info(args ...interface{}) {
	if c.StackTrace {
		args = append([]interface{}{getCallerInfo()}, args...)
	}
	c.Logger.Info(args...)
}

// Print logs a message at the Print level.
func (c *ChaincodeLogger) Print(args ...interface{}) {
	if c.StackTrace {
		args = append([]interface{}{getCallerInfo()}, args...)
	}
	c.Logger.Print(args)
}

// Warn logs a message at the Warn level.
func (c *ChaincodeLogger) Warn(args ...interface{}) {
	if c.StackTrace {
		args = append([]interface{}{getCallerInfo()}, args...)
	}
	c.Logger.Warn(args...)
}

// Warning logs a message at the Warning level.
func (c *ChaincodeLogger) Warning(args ...interface{}) {
	if c.StackTrace {
		args = append([]interface{}{getCallerInfo()}, args...)
	}
	c.Logger.Warning(args...)
}

// Error logs a message at the Error level.
func (c *ChaincodeLogger) Error(args ...interface{}) {
	if c.StackTrace {
		args = append([]interface{}{getCallerInfo()}, args...)
	}
	c.Logger.Error(args...)
}

// Fatal logs a message at the Fatal level.
func (c *ChaincodeLogger) Fatal(args ...interface{}) {
	if c.StackTrace {
		args = append([]interface{}{getCallerInfo()}, args...)
	}
	c.Logger.Fatal(args...)
}

// Panic logs a message at the Panic level.
func (c *ChaincodeLogger) Panic(args ...interface{}) {
	if c.StackTrace {
		args = append([]interface{}{getCallerInfo()}, args...)
	}
	c.Logger.Panic(args...)
}

// ------------------------------------------------------------------

// Tracef logs a formatted message at the Trace level.
func (c *ChaincodeLogger) Tracef(format string, args ...interface{}) {
	if c.StackTrace {
		args = append([]interface{}{getCallerInfo()}, args...)
	}
	c.Logger.Tracef(format, args...)
}

// Debugf logs a formatted message at the Debug level.
func (c *ChaincodeLogger) Debugf(format string, args ...interface{}) {
	if c.StackTrace {
		args = append([]interface{}{getCallerInfo()}, args...)
	}
	c.Logger.Debugf(format, args...)
}

// Infof logs a formatted message at the Info level.
func (c *ChaincodeLogger) Infof(format string, args ...interface{}) {
	if c.StackTrace {
		args = append([]interface{}{getCallerInfo()}, args...)
	}
	c.Logger.Infof(format, args...)
}

// Printf logs a formatted message at the Print level.
func (c *ChaincodeLogger) Printf(format string, args ...interface{}) {
	if c.StackTrace {
		args = append([]interface{}{getCallerInfo()}, args...)
	}
	c.Logger.Printf(format, args...)
}

// Warnf logs a formatted message at the Warn level.
func (c *ChaincodeLogger) Warnf(format string, args ...interface{}) {
	if c.StackTrace {
		args = append([]interface{}{getCallerInfo()}, args...)
	}
	c.Logger.Warnf(format, args...)
}

// Warningf logs a formatted message at the Warning level.
func (c *ChaincodeLogger) Warningf(format string, args ...interface{}) {
	if c.StackTrace {
		args = append([]interface{}{getCallerInfo()}, args...)
	}
	c.Logger.Warningf(format, args...)
}

// Errorf logs a formatted message at the Error level.
func (c *ChaincodeLogger) Errorf(format string, args ...interface{}) {
	if c.StackTrace {
		args = append([]interface{}{getCallerInfo()}, args...)
	}
	c.Logger.Errorf(format, args...)
	// str := fmt.Sprintf(format, args...)
	// var err error
	// err = errors.New(str)
	// return err
}

// Fatalf logs a formatted message at the Fatal level.
func (c *ChaincodeLogger) Fatalf(format string, args ...interface{}) {
	if c.StackTrace {
		args = append([]interface{}{getCallerInfo()}, args...)
	}
	c.Logger.Fatalf(format, args...)
}

// Panicf logs a formatted message at the Panic level.
func (c *ChaincodeLogger) Panicf(format string, args ...interface{}) {
	if c.StackTrace {
		args = append([]interface{}{getCallerInfo()}, args...)
	}
	c.Logger.Panicf(format, args...)
}

// ------------------------------------------------------------------

// Traceln logs a message with a new line at the Trace level.
func (c *ChaincodeLogger) Traceln(args ...interface{}) {
	if c.StackTrace {
		args = append([]interface{}{getCallerInfo()}, args...)
	}
	c.Logger.Traceln(args...)
}

// Debugln logs a message with a new line at the Debug level.
func (c *ChaincodeLogger) Debugln(args ...interface{}) {
	if c.StackTrace {
		args = append([]interface{}{getCallerInfo()}, args...)
	}
	c.Logger.Debugln(args...)
}

// Infoln logs a message with a new line at the Info level.
func (c *ChaincodeLogger) Infoln(args ...interface{}) {
	if c.StackTrace {
		args = append([]interface{}{getCallerInfo()}, args...)
	}
	c.Logger.Infoln(args...)
}

// Println logs a message with a new line at the Print level.
func (c *ChaincodeLogger) Println(args ...interface{}) {
	if c.StackTrace {
		args = append([]interface{}{getCallerInfo()}, args...)
	}
	c.Logger.Println(args...)
}

// Warnln logs a message with a new line at the Warn level.
func (c *ChaincodeLogger) Warnln(args ...interface{}) {
	if c.StackTrace {
		args = append([]interface{}{getCallerInfo()}, args...)
	}
	c.Logger.Warnln(args...)
}

// Warningln logs a message with a new line at the Warning level.
func (c *ChaincodeLogger) Warningln(args ...interface{}) {
	if c.StackTrace {
		args = append([]interface{}{getCallerInfo()}, args...)
	}
	c.Logger.Warningln(args...)
}

// Errorln logs a message with a new line at the Error level.
func (c *ChaincodeLogger) Errorln(args ...interface{}) {
	if c.StackTrace {
		args = append([]interface{}{getCallerInfo()}, args...)
	}
	c.Logger.Errorln(args...)
}

// Fatalln logs a message with a new line at the Fatal level.
func (c *ChaincodeLogger) Fatalln(args ...interface{}) {
	if c.StackTrace {
		args = append([]interface{}{getCallerInfo()}, args...)
	}
	c.Logger.Fatalln(args...)
}

// Panicln logs a message with a new line at the Panic level.
func (c *ChaincodeLogger) Panicln(args ...interface{}) {
	if c.StackTrace {
		args = append([]interface{}{getCallerInfo()}, args...)
	}
	c.Logger.Panicln(args...)
}

// Exit calls the logger's Exit method.
func (c *ChaincodeLogger) Exit(code int) {
	c.Logger.Exit(code)
}
