package logger

import (
	"fmt"
	"github.com/Eugene-Usachev/fastbytes"
	"io"
	"os"
	"strconv"
	"time"
)

var std = os.Stderr

type StandardLogger struct {
	console io.Writer

	// errorWriter is the writer to which errors will be written.
	errorWriter io.Writer
	// fatalWriter is the writer to which fatal errors will be written.
	fatalWriter io.Writer
	// warningWriter is the writer to which warnings will be written.
	warningWriter io.Writer
	// infoWriter is the writer to which information will be written.
	infoWriter io.Writer
	// successWriter is the writer to which successes will be written.
	successWriter io.Writer
	// recordWriter is the writer to which records will be written.
	recordWriter io.Writer
	// rawWriter is the writer to which raw logs will be written.
	rawWriter io.Writer

	showDate bool
}

type StandardLoggerConfig struct {
	// IsWritingToTheConsole indicates whether the logger should write to the console.
	IsWritingToTheConsole bool
	// ErrorWriter is the writer to which errors will be written.
	ErrorWriter io.Writer
	// WarningWriter is the writer to which warnings will be written.
	WarningWriter io.Writer
	// InfoWriter is the writer to which information will be written.
	InfoWriter io.Writer
	// SuccessWriter is the writer to which successes will be written.
	SuccessWriter io.Writer
	// FatalWriter is the writer to which errors will be written.
	FatalWriter io.Writer
	// RecordWriter is the writer to which records will be written.
	RecordWriter io.Writer
	// RawWriter is the writer to which raw logs will be written.
	RawWriter io.Writer

	ShowDate bool
}

// NewStandardLogger creates a new StandardLogger.
func NewStandardLogger(cfg *StandardLoggerConfig) *StandardLogger {
	logger := &StandardLogger{}
	if cfg.IsWritingToTheConsole {
		logger.console = std
	}

	logger.errorWriter = cfg.ErrorWriter
	logger.warningWriter = cfg.WarningWriter
	logger.infoWriter = cfg.InfoWriter
	logger.successWriter = cfg.SuccessWriter
	logger.fatalWriter = cfg.FatalWriter
	logger.recordWriter = cfg.RecordWriter
	logger.rawWriter = cfg.RawWriter

	logger.showDate = cfg.ShowDate

	if len(Now.Load().([]byte)) == 0 {
		go func() {
			buf := make([]byte, 0, 19)
			for {
				time.Sleep(300 * time.Millisecond)
				date := time.Now()
				buf = make([]byte, 0, 19)
				year, month, day := date.Date()
				hour, min, sec := date.Clock()
				buf = strconv.AppendInt(buf, int64(year), 10)
				buf = append(buf, '/')
				buf = strconv.AppendInt(buf, int64(month), 10)
				buf = append(buf, '/')
				buf = strconv.AppendInt(buf, int64(day), 10)
				buf = append(buf, ' ')
				buf = strconv.AppendInt(buf, int64(hour), 10)
				buf = append(buf, ':')
				buf = strconv.AppendInt(buf, int64(min), 10)
				buf = append(buf, ':')
				if sec < 10 {
					buf = append(buf, '0')
				}
				buf = strconv.AppendInt(buf, int64(sec), 10)
				buf = append(buf, ' ')
				Now.Store(buf)
			}
		}()
	}

	return logger
}

func (logger *StandardLogger) log(buf []byte, writer io.Writer) {
	if logger.console != nil {
		logger.console.Write(buf)
	}
	if writer != nil {
		writer.Write(buf)
	}
}

func (logger *StandardLogger) info(buf []byte) {
	logger.log(buf, logger.infoWriter)
}

func (logger *StandardLogger) error(buf []byte) {
	logger.log(buf, logger.errorWriter)
}

func (logger *StandardLogger) warning(buf []byte) {
	logger.log(buf, logger.warningWriter)
}

func (logger *StandardLogger) success(buf []byte) {
	logger.log(buf, logger.successWriter)
}

func (logger *StandardLogger) fatal(buf []byte) {
	logger.log(buf, logger.fatalWriter)
	os.Exit(1)
}

func (logger *StandardLogger) record(buf []byte) {
	logger.log(buf, logger.recordWriter)
}

func (logger *StandardLogger) raw(buf []byte) {
	logger.log(buf, logger.rawWriter)
}

// Raw logs a raw log to the logger.rawWriter.
func (logger *StandardLogger) Raw(record []byte) {
	logger.raw(record)
}

// RawWithWriter logs a raw log to the writer.
func (logger *StandardLogger) RawWithWriter(record []byte, writer io.Writer) {
	logger.log(record, writer)
}

// Record logs a record to the logger.recordWriter. You can create a record with Builder(). Will reset the record.
func (logger *StandardLogger) Record(record *Record) {
	logger.record(record.rec)
	record.Reset()
	if record.wasGot {
		recordPool.Put(record)
	}
}

// RecordWithWriter logs a record to the writer. You can create a record with Builder().
func (logger *StandardLogger) RecordWithWriter(record Record, writer io.Writer) {
	logger.log(record.rec, writer)
	record.Reset()
	if record.wasGot {
		recordPool.Put(record)
	}
}

// Info logs a message to the logger.infoWriter.
func (logger *StandardLogger) Info(args ...interface{}) {
	buf := make([]byte, 0, 70)
	if logger.showDate {
		buf = append(Now.Load().([]byte), buf...)
	}
	buf = addArgsToLog(buf, args...)
	logger.info(append(buf, '\n'))
	buf = nil
}

// FormatInfo logs a message with format to the logger.infoWriter.
func (logger *StandardLogger) FormatInfo(f string, args ...interface{}) {
	if logger.showDate {
		buf := make([]byte, 0, 70)
		buf = append(Now.Load().([]byte), buf...)
		buf = append(buf, fastbytes.S2B(fmt.Sprintf(f, args...))...)
		logger.info(buf)
	} else {
		logger.info(fastbytes.S2B(fmt.Sprintf(f, args...)))
	}
}

// InfoPrepare logs a prepared record to the logger.infoWriter. Will not reset the record.
func (logger *StandardLogger) InfoPrepare(record *Record) {
	if record.isShowDate {
		copy(record.rec[:19], Now.Load().([]byte))
	}
	logger.info(record.rec)
}

// Error logs a message to the logger.errorWriter.
func (logger *StandardLogger) Error(args ...interface{}) {
	buf := make([]byte, 0, 70)
	if logger.showDate {
		buf = append(Now.Load().([]byte), buf...)
	}
	buf = addArgsToLog(buf, args...)
	logger.error(append(buf, '\n'))
	buf = nil
}

// FormatError logs a message with format to the logger.errorWriter.
func (logger *StandardLogger) FormatError(f string, args ...interface{}) {
	if logger.showDate {
		buf := make([]byte, 0, 70)
		buf = append(Now.Load().([]byte), buf...)
		buf = append(buf, fastbytes.S2B(fmt.Sprintf(f, args...))...)
		logger.error(buf)
	} else {
		logger.error(fastbytes.S2B(fmt.Sprintf(f, args...)))
	}
}

// ErrorPrepare logs a prepared record to the logger.errorWriter. Will not reset the record.
func (logger *StandardLogger) ErrorPrepare(record *Record) {
	if record.isShowDate {
		copy(record.rec[:19], Now.Load().([]byte))
	}
	logger.error(record.rec)
}

// Warning logs a message to the logger.warningWriter.
func (logger *StandardLogger) Warning(args ...interface{}) {
	buf := make([]byte, 0, 70)
	if logger.showDate {
		buf = append(Now.Load().([]byte), buf...)
	}
	buf = addArgsToLog(buf, args...)
	logger.warning(append(buf, '\n'))
	buf = nil
}

// FormatWarning logs a message with format to the logger.warningWriter.
func (logger *StandardLogger) FormatWarning(f string, args ...interface{}) {
	if logger.showDate {
		buf := make([]byte, 0, 70)
		buf = append(Now.Load().([]byte), buf...)
		buf = append(buf, fastbytes.S2B(fmt.Sprintf(f, args...))...)
		logger.warning(buf)
	} else {
		logger.warning(fastbytes.S2B(fmt.Sprintf(f, args...)))
	}
}

// WarningPrepare logs a prepared record to the logger.warningWriter. Will not reset the record.
func (logger *StandardLogger) WarningPrepare(record *Record) {
	if record.isShowDate {
		copy(record.rec[:19], Now.Load().([]byte))
	}
	logger.warning(record.rec)
}

// Success logs a message to the logger.successWriter.
func (logger *StandardLogger) Success(args ...interface{}) {
	buf := make([]byte, 0, 70)
	if logger.showDate {
		buf = append(Now.Load().([]byte), buf...)
	}
	buf = addArgsToLog(buf, args...)
	logger.success(append(buf, '\n'))
	buf = nil
}

// FormatSuccess logs a message with format to the logger.successWriter.
func (logger *StandardLogger) FormatSuccess(f string, args ...interface{}) {
	if logger.showDate {
		buf := make([]byte, 0, 70)
		buf = append(Now.Load().([]byte), buf...)
		buf = append(buf, fastbytes.S2B(fmt.Sprintf(f, args...))...)
		logger.success(buf)
	} else {
		logger.success(fastbytes.S2B(fmt.Sprintf(f, args...)))
	}
}

// SuccessPrepare logs a prepared record to the logger.successWriter. Will not reset the record.
func (logger *StandardLogger) SuccessPrepare(record *Record) {
	if record.isShowDate {
		copy(record.rec[:19], Now.Load().([]byte))
	}
	logger.success(record.rec)
}

// Fatal logs a message to the logger.fatalWriter. It will exit the program.
func (logger *StandardLogger) Fatal(args ...interface{}) {
	buf := make([]byte, 0, 70)
	if logger.showDate {
		buf = append(Now.Load().([]byte), buf...)
	}
	buf = addArgsToLog(buf, args...)
	logger.fatal(append(buf, '\n'))
	buf = nil
}

// FormatFatal logs a message with format to the logger.fatalWriter. It will exit the program.
func (logger *StandardLogger) FormatFatal(f string, args ...interface{}) {
	if logger.showDate {
		buf := make([]byte, 0, 70)
		buf = append(Now.Load().([]byte), buf...)
		buf = append(buf, fastbytes.S2B(fmt.Sprintf(f, args...))...)
		logger.fatal(buf)
	} else {
		logger.fatal(fastbytes.S2B(fmt.Sprintf(f, args...)))
	}
}

// FatalPrepare logs a prepared record to the logger.fatalWriter. Will not reset the record.
func (logger *StandardLogger) FatalPrepare(record *Record) {
	if record.isShowDate {
		copy(record.rec[:19], Now.Load().([]byte))
	}
	logger.fatal(record.rec)
}
