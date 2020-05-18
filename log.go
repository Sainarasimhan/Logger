package log

import (
	stdlog "log"
	"os"
	"strings"
)

//Log Levels and Strings
const (
	//LevelError - Use to Print Errors
	LevelError = "Error"
	//LevelInfo - Use to Print Info (Default)
	LevelInfo = "Info"
	//LevelDebug - Use to Print Debug messages
	LevelDebug = "Debug"

	//ErrStr - String representing Error in logs
	ErrStr = "ER"
	//InfoStr - String representing info in logs
	InfoStr = "INF"
	//DebugStr - String representing Debug in logs
	DebugStr = "DBG"
)

//Copied from std Log package, can be modified if req
const (
	Ldate         = 1 << iota     // the date in the local time zone: 2009/01/23
	Ltime                         // the time in the local time zone: 01:23:23
	Lmicroseconds                 // microsecond resolution: 01:23:23.123123.  assumes Ltime.
	Llongfile                     // full file name and line number: /a/b/c/d.go:23
	Lshortfile                    // final file name element and line number: d.go:23. overrides Llongfile
	LUTC                          // if Ldate or Ltime is set, use UTC rather than the local time zone
	LstdFlags     = Ldate | Ltime // initial values for the standard logger
)

//Cfg - Configuration of Logger
type Cfg struct {
	// Log Level
	Level string `yaml:"Level" valid:"Level,required"`
	// Prefix to appear in Logs
	Prefix string `yaml:"Prefix" valid:"ascii, required"`
	// Log Flags
	Flags int `yaml:"Flags" valid:"numeric,required"`
}

//Logger - Logger with levels; uses standard logger internally
type Logger struct {
	log                      *stdlog.Logger
	errorFn, debugFn, infoFn LoggerFunc
}

//PrintFunc - Func type returned by Log Methods
type PrintFunc func(string, ...interface{})

//LoggerFunc - Func type of Log Write
type LoggerFunc func(string, ...string) PrintFunc

//NoOpFn Logger that does nothing
var (
	NoOpFn = func(_ string, _ ...string) PrintFunc {
		return func(_ string, _ ...interface{}) {}
	}
)

//New - Creates new Logger with requested prefix and flags
func New(c Cfg) *Logger {
	l := stdlog.New(
		os.Stdout, // By Default write to std out
		c.Prefix+" ",
		c.Flags)

	//Set All Methods to NoOpLogger, by Default
	logger := &Logger{log: l,
		errorFn: NoOpFn,
		debugFn: NoOpFn,
		infoFn:  NoOpFn}

	//Set Level for logger
	logger.Level(c.Level)

	return logger
}

//Level - Change Logger Level Supported Levels (Error/Info/Debug)
//Can be called dynamically to change log level
func (l *Logger) Level(level string) *Logger {
	switch level {
	case LevelDebug:
		l.debugFn = l.logWrite
		fallthrough //If Debug is requested, enable All Levels
	case LevelInfo:
		l.infoFn = l.logWrite
		fallthrough //If Info is requested, enable Info and Error
	case LevelError:
		l.errorFn = l.logWrite
	}
	return l
}

//Error - Error logging.
//Usage - log.Error(context info)(Log Message)
//Pass Contexts along with actual error messages
//Verifies if Error is allowed, if so adds  Error String
func (l *Logger) Error(ctx ...string) PrintFunc {
	return l.errorFn(ErrStr, ctx...)
}

//Info - Info logging.
//Usage - log.Info(context info)(Log Message)
//Pass Contexts along with actual info messages
//Verifies if Info is allowed, if so adds Info String
func (l *Logger) Info(ctx ...string) PrintFunc {
	return l.infoFn(InfoStr, ctx...)
}

//Debug - Debug logging.
//Usage - log.Debug(context info)(Log Message)
//Pass Contexts along with actual debug messages
//Verifies if Debug is allowed, if so adds Debug String
func (l *Logger) Debug(ctx ...string) PrintFunc {
	return l.debugFn(DebugStr, ctx...)
}

//logWrite - function to write Log , Joins context and actual log.
//uses square braces [] to wrap contexts
//uses | to split contexts
//uses <> to wrap actual log message
func (l *Logger) logWrite(level string, ctx ...string) PrintFunc {
	prefix := level + "[" + strings.Join(ctx, "|") + "]" // Add Logging Level and [headers]
	return func(format string, args ...interface{}) {
		l.log.Printf(prefix+"<"+format+">", args...) // Add Log message within <>
	}
}

//Fatal - calls stdlog fatal func
func Fatal(args ...interface{}) {
	stdlog.Fatal(args...)
}
