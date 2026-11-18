package realtime

func (e *Error) Error() string {
	if e.Message != "" {
		return "realtime: " + e.Code.String() + ": " + e.Message
	}
	return "realtime: " + e.Code.String()
}
