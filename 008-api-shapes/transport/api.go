package transport

import (
	"api-shapes/pkg/router"
	"api-shapes/store"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/google/uuid"
)

func List(w http.ResponseWriter, r *http.Request) {
	var res []UserRes
	l := store.List()
	for _, u := range l {
		var ur UserRes
		ur.Bind(u)
		res = append(res, ur)
	}

	router.JsonResponse(w, res, http.StatusOK)
}

func Retrieve(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/users/")
	u := store.Retrieve(id)

	var res UserRes
	res.Bind(u)

	router.JsonResponse(w, res, http.StatusOK)
}

func Create(w http.ResponseWriter, r *http.Request) {
	var req UserReq
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	u := store.Create(store.User{ID: uuid.NewString(), Name: req.Name})

	var res UserRes
	res.Bind(u)

	router.JsonResponse(w, res, http.StatusCreated)
}

func Update(w http.ResponseWriter, r *http.Request) {
	var req UserReq
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	id := strings.TrimPrefix(r.URL.Path, "/users/")

	u := store.Retrieve(id)
	u.Name = req.Name
	u = store.Update(u)

	var res UserRes
	res.Bind(u)

	router.JsonResponse(w, res, http.StatusOK)
}

func Delete(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/users/")
	store.Delete(id)

	w.WriteHeader(http.StatusNoContent)
}
