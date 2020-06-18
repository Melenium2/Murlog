package murlog

import (
	"fmt"
	"runtime"
	"strings"
	"time"
)

type Prefix func() interface{}

type Config struct {
	prefixes  []Prefix
	iLogger   bool
	loggerUrl string
}

func NewConfig() *Config {
	return &Config{
		iLogger: false,
	}
}

func (c *Config) Pref(prefix Prefix) {
	c.prefixes = append(c.prefixes, prefix)
}

func (c *Config) InternalLogStore(url string) {
	c.iLogger = true
	c.loggerUrl = url
}

func (c *Config) TimePref(format string) {
	c.Pref(func() interface{} {
		return fmt.Sprintf("ts=%s", time.Now().UTC().Format(format))
	})
}

func (c *Config) CallerPref() {
	c.Pref(func() interface{} {
		_, file, row, _ := runtime.Caller(4)
		idx := strings.LastIndexByte(file, '/')
		return fmt.Sprintf("caller=%s:%d", file[idx+1:], row)
	})
}
