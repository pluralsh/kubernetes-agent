package structerr

import (
	"errors"
	"fmt"
	"sort"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"
)

// MetadataItem is an item that associated a metadata key with an arbitrary value.
type MetadataItem struct {
	// Key is the key of the metadata item that will be used as the logging key.
	Key string
	// Value is the value of the metadata item that will be formatted as the logging value.
	Value any
}

// Error is a structured error that contains additional details.
type Error struct {
	err     error
	code    codes.Code
	details []proto.Message
	// metadata is the array of metadata items added to this error. Note that we explicitly
	// don't use a map here so that we don't have any allocation overhead here in the general
	// case where there is no metadata.
	metadata []MetadataItem
}

type grpcStatuser interface {
	GRPCStatus() *status.Status
}

func newError(code codes.Code, format string, a ...any) Error {
	for i, arg := range a {
		err, ok := arg.(error)
		if !ok {
			continue
		}

		if errors.As(err, &(Error{})) {
			// We need to explicitly handle this, otherwise `status.FromError()` would
			// return these because we implement `GRPCStatus()`.
			continue
		}

		// If we see any wrapped gRPC error, then we retain its error code and details.
		// Note that we cannot use `status.FromError()` here, as that would only return an
		// error in case the immediate error is a gRPC status error.
		var wrappedGRPCStatus grpcStatuser
		if errors.As(err, &wrappedGRPCStatus) {
			grpcStatus := wrappedGRPCStatus.GRPCStatus()

			// The error message from gRPC errors is awkward because they include
			// RPC-specific constructs. This is awkward especially in the case where
			// these are embedded in the middle of an error message.
			//
			// So if we see that the top-level error is a gRPC error, then we only use
			// the status message as error message. But otherwise, we use the top-level
			// error message.
			message := err.Error()
			if st, ok := status.FromError(err); ok {
				message = st.Message()
			}

			var details []proto.Message
			for _, detail := range grpcStatus.Details() {
				if detailProto, ok := detail.(proto.Message); ok {
					details = append(details, detailProto)
				}
			}

			a[i] = Error{
				err:     errors.New(message),
				code:    grpcStatus.Code(),
				details: details,
			}
		}
	}

	formattedErr := fmt.Errorf(format, a...)

	// When we wrap an Error, we retain its error code. The intent of this is to retain the most
	// specific error code we have in the general case. As `Unknown` does not really count as a
	// specific error code, we will ignore these errors.
	//
	// Note that this impacts our middleware status handler, where we wrap non-context-errors
	// via `structerr.NewInternal()`. The result is that the caller should never see any
	// `Unknown` errors.
	var wrappedErr Error
	if errors.As(formattedErr, &wrappedErr) {
		if wrappedErr.code != codes.Unknown {
			code = wrappedErr.code
		}
	}

	return Error{
		err:  formattedErr,
		code: code,
	}
}

// New returns a new Error with an Unknown error code. This constructor should be used in the
// general case where it is not clear what the specific error category is. As Unknown errors get
// treated specially, they will be overridden when wrapped with an error that has a more specific
// error code.
func New(format string, a ...any) Error {
	return newError(codes.Unknown, format, a...)
}

// NewAborted constructs a new error code with the Aborted error code. Please refer to New for
// further details.
func NewAborted(format string, a ...any) Error {
	return newError(codes.Aborted, format, a...)
}

// NewAlreadyExists constructs a new error code with the AlreadyExists error code. Please refer to
// New for further details.
func NewAlreadyExists(format string, a ...any) Error {
	return newError(codes.AlreadyExists, format, a...)
}

// NewCanceled constructs a new error code with the Canceled error code. Please refer to New for
// further details.
func NewCanceled(format string, a ...any) Error {
	return newError(codes.Canceled, format, a...)
}

// NewDataLoss constructs a new error code with the DataLoss error code. Please refer to New for
// further details.
func NewDataLoss(format string, a ...any) Error {
	return newError(codes.DataLoss, format, a...)
}

// NewDeadlineExceeded constructs a new error code with the DeadlineExceeded error code. Please
// refer to New for further details.
func NewDeadlineExceeded(format string, a ...any) Error {
	return newError(codes.DeadlineExceeded, format, a...)
}

// NewFailedPrecondition constructs a new error code with the FailedPrecondition error code. Please
// refer to New for further details.
func NewFailedPrecondition(format string, a ...any) Error {
	return newError(codes.FailedPrecondition, format, a...)
}

// NewInternal constructs a new error code with the Internal error code. Please refer to New for
// further details.
func NewInternal(format string, a ...any) Error {
	return newError(codes.Internal, format, a...)
}

// NewInvalidArgument constructs a new error code with the InvalidArgument error code. Please refer
// to New for further details.
func NewInvalidArgument(format string, a ...any) Error {
	return newError(codes.InvalidArgument, format, a...)
}

// NewNotFound constructs a new error code with the NotFound error code. Please refer to New for
// further details.
func NewNotFound(format string, a ...any) Error {
	return newError(codes.NotFound, format, a...)
}

// NewPermissionDenied constructs a new error code with the PermissionDenied error code. Please
// refer to New for further details.
func NewPermissionDenied(format string, a ...any) Error {
	return newError(codes.PermissionDenied, format, a...)
}

// NewResourceExhausted constructs a new error code with the ResourceExhausted error code. Please
// refer to New for further details.
func NewResourceExhausted(format string, a ...any) Error {
	return newError(codes.ResourceExhausted, format, a...)
}

// NewUnavailable constructs a new error code with the Unavailable error code. Please refer to New for further details.
// Please note that the Unavailable status code is a signal telling clients to retry automatically. This auto-retry
// mechanism is handled at the library layer, without client consensus. Typically, it is used for the situations where
// the gRPC is not available or is not responding. Here are some discrete examples:
//
//   - Server downtime: The server hosting the gRPC service is down for maintenance or has crashed.
//   - Network issues: Connectivity problems between the client and server, like network congestion or a broken connection,
//     can cause the service to appear unavailable.
//   - Load balancing failure: In a distributed system, the load balancer may be unable to route the client's request to a
//     healthy instance of the gRPC service. This can happen if all instances are down or if the load balancer is
//     misconfigured.
//   - TLS/SSL handshake failure: If there's a problem during the TLS/SSL handshake between the client and the server, the
//     connection may fail, leading to an UNAVAILABLE status code.
//
// Thus, this status code should be used by interceptors or network-related components. gRPC handlers should use another
// status code instead.
func NewUnavailable(format string, a ...any) Error {
	return newError(codes.Unavailable, format, a...)
}

// NewUnauthenticated constructs a new error code with the Unauthenticated error code. Please refer
// to New for further details.
func NewUnauthenticated(format string, a ...any) Error {
	return newError(codes.Unauthenticated, format, a...)
}

// NewUnimplemented constructs a new error code with the Unimplemented error code. Please refer to
// New for further details.
func NewUnimplemented(format string, a ...any) Error {
	return newError(codes.Unimplemented, format, a...)
}

// Error returns the error message of the Error.
func (e Error) Error() string {
	return e.err.Error()
}

// Unwrap returns the wrapped error if any, otherwise it returns nil.
func (e Error) Unwrap() error {
	return errors.Unwrap(e.err)
}

// Is checks whether the error is equivalent to the target error. Errors are only considered
// equivalent if the GRPC representation of this error is the same.
func (e Error) Is(targetErr error) bool {
	target, ok := targetErr.(Error) // nolint: errorlint
	if !ok {
		return false
	}

	return errors.Is(e.GRPCStatus().Err(), target.GRPCStatus().Err())
}

// Code returns the error code of the Error.
func (e Error) Code() codes.Code {
	return e.code
}

// GRPCStatus returns the gRPC status of this error.
func (e Error) GRPCStatus() *status.Status {
	st := status.New(e.Code(), e.Error())

	if details := e.Details(); len(details) > 0 {
		proto := st.Proto()

		for _, detail := range details {
			marshaled, err := anypb.New(detail)
			if err != nil {
				return status.New(codes.Internal, fmt.Sprintf("marshaling error details: %v", err))
			}

			proto.Details = append(proto.Details, marshaled)
		}

		st = status.FromProto(proto)
	}

	return st
}

// errorChain returns the complete chain of `structerr.Error`s wrapped by this error, including the
// error itself.
func (e Error) errorChain() []Error {
	var result []Error
	for err := error(e); err != nil; err = errors.Unwrap(err) {
		if structErr, ok := err.(Error); ok { // nolint: errorlint
			result = append(result, structErr)
		}
	}
	return result
}

// Metadata returns the Error's metadata. Please refer to `ExtractMetadata()` for the exact semantics of this function.
func (e Error) Metadata() map[string]any {
	return ExtractMetadata(e)
}

// MetadataItems returns a copy of all metadata items added to this error. This function has the
// same semantics as `Metadata()`. The results are sorted by their metadata key.
func (e Error) MetadataItems() []MetadataItem {
	metadata := e.Metadata()

	items := make([]MetadataItem, 0, len(metadata))
	for key, value := range metadata {
		items = append(items, MetadataItem{Key: key, Value: value})
	}

	sort.Slice(items, func(i, j int) bool {
		return items[i].Key < items[j].Key
	})

	return items
}

// WithMetadata adds an additional metadata item to the Error. The metadata has the intent to
// provide more context around errors to the consumer of the Error. Calling this function multiple
// times with the same key will override any previous values.
func (e Error) WithMetadata(key string, value any) Error {
	for i, metadataItem := range e.metadata {
		// In case the key already exists we override it.
		if metadataItem.Key == key {
			e.metadata[i].Value = value
			return e
		}
	}

	// Otherwise we append a new metadata item.
	e.metadata = append(e.metadata, MetadataItem{
		Key: key, Value: value,
	})
	return e
}

// WithMetadataItems is a convenience function to append multiple metadata items to an error. It
// behaves the same as if `WithMetadata()` was called for each of the items separately.
func (e Error) WithMetadataItems(items ...MetadataItem) Error {
	for _, item := range items {
		e = e.WithMetadata(item.Key, item.Value)
	}
	return e
}

// Details returns the chain error details set by this and any wrapped Error. The returned array
// contains error details ordered from top-level error details to bottom-level error details.
func (e Error) Details() []proto.Message {
	var details []proto.Message
	for _, err := range e.errorChain() {
		details = append(details, err.details...)
	}
	return details
}

// WithDetail sets the Error detail that provides additional structured information about the error
// via gRPC so that callers can programmatically determine the exact circumstances of an error.
func (e Error) WithDetail(detail proto.Message) Error {
	e.details = append(e.details, detail)
	return e
}

// WithGRPCCode overrides the gRPC code embedded into the error.
func (e Error) WithGRPCCode(code codes.Code) Error {
	e.code = code
	return e
}

// ErrorMetadater is an interface that can be implemented by error types in order to provide custom metadata items
// without itself being a `structerr.Error`.
type ErrorMetadater interface {
	// ErrorMetadata returns the list of metadata items attached to this error.
	ErrorMetadata() []MetadataItem
}

// ExtractMetadata extracts metadata from the given error if any of the errors in its chain contain any. Errors may
// contain in case they are either a `structerr.Error` or in case they implement the `ErrorMetadater` interface. The
// metadata will contain the combination of all added metadata of this error as well as any wrapped Errors.
//
// When the same metada key exists multiple times in the error chain, then the value that is
// highest up the callchain will be returned. This is done because in general, the higher up the
// callchain one is the more context is available.
func ExtractMetadata(err error) map[string]any {
	metadata := map[string]any{}

	for ; err != nil; err = errors.Unwrap(err) {
		var metadataItems []MetadataItem
		if structerr, ok := err.(Error); ok { // nolint: errorlint
			metadataItems = structerr.metadata
		} else if errorMetadater, ok := err.(ErrorMetadater); ok { // nolint: errorlint
			metadataItems = errorMetadater.ErrorMetadata()
		} else {
			continue
		}

		for _, m := range metadataItems {
			if _, exists := metadata[m.Key]; !exists {
				metadata[m.Key] = m.Value
			}
		}
	}

	return metadata
}
