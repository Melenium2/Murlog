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
	eLogger   bool
	loggerUrl string
}

// Create new config fo logger. eLogger means 'ExternalLogger', default false
func NewConfig() *Config {
	return &Config{
		eLogger: false,
	}
}

/*
	Add custom prefix to logger output. All the prefixes will be displayed
	sequentially.
 */
func (c *Config) Pref(prefix Prefix) {
	c.prefixes = append(c.prefixes, prefix)
}

/*
	Add external logger url. Additionally all logs will be sent to that url
 */
func (c *Config) InternalLogStore(url string) {
	c.eLogger = true
	c.loggerUrl = url
}

/*
	Timestamp prefix to the log. Need to provide format of ts for the log.
	For example: '2006.01.02 15-04-05'
 */
func (c *Config) TimePref(format string) {
	c.Pref(func() interface{} {
		return fmt.Sprintf("ts=%s", time.Now().UTC().Format(format))
	})
}

/*
	Prefix will add caller func to the log from where the log was called.
 */
func (c *Config) CallerPref() {
	c.CallerCustomPref(4)
}

/*
	Prefix will add caller func to the log from where the log was called.
	You can provide the custom number of func's skip.
*/
func (c *Config) CallerCustomPref(n int) {
	c.Pref(func() interface{} {
		_, file, row, _ := runtime.Caller(n)
		idx := strings.LastIndexByte(file, '/')
		return fmt.Sprintf("caller=%s:%d", file[idx+1:], row)
	})
}
