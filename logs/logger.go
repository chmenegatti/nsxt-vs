package logs

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

func InitLogger() *zap.Logger {
	// Cria o diretório de logs se não existir
	if err := os.MkdirAll("logs", 0744); err != nil {
		panic(err)
	}

	// Configuração do arquivo de log com rotação
	fileRotate := &lumberjack.Logger{
		Filename:   "logs/app.log", // Nome do arquivo de log
		MaxSize:    10,             // Tamanho máximo em megabytes
		MaxBackups: 5,              // Número máximo de arquivos de backup
		MaxAge:     30,             // Dias máximos para manter os logs
		Compress:   true,           // Comprimir arquivos de backup
	}

	// Configuração do encoder
	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "ts",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		FunctionKey:    zapcore.OmitKey,
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}

	// Criar core para arquivo
	fileEncoder := zapcore.NewJSONEncoder(encoderConfig)
	fileWriter := zapcore.AddSync(fileRotate)
	fileLevel := zap.NewAtomicLevelAt(zap.InfoLevel)
	fileCore := zapcore.NewCore(fileEncoder, fileWriter, fileLevel)

	// Criar core para console
	consoleEncoder := zapcore.NewConsoleEncoder(encoderConfig)
	consoleWriter := zapcore.AddSync(os.Stdout)
	consoleLevel := zap.NewAtomicLevelAt(zap.InfoLevel)
	consoleCore := zapcore.NewCore(consoleEncoder, consoleWriter, consoleLevel)

	// Combinar os cores
	core := zapcore.NewTee(
		fileCore,
		consoleCore,
	)

	// Criar o logger
	logger := zap.New(
		core,
		zap.AddCaller(),
		zap.AddStacktrace(zap.ErrorLevel),
	)

	return logger
}
