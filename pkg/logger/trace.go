package logger

import "github.com/google/uuid"

const TraceID = "TraceID"

// generateTraceID 生成一個新的TraceID
func GenerateTraceID() string {
	return uuid.New().String()
}
