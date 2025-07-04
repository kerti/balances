package failure

import "fmt"

// Failure is a wrapper for error messages and codes
type Failure struct {
	Code      Code    `json:"code"`
	Operation *string `json:"operation,omitempty"`
	Entity    *string `json:"entity,omitempty"`
	Message   string  `json:"message"`
}

// Error returns the error code and message in a formatted string
func (e *Failure) Error() string {
	if e.Operation != nil {
		if e.Entity != nil {
			return fmt.Sprintf("[%s] %s on %s: %s", e.Code, *e.Operation, *e.Entity, e.Message)
		} else {
			return fmt.Sprintf("[%s] %s: %s", e.Code, *e.Operation, e.Message)
		}
	} else {
		if e.Entity != nil {
			return fmt.Sprintf("[%s] %s: %s", e.Code, *e.Entity, e.Message)
		} else {
			return fmt.Sprintf("%s: %s", e.Code, e.Message)
		}
	}
}

// BadRequest returns a new Failure with code for bad requests
func BadRequest(err error) error {
	if err != nil {
		return &Failure{
			Code:    CodeBadRequest,
			Message: err.Error(),
		}
	}
	return nil
}

// BadRequestFromString returns a new Failure with code for bad requests with message set from string
func BadRequestFromString(msg string) error {
	return &Failure{
		Code:    CodeBadRequest,
		Message: msg,
	}
}

// Unauthorized returns a new Failure with code for unauthorized requests
func Unauthorized(msg string) error {
	return &Failure{
		Code:    CodeUnauthorized,
		Message: msg,
	}
}

// InternalError returns a new Failure with code for internal error and message derived from an error interface
func InternalError(operationName, entityName string, err error) error {
	if err != nil {
		return &Failure{
			Code:      CodeInternalError,
			Operation: &operationName,
			Entity:    &entityName,
			Message:   err.Error(),
		}
	}
	return nil
}

// Unimplemented returns a new Failure with code for unimplemented method
func Unimplemented(methodName string) error {
	return &Failure{
		Code:    CodeUnimplemented,
		Message: methodName,
	}
}

// EntityNotFound returns a new Failure with code for entity not found
func EntityNotFound(operationName, entityName string) error {
	return &Failure{
		Code:      CodeEntityNotFound,
		Operation: &operationName,
		Entity:    &entityName,
		Message:   "Record not found.",
	}
}

// OperationNotPermitted returns a new Failure with code for operation not permitted
func OperationNotPermitted(operationName, entityName string, message string) error {
	return &Failure{
		Code:      CodeOperationNotPermitted,
		Operation: &operationName,
		Entity:    &entityName,
		Message:   message,
	}
}

// GetCode returns the error code of an error interface
func GetCode(err error) Code {
	if f, ok := err.(*Failure); ok {
		return f.Code
	}
	return CodeInternalError
}

// Code represents a failure code string
type Code string

var (
	// CodeBadRequest is the string code for errors related to bad requests
	CodeBadRequest Code = "BadRequest"
	// CodeUnauthorized us the string code for unauthorized requests
	CodeUnauthorized Code = "Unauthorized"
	// CodeInternalError is the string code for internal errors
	CodeInternalError Code = "InternalError"
	// CodeUnimplemented is the string code for errors caused by unimplemented methods
	CodeUnimplemented Code = "Unimplemented"
	// CodeEntityNotFound is the string code for indicating an entity is not found
	CodeEntityNotFound Code = "EntityNotFound"
	// CodeOperationNotPermitted is the string code for indicating that an operation is not permitted
	CodeOperationNotPermitted Code = "OperationNotPermitted"
)
