package main

import (
	"errors"
	"github.com/Eugene-Usachev/fastbytes"
	testLogger "github.com/Eugene-Usachev/logger"
	"os"
	"path/filepath"
	"time"
)

func main() {

	logsDir := filepath.Join("../", "logs")
	err := os.Mkdir(logsDir, 0777)
	if err != nil && !errors.Is(err, os.ErrExist) {
		return
	}

	filePathInfo := filepath.Join(logsDir, "info.txt")
	infoFile, _ := os.OpenFile(filePathInfo, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	defer infoFile.Close()

	filePathError := filepath.Join(logsDir, "error.txt")
	errorFile, _ := os.OpenFile(filePathError, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	defer errorFile.Close()

	filePathWarning := filepath.Join(logsDir, "warning.txt")
	warningFile, _ := os.OpenFile(filePathWarning, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	defer warningFile.Close()

	filePathSuccess := filepath.Join(logsDir, "success.txt")
	successFile, _ := os.OpenFile(filePathSuccess, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	defer successFile.Close()

	filePathRecord := filepath.Join(logsDir, "record.txt")
	recordFile, _ := os.OpenFile(filePathRecord, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	defer recordFile.Close()

	filePathRaw := filepath.Join(logsDir, "raw.txt")
	rawFile, _ := os.OpenFile(filePathRaw, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	defer rawFile.Close()

	filePathFatal := filepath.Join(logsDir, "fatal.txt")
	fatalFile, _ := os.OpenFile(filePathFatal, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	defer fatalFile.Close()

	cfg := &testLogger.StandardLoggerConfig{
		IsWritingToTheConsole: true,
		ErrorWriter:           errorFile,
		WarningWriter:         warningFile,
		InfoWriter:            infoFile,
		SuccessWriter:         successFile,
		FatalWriter:           fatalFile,
		RecordWriter:          recordFile,
		RawWriter:             rawFile,
		ShowDate:              true,
	}

	logger := testLogger.NewStandardLogger(cfg)

	logger.Info("[info std] Hello, World!")
	logger.Error("[error std] Hello, World!")
	logger.Warning("[warning std] Hello, World!")
	logger.Success("[success std] Hello, World!")
	logger.Record(testLogger.Builder().ShowDate().Prefix("[record std] ").NewLine(true).AppendArgs("Hello, World!").Build())
	logger.Raw(fastbytes.S2B("[raw std] Hello, World!\n"))

	fastLogger := testLogger.NewFastLogger(&testLogger.FastLoggerConfig{
		StandardLoggerConfig: *cfg,
		FlushInterval:        1 * time.Second,
		FatalFunc:            nil,
	})

	fastLogger.Info("[info fast] Hello, World!")
	fastLogger.Error("[error fast] Hello, World!")
	fastLogger.Warning("[warning fast] Hello, World!")
	fastLogger.Success("[success fast] Hello, World!")
	fastLogger.Record(testLogger.Builder().ShowDate().Prefix("[record fast] ").NewLine(true).AppendArgs("Hello, World!").Build())
	fastLogger.Raw(fastbytes.S2B("[raw fast] Hello, World!\n"))

	prepared := testLogger.Builder().Prefix("[my prefix] ").NewLine(true).AppendArgs("Hello, World!").Prepare()
	fastLogger.InfoPrepare(prepared)

	fastLogger.Fatal("[fatal fast] Hello, World!")
}
