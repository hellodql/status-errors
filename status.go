package status_errors


func GetErrorType(err error) int {
	err = Unwrap(err)
	if customErr, ok := err.(*httpError); ok {
		errorType := customErr.errorType
		return int(errorType)
	}

	return int(UndefinedErr)
}
