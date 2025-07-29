package dtos

// Defining the query parameters for the `Logger` struct.
type LoggerListReq struct {
	Pager
	Level     string `form:"level"`     // log level
	TraceID   string `form:"traceID"`   // trace ID
	UserName  string `form:"username"`  // user name
	Tag       string `form:"tag"`       // log tag
	Message   string `form:"message"`   // log message
	StartTime string `form:"startTime"` // start time
	EndTime   string `form:"endTime"`   // end time
}
