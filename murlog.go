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

/*
	Create instance of logger interface.
 */
func NewLogger(murfig *Config) Logger {
	return &murlogimpl{
		config: murfig,
	}
}

/*
	Print values from params to os.stdout.
	If external logger was provide, func will sent log to the external logger.
 */
func (m murlogimpl) Log(keyvals ...interface{}) error {
	l := m.log(keyvals...)
	fmt.Fprintf(os.Stdout, "%s\n", "[Info]\t" + l)

	if m.config.eLogger {
		go m.sendInternal("Info", l)
	}

	return nil
}

/*
	Print values from params to os.stderr.
	If external logger was provide, func will sent log to the external logger.
 */
func (m murlogimpl) ErrorLog(keyvals ...interface{}) error {
	e := m.log(keyvals...)
	fmt.Fprintf(os.Stderr, "%s\n", "[Error]\t" + e)

	if m.config.eLogger {
		go m.sendInternal("Error", e)
	}

	return nil
}

/*
	Send log to external logger. Func will be restarted if connection problem will be represent.
 */
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

/*
	Func call prefix func and get his values as a string.
 */
func (m murlogimpl) constructPrefixes() string {
	log := ""
	for _, v := range m.config.prefixes {
		log += fmt.Sprintf("%v ", v())
	}
	return log
}

/*
	Func construct full log string and return her.
 */
func (m murlogimpl) log(keyvals ...interface{}) string {
	log := m.constructPrefixes()

	for i, s := range keyvals {
		sep := "="
		if (i+1)%2 == 0 {
			sep = " "
		}
		log += fmt.Sprintf("%v%s", s, sep)
	}

	return log
}

type emptymurlogger struct {}

// Empty logger instance for tests
func NewNopLogger() Logger {
	return emptymurlogger{}
}

func (m emptymurlogger) Log(keyvals ...interface{}) error {
	fmt.Fprintf(os.Stdout, "%v\n", keyvals)
	return nil
}

func (m emptymurlogger) ErrorLog(keyvals ...interface{}) error {
	fmt.Fprintf(os.Stdout, "%v\n", keyvals)
	return nil
}
