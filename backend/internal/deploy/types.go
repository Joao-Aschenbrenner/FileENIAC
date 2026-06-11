package errors

import (
	"fmt"
)

type DeployError struct {
	Code    string
	Message string
	Cause   error
}

func (e *DeployError) Error() string {
	if e.Cause != nil {
		return fmt.Sprintf("[%s] %s: %v", e.Code, e.Message, e.Cause)
	}
	return fmt.Sprintf("[%s] %s", e.Code, e.Message)
}

func (e *DeployError) Unwrap() error {
	return e.Cause
}

func NewDeployError(code, message string, cause error) *DeployError {
	return &DeployError{
		Code:    code,
		Message: message,
		Cause:   cause,
	}
}

var (
	ErrConfigNotFound     = &DeployError{Code: "CONFIG_NOT_FOUND", Message: "Configuration file not found"}
	ErrProjectNotFound    = &DeployError{Code: "PROJECT_NOT_FOUND", Message: "Project not found in configuration"}
	ErrFTPSConnectFailed  = &DeployError{Code: "FTPS_CONNECT_FAILED", Message: "Failed to connect to FTPS server"}
	ErrFTPSAuthFailed     = &DeployError{Code: "FTPS_AUTH_FAILED", Message: "Failed to authenticate to FTPS server"}
	ErrTokenInvalid       = &DeployError{Code: "TOKEN_INVALID", Message: "Invalid deployment token"}
	ErrTokenExpired       = &DeployError{Code: "TOKEN_EXPIRED", Message: "Deployment token has expired"}
	ErrUploadFailed       = &DeployError{Code: "UPLOAD_FAILED", Message: "Failed to upload file to server"}
	ErrTriggerFailed      = &DeployError{Code: "TRIGGER_FAILED", Message: "Failed to trigger deployment on server"}
	ErrMigrationFailed    = &DeployError{Code: "MIGRATION_FAILED", Message: "Migration execution failed"}
	ErrVerifyFailed       = &DeployError{Code: "VERIFY_FAILED", Message: "Deployment verification failed"}
	ErrHistoryReadFailed  = &DeployError{Code: "HISTORY_READ_FAILED", Message: "Failed to read deployment history"}
	ErrHistoryWriteFailed = &DeployError{Code: "HISTORY_WRITE_FAILED", Message: "Failed to write to deployment history"}
	ErrRollbackNotFound   = &DeployError{Code: "ROLLBACK_NOT_FOUND", Message: "No successful deployment found for rollback"}
	ErrPackFailed         = &DeployError{Code: "PACK_FAILED", Message: "Failed to create deployment package"}
)

func IsDeployError(err error, code string) bool {
	if de, ok := err.(*DeployError); ok {
		return de.Code == code
	}
	return false
}