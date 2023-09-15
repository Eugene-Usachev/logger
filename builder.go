package logger

import (
	"fmt"
	"github.com/Eugene-Usachev/fastbytes"
	"sync"
)

type Record struct {
	// prefix is the prefix of the log. nil by default.
	prefix []byte
	// isShowDate indicates will log have a date. By default, it's true
	isShowDate bool
	// isNewLine indicates will log create a new line ('\n'). By default, it's true
	isNewLine  bool
	rec        []byte
	wasGot     bool
	wasPrepare bool
}

func NewBuilder() *Record {
	return &Record{
		rec:        make([]byte, 0),
		isShowDate: true,
		isNewLine:  true,
		wasGot:     false,
	}
}

var recordPool = sync.Pool{
	New: func() any {
		return NewBuilder()
	},
}

func Builder() *Record {
	return recordPool.Get().(*Record)
}

func (r *Record) Prefix(prefix string) *Record {
	r.prefix = fastbytes.S2B(prefix)
	return r
}

func (r *Record) ShowDate() *Record {
	r.isShowDate = true
	return r
}

func (r *Record) NoDate() *Record {
	r.isShowDate = false
	return r
}

func (r *Record) AppendArgs(args ...interface{}) *Record {
	r.rec = addArgsToLog(r.rec, args...)
	return r
}

func (r *Record) AppendFormat(format string, args ...interface{}) *Record {
	r.rec = append(r.rec, fmt.Sprintf(format, args...)...)
	return r
}

func (r *Record) NewLine(flag bool) *Record {
	r.isNewLine = flag
	return r
}

func (r *Record) Build() *Record {
	if r.prefix != nil {
		r.rec = append(r.prefix, r.rec...)
	}
	if r.isShowDate {
		r.rec = append(Now.Load().([]byte), r.rec...)
	}
	if r.isNewLine {
		r.rec = append(r.rec, '\n')
	}
	r.wasPrepare = true
	return r
}

// Prepare a record to use with Prepare functions.
func (r *Record) Prepare() *Record {
	if r.wasPrepare {
		return r
	}
	if r.prefix != nil {
		r.rec = append(r.prefix, r.rec...)
	}
	if r.isShowDate {
		r.rec = append(Now.Load().([]byte), r.rec...)
	}
	if r.isNewLine {
		r.rec = append(r.rec, '\n')
	}
	r.wasPrepare = true
	return r
}

func (r *Record) Reset() {
	r.rec = r.rec[:0]
	r.isShowDate = true
	r.isNewLine = true
	r.wasPrepare = false
}
