package logger

type LogLevel int

const (
	All LogLevel = iota
	NoRequests
	NoProxy
	NoInfo
	None
)