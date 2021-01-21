package murlog

import (
	"github.com/Melenium2/Murlog/colorable"
	"github.com/mattn/go-isatty"
	"io"
	"os"
	"strings"
	"time"
)

type Config struct {
	// Format define the string with logging tags
	// Optional.
	// Default: ${red}[${time}] ${white}- ${default}
	Format string
	// TimeZone can be specified, such as "UTC" or "Asia/Shanghai", etc
	// Optional.
	// Default: Local
	TimeZone string
	// // TimeFormat
	//	Optional.
	//	Default: 15:04:05
	TimeFormat   string
	// TimeInterval is the delay before the timestamp is updated
	// Optional.
	// Default: 500 * time.Millisecond
	TimeInterval time.Duration
	// Output is a writter where logs are written
	// Optional.
	// Default: os.Stderr
	Output       io.Writer

	enableColors     bool
	enableLatency    bool
	timeZoneLocation *time.Location
}

var DefaultConfig = Config{
	Format:       "${red}[${time}] ${white}- ${default}\n",
	TimeZone:     "Local",
	TimeFormat:   "15:04:05",
	TimeInterval: time.Millisecond * 500,
	Output:       os.Stderr,

	enableColors: true,
}

func defaultConfig(conf ...Config) Config {
	var c Config
	if len(conf) == 0 {
		c = DefaultConfig
	} else {
		c = conf[0]
		if c.Format == "" {
			c.Format = DefaultConfig.Format
		}
		if c.Output == nil {
			c.Output = DefaultConfig.Output
		}
		if c.TimeFormat == "" {
			c.TimeFormat = DefaultConfig.TimeFormat
		}
		if c.TimeZone == "" {
			c.TimeZone = DefaultConfig.TimeZone
		}
		if int(c.TimeInterval) <= 0 {
			c.TimeInterval = DefaultConfig.TimeInterval
		}
		if c.Format == "" || c.Output == nil {
			c.enableColors = true
		}
	}

	tz, err := time.LoadLocation(c.TimeZone)
	if err != nil || tz == nil {
		c.timeZoneLocation = time.Local
	} else {
		c.timeZoneLocation = tz
	}

	c.enableLatency = strings.Contains(c.Format, "${latency}")

	if c.enableColors {
		c.Output = colorable.NewColorableStderr()
		if !isatty.IsTerminal(os.Stderr.Fd()) && !isatty.IsCygwinTerminal(os.Stderr.Fd()) {
			c.Output = colorable.NewNonColorable(os.Stderr)
		}
	}

	return c
}
