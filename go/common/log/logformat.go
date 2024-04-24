package log

import (
	"sync"
	"sync/atomic"
)

// IMPORTANT - the utilities and constants in this format class are copied from the geth log package
// (go-ethereum/log/format.go) and modified to suit our needs. Unfortunately geth log format is not configurable
// and implementation is all private so had to duplicate a lot of it.
//
// The code was modified in places to remove things we don't need, e.g. color, but other than that it's identical.

const (
	timeFormat        = "2006-01-02T15:04:05-0700"
	termTimeFormat    = "01-02|15:04:05.000"
	floatFormat       = 'f'
	termMsgJust       = 40
	termCtxMaxPadding = 40
)

// locationTrims are trimmed for display to avoid unwieldy log lines.
var locationTrims = []string{
	"github.com/ethereum/go-ethereum/",
	"github.com/ten-protocol/go-ten/",
}

// locationEnabled is an atomic flag controlling whether the terminal formatter
// should append the log locations too when printing entries.
var locationEnabled atomic.Bool

// locationLength is the maxmimum path length encountered, which all logs are
// padded to aid in alignment.
var locationLength atomic.Uint32

// fieldPadding is a global map with maximum field value lengths seen until now
// to allow padding log contexts in a bit smarter way.
var fieldPadding = make(map[string]int)

// fieldPaddingLock is a global mutex protecting the field padding map.
var fieldPaddingLock sync.RWMutex

/*
// TenLogFormat - returns a log format that is used by the Ten logger for both console and file logging.
// Note: this is mostly a copy of gethlog.TerminalFormat but putting it here gives us control and
// means we aren't forced to use the shortened hash format
func TenLogFormat() gethlog.Format {
	return gethlog.FormatFunc(func(r *gethlog.Record) []byte {
		msg := escapeMessage(r.Msg)

		b := &bytes.Buffer{}
		lvl := r.Lvl.AlignedString()
		if locationEnabled.Load() {
			// Log origin printing was requested, format the location path and line number
			location := fmt.Sprintf("%+v", r.Call)
			for _, prefix := range locationTrims {
				location = strings.TrimPrefix(location, prefix)
			}
			// Maintain the maximum location length for fancyer alignment
			align := int(locationLength.Load())
			if align < len(location) {
				align = len(location)
				locationLength.Store(uint32(align))
			}
			padding := strings.Repeat(" ", align-len(location))

			// Assemble and print the log heading
			fmt.Fprintf(b, "%s[%s|%s]%s %s ", lvl, r.Time.Format(termTimeFormat), location, padding, msg)
		} else {
			fmt.Fprintf(b, "%s[%s] %s ", lvl, r.Time.Format(termTimeFormat), msg)
		}
		// try to justify the log output for short messages
		length := utf8.RuneCountInString(msg)
		if len(r.Ctx) > 0 && length < termMsgJust {
			b.Write(bytes.Repeat([]byte{' '}, termMsgJust-length))
		}
		// print the keys logfmt style
		logfmt(b, r.Ctx)
		return b.Bytes()
	})
}

func logfmt(buf *bytes.Buffer, ctx []interface{}) {
	for i := 0; i < len(ctx); i += 2 {
		if i != 0 {
			buf.WriteByte(' ')
		}

		k, ok := ctx[i].(string)
		v := formatLogfmtValue(ctx[i+1])
		if !ok {
			k, v = ErrKey, fmt.Sprintf("%+T is not a string key", ctx[i])
		} else {
			k = escapeString(k)
		}

		// XXX: we should probably check that all of your key bytes aren't invalid
		fieldPaddingLock.RLock()
		padding := fieldPadding[k]
		fieldPaddingLock.RUnlock()

		length := utf8.RuneCountInString(v)
		if padding < length && length <= termCtxMaxPadding {
			padding = length

			fieldPaddingLock.Lock()
			fieldPadding[k] = padding
			fieldPaddingLock.Unlock()
		}
		buf.WriteString(k)
		buf.WriteByte('=')
		buf.WriteString(v)
		if i < len(ctx)-2 && padding > length {
			buf.Write(bytes.Repeat([]byte{' '}, padding-length))
		}
	}
	buf.WriteByte('\n')
}

func formatShared(value interface{}) (result interface{}) {
	defer func() {
		if err := recover(); err != nil {
			if v := reflect.ValueOf(value); v.Kind() == reflect.Ptr && v.IsNil() {
				result = "nil"
			} else {
				panic(err)
			}
		}
	}()

	switch v := value.(type) {
	case time.Time:
		return v.Format(timeFormat)

	case error:
		return v.Error()

	case fmt.Stringer:
		return v.String()

	default:
		return v
	}
}

// formatValue formats a value for serialization
func formatLogfmtValue(value interface{}) string {
	if value == nil {
		return "nil"
	}

	switch v := value.(type) {
	case time.Time:
		// Performance optimization: No need for escaping since the provided
		// timeFormat doesn't have any escape characters, and escaping is
		// expensive.
		return v.Format(timeFormat)

	case *big.Int:
		// Big ints get consumed by the Stringer clause, so we need to handle
		// them earlier on.
		if v == nil {
			return "<nil>"
		}
		return formatLogfmtBigInt(v)

	case *uint256.Int:
		// Uint256s get consumed by the Stringer clause, so we need to handle
		// them earlier on.
		if v == nil {
			return "<nil>"
		}
		return formatLogfmtUint256(v)
	}

	value = formatShared(value)
	switch v := value.(type) {
	case bool:
		return strconv.FormatBool(v)
	case float32:
		return strconv.FormatFloat(float64(v), floatFormat, 3, 64)
	case float64:
		return strconv.FormatFloat(v, floatFormat, 3, 64)
	case int8:
		return strconv.FormatInt(int64(v), 10)
	case uint8:
		return strconv.FormatInt(int64(v), 10)
	case int16:
		return strconv.FormatInt(int64(v), 10)
	case uint16:
		return strconv.FormatInt(int64(v), 10)
	// Larger integers get thousands separators.
	case int:
		return gethlog.FormatLogfmtInt64(int64(v))
	case int32:
		return gethlog.FormatLogfmtInt64(int64(v))
	case int64:
		return gethlog.FormatLogfmtInt64(v)
	case uint:
		return gethlog.FormatLogfmtUint64(uint64(v))
	case uint32:
		return gethlog.FormatLogfmtUint64(uint64(v))
	case uint64:
		return gethlog.FormatLogfmtUint64(v)
	case string:
		return escapeString(v)
	default:
		return escapeString(fmt.Sprintf("%+v", value))
	}
}

// formatLogfmtBigInt formats n with thousand separators.
func formatLogfmtBigInt(n *big.Int) string {
	if n.IsUint64() {
		return gethlog.FormatLogfmtUint64(n.Uint64())
	}
	if n.IsInt64() {
		return gethlog.FormatLogfmtInt64(n.Int64())
	}

	var (
		text  = n.String()
		buf   = make([]byte, len(text)+len(text)/3)
		comma = 0
		i     = len(buf) - 1
	)
	for j := len(text) - 1; j >= 0; j, i = j-1, i-1 {
		c := text[j]

		switch {
		case c == '-':
			buf[i] = c
		case comma == 3:
			buf[i] = ','
			i--
			comma = 0
			fallthrough
		default:
			buf[i] = c
			comma++
		}
	}
	return string(buf[i+1:])
}

// formatLogfmtUint256 formats n with thousand separators.
func formatLogfmtUint256(n *uint256.Int) string {
	if n.IsUint64() {
		return gethlog.FormatLogfmtUint64(n.Uint64())
	}
	var (
		text  = n.Dec()
		buf   = make([]byte, len(text)+len(text)/3)
		comma = 0
		i     = len(buf) - 1
	)
	for j := len(text) - 1; j >= 0; j, i = j-1, i-1 {
		c := text[j]

		switch {
		case c == '-':
			buf[i] = c
		case comma == 3:
			buf[i] = ','
			i--
			comma = 0
			fallthrough
		default:
			buf[i] = c
			comma++
		}
	}
	return string(buf[i+1:])
}

// escapeString checks if the provided string needs escaping/quoting, and
// calls strconv.Quote if needed
func escapeString(s string) string {
	needsQuoting := false
	for _, r := range s {
		// We quote everything below " (0x22) and above~ (0x7E), plus equal-sign
		if r <= '"' || r > '~' || r == '=' {
			needsQuoting = true
			break
		}
	}
	if !needsQuoting {
		return s
	}
	return strconv.Quote(s)
}

// escapeMessage checks if the provided string needs escaping/quoting, similarly
// to escapeString. The difference is that this method is more lenient: it allows
// for spaces and linebreaks to occur without needing quoting.
func escapeMessage(s string) string {
	needsQuoting := false
	for _, r := range s {
		// Allow CR/LF/TAB. This is to make multi-line messages work.
		if r == '\r' || r == '\n' || r == '\t' {
			continue
		}
		// We quote everything below <space> (0x20) and above~ (0x7E),
		// plus equal-sign
		if r < ' ' || r > '~' || r == '=' {
			needsQuoting = true
			break
		}
	}
	if !needsQuoting {
		return s
	}
	return strconv.Quote(s)
}
*/
