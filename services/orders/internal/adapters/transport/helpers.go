package transport

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"mime"
	"net/http"
	"strings"
)

const maxRequestBodySize int64 = 1 << 20 // 1 MB

type APIError struct {
	Status  int
	Message string
}

func (e *APIError) Error() string {
	return e.Message
}

type ErrorResponse struct {
	Error string `json:"error"`
}

func readJSON(w http.ResponseWriter, r *http.Request, dst any) *APIError {
	contentType := r.Header.Get("Content-Type")

	if contentType == "" {
		return &APIError{
			Status:  http.StatusUnsupportedMediaType,
			Message: "Content-Type must be application/json",
		}
	}

	mediaType, _, err := mime.ParseMediaType(contentType)
	if err != nil {
		return &APIError{
			Status:  http.StatusUnsupportedMediaType,
			Message: "invalid Content-Type header",
		}
	}

	if mediaType != "application/json" {
		return &APIError{
			Status:  http.StatusUnsupportedMediaType,
			Message: "Content-Type must be application/json",
		}
	}

	// limit body
	r.Body = http.MaxBytesReader(w, r.Body, maxRequestBodySize)

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	err = decoder.Decode(dst)
	if err != nil {
		var syntaxError *json.SyntaxError
		var typeError *json.UnmarshalTypeError
		var maxBytesError *http.MaxBytesError

		switch {
		case errors.As(err, &maxBytesError):
			return &APIError{
				Status:  http.StatusRequestEntityTooLarge,
				Message: fmt.Sprintf("request body must not be larger than %d bytes", maxRequestBodySize),
			}

		case errors.Is(err, io.EOF):
			return &APIError{
				Status:  http.StatusBadRequest,
				Message: "request body must not be empty",
			}

		case errors.As(err, &syntaxError):
			return &APIError{
				Status: http.StatusBadRequest,
				Message: fmt.Sprintf(
					"request body contains malformed JSON near character %d",
					syntaxError.Offset,
				),
			}

		case errors.Is(err, io.ErrUnexpectedEOF):
			return &APIError{
				Status:  http.StatusBadRequest,
				Message: "request body contains malformed JSON",
			}

		case errors.As(err, &typeError):
			if typeError.Field != "" {
				return &APIError{
					Status: http.StatusBadRequest,
					Message: fmt.Sprintf(
						"incorrect type for field %q",
						typeError.Field,
					),
				}
			}

			return &APIError{
				Status:  http.StatusBadRequest,
				Message: "request body contains an incorrect JSON type",
			}

		case strings.HasPrefix(err.Error(), "json: unknown field "):
			field := strings.TrimPrefix(err.Error(), "json: unknown field ")

			return &APIError{
				Status: http.StatusBadRequest,
				Message: fmt.Sprintf(
					"request body contains unknown field %s",
					field,
				),
			}

		default:
			return &APIError{
				Status:  http.StatusBadRequest,
				Message: "request body contains invalid JSON",
			}
		}
	}

	// after json next only eof
	//
	// wrong:
	// {"product_id": 1, "quantity": 2}
	// {"product_id": 2, "quantity": 3}
	//
	// {"product_id": 1, "quantity": 2}
	// aboba
	if err := decoder.Decode(&struct{}{}); err != io.EOF {
		return &APIError{
			Status:  http.StatusBadRequest,
			Message: "request body must contain only one JSON value",
		}
	}

	return nil
}

func writeJSON(w http.ResponseWriter, status int, data any) error {
	response, err := json.Marshal(data)
	if err != nil {
		return err
	}

	response = append(response, '\n')

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	_, err = w.Write(response)
	return err
}

func writeError(w http.ResponseWriter, status int, message string) {
	err := writeJSON(w, status, ErrorResponse{
		Error: message,
	})

	if err != nil {
		slog.Error(
			"failed to write error response",
			"error", err,
		)
	}
}
