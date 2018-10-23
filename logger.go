package main

import (
	"encoding/json"
	"fmt"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"io/ioutil"
	"os"
	"strings"
)

type logConfig struct {
	Desc      string
	Level     string
	Stdout    bool
	Encoding  string
	AddCaller bool
	Color     bool
	FilesOut  bool
	LogsPath  []*logFilePath
}

type logFilePath struct {
	Level string
	Hook  *lumberjack.Logger
}

var logger *zap.Logger

func init() {
	logcfg := logConfig{
		Desc:      "development",
		Level:     "debug",
		Stdout:    true,
		Encoding:  "console",
		AddCaller: true,
		Color:     false,
		FilesOut:  true,
		/*
			Hook: &lumberjack.Logger{
				Filename:   "./logs/gaven.log", // Filename is the file to write logs to.
				MaxSize:    1024,               // megabytes
				MaxAge:     7,                  // days
				MaxBackups: 3,                  // the maximum number of old log files to retain.
			},
		*/
	}
	exists, err := pathExists("./logs.json")
	if err != nil {
		fmt.Println(err)
	}
	if exists {
		file, err := ioutil.ReadFile("./logs.json")
		if err != nil {
			panic(err)
		}
		if err := json.Unmarshal(file, &logcfg); err != nil {
			panic(err)
		}
	}

	// Print only logs above the loggerLevel.
	loggerLevel := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		switch logcfg.Level {
		case "debug":
			return lvl >= zap.DebugLevel
		case "info":
			return lvl >= zap.InfoLevel
		case "error":
			return lvl >= zap.ErrorLevel
		default:
			return lvl >= zap.InfoLevel
		}
	})
	fileWriter := zapcore.AddSync(logcfg.LogsPath[1].Hook)
	// Output should also go to standard out.
	consoleDebugging := zapcore.Lock(os.Stdout)
	var encoderConfig zapcore.EncoderConfig
	var fileEncoder, consoleEncoder zapcore.Encoder
	if strings.EqualFold(logcfg.Desc, "production") {
		encoderConfig = zap.NewProductionEncoderConfig()
	}
	if strings.EqualFold(logcfg.Desc, "development") {
		encoderConfig = zap.NewDevelopmentEncoderConfig()
	}
	if strings.EqualFold(logcfg.Encoding, "json") {
		fileEncoder = zapcore.NewJSONEncoder(encoderConfig)
		consoleEncoder = zapcore.NewJSONEncoder(encoderConfig)
	}
	if strings.EqualFold(logcfg.Encoding, "console") {
		fileEncoder = zapcore.NewConsoleEncoder(encoderConfig)
		if logcfg.Color {
			encoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
		}
		consoleEncoder = zapcore.NewConsoleEncoder(encoderConfig)
	}
	var cores []zapcore.Core
	if !logcfg.Stdout && !logcfg.FilesOut || logcfg.Stdout {
		cores := append(cores, zapcore.NewCore(consoleEncoder, consoleDebugging, loggerLevel))
	}
	core := zapcore.NewTee(cores...)
	/*
		core := zapcore.NewTee(
			// Printed in the console.
			zapcore.NewCore(consoleEncoder, consoleDebugging, loggerLevel),
			// Printed in file.
			zapcore.NewCore(fileEncoder, fileWriter, loggerLevel),
		)
	*/
	// From a zapcore.Core to construct a Logger.
	if logcfg.AddCaller {
		logger = zap.New(core, zap.AddCaller())
	} else {
		logger = zap.New(core)
	}
}
