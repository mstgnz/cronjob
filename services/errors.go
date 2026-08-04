package services

import (
	"github.com/mstgnz/cronjob/pkg/logger"
	"github.com/mstgnz/cronjob/pkg/response"
)

// serverError records the real failure and hands the caller a generic message.
// Driver errors carry table names, column names and fragments of SQL, which is
// both an information leak and an oracle for anyone probing the API.
func serverError(err error) response.Response {
	logger.Warn("Service Error", err.Error())
	return response.Response{Status: false, Message: "An unexpected error occurred"}
}

// badRequestError is the same idea for malformed input: the caller learns the
// request was rejected, the details go to the log.
func badRequestError(err error) response.Response {
	logger.Info("Bad Request", err.Error())
	return response.Response{Status: false, Message: "Invalid request body"}
}
