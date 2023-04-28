package gitlab

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v16/internal/tool/errz"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v16/internal/tool/httpz"
	"google.golang.org/protobuf/encoding/protojson"
)

type ResponseHandlerStruct struct {
	AcceptHeader string
	HandleFunc   func(*http.Response, error) error
}

func (r ResponseHandlerStruct) Handle(resp *http.Response, err error) error {
	return r.HandleFunc(resp, err)
}

func (r ResponseHandlerStruct) Accept() string {
	return r.AcceptHeader
}

func NakedResponseHandler(response **http.Response) ResponseHandler {
	return ResponseHandlerStruct{
		HandleFunc: func(r *http.Response, err error) error {
			if err != nil {
				return err
			}
			*response = r
			return nil
		},
	}
}

func JsonResponseHandler(response interface{}) ResponseHandler {
	return ResponseHandlerStruct{
		AcceptHeader: "application/json",
		HandleFunc: handleOkResponse(func(body []byte) error {
			if err := json.Unmarshal(body, response); err != nil {
				return fmt.Errorf("json.Unmarshal: %w", err)
			}
			return nil
		}),
	}
}

func ProtoJsonResponseHandler(response ValidatableMessage) ResponseHandler {
	return ResponseHandlerStruct{
		AcceptHeader: "application/json",
		HandleFunc: handleOkResponse(func(body []byte) error {
			err := protojson.UnmarshalOptions{
				DiscardUnknown: true,
			}.Unmarshal(body, response)
			if err != nil {
				return fmt.Errorf("protojson.Unmarshal: %w", err)
			}
			if err = response.ValidateAll(); err != nil {
				return fmt.Errorf("ValidateAll: %w", err)
			}
			return nil
		}),
	}
}

func handleOkResponse(h func(body []byte) error) func(*http.Response, error) error {
	return func(resp *http.Response, err error) (retErr error) {
		if err != nil {
			return err
		}
		defer errz.SafeClose(resp.Body, &retErr)
		switch resp.StatusCode {
		case http.StatusOK, http.StatusCreated:
			contentType := resp.Header.Get(httpz.ContentTypeHeader)
			if !httpz.IsContentType(contentType, "application/json") {
				return fmt.Errorf("unexpected %s in response: %q", httpz.ContentTypeHeader, contentType)
			}
			body, err := io.ReadAll(resp.Body)
			if err != nil {
				return fmt.Errorf("response body read: %w", err)
			}
			return h(body)
		default: // Unexpected status
			path := ""
			if resp.Request != nil && resp.Request.URL != nil {
				path = resp.Request.URL.Path
			}
			return &ClientError{
				StatusCode: int32(resp.StatusCode),
				Path:       path,
			}
		}
	}
}

// NoContentResponseHandler can be used when no response is expected or response must be discarded.
func NoContentResponseHandler() ResponseHandler {
	return ResponseHandlerStruct{
		HandleFunc: func(resp *http.Response, err error) (retErr error) {
			if err != nil {
				return err
			}
			defer errz.SafeClose(resp.Body, &retErr)
			switch resp.StatusCode {
			case http.StatusOK, http.StatusNoContent:
				const maxBodySlurpSize = 8 * 1024
				_, err = io.CopyN(io.Discard, resp.Body, maxBodySlurpSize)
				if err == io.EOF { // nolint:errorlint
					err = nil
				}
				return err
			default: // Unexpected status
				path := ""
				if resp.Request != nil && resp.Request.URL != nil {
					path = resp.Request.URL.Path
				}
				return &ClientError{
					StatusCode: int32(resp.StatusCode),
					Path:       path,
				}
			}
		},
	}
}
