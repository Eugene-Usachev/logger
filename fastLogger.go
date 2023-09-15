package logger

import (
	"fmt"
	"sync"
	"sync/atomic"
	"time"
)

/*
FastLogger is a logger that writes to the standard logger with buffers. FastLogger is thread-safe and faster than StandardLogger by many orders.
FastLogger will not log buffers if the application crashes before the recording interval has passed!
If you don't use panics outside FastLogger, FastLogger will only not record if the machine is turned off.
Use FastLogger.Fatal() or FastLogger.FormatFatal() instead of panic, or use FastLogger.Flush() before shutting down the application.

Example:

	logger := NewFastLogger(&FastLoggerConfig{
		StandardLoggerConfig: StandardLoggerConfig{
			IsWritingToTheConsole: true,
			ErrorWriter:           errorFile,
			WarningWriter:         warningFile,
			InfoWriter:            infoFile,
			SuccessWriter:         successFile,
			FatalWriter:           fatalFile,
			RecordWriter:          recordFile,
			RawWriter:             rawFile,
		},
		FlushInterval: 1 * time.Second,
		FatalFunc:     func(reason any) {
		panic(fmt.Sprintf("logger could not write to the file reason: %v", reason))
		},
	})
*/
type FastLogger struct {
	stdLogger *StandardLogger

	isRunning atomic.Bool

	infoLogs     []byte
	infoMutex    sync.Mutex
	errorLogs    []byte
	errorMutex   sync.Mutex
	warningLogs  []byte
	warningMutex sync.Mutex
	successLogs  []byte
	successMutex sync.Mutex
	recordLogs   []byte
	recordMutex  sync.Mutex
	rawLogs      []byte
	rawMutex     sync.Mutex

	fatalFunc func(reason any)
}

type FastLoggerConfig struct {
	StandardLoggerConfig

	// FlushInterval is the interval between flushes to the writers.
	FlushInterval time.Duration
	// FatalFunc is the function to call when a fatal error occurs in the logger.
	FatalFunc func(reason any)
}

// NewFastLogger creates a new FastLogger.
func NewFastLogger(cfg *FastLoggerConfig) *FastLogger {
	logger := &FastLogger{
		stdLogger:    NewStandardLogger(&cfg.StandardLoggerConfig),
		isRunning:    atomic.Bool{},
		infoLogs:     make([]byte, 0),
		infoMutex:    sync.Mutex{},
		errorLogs:    make([]byte, 0),
		errorMutex:   sync.Mutex{},
		warningLogs:  make([]byte, 0),
		warningMutex: sync.Mutex{},
		successLogs:  make([]byte, 0),
		successMutex: sync.Mutex{},
		recordLogs:   make([]byte, 0),
		recordMutex:  sync.Mutex{},
		rawLogs:      make([]byte, 0),
		rawMutex:     sync.Mutex{},
		fatalFunc:    cfg.FatalFunc,
	}

	logger.isRunning.Store(true)
	interval := cfg.FlushInterval
	go func() {
		for {
			time.Sleep(interval)
			logger.Flush()
			if !logger.isRunning.Load() {
				break
			}
		}
	}()

	return logger
}

// Flush flushes all logs to the logger.
func (logger *FastLogger) Flush() {
	defer func() {
		if err := recover(); err != nil {
			logger.fatalFunc(err)
			logger.infoLogs = logger.infoLogs[:0]
			logger.infoMutex.Unlock()
			logger.errorLogs = logger.errorLogs[:0]
			logger.errorMutex.Unlock()
			logger.warningLogs = logger.warningLogs[:0]
			logger.warningMutex.Unlock()
			logger.successLogs = logger.successLogs[:0]
			logger.successMutex.Unlock()
		}
	}()
	{
		logger.infoMutex.Lock()
		if len(logger.infoLogs) > 0 {
			logger.stdLogger.info(logger.infoLogs)
			logger.infoLogs = logger.infoLogs[:0]
		}
		logger.infoMutex.Unlock()
	}
	{
		logger.errorMutex.Lock()
		if len(logger.errorLogs) > 0 {
			logger.stdLogger.error(logger.errorLogs)
			logger.errorLogs = logger.errorLogs[:0]
		}
		logger.errorMutex.Unlock()
	}
	{
		logger.warningMutex.Lock()
		if len(logger.warningLogs) > 0 {
			logger.stdLogger.warning(logger.warningLogs)
			logger.warningLogs = logger.warningLogs[:0]
		}
		logger.warningMutex.Unlock()
	}
	{
		logger.successMutex.Lock()
		if len(logger.successLogs) > 0 {
			logger.stdLogger.success(logger.successLogs)
			logger.successLogs = logger.successLogs[:0]
		}
		logger.successMutex.Unlock()
	}
	{
		logger.recordMutex.Lock()
		if len(logger.recordLogs) > 0 {
			logger.stdLogger.record(logger.recordLogs)
			logger.recordLogs = logger.recordLogs[:0]
		}
		logger.recordMutex.Unlock()
	}
	{
		logger.rawMutex.Lock()
		if len(logger.rawLogs) > 0 {
			logger.stdLogger.raw(logger.rawLogs)
			logger.rawLogs = logger.rawLogs[:0]
		}
		logger.rawMutex.Unlock()
	}
}

// Stop stops the logger.
func (logger *FastLogger) Stop() {
	logger.isRunning.Store(false)
}

// StopWithoutFlush stops the logger without flushing. WILL CLEAR NOT FLUSHED LOGS!
func (logger *FastLogger) StopWithoutFlush() {
	logger.infoMutex.Lock()
	logger.errorMutex.Lock()
	logger.warningMutex.Lock()
	logger.successMutex.Lock()
	logger.recordMutex.Lock()
	logger.rawMutex.Lock()
	logger.infoLogs = logger.infoLogs[:0]
	logger.errorLogs = logger.errorLogs[:0]
	logger.warningLogs = logger.warningLogs[:0]
	logger.successLogs = logger.successLogs[:0]
	logger.recordLogs = logger.recordLogs[:0]
	logger.rawLogs = logger.rawLogs[:0]
	logger.infoMutex.Unlock()
	logger.errorMutex.Unlock()
	logger.warningMutex.Unlock()
	logger.successMutex.Unlock()
	logger.recordMutex.Unlock()
	logger.rawMutex.Unlock()
	logger.isRunning.Store(false)
}

// Record logs a record to the logger.stdLogger.recordWriter.
func (logger *FastLogger) Record(record *Record) {
	logger.recordMutex.Lock()
	logger.recordLogs = append(logger.recordLogs, record.rec...)
	record.Reset()
	if record.wasGot {
		recordPool.Put(record)
	}
	logger.recordMutex.Unlock()
}

// Raw logs a raw log to the logger.stdLogger.rawWriter.
func (logger *FastLogger) Raw(data []byte) {
	logger.rawMutex.Lock()
	logger.rawLogs = append(logger.rawLogs, data...)
	logger.rawMutex.Unlock()
}

// Info logs a message to the logger.stdLogger.infoWriter.

func (logger *FastLogger) Info(args ...interface{}) {
	logger.infoMutex.Lock()
	if logger.stdLogger.showDate {
		logger.infoLogs = append(logger.infoLogs, Now.Load().([]byte)...)
	}
	logger.infoLogs = addArgsToLog(logger.infoLogs, args...)
	logger.infoLogs = append(logger.infoLogs, '\n')
	logger.infoMutex.Unlock()
}

// FormatInfo logs a message with format to the logger.stdLogger.infoWriter.
func (logger *FastLogger) FormatInfo(f string, args ...interface{}) {
	logger.infoMutex.Lock()
	if logger.stdLogger.showDate {
		logger.infoLogs = append(logger.infoLogs, Now.Load().([]byte)...)
	}
	logger.infoLogs = append(logger.infoLogs, fmt.Sprintf(f, args...)...)
	logger.infoMutex.Unlock()
}

// InfoPrepare logs a prepared record to the logger.infoWriter. Will not reset the record.
func (logger *FastLogger) InfoPrepare(record *Record) {
	logger.infoMutex.Lock()
	if record.isShowDate {
		copy(record.rec[:19], Now.Load().([]byte))
	}
	logger.infoLogs = append(logger.infoLogs, record.rec...)
	logger.infoMutex.Unlock()
}

// Error logs a message to the logger.stdLogger.errorWriter.
func (logger *FastLogger) Error(args ...interface{}) {
	logger.errorMutex.Lock()
	if logger.stdLogger.showDate {
		logger.errorLogs = append(logger.errorLogs, Now.Load().([]byte)...)
	}
	logger.errorLogs = addArgsToLog(logger.errorLogs, args...)
	logger.errorLogs = append(logger.errorLogs, '\n')
	logger.errorMutex.Unlock()
}

// FormatError logs a message with format to the logger.stdLogger.errorWriter.
func (logger *FastLogger) FormatError(f string, args ...interface{}) {
	logger.errorMutex.Lock()
	if logger.stdLogger.showDate {
		logger.errorLogs = append(logger.errorLogs, Now.Load().([]byte)...)
	}
	logger.errorLogs = append(logger.errorLogs, fmt.Sprintf(f, args...)...)
	logger.errorMutex.Unlock()
}

// ErrorPrepare logs a prepared record to the logger.errorWriter. Will not reset the record.
func (logger *FastLogger) ErrorPrepare(record *Record) {
	logger.errorMutex.Lock()
	if record.isShowDate {
		copy(record.rec[:19], Now.Load().([]byte))
	}
	logger.errorLogs = append(logger.errorLogs, record.rec...)
	logger.errorMutex.Unlock()
}

// Warning logs a message to the logger.stdLogger.warningWriter.
func (logger *FastLogger) Warning(args ...interface{}) {
	logger.warningMutex.Lock()
	if logger.stdLogger.showDate {
		logger.warningLogs = append(logger.warningLogs, Now.Load().([]byte)...)
	}
	logger.warningLogs = addArgsToLog(logger.warningLogs, args...)
	logger.warningLogs = append(logger.warningLogs, '\n')
	logger.warningMutex.Unlock()
}

// FormatWarning logs a message with format to the logger.stdLogger.warningWriter.
func (logger *FastLogger) FormatWarning(f string, args ...interface{}) {
	logger.warningMutex.Lock()
	if logger.stdLogger.showDate {
		logger.warningLogs = append(logger.warningLogs, Now.Load().([]byte)...)
	}
	logger.warningLogs = append(logger.warningLogs, fmt.Sprintf(f, args...)...)
	logger.warningMutex.Unlock()
}

// WarningPrepare logs a prepared record to the logger.warningWriter. Will not reset the record.
func (logger *FastLogger) WarningPrepare(record *Record) {
	logger.warningMutex.Lock()
	if record.isShowDate {
		copy(record.rec[:19], Now.Load().([]byte))
	}
	logger.warningLogs = append(logger.warningLogs, record.rec...)
	logger.warningMutex.Unlock()
}

// Success logs a message to the logger.stdLogger.successWriter.
func (logger *FastLogger) Success(args ...interface{}) {
	logger.successMutex.Lock()
	if logger.stdLogger.showDate {
		logger.successLogs = append(logger.successLogs, Now.Load().([]byte)...)
	}
	logger.successLogs = addArgsToLog(logger.successLogs, args...)
	logger.successLogs = append(logger.successLogs, '\n')
	logger.successMutex.Unlock()
}

// FormatSuccess logs a message with format to the logger.stdLogger.successWriter.
func (logger *FastLogger) FormatSuccess(f string, args ...interface{}) {
	logger.successMutex.Lock()
	if logger.stdLogger.showDate {
		logger.successLogs = append(logger.successLogs, Now.Load().([]byte)...)
	}
	logger.successLogs = append(logger.successLogs, fmt.Sprintf(f, args...)...)
	logger.successMutex.Unlock()
}

// SuccessPrepare logs a prepared record to the logger.successWriter. Will not reset the record.
func (logger *FastLogger) SuccessPrepare(record *Record) {
	logger.successMutex.Lock()
	if record.isShowDate {
		copy(record.rec[:19], Now.Load().([]byte))
	}
	logger.successLogs = append(logger.successLogs, record.rec...)
	logger.successMutex.Unlock()
}

// Fatal logs a message to the logger.stdLogger.fatalWriter.
func (logger *FastLogger) Fatal(args ...interface{}) {
	logger.Flush()
	logger.stdLogger.Fatal(args...)
}

// FormatFatal logs a message with format to the logger.stdLogger.fatalWriter.
func (logger *FastLogger) FormatFatal(f string, args ...interface{}) {
	logger.Flush()
	logger.stdLogger.FormatFatal(f, args...)
}

// FatalPrepare logs a prepared record to the logger.fatalWriter. Will not reset the record.
func (logger *FastLogger) FatalPrepare(record *Record) {
	logger.Flush()
	logger.stdLogger.FatalPrepare(record)
}
