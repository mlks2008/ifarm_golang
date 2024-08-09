package zaplogger

import (
	"components/log/zaplogger/rotatelogs"
	"fmt"
	"github.com/go-kratos/kratos/v2/log"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
	"sync"
	"time"
)

var (
	rw            sync.RWMutex
	DefaultLogger *ZLogger
	loggers       map[string]*ZLogger
	appId         string
)

func init() {
	DefaultLogger = NewConfigLogger(DefaultConsoleConfig())
	loggers = make(map[string]*ZLogger)
}

type ZLogger struct {
	*zap.Logger
	*Config
}

func Flush() {
	_ = DefaultLogger.Sync()

	for _, logger := range loggers {
		_ = logger.Sync()
	}
}

func NewLogger(debug bool, opts ...zap.Option) *ZLogger {
	return NewLoggerWithName("", "", "", debug, opts...)
}

func NewLoggerWithName(logPath, serviceName, refLoggerName string, debug bool, opts ...zap.Option) *ZLogger {
	defer rw.Unlock()
	rw.Lock()

	if logger, found := loggers[refLoggerName]; found {
		return logger
	}

	if refLoggerName == "" {
		refLoggerName = "logstash"
	}

	if serviceName == "" {
		serviceName, _ = os.Hostname()
	}
	appId = serviceName

	config := DefaultConfig(logPath, serviceName, refLoggerName)
	if debug {
		config.Level = "debug"
		config.StackLevel = "error"
		config.EnableConsole = true
	}

	logger := NewConfigLogger(config, opts...)
	loggers[refLoggerName] = logger

	return logger
}

func NewConfigLogger(config *Config, opts ...zap.Option) *ZLogger {
	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "ts",
		LevelKey:       "level",
		CallerKey:      "caller",
		MessageKey:     "msg",
		NameKey:        "name",
		StacktraceKey:  "stack",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeDuration: zapcore.StringDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}

	encoderConfig.EncodeLevel = func(level zapcore.Level, encoder zapcore.PrimitiveArrayEncoder) {
		if appId != "" {
			encoder.AppendString(fmt.Sprintf("%s  %-5s", appId, level.CapitalString()))
		} else {
			encoder.AppendString(level.CapitalString())
		}
	}

	if config.PrintCaller {
		encoderConfig.EncodeTime = config.TimeEncoder()
		encoderConfig.EncodeName = zapcore.FullNameEncoder
		encoderConfig.FunctionKey = zapcore.OmitKey
		opts = append(opts, zap.AddCallerSkip(1))
		opts = append(opts, zap.AddCaller())
	}

	opts = append(opts, zap.AddStacktrace(GetLevel(config.StackLevel)))

	var writers []zapcore.WriteSyncer

	if config.EnableWriteFile {
		hook, err := rotatelogs.New(
			config.FilePathFormat,
			rotatelogs.WithLinkName(config.FileLinkPath),
			rotatelogs.WithMaxAge(time.Hour*24*time.Duration(config.MaxAge)),
			rotatelogs.WithRotationTime(time.Second*time.Duration(config.RotationTime)),
		)

		if err != nil {
			panic(err)
		}

		writers = append(writers, zapcore.AddSync(hook))
	}

	if config.EnableConsole {
		writers = append(writers, zapcore.AddSync(os.Stderr))
	}

	if config.IncludeStdout {
		writers = append(writers, zapcore.Lock(os.Stdout))
	}

	if config.IncludeStderr {
		writers = append(writers, zapcore.Lock(os.Stderr))
	}

	core := zapcore.NewCore(
		zapcore.NewConsoleEncoder(encoderConfig),
		zapcore.AddSync(zapcore.NewMultiWriteSyncer(writers...)),
		zap.NewAtomicLevelAt(GetLevel(config.Level)),
	)

	logger := &ZLogger{
		Logger: zap.New(core, opts...),
		Config: config,
	}

	return logger
}

func (l *ZLogger) Log(level log.Level, kv ...interface{}) error {
	if len(kv) == 0 || len(kv)%2 != 0 {
		l.Logger.Warn(fmt.Sprint("kv must appear in pairs: ", kv))
		return nil
	}
	// Zap.Field is used when keyvals pairs appear
	var data []zap.Field
	for i := 0; i < len(kv); i += 2 {
		data = append(data, zap.Any(fmt.Sprint(kv[i]), fmt.Sprint(kv[i+1])))
	}
	switch level {
	case log.LevelInfo:
		l.Logger.Info("", data...)
	case log.LevelWarn:
		l.Logger.Warn("", data...)
	case log.LevelError:
		l.Logger.Error("", data...)
	case log.LevelFatal:
		l.Logger.Fatal("", data...)
	default:
		l.Logger.Debug("", data...)
	}
	return nil
}

func (l *ZLogger) Info(msg string, kv ...interface{}) {
	if len(kv) == 0 || len(kv)%2 != 0 {
		l.Logger.Warn(fmt.Sprint("kv must appear in pairs: ", kv))
		return
	}
	// Zap.Field is used when keyvals pairs appear
	var data []zap.Field
	for i := 0; i < len(kv); i += 2 {
		data = append(data, zap.Any(fmt.Sprint(kv[i]), fmt.Sprint(kv[i+1])))
	}

	l.Logger.Info(msg, data...)
}

func (l *ZLogger) Error(err error, msg string, kv ...interface{}) {
	if len(kv) == 0 || len(kv)%2 != 0 {
		l.Logger.Warn(fmt.Sprint("kv must appear in pairs: ", kv))
		return
	}
	// Zap.Field is used when keyvals pairs appear
	var data []zap.Field
	for i := 0; i < len(kv); i += 2 {
		data = append(data, zap.Any(fmt.Sprint(kv[i]), fmt.Sprint(kv[i+1])))
	}

	l.Logger.Error(msg, data...)
}
