package services

//ErrorNotFound can be use when a particular resource
// is not found
type ErrorNotFound struct {
	message string
}

func (e *ErrorNotFound) Error() string {
	if e.message == "" {
		return "Not Found"
	}
	return e.message
}

//New returns a new ErrorNotFound
func (ErrorNotFound) New(message string) *ErrorNotFound {
	return &ErrorNotFound{message}
}

//ErrorParsing can be use to create errors
//For cusome parsing
type ErrorParsing struct {
	message string
}

func (e *ErrorParsing) Error() string {
	if e.message == "" {
		return "Error Parsing"
	}
	return e.message
}

//New return a new ErrorParsing
func (ErrorParsing) New(message string) *ErrorParsing {
	return &ErrorParsing{message}
}
