package grpctool

import (
	"context"
	"errors"
	"fmt"
	"net"
	"strings"

	"github.com/ash2k/stager"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v14/internal/tool/errz"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v14/internal/tool/logz"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v14/internal/tool/prototool"
	grpccorrelation "gitlab.com/gitlab-org/labkit/correlation/grpc"
	"go.uber.org/zap"
	statuspb "google.golang.org/genproto/googleapis/rpc/status"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

func HandleSendError(log *zap.Logger, msg string, err error) error {
	// The problem is almost certainly with the client's connection.
	// Still log it on Debug.
	log.Debug(msg, logz.Error(err))
	if IsStatusError(err) {
		s := status.Convert(err).Proto()
		if isErrIllegalHeaderWriteError(s) {
			s.Code = int32(codes.Canceled)
		}
		s.Message = fmt.Sprintf("%s: %s", msg, s.Message)
		err = status.ErrorProto(s)
	} else {
		err = status.Errorf(codes.Canceled, "%s: %v", msg, err)
	}
	return err
}

// isErrIllegalHeaderWriteError checks if an error is google.golang.org/grpc/internal/transport.ErrIllegalHeaderWrite
// We cannot/shouldn't import an internal package, so we just do an error string check :facepalm:
// See https://github.com/grpc/grpc-go/blob/v1.40.0/internal/transport/http2_server.go#L999
// See https://github.com/grpc/grpc-go/issues/3575
// See https://github.com/grpc/grpc-go/issues/4696
func isErrIllegalHeaderWriteError(s *statuspb.Status) bool {
	return s.Code == int32(codes.Internal) && strings.HasSuffix(s.Message, "transport: the stream is done or WriteHeader was already called")
}

func RequestCanceledOrTimedOut(err error) bool {
	return RequestCanceled(err) || RequestTimedOut(err)
}

func RequestCanceled(err error) bool {
	for err != nil {
		if err == context.Canceled { // nolint:errorlint
			return true
		}
		code := status.Code(err)
		if code == codes.Canceled {
			return true
		}
		err = errors.Unwrap(err)
	}
	return false
}

func RequestTimedOut(err error) bool {
	for err != nil {
		if err == context.DeadlineExceeded { // nolint:errorlint
			return true
		}
		code := status.Code(err)
		if code == codes.DeadlineExceeded {
			return true
		}
		err = errors.Unwrap(err)
	}
	return false
}

func StartServer(stage stager.Stage, server *grpc.Server, listener func() (net.Listener, error)) {
	stage.Go(func(ctx context.Context) error {
		// gRPC listener
		lis, err := listener()
		if err != nil {
			return err
		}
		return server.Serve(lis)
	})
	stage.Go(func(ctx context.Context) error {
		<-ctx.Done() // can be cancelled because Serve() failed or main ctx was canceled or some stage failed
		server.GracefulStop()
		return nil
	})
}

func IsStatusError(err error) bool {
	_, ok := err.(interface { // nolint:errorlint
		GRPCStatus() *status.Status
	})
	return ok
}

func MetaToValuesMap(meta metadata.MD) map[string]*prototool.Values {
	if len(meta) == 0 {
		return nil
	}
	result := make(map[string]*prototool.Values, len(meta))
	for k, v := range meta {
		val := make([]string, len(v))
		copy(val, v) // metadata may be mutated, so copy
		result[k] = &prototool.Values{
			Value: val,
		}
	}
	return result
}

func ValuesMapToMeta(vals map[string]*prototool.Values) metadata.MD {
	result := make(metadata.MD, len(vals))
	for k, v := range vals {
		val := make([]string, len(v.Value))
		copy(val, v.Value) // metadata may be mutated, so copy
		result[k] = val
	}
	return result
}

func MaybeWrapWithCorrelationId(err error, client grpc.ClientStream) error {
	md, headerErr := client.Header()
	if headerErr != nil {
		return err
	}
	return errz.MaybeWrapWithCorrelationId(err, grpccorrelation.CorrelationIDFromMetadata(md))
}

func DeferMaybeWrapWithCorrelationId(err *error, client grpc.ClientStream) {
	if *err == nil {
		return
	}
	*err = MaybeWrapWithCorrelationId(*err, client)
}

func SplitGrpcMethod(fullMethodName string) (string /* service */, string /* method */) {
	if fullMethodName != "" && fullMethodName[0] == '/' {
		fullMethodName = fullMethodName[1:]
	}
	pos := strings.LastIndex(fullMethodName, "/")
	if pos == -1 {
		return "unknown", fullMethodName
	}
	service := fullMethodName[:pos]
	method := fullMethodName[pos+1:]
	return service, method
}
