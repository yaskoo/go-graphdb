package graphdb

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
)

// CombinedResponseHandler creates a single response handler from multiple.
// It iterates all handlers, but stops on the first one that returns an error, returning that error to the caller.
// Useful for combining error checkers and response parsers.
func CombinedResponseHandler(handlers ...ResponseHandler) ResponseHandler {
	return func(resp *http.Response) error {
		for _, handler := range handlers {
			if err := handler(resp); err != nil {
				return err
			}
		}
		return nil
	}
}

func ExpectStatusCode(code int) ResponseHandler {
	return func(resp *http.Response) error {
		if resp.StatusCode == code {
			return nil
		}

		return fmt.Errorf("expected status code %d but got %d", code, resp.StatusCode)
	}
}

func ExpectOneOfStatusCode(codes ...int) ResponseHandler {
	return func(resp *http.Response) error {
		for _, code := range codes {
			if resp.StatusCode == code {
				return nil
			}
		}
		return fmt.Errorf("expected one of %v status codes but got %d", codes, resp.StatusCode)
	}
}

func UnmarshalJson(v any) ResponseHandler {
	return func(resp *http.Response) error {
		dec := json.NewDecoder(resp.Body)
		if err := dec.Decode(&v); err != nil {
			return err
		}
		return nil
	}
}

func ErrNotStatus(status int, message string, resp *http.Response) error {
	if resp.StatusCode == status {
		return nil
	}

	all, _ := io.ReadAll(resp.Body)
	if len(all) > 0 && message != "" {
		return fmt.Errorf("%s: %s", message, string(all))
	}

	if message != "" {
		return errors.New(message)
	}
	return fmt.Errorf("unexpected status code: %d", status)
}
