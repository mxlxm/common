package utils

import (
	"path"
	"strings"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

type LogConfig struct {
	lumberjack.Logger
	zap.Config
	Encoder zapcore.Encoder
	Fields  []zap.Field
	Options []zap.Option
}

const (
	// lumberjack
	defaultSize      = 1024
	defaultAge       = 7
	defaultBackup    = 3
	defaultLocalTime = true
	defaultCompress  = true
	// zap
	defaultDevelopment       = false
	defaultDisableCaller     = false
	defaultDisableStacktrace = false
	defaultEncoding          = "json"
)

var (
	defaultLevel         = zap.NewAtomicLevelAt(zap.DebugLevel)
	defaultSampling      *zap.SamplingConfig
	defaultEncoderConfig = zap.NewProductionEncoderConfig()
	defaultInitialFields = map[string]interface{}{}
)

func InitLogConfig(c LogConfig) (l *LogConfig) {
	l = &LogConfig{
		lumberjack.Logger{
			MaxSize:    defaultSize,
			MaxAge:     defaultAge,
			MaxBackups: defaultBackup,
			LocalTime:  true,
			Compress:   true,
		},
		zap.Config{
			Level:             defaultLevel,
			DisableCaller:     defaultDisableCaller,
			DisableStacktrace: defaultDisableStacktrace,
			Encoding:          defaultEncoding,
			EncoderConfig:     defaultEncoderConfig,
			InitialFields:     defaultInitialFields,
		},
		zapcore.NewJSONEncoder(defaultEncoderConfig),
		[]zap.Field{},
		[]zap.Option{},
	}
	if c.Filename != "" {
		l.Filename = c.Filename
	} else {
		l.Filename = "default.log"
	}
	if c.Level != (zap.AtomicLevel{}) {
		l.Level = c.Level
	}
	if c.MaxSize != 0 {
		l.MaxSize = c.MaxSize
	}
	if c.MaxAge != 0 {
		l.MaxAge = c.MaxAge
	}
	if c.MaxBackups != 0 {
		l.MaxBackups = c.MaxBackups
	}
	if !c.DisableCaller {
		l.Options = append(l.Options, zap.AddCaller())
	}
	if !c.DisableStacktrace {
		l.Options = append(l.Options, zap.AddStacktrace(zap.ErrorLevel))
	}
	if c.Sampling != nil {
		l.Options = append(l.Options, zap.WrapCore(func(core zapcore.Core) zapcore.Core {
			return zapcore.NewSampler(core, time.Second, int(c.Sampling.Initial), int(c.Sampling.Thereafter))
		}))
	}
	switch c.Encoding {
	case "console":
		encoderCfg := zap.NewProductionEncoderConfig()
		encoderCfg.TimeKey = "timestamp"
		encoderCfg.EncodeTime = zapcore.ISO8601TimeEncoder
		encoderCfg.EncodeLevel = zapcore.CapitalLevelEncoder
		encoderCfg.EncodeDuration = zapcore.StringDurationEncoder
		encoderCfg.EncodeCaller = zapcore.ShortCallerEncoder
		l.Encoder = zapcore.NewConsoleEncoder(encoderCfg)
	default:
		encoderCfg := zap.NewProductionEncoderConfig()
		encoderCfg.TimeKey = "timestamp"
		encoderCfg.EncodeTime = zapcore.ISO8601TimeEncoder
		encoderCfg.EncodeLevel = zapcore.CapitalLevelEncoder
		encoderCfg.EncodeDuration = zapcore.StringDurationEncoder
		encoderCfg.EncodeCaller = zapcore.ShortCallerEncoder
		l.Encoder = zapcore.NewJSONEncoder(encoderCfg)
	}
	if len(c.InitialFields) > 0 {
		for k, v := range c.InitialFields {
			l.Fields = append(l.Fields, zap.Field(zap.Any(k, v)))
		}
	} else {
		l.Fields = append(l.Fields, zap.Field(zap.Any("app", GetProgramName())))
	}
	l.Options = append(l.Options, zap.Fields(l.Fields...))
	return l
}

func LogInit(c *LogConfig) (l *zap.Logger, err error) {
	if !strings.HasPrefix(c.Filename, "/") {
		c.Filename = path.Clean(ExecDir() + "/../logs/" + c.Filename)
	}
	w := zapcore.AddSync(&lumberjack.Logger{
		Filename:   c.Filename,
		MaxSize:    c.MaxSize,
		MaxBackups: c.MaxBackups,
		MaxAge:     c.MaxAge,
		LocalTime:  true,
	})
	l = zap.New(zapcore.NewCore(c.Encoder, w, c.Level), c.Options...)
	return
}

func SugarInit(c *LogConfig) (s *zap.SugaredLogger, err error) {
	var l *zap.Logger
	if l, err = LogInit(c); err != nil {
		return
	}
	s = l.Sugar()
	return
}
