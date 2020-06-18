package murlog

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"
)

type Logger interface {
	Log(keyvals ...interface{}) error
	ErrorLog(keyvals ...interface{}) error
}

type murlogimpl struct {
	config *Config
}

func NewLogger(murfig *Config) Logger {
	return &murlogimpl{
		config: murfig,
	}
}

func (m murlogimpl) Log(keyvals ...interface{}) error {
	l := m.log(keyvals...)
	fmt.Fprintf(os.Stdout, "%s\n", "[Info]\t" + l)

	if m.config.iLogger {
		go m.sendInternal("Info", l)
	}

	return nil
}

func (m murlogimpl) ErrorLog(keyvals ...interface{}) error {
	e := m.log(keyvals...)
	fmt.Fprintf(os.Stderr, "%s\n", "[Error]\t" + e)

	if m.config.iLogger {
		go m.sendInternal("Error", e)
	}

	return nil
}

func (m murlogimpl) sendInternal(msgType, log string) {
	checkpoint := 0
	b := bytes.Buffer{}
	json.NewEncoder(&b).Encode(map[string]string{
		"message_type": msgType,
		"log":          log,
	})
restart:
	resp, err := http.Post(m.config.loggerUrl, "application/json", &b)
	if err != nil {
		checkpoint++
		time.Sleep(time.Second * 15)
		if checkpoint > 5 {
			m.log("error", fmt.Sprintf("internal server say %v", err))
			return
		}
		goto restart
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		checkpoint++
		time.Sleep(time.Second * 15)
		if checkpoint > 5 {
			m.log("error", fmt.Sprintf("internal server response with %v code", resp.StatusCode))
			return
		}
		goto restart
	}
}

func (m murlogimpl) constructPrefixes() string {
	log := ""
	for _, v := range m.config.prefixes {
		log += fmt.Sprintf("%v \t", v())
	}

	return log
}

func (m murlogimpl) log(keyvals ...interface{}) string {
	log := m.constructPrefixes()

	for i, s := range keyvals {
		sep := "="
		if (i+1)%2 == 0 {
			sep = " "
		}
		log += fmt.Sprintf("%s%s", s, sep)
	}

	return log
}
