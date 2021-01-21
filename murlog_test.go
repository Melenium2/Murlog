package murlog_test

import (
	murlog "github.com/Melenium2/Murlog"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"
)

func TestNewLogger_Default(t *testing.T) {
	var tt = []struct {
		name    string
		config  murlog.Config
		message string
	}{
		{
			name:    "contains default formatting and simple message",
			config:  murlog.Config{},
			message: "simple message",
		},
		{
			name: "contains multicolor output with simple message",
			config: murlog.Config{
				Format: "${cyan}msg = ${white}${default}\n",
			},
			message: "simple message",
		},
	}

	for _, test := range tt {
		t.Run(test.name, func(t *testing.T) {
			rescueStdout := os.Stdout
			r, w, _ := os.Pipe()
			os.Stdout = w

			test.config.Output = os.Stdout
			l := murlog.New(test.config)

			l(test.message)

			_ = w.Close()
			out, _ := ioutil.ReadAll(r)
			os.Stdout = rescueStdout

			assert.Contains(t, string(out), test.message)
		})
	}
}

func TestNewLogger_Middleware_ShouldPrintTextIfRequestComplete(t *testing.T) {
	log := murlog.NewMiddleware(murlog.Config{
		Format: "${red}${time} ${cyan}${method} ${path} ${magenta}${code}${reset} ${latency} ${default}\n",
	})

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(time.Millisecond * 500)
		w.WriteHeader(200)
	})
	req, err := http.NewRequest("GET", "/test-check", nil)
	assert.NoError(t, err)

	res := httptest.NewRecorder()

	log(handler).ServeHTTP(res, req)
}

func TestNewLogger_Log_ShouldPrintTextIfRequestCompleteWithDefaultConfig(t *testing.T) {
	log := murlog.NewMiddleware()

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(time.Millisecond * 200)
		panic("Error message")
	})
	req, err := http.NewRequest("GET", "/test-check", nil)
	assert.NoError(t, err)

	res := httptest.NewRecorder()

	log(handler).ServeHTTP(res, req)
}
