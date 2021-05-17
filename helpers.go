package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"go.uber.org/zap"
)

// DecodeJSONBody decodes the JSON from Request.Body and checks for common errors
// ToDo: Extend with all common errors and maybe create a middleware?
func DecodeJSONBody(w http.ResponseWriter, r *http.Request, dst interface{}) interface{} {
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	err := decoder.Decode(&dst)
	if err != nil {
		var syntaxError *json.SyntaxError
		var unmarshalTypeError *json.UnmarshalTypeError

		switch {
		case errors.As(err, &syntaxError):
			http.Error(w, err.Error(), http.StatusBadRequest)

		case errors.As(err, &unmarshalTypeError):
			http.Error(w, err.Error(), http.StatusBadRequest)

		case strings.HasPrefix(err.Error(), "json: unknown field"):
			fieldName := strings.TrimPrefix(err.Error(), "json: unknown field")
			msg := fmt.Sprintf("Request body contains unknown field %s", fieldName)
			http.Error(w, msg, http.StatusBadRequest)

		default:
			http.Error(w, err.Error(), http.StatusBadRequest)
		}

	}

	return dst
}

// ifError very.IsBadPractice()! Don't copy this in actual app
func ifError(err error) {
	if err != nil {
		zap.L().Error(err.Error())
	}
}
