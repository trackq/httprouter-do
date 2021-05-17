package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/gobuffalo/pop/v5"
	"github.com/gofrs/uuid"
	"github.com/julienschmidt/httprouter"
	"yesplease.ai/httprouter-do/models"
)

func main() {
	router := httprouter.New()
	router.GET("/healthcheck/ping", Ping)
	router.GET("/", Index)
	router.POST("/", Create)
	router.PATCH("/:id", Update)
	router.DELETE("/:id", Delete)

	log.Println("Ready on port: 8080")
	log.Fatal(http.ListenAndServe(":8080", router))
}

func Index(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	db, err := pop.Connect("development")
	ifError(err)

	todoes := []models.Todo{}
	err = db.Order("created_at desc").All(&todoes)
	ifError(err)

	json.NewEncoder(w).Encode(todoes)
}

func Create(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	db, err := pop.Connect("development")
	ifError(err)

	todo := &models.Todo{}

	DecodeJSONBody(w, r, todo)

	vErr, err := db.ValidateAndCreate(todo)
	ifError(err)
	ifError(vErr)

	json.NewEncoder(w).Encode(todo)
}

func Update(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	db, err := pop.Connect("development")
	ifError(err)

	id := ps.ByName("id")
	uId, err := uuid.FromString(id)
	ifError(err)

	todo := &models.Todo{}
	DecodeJSONBody(w, r, todo)

	todo.ID = uId

	vErr, err := db.ValidateAndUpdate(todo)
	ifError(err)
	ifError(vErr)

	json.NewEncoder(w).Encode(todo)
}

func Delete(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	db, err := pop.Connect("development")
	ifError(err)

	id := ps.ByName("id")
	uId, err := uuid.FromString(id)
	ifError(err)

	todo := &models.Todo{ID: uId}
	destroyed := db.Destroy(todo)

	if destroyed != nil {
		fmt.Println(err)
		return
	}

	w.WriteHeader(http.StatusAccepted)
}

func Ping(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	fmt.Fprintf(w, "Pong")
}

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
		fmt.Println(err.Error())
	}
}
