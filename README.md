# logger

This package provides a convenient tool for logging Go programs. It includes an already configured `FastLogger` that offers fast logging performance, with each operation taking less than 1 microsecond.

## Installation

To install the `logger` package, use the following command:
`go get github.com/Eugene-Usachev/logger`

## Usage

Here's an example of how to use the `FastLogger`:
```go
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
	fastLogger := testLogger.NewFastLogger(&testLogger.FastLoggerConfig{
		StandardLoggerConfig: *cfg,
		FlushInterval:        1 * time.Second,
		FatalFunc:            nil,
	})

	fastLogger.Info("[info fast] Hello, World!")
	fastLogger.Error("[error fast] Error message")
	fastLogger.Warning("[warning fast] Warning message")
	fastLogger.Success("[success fast] Success message")
	fastLogger.Record(testLogger.Builder().ShowDate().Prefix("[record fast] ").NewLine(true).AppendArgs("Build message").Build())
	fastLogger.Raw(fastbytes.S2B("[raw fast] Raw message!\n"))

	prepared := testLogger.Builder().Prefix("[my prefix] ").NewLine(true).AppendArgs("Hello,").Prepare()
	fastLogger.InfoPrepare(prepared)

	fastLogger.Fatal("[fatal fast] Fatal message")
}
```

## Contributing

We welcome contributions to the logger package! If you encounter any issues or have suggestions for improvements, please feel free to open an issue or contribute directly to the codebase. Your feedback and contributions are valuable in making this package even better.
## License

This project is licensed under the MIT License. See the [LICENSE](https://github.com/Eugene-Usachev/logger/blob/main/LICENSE) file for details.