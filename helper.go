package logger

import (
	"strconv"
	"sync/atomic"
)

var Now = func() atomic.Value {
	var v atomic.Value
	v.Store(make([]byte, 0))
	return v
}()

func addArgsToLog(buf []byte, args ...interface{}) []byte {
	for i := 0; i < len(args); i++ {
		switch args[i].(type) {
		case string:
			buf = append(buf, args[i].(string)...)
			break
		case int8:
			buf = strconv.AppendInt(buf, int64(args[i].(int8)), 10)
			break
		case int16:
			buf = strconv.AppendInt(buf, int64(args[i].(int16)), 10)
			break
		case int32:
			buf = strconv.AppendInt(buf, int64(args[i].(int32)), 10)
			break
		case int:
			buf = strconv.AppendInt(buf, int64(args[i].(int)), 10)
			break
		case int64:
			buf = strconv.AppendInt(buf, args[i].(int64), 10)
			break
		case uint:
			buf = strconv.AppendUint(buf, uint64(args[i].(uint)), 10)
			break
		case uint8:
			buf = strconv.AppendUint(buf, uint64(args[i].(uint8)), 10)
			break
		case uint16:
			buf = strconv.AppendUint(buf, uint64(args[i].(uint16)), 10)
			break
		case uint32:
			buf = strconv.AppendUint(buf, uint64(args[i].(uint32)), 10)
			break
		case uint64:
			buf = strconv.AppendUint(buf, args[i].(uint64), 10)
			break
		case float32:
			buf = strconv.AppendFloat(buf, float64(args[i].(float32)), 'f', -1, 64)
			break
		case float64:
			buf = strconv.AppendFloat(buf, args[i].(float64), 'f', -1, 64)
			break
		case bool:
			buf = strconv.AppendBool(buf, args[i].(bool))
			break
		case complex64:
			buf = append(buf, strconv.FormatComplex(complex128(args[i].(complex64)), 'f', -1, 64)...)
			break
		case complex128:
			buf = append(buf, strconv.FormatComplex(args[i].(complex128), 'f', -1, 64)...)
			break
		case []byte:
			buf = append(buf, args[i].([]byte)...)
			break
		case uintptr:
			buf = strconv.AppendUint(buf, uint64(args[i].(uintptr)), 10)
			break
		default:
			buf = append(buf, args[i].(string)...)
			break
		}

	}
	return buf
}
