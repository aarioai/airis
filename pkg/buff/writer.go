package buff

import (
	"github.com/aarioai/airis/core/atype"
	"strconv"
	"strings"
)

// StringWithBuffer 使用提供的 buffer 转换字符串，减少内存分配
// var buf strings.Builder
// atype.StringWithBuffer(123, &buf)
// str = buf.String()  // "123"
func WriteInt(d any, buf *strings.Builder) {
	if d == nil {
		return
	}
	switch v := d.(type) {
	case string:
		buf.WriteString(v)
	case []byte:
		buf.Write(v)
	case atype.AByte:
		buf.WriteByte(byte(v))
	case atype.Date:
		buf.WriteString(string(v))
	case atype.Datetime:
		buf.WriteString(string(v))
	case bool:
		if v {
			buf.WriteString("true")
		} else {
			buf.WriteString("false")
		}
	case int:
		writeIntToBuffer(int64(v), buf)
	case int64:
		writeIntToBuffer(v, buf)
	case uint:
		writeUintToBuffer(uint64(v), buf)
	case uint64:
		writeUintToBuffer(v, buf)
	default:
		buf.WriteString(atype.String(v))
	}
}

// writeIntToBuffer 写入整数到 buffer
func writeIntToBuffer(v int64, buf *strings.Builder) {
	if v >= 0 && v < 100 {
		if v < 10 {
			buf.WriteByte(byte(v + '0'))
			return
		}
		buf.WriteByte(byte(v/10 + '0'))
		buf.WriteByte(byte(v%10 + '0'))
		return
	}
	buf.WriteString(strconv.FormatInt(v, 10))
}

// writeUintToBuffer 写入无符号整数到 buffer
func writeUintToBuffer(v uint64, buf *strings.Builder) {
	if v < 100 {
		if v < 10 {
			buf.WriteByte(byte(v + '0'))
			return
		}
		buf.WriteByte(byte(v/10 + '0'))
		buf.WriteByte(byte(v%10 + '0'))
		return
	}
	buf.WriteString(strconv.FormatUint(v, 10))
}
