package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gobuffalo/pop/v5"
	"github.com/gofrs/uuid"
	"github.com/julienschmidt/httprouter"
	"go.uber.org/zap"
	"yesplease.ai/httprouter-do/models"
)

func main() {
	router := httprouter.New()
	router.GET("/healthcheck/ping", Ping)
	router.GET("/", Index)
	router.POST("/", Create)
	router.PATCH("/:id", Update)
	router.DELETE("/:id", Delete)

	lm := NewZapLogger(router)

	lm.logger.Info("Server started.",
		// Structured context as strongly typed fields.
		zap.Int("port", 8080),
		zap.String("host", "localhost"),
		zap.String("Status", "Ready"),
	)

	log.Fatal(http.ListenAndServe(":8080", lm))
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
