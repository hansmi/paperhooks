package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-resty/resty/v2"
)

type RequestError struct {
	StatusCode int
	Message    string
}

func (e *RequestError) Error() string {
	return fmt.Sprintf("HTTP status %d (%s): %s", e.StatusCode, http.StatusText(e.StatusCode), e.Message)
}

func (e *RequestError) Is(other error) bool {
	err, ok := other.(*RequestError)

	return ok && e.StatusCode == err.StatusCode && e.Message == err.Message
}

type requestError struct {
	json.RawMessage
}

func convertError(requestErr error, resp *resty.Response) error {
	if requestErr != nil {
		return requestErr
	}

	if resp.IsSuccess() {
		switch resp.StatusCode() {
		case http.StatusOK, http.StatusNoContent:
			return nil
		}
	}

	err := &RequestError{
		StatusCode: resp.StatusCode(),
	}

	switch respErr := resp.Error().(type) {
	case *requestError:
		var buf bytes.Buffer

		if compactErr := json.Compact(&buf, respErr.RawMessage); compactErr != nil {
			err.Message = string(respErr.RawMessage)
		} else {
			err.Message = buf.String()
		}
	}

	if err.Message == "" {
		err.Message = resp.Status()
	}

	if err.Message == "" {
		err.Message = "unknown error"
	}

	return err
}
