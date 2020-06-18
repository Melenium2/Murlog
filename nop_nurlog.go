package murlog

type emptymurlogger struct {}

func NewNopLogger() Logger {
	return emptymurlogger{}
}

func (m emptymurlogger) Log(keyvals ...interface{}) error {
	return nil
}

func (m emptymurlogger) ErrorLog(keyvals ...interface{}) error {
	return nil
}
