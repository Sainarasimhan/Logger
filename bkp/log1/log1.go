package log

import (
	"io"
	stdlog "log"
	"strings"
)

//Log Levels
const (
	//LevelError - Use to Print Errors
	LevelError = iota
	//LevelInfo - Use to Print Info (Default)
	LevelInfo
	//LevelDebug - Use to Print Debug messages
	LevelDebug
)

//Log Strings
const (
	//ErrStr - String representing Error logs
	ErrStr string = "[Error]"
	//InfoStr - String representing logs
	InfoStr string = "[Info ]"
	//DebugStr - String representing Debug logs
	DebugStr string = "[Debug]"
)

var (
	//NoOpLogger Logger that does nothing
	NoOpLogger = func(_ string, _ ...string) PrintFunc {
		return func(_ string, _ ...interface{}) {}
	}
)

//Logger - Logger with levels; uses standard logger internally
type Logger struct {
	log   *stdlog.Logger
	level int8
}

//PrintFunc - Func type returned by Log Methods
type PrintFunc func(string, ...interface{})

//LoggerFunc - Func type of Log Write
type LoggerFunc func(string, ...string) PrintFunc

//New - Creates new Logger with requested prefix and flags
func New(w io.Writer, prefix string, flag int) *Logger {
	l := stdlog.New(w, prefix+" ", flag)

	//Set All Methods to NoOpLogger, by Default
	return &Logger{log: l}
}

//Level - Change Logger Level Supported Levels (Error/Info/Debug)
//Can be called dynamically to change log level
func (l *Logger) Level(level int8) *Logger {
	l.level = level
	return l
}

//Error - Error logging.
//Usage - log.Error(context info)(Log Message)
//Pass Contexts along with actual error messages
//Verifies if Error is allowed, if so adds  Error String
func (l *Logger) Error(ctx ...string) PrintFunc {
	if l.level >= LevelError {
		return l.logWrite(ErrStr, ctx...)
	}
	return NoOpLogger("", "")
}

//Info - Info logging.
//Usage - log.Info(context info)(Log Message)
//Pass Contexts along with actual info messages
//Verifies if Info is allowed, if so adds Info String
func (l *Logger) Info(ctx ...string) PrintFunc {
	if l.level >= LevelError {
		return l.logWrite(InfoStr, ctx...)
	}
	return NoOpLogger("", "")
}

//Debug - Debug logging.
//Usage - log.Debug(context info)(Log Message)
//Pass Contexts along with actual debug messages
//Verifies if Debug is allowed, if so adds Debug String
func (l *Logger) Debug(ctx ...string) PrintFunc {
	if l.level >= LevelError {
		return l.logWrite(DebugStr, ctx...)
	}
	return NoOpLogger("", "")
}

//logWrite - function to write Log , Joins context and actual log.
//uses square braces [] to wrap contexts
//uses | to split contexts
//uses <> to wrap actual log message
func (l *Logger) logWrite(level string, ctx ...string) PrintFunc {
	prefix := level + "[" + strings.Join(ctx, "|") + "]" // Add Logging Level and [headers]
	return func(format string, args ...interface{}) {
		l.log.Printf(prefix+"<"+format+">", args) // Add Log message within <>
	}
}
