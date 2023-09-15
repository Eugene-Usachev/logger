package logger

import (
	zerolog "github.com/rs/zerolog/log"
	"log"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

var (
	logsDir       = "logs"
	s       int8  = 10
	l       int64 = 9223372036854775807

	str_s = func() string {
		s := ""
		for i := 0; i < 48; i++ {
			s += "1"
		}
		return s
	}()
	str_l = func() string {
		s := strings.Builder{}
		for i := 0; i < 1000000; i++ {
			s.WriteString("1")
		}
		return s.String()
	}()
)

func BenchmarkLog__File(b *testing.B) {
	filePath := filepath.Join(logsDir, "log.txt")
	file, err := os.Create(filePath)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	log.SetOutput(file)
	for i := 0; i < b.N; i++ {
		log.Print(str_s, l)
		if s < 0 {
		}
	}
	err = file.Sync()
	if err != nil {
		log.Fatal(err)
	}
}

func BenchmarkZerolog__File(b *testing.B) {
	filePath := filepath.Join(logsDir, "zerolog.txt")
	file, err := os.Create(filePath)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	logger := zerolog.Output(file)
	for i := 0; i < b.N; i++ {
		logger.Debug().Int64("l", l).Msg(str_s)
		if s < 0 {
		}
	}
	err = file.Sync()
	if err != nil {
		log.Fatal(err)
	}
}

func BenchmarkStandardLogger_Log__File(b *testing.B) {
	filePath := filepath.Join(logsDir, "StandartLoggerLog.txt")
	file, err := os.Create(filePath)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	var logger = NewStandardLogger(&StandardLoggerConfig{
		IsWritingToTheConsole: false,
		ErrorWriter:           nil,
		WarningWriter:         nil,
		InfoWriter:            file,
		SuccessWriter:         nil,
		FatalWriter:           nil,
	})
	for i := 0; i < b.N; i++ {
		logger.Info(str_s, l)
		if s < 0 {
		}
	}
	err = file.Sync()
	if err != nil {
		log.Fatal(err)
	}
}

func BenchmarkStandardLogger_Log__File_Date(b *testing.B) {
	filePath := filepath.Join(logsDir, "StandartLoggerLogDate.txt")
	file, err := os.Create(filePath)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	var logger = NewStandardLogger(&StandardLoggerConfig{
		IsWritingToTheConsole: false,
		ErrorWriter:           nil,
		WarningWriter:         nil,
		InfoWriter:            file,
		SuccessWriter:         nil,
		FatalWriter:           nil,
		ShowDate:              true,
	})
	for i := 0; i < b.N; i++ {
		logger.Info(str_s, l)
		if s < 0 {
		}
	}
	err = file.Sync()
	if err != nil {
		log.Fatal(err)
	}
}

func BenchmarkFastLogger_Log__File(b *testing.B) {
	filePath := filepath.Join(logsDir, "FastLoggerLog.txt")
	file, err := os.Create(filePath)
	if err != nil {
		log.Fatal(err)
	}
	var logger = NewFastLogger(&FastLoggerConfig{
		StandardLoggerConfig: StandardLoggerConfig{
			IsWritingToTheConsole: false,
			ErrorWriter:           nil,
			WarningWriter:         nil,
			InfoWriter:            file,
			SuccessWriter:         nil,
			FatalWriter:           nil,
		},
		FlushInterval: 1 * time.Second,
		FatalFunc:     nil,
	})
	for i := 0; i < b.N; i++ {
		logger.Info(str_s, l)
		if s < 0 {
		}
	}
	err = file.Sync()
	if err != nil {
		log.Fatal(err)
	}
	logger.Stop()
}

func BenchmarkFastLogger_Log__File_Date(b *testing.B) {
	filePath := filepath.Join(logsDir, "FastLoggerLogDate.txt")
	file, err := os.Create(filePath)
	if err != nil {
		log.Fatal(err)
	}
	var logger = NewFastLogger(&FastLoggerConfig{
		StandardLoggerConfig: StandardLoggerConfig{
			IsWritingToTheConsole: false,
			ErrorWriter:           nil,
			WarningWriter:         nil,
			InfoWriter:            file,
			SuccessWriter:         nil,
			FatalWriter:           nil,
			ShowDate:              true,
		},
		FlushInterval: 1 * time.Second,
		FatalFunc:     nil,
	})
	for i := 0; i < b.N; i++ {
		logger.Info(str_s, l)
		if s < 0 {
		}
	}
	err = file.Sync()
	if err != nil {
		log.Fatal(err)
	}
	logger.Stop()
}

func BenchmarkFastLogger_Log__File_Builder(b *testing.B) {
	filePath := filepath.Join(logsDir, "FastLoggerLogBuilder.txt")
	file, err := os.Create(filePath)
	if err != nil {
		log.Fatal(err)
	}
	var logger = NewFastLogger(&FastLoggerConfig{
		StandardLoggerConfig: StandardLoggerConfig{
			IsWritingToTheConsole: false,
			ErrorWriter:           nil,
			WarningWriter:         nil,
			InfoWriter:            nil,
			SuccessWriter:         nil,
			FatalWriter:           nil,
			RecordWriter:          file,
			RawWriter:             nil,
		},
		FlushInterval: 1 * time.Second,
		FatalFunc:     nil,
	})
	for i := 0; i < b.N; i++ {
		logger.Record(Builder().AppendArgs(str_s, l).
			Build())
		if s < 0 {
		}
	}
	err = file.Sync()
	if err != nil {
		log.Fatal(err)
	}
	logger.Stop()
}

func BenchmarkFastLogger_Log__File_Builder_Date(b *testing.B) {
	filePath := filepath.Join(logsDir, "FastLoggerLogBuilderDate.txt")
	file, err := os.Create(filePath)
	if err != nil {
		log.Fatal(err)
	}
	var logger = NewFastLogger(&FastLoggerConfig{
		StandardLoggerConfig: StandardLoggerConfig{
			IsWritingToTheConsole: false,
			ErrorWriter:           nil,
			WarningWriter:         nil,
			InfoWriter:            nil,
			SuccessWriter:         nil,
			FatalWriter:           nil,
			RecordWriter:          file,
			RawWriter:             nil,
		},
		FlushInterval: 1 * time.Second,
		FatalFunc:     nil,
	})
	for i := 0; i < b.N; i++ {
		logger.Record(Builder().
			ShowDate().
			AppendArgs(str_s, l).
			Build())
		if s < 0 {
		}
	}
	err = file.Sync()
	if err != nil {
		log.Fatal(err)
	}
	logger.Stop()
}

func BenchmarkFastLogger_Log__File_Builder_DateAndPrefix(b *testing.B) {
	filePath := filepath.Join(logsDir, "FastLoggerLogBuilderDateAndPrefix.txt")
	file, err := os.Create(filePath)
	if err != nil {
		log.Fatal(err)
	}
	var logger = NewFastLogger(&FastLoggerConfig{
		StandardLoggerConfig: StandardLoggerConfig{
			IsWritingToTheConsole: false,
			ErrorWriter:           nil,
			WarningWriter:         nil,
			InfoWriter:            nil,
			SuccessWriter:         nil,
			FatalWriter:           nil,
			RecordWriter:          file,
			RawWriter:             nil,
		},
		FlushInterval: 1 * time.Second,
		FatalFunc:     nil,
	})
	for i := 0; i < b.N; i++ {
		logger.Record(Builder().ShowDate().AppendArgs(str_s, l).Prefix("[prefix]").Build())
		if s < 0 {
		}
	}
	err = file.Sync()
	if err != nil {
		log.Fatal(err)
	}
	logger.Stop()
}

func BenchmarkFastLogger_Log__File_Builder_DateAndPrefix_Prepare(b *testing.B) {
	filePath := filepath.Join(logsDir, "FastLoggerLogBuilderDateAndPrefix_Prepare.txt")
	file, err := os.Create(filePath)
	if err != nil {
		log.Fatal(err)
	}
	var logger = NewFastLogger(&FastLoggerConfig{
		StandardLoggerConfig: StandardLoggerConfig{
			IsWritingToTheConsole: false,
			ErrorWriter:           nil,
			WarningWriter:         nil,
			InfoWriter:            nil,
			SuccessWriter:         nil,
			FatalWriter:           nil,
			RecordWriter:          file,
			RawWriter:             nil,
		},
		FlushInterval: 1 * time.Second,
		FatalFunc:     nil,
	})
	prepared := Builder().ShowDate().AppendArgs(str_s, l).Prefix("[prefix]").Prepare()
	for i := 0; i < b.N; i++ {
		logger.InfoPrepare(prepared)
		if s < 0 {
		}
	}
	err = file.Sync()
	if err != nil {
		log.Fatal(err)
	}
	logger.Stop()
}

//func BenchmarkLog(b *testing.B) {
//	var buf bytes.Buffer
//	log.SetOutput(&buf)
//	for i := 0; i < b.N; i++ {
//		log.Print(str_s, l)
//		if s < 0 {
//		}
//	}
//	buf.Reset()
//}
//
//func BenchmarkFmtLog(b *testing.B) {
//	var buf bytes.Buffer
//	log.SetOutput(&buf)
//	for i := 0; i < b.N; i++ {
//		log.Printf("%s %d", str_s, l)
//		if s < 0 {
//		}
//	}
//	buf.Reset()
//}
//
//func BenchmarkStandardLogger_Log(b *testing.B) {
//	var buf bytes.Buffer
//	var logger = NewStandardLoggerWithCustomWriter("", &buf)
//	for i := 0; i < b.N; i++ {
//		logger.Info(str_s, l)
//		if s < 0 {
//		}
//	}
//	buf.Reset()
//}
//
//func BenchmarkStandardLogger_LogWithFmt(b *testing.B) {
//	var buf bytes.Buffer
//	var logger = NewStandardLoggerWithCustomWriter("", &buf)
//	for i := 0; i < b.N; i++ {
//		logger.FormatInfo("%s %d", str_s, l)
//		if s < 0 {
//		}
//	}
//	buf.Reset()
//}
//
//func BenchmarkFastLogger_Log(b *testing.B) {
//	var buf bytes.Buffer
//	var logger = NewFastLoggerWithCustomWriter("", &buf)
//	for i := 0; i < b.N; i++ {
//		logger.Info(str_s, l)
//		if s < 0 {
//		}
//	}
//	buf.Reset()
//}
//
//func BenchmarkFastLogger_LogWithFmt(b *testing.B) {
//	var buf bytes.Buffer
//	var logger = NewFastLoggerWithCustomWriter("", &buf)
//	for i := 0; i < b.N; i++ {
//		logger.FormatInfo("%s %d", str_s, l)
//		if s < 0 {
//		}
//	}
//	buf.Reset()
//}
