package murlog

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestNewLogger_Log_ShouldReturnOnlyOneCustomPrefix(t *testing.T) {
	c := NewConfig()
	c.Pref(func() interface{} {
		return "[Info]"
	})
	l := NewLogger(c)

	assert.NoError(t, l.Log())
}

func TestNewLogger_Log_ShouldReturnCallerPrefixAndTimestampPrefix(t *testing.T) {
	c := NewConfig()
	c.TimePref(time.ANSIC)
	c.CallerPref()
	l := NewLogger(c)

	assert.NoError(t, l.Log())
}

func TestNewLogger_Log_ShouldReturnCallerPrefixAndTimestampPrefixAndCustomPrefix(t *testing.T) {
	c := NewConfig()
	c.TimePref(time.ANSIC)
	c.CallerPref()
	c.Pref(func() interface{} {
		return "[Error]"
	})
	l := NewLogger(c)

	assert.NoError(t, l.Log())
}

func TestNewLogger_Log_ShouldReturnDefaultPrefixesAndMessage(t *testing.T) {
	c := NewConfig()
	c.TimePref(time.ANSIC)
	c.CallerPref()
	c.Pref(func() interface{} {
		return "service=guard"
	})
	l := NewLogger(c)

	assert.NoError(t, l.Log("action", "start"))
}

func TestNewLogger_Log_ShouldReturnErrorMessageWithDefaultPrefixes(t *testing.T) {
	c := NewConfig()
	c.TimePref(time.ANSIC)
	c.CallerPref()
	c.Pref(func() interface{} {
		return "service=guard"
	})
	l := NewLogger(c)

	assert.NoError(t, l.ErrorLog("error", "can not start"))
}


