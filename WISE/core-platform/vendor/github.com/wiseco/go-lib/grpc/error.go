package grpc

import (
	"net/http"

	grpcRoot "github.com/wiseco/protobuf/golang"
	grpcCodes "google.golang.org/grpc/codes"
	grpcStatus "google.golang.org/grpc/status"
)

// StatusWithErrorCode returns a GRPC compatible error that can be passed back from services to callers
func StatusWithErrorCode(code grpcRoot.ErrorCode, grpcCode grpcCodes.Code, message string) error {
	status, _ := grpcStatus.New(grpcCode, message).WithDetails(&grpcRoot.Error{Code: code})
	return status.Err()
}

// NewHttpError is a http standard error for api responses
func NewHttpError(requestID string, httpMethod string, code grpcCodes.Code, errorCode grpcRoot.ErrorCode, message string) *grpcRoot.HttpError {
	httpStatus, ok := httpStatusFromStatusCode[code]
	if !ok {
		httpStatus = 500
	}

	return &grpcRoot.HttpError{
		RequestId: requestID,
		Method:    httpMethod,
		Status:    int32(httpStatus),
		Code:      errorCode,
		Message:   message,
	}
}

// httpStatusFromStatusCode maps grpc error codes to http status
var httpStatusFromStatusCode = map[grpcCodes.Code]int{
	grpcCodes.OK:                 http.StatusOK,
	grpcCodes.Canceled:           http.StatusInternalServerError,
	grpcCodes.Unknown:            http.StatusInternalServerError,
	grpcCodes.InvalidArgument:    http.StatusBadRequest,
	grpcCodes.DeadlineExceeded:   http.StatusInternalServerError,
	grpcCodes.NotFound:           http.StatusNotFound,
	grpcCodes.AlreadyExists:      http.StatusBadRequest,
	grpcCodes.PermissionDenied:   http.StatusForbidden,
	grpcCodes.ResourceExhausted:  http.StatusInternalServerError,
	grpcCodes.FailedPrecondition: http.StatusPreconditionFailed,
	grpcCodes.Aborted:            http.StatusInternalServerError,
	grpcCodes.OutOfRange:         http.StatusRequestedRangeNotSatisfiable,
	grpcCodes.Unimplemented:      http.StatusNotImplemented,
	grpcCodes.Internal:           http.StatusInternalServerError,
	grpcCodes.Unavailable:        http.StatusServiceUnavailable,
	grpcCodes.DataLoss:           http.StatusInternalServerError,
	grpcCodes.Unauthenticated:    http.StatusUnauthorized,
}

// NewHttpErrorFromError creates a http error from a generic error - optimized for grpc status errors
func NewHttpErrorFromError(requestID string, httpMethod string, err error) *grpcRoot.HttpError {
	httpStatus, ok := httpStatusFromStatusCode[grpcStatus.Code(err)]
	if !ok {
		httpStatus = 500
	}

	code := grpcRoot.ErrorCode_EC_UNSPECIFIED
	message := http.StatusText(httpStatus)
	status := grpcStatus.Convert(err)
	if status != nil {
		message = status.Message()

		details := status.Details()
		if len(details) > 0 {
			rootError, ok := details[0].(*grpcRoot.Error)
			if ok {
				code = rootError.Code
			}
		}
	}

	return &grpcRoot.HttpError{
		RequestId: requestID,
		Method:    httpMethod,
		Status:    int32(httpStatus),
		Code:      code,
		Message:   message,
	}
}
