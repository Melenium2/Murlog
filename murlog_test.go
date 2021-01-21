package murlog_test

import (
	murlog "github.com/Melenium2/Murlog"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestNewLogger_Log_ShouldPrintTextWithFormat(t *testing.T) {
	l := murlog.New(murlog.Config{
		Format: "${cyan}msg = ${white}${default}\n",
	})

	 l("privet sosed")
}

func TestNewLogger_Log_ShouldPrintDefaultText(t *testing.T) {
	l := murlog.New()

	l("privet sosed")
}

func TestNewLogger_Log_ShouldPrintTextIfRequestComplete(t *testing.T) {
	log := murlog.NewMiddleware(murlog.Config{
		Format: "${red}${time} ${cyan}${method} ${path} ${magenta}${code}${reset} ${latency} ${default}\n",
	})

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
	})
	req, err := http.NewRequest("GET", "/test-check", nil)
	assert.NoError(t, err)

	res := httptest.NewRecorder()

	log(handler).ServeHTTP(res, req)
}

func TestNewLogger_Log_ShouldPrintTextIfRequestCompleteWithoutConfig(t *testing.T) {
	log := murlog.NewMiddleware()

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
	})
	req, err := http.NewRequest("GET", "/test-check", nil)
	assert.NoError(t, err)

	res := httptest.NewRecorder()

	log(handler).ServeHTTP(res, req)
}

