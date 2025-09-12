package logger

type LogLevel uint8

const (
	All LogLevel = iota
	NoRequests
	NoProxy
	NoInfo
	None
)
