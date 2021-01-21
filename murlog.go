package murlog

import (
	"fmt"
	"github.com/Melenium2/Murlog/fasttemplate"
	"io"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

type LogFunc func(v string, kv ...map[string]string)
type LogMiddlewareFunc func(next http.Handler) http.Handler

func New(config ...Config) LogFunc {
	cfg := defaultConfig(config...)

	var timestamp atomic.Value
	timestamp.Store(time.Now().In(cfg.timeZoneLocation).Format(cfg.TimeFormat))

	if strings.Contains(cfg.Format, "${time}") {
		go func() {
			time.Sleep(cfg.TimeInterval)
			timestamp.Store(time.Now().In(cfg.timeZoneLocation).Format(cfg.TimeFormat))
		}()
	}

	tmpl := fasttemplate.New(cfg.Format, "${", "}")

	var (
		mu sync.Mutex
	)

	return func(v string, kv ...map[string]string) {
		var (
			err       error
			keyValues map[string]string
		)
		if len(kv) > 0 {
			keyValues = kv[0]
		}

		buf := fasttemplate.Get()

		if cfg.enableColors {
			_, _ = buf.WriteString(fmt.Sprintf("%s - %s\n",
				timestamp.Load().(string),
				v,
			))

			_, _ = cfg.Output.Write(buf.Bytes())

			fasttemplate.Put(buf)

			return
		}

		{
			if keyValues == nil {
				keyValues = make(map[string]string)
			}
			keyValues["default"] = v
			keyValues["time"] = timestamp.Load().(string)
		}

		err = fillBuffer(tmpl, buf, keyValues)

		if err != nil {
			_, _ = buf.WriteString(err.Error())
		}

		mu.Lock()

		if _, err := cfg.Output.Write(buf.Bytes()); err != nil {
			if _, err := cfg.Output.Write([]byte(err.Error())); err != nil {
				// ??? =)
			}
		}

		mu.Unlock()

		fasttemplate.Put(buf)
	}
}

func NewMiddleware(config ...Config) LogMiddlewareFunc {
	defaultFormat := "${red}[${time}] ${cyan}${method} ${path} - ${magenta}${code}${reset} ${latency} ${default}\n"
	var c Config
	if len(config) == 0 {
		c.Format = defaultFormat
	} else if config[0].Format == "" {
		c = config[0]
		c.Format = defaultFormat
	}

	log := New(c)

	return func(next http.Handler) http.Handler {
		var (
			start, end time.Time
		)

		fn := func(w http.ResponseWriter, r *http.Request) {
			wrappedHeader := WrapResponseWriter(w)
			var kvs = map[string]string{
				"path":   r.URL.EscapedPath(),
				"method": r.Method,
			}

			defer func() {
				if err := recover(); err != nil {
					wrappedHeader.WriteHeader(http.StatusInternalServerError)
					kvs["code"] = strconv.Itoa(wrappedHeader.Status)

					log(err.(string), kvs)
				}
			}()

			{
				start = time.Now()
				next.ServeHTTP(w, r)
				end = time.Now()

				kvs["latency"] = end.Sub(start).String()
				kvs["code"] = strconv.Itoa(wrappedHeader.Status)
			}
			log("", kvs)
		}

		return http.HandlerFunc(fn)
	}
}

func fillBuffer(tmpl *fasttemplate.Template, buf *fasttemplate.ByteBuffer, kv map[string]string) error {
	_, err := tmpl.ExecuteFunc(buf, func(w io.Writer, tag string) (int, error) {
		switch tag {
		case TimeTag:
			return buf.WriteString(kv["time"])
		case DefaultTag:
			return buf.WriteString(kv["default"])
		case LatencyTag:
			return buf.WriteString(kv["latency"])
		case MethodTag:
			return buf.WriteString(kv["method"])
		case CodeTag:
			return buf.WriteString(kv["code"])
		case PathTag:
			return buf.WriteString(kv["path"])
		case Black:
			return buf.WriteString(cBlack)
		case Red:
			return buf.WriteString(cRed)
		case Green:
			return buf.WriteString(cGreen)
		case Yellow:
			return buf.WriteString(cYellow)
		case Blue:
			return buf.WriteString(cBlue)
		case Magenta:
			return buf.WriteString(cMagenta)
		case Cyan:
			return buf.WriteString(cCyan)
		case White:
			return buf.WriteString(cWhite)
		case Reset:
			return buf.WriteString(cReset)
		default:
			return 0, nil
		}
	})

	return err
}
