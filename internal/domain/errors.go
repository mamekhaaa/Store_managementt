package domain

type ErrorCode string

const (
	CodeValidation   ErrorCode = "validation_error"
	CodeNotFound     ErrorCode = "not_found"
	CodeForbidden    ErrorCode = "forbidden"
	CodeUnauthorized ErrorCode = "unauthorized"
	CodeConflict     ErrorCode = "conflict"
	CodeInternal     ErrorCode = "internal_error"
)

type AppError struct {
	Code    ErrorCode
	Message string
}

func (e *AppError) Error() string { return e.Message }

func NewValidationError(message string) *AppError {
	return &AppError{Code: CodeValidation, Message: message}
}

func NewNotFoundError(message string) *AppError {
	return &AppError{Code: CodeNotFound, Message: message}
}

func NewForbiddenError(message string) *AppError {
	return &AppError{Code: CodeForbidden, Message: message}
}

func NewUnauthorizedError(message string) *AppError {
	return &AppError{Code: CodeUnauthorized, Message: message}
}

func NewConflictError(message string) *AppError {
	return &AppError{Code: CodeConflict, Message: message}
}

func NewInternalError(message string) *AppError {
	return &AppError{Code: CodeInternal, Message: message}
}
