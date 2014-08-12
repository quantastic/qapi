package main

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gorilla/mux"
	"github.com/quantastic/qapi"
)

func Router() http.Handler {
	mux := mux.NewRouter()
	mux.Methods("GET").Path("/").HandlerFunc(Home)
	mux.Methods("POST").Path("/times").HandlerFunc(CreateTime)
	return mux
}

func Home(w http.ResponseWriter, r *http.Request) {
	jsonWrite(w, http.StatusOK, map[string]interface{}{})
}

func CreateTime(w http.ResponseWriter, r *http.Request) {
	t := qapi.NewTime()
	if err := mapTime(r, &t); err != nil {
		badRequest(w, err)
		return
	}
	if err := t.Valid(); err != nil {
		badRequest(w, err)
		return
	}
	if err := db.SaveTime(&t); err != nil {
		jsonWrite(w, http.StatusInternalServerError, err)
		return
	}
	res := map[string]interface{}{"time": NewTime(t)}
	jsonWrite(w, http.StatusCreated, res)
}

func mapTime(r *http.Request, t *qapi.Time) error {
	m := NewMap()
	m.Required("category", &t.Category)
	m.Required("start", &t.Start)
	m.Optional("end", &t.End)
	m.Optional("note", &t.Note)
	return JSONMap(r, m)
}

func NewTime(t qapi.Time) Time {
	return Time{
		Id:       t.Id,
		URL:      fmt.Sprintf("%s/times/%s", config.Url, t.Id),
		Category: NewTimeCategory(t.Category),
		End:      t.End,
		Start:    t.Start,
		Note:     t.Note,
		Created:  t.Created,
		Updated:  t.Updated,
	}
}

type Time struct {
	Id       string       `json:"id"`
	URL      string       `json:"url"`
	Category TimeCategory `json:"category"`
	End      *time.Time   `json:"end"`
	Start    time.Time    `json:"start"`
	Note     string       `json:"note"`
	Created  time.Time    `json:"created"`
	Updated  time.Time    `json:"updated"`
}

func NewTimeCategory(category []string) TimeCategory {
	return TimeCategory{
		Name: category,
		URL:  fmt.Sprintf("%s/times/c/%s", config.Url, strings.Join(category, "/")),
	}
}

type TimeCategory struct {
	Name []string `json:"name"`
	URL  string   `json:"url"`
}
