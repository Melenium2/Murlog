package murlog

const (
	// Formatted timestamp
	TimeTag    = "time"
	// Indicates a response delay
	LatencyTag = "latency"
	// Request method
	MethodTag  = "method"
	// Response code
	CodeTag    = "code"
	// Request uri
	PathTag    = "path"
	// Tag for default output
	// For example this is just text or error
	DefaultTag = "default"

	// CMD Colors
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

const (
	cBlack   = "\u001b[90m"
	cRed     = "\u001b[91m"
	cGreen   = "\u001b[92m"
	cYellow  = "\u001b[93m"
	cBlue    = "\u001b[94m"
	cMagenta = "\u001b[95m"
	cCyan    = "\u001b[96m"
	cWhite   = "\u001b[97m"
	cReset   = "\u001b[0m"
)