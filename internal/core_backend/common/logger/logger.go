package logger

import (
	"encoding/json"
	"log"
	"os"

	"github.com/fatih/color"
)

// RequestErrorLog Error with x-request-id
type RequestErrorLog struct {
	XRequestID   string `json:"x-request-id"`
	ErrorMessage string `json:"error_message"`
}

// OutputLog Output log with Level
type OutputLog struct {
	Level   string `json:"level"`
	Message string `json:"message"`
}

// LogError Write ERROR log stdout
func LogError(message string) {
	errLog := OutputLog{
		Message: message,
		Level:   "ERROR",
	}
	output, _ := json.Marshal(errLog)
	log.SetOutput(os.Stderr)
	color.Red(string(output))
}

// LogSuccess Write SUCCESS log stdout
func LogSuccess(message string) {
	successLog := OutputLog{
		Message: message,
		Level:   "SUCCESS",
	}
	output, _ := json.Marshal(successLog)
	log.SetOutput(os.Stdout)
	color.HiGreen(string(output))
}

// LogInfo Write INFO log stdout
func LogInfo(message string) {
	infoLog := OutputLog{
		Message: message,
		Level:   "INFO",
	}
	output, _ := json.Marshal(infoLog)
	log.SetOutput(os.Stdout)
	color.Yellow(string(output))
}

// LogRequestError logs error with x-request-id
func LogRequestError(xRequestID, message string) {
	var errlog RequestErrorLog
	errlog.ErrorMessage = message
	errlog.XRequestID = xRequestID
	output, _ := json.Marshal(errlog)
	log.SetOutput(os.Stderr)
	color.Red(string(output))
}
