package entity

// ProcessStatus is the status of the process
type ProcessStatus string

var (
	ProcessStatus_NEW        ProcessStatus = "New"
	ProcessStatus_PROCESSING ProcessStatus = "Processing"
	ProcessStatus_SUCCESS    ProcessStatus = "Success"
	ProcessStatus_FAIL       ProcessStatus = "Fail"
)
