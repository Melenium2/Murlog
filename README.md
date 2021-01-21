# Murlog

Simple logger (logger middleware) for Golang.

## Installation

```shell
go get -u github.com/Melenium2/Murlog
```

## Signatures

```go
func New(config ...Config) LogFunc
```

or for middleware

```go
func NewMiddleware(config ...Config) LogMiddlewareFunc
```

## Config

Just call

```go
New()
```

or

```go
NewMiddleware()
```

### Changing Format of output

```go
murlog.New(murlog.Config{
    // For more option, see Default config section
    Format: "[${time}] --- ${default} =>\n",
})
```

### Changing TimeZone & TimeFormat

```go
murlog.New(murlog.Config{
    // For more option, see Default config section
    Format: "[${time}] - ${default}\n",
    TimeFormat: "02/Jan/2006",
    TimeZone: "Asia/Shanghai",
})
```

### Changing Output

For example, you can save logs to the file

```go
file, _ := os.OpenFile("./mylogs.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0777)
defer file.Close()

murlog.New(murlog.Config{
    Output: file,
})
```

For `NewMiddleware` all options the same

### Default Config

```go
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
}
```

```go
var DefaultConfig = Config{
	Format:       "${red}[${time}] ${white}- ${default}\n",
	TimeZone:     "Local",
	TimeFormat:   "15:04:05",
	TimeInterval: time.Millisecond * 500,
	Output:       os.Stderr,
}
```

### Logger variables
You can specify your own format of log with special variables.

All variables must be of the form `${variable}`


Example:
```go
config := Config{
    Format: "${red}[${time}] ${cyan}${method} ${path} - ${magenta}${code}${reset} ${latency} ${default}",       
}
```
and result

```shell
[16:00:14] GET /test-check - 500  Error message
```
The output will be colored if your CMD supports it 

### Variables
```go
const (
	TimeTag    = "time"
	LatencyTag = "latency"
	MethodTag  = "method"
	CodeTag    = "code"
	PathTag    = "path"
	DefaultTag = "default"

	Black   = "black"
	Red     = "red"
	Green   = "green"
	Yellow  = "yellow"
	Blue    = "blue"
	Magenta = "magenta"
	Cyan    = "cyan"
	White   = "white"
	Reset   = "reset"
)
```


