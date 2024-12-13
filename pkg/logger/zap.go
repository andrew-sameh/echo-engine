package logger

import (
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/andrew-sameh/echo-engine/internal/config"
	"github.com/andrew-sameh/echo-engine/pkg/file"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

// Zap SugaredLogger by default
// DesugarZap performance-sensitive code
type Logger struct {
	Zap        *zap.SugaredLogger
	DesugarZap *zap.Logger
}

// Logger singleton
var (
	ZLogger Logger
)

func (l *Logger) LogWithFields(msg string, requestID string, serviceName string) {
	l.Zap.With(zap.String("request_id", requestID), zap.String("service_name", serviceName)).Info(msg)
}
func NewLogger(config config.LoggerConfig) *Logger {
	var options []zap.Option
	var encoder zapcore.Encoder

	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "ts",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		FunctionKey:    zapcore.OmitKey,
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.CapitalLevelEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
		EncodeTime:     prodTimeEncoder,
	}

	if config.Format == "json" {
		encoder = zapcore.NewJSONEncoder(encoderConfig)
	} else {
		encoder = zapcore.NewConsoleEncoder(encoderConfig)
	}

	level := zap.NewAtomicLevelAt(toLevel(config.Level))

	core := zapcore.NewCore(encoder, toWriter(config), level)

	stackLevel := zap.NewAtomicLevel()
	stackLevel.SetLevel(zap.WarnLevel)
	options = append(options,
		zap.AddCaller(),
		zap.AddCallerSkip(1),
		zap.AddStacktrace(stackLevel),
	)

	logger := zap.New(core, options...)
	ZLogger = Logger{Zap: logger.Sugar(), DesugarZap: logger}
	return &ZLogger
}

func localTimeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(t.Format(time.Stamp))
}
func prodTimeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(t.UTC().Format("2006-01-02T15:04:05Z0700"))
	// enc.AppendString(t.Format(constants.TimeFormat))
}

func toLevel(level string) zapcore.Level {
	switch strings.ToLower(level) {
	case "debug":
		return zap.DebugLevel
	case "info":
		return zap.InfoLevel
	case "warn":
		return zap.WarnLevel
	case "error":
		return zap.ErrorLevel
	case "dpanic":
		return zap.DPanicLevel
	case "panic":
		return zap.PanicLevel
	case "fatal":
		return zap.FatalLevel
	default:
		return zap.InfoLevel
	}
}

func toWriter(config config.LoggerConfig) zapcore.WriteSyncer {
	fp := ""
	sp := string(filepath.Separator)

	fp, _ = filepath.Abs(filepath.Dir(filepath.Join(".")))
	fp += sp + "logs" + sp

	if config.Directory != "" {
		if err := file.EnsureDirRW(config.Directory); err != nil {
			fp = config.Directory
		}
	}
	// Create a slice of WriteSyncer
	var writers []zapcore.WriteSyncer

	// Always add os.Stdout to writers
	writers = append(writers, zapcore.AddSync(os.Stdout))

	// Only add local file logging if config.Local is true
	if config.Local {
		writers = append(writers, zapcore.AddSync(&lumberjack.Logger{ // 文件切割
			Filename:   filepath.Join(fp, config.Name) + ".log",
			MaxSize:    100,
			MaxAge:     0,
			MaxBackups: 0,
			LocalTime:  false,
			Compress:   true,
		}))
	}

	return zapcore.NewMultiWriteSyncer(writers...)
}
