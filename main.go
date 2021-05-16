package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/julienschmidt/httprouter"
	"yesplease.ai/httprouter-do/models"
)

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

func Ping(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	fmt.Fprintf(w, "Pong")
}

func Index(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	fmt.Fprintf(w, "httprouter-do")
}

func Create(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var t models.Todo
	DecodeJSONBody(w, r, &t)

	fmt.Fprintf(w, "%+v", t)
}

func Update(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
}

func Delete(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

}

func main() {
	router := httprouter.New()
	router.GET("/healthcheck/ping", Ping)
	router.GET("/", Index)
	router.POST("/", Create)
	router.PATCH("/", Update)
	router.DELETE("/:id", Delete)

	log.Println("Ready on port: 8080")
	log.Fatal(http.ListenAndServe(":8080", router))
}
