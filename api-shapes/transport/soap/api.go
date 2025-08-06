package soap

import (
	"api-shapes/pkg/router"
	"api-shapes/store"
	"api-shapes/transport"
	"encoding/xml"
	"io"
	"net/http"
	"strings"

	"github.com/google/uuid"
)

type userAPI struct{}

func NewUserAPI() transport.UserAPI {
	return &userAPI{}
}

func (u *userAPI) List(w http.ResponseWriter, r *http.Request) {
	var res transport.ListRes
	l := store.List()
	for _, u := range l {
		var ur transport.UserRes
		ur.Bind(u)
		res.Users = append(res.Users, ur)
	}

	router.XMLResponse(w, res, http.StatusOK)
}

func (u *userAPI) Retrieve(w http.ResponseWriter, r *http.Request) {
	id := strings.Split(r.URL.Path, "/")[3]
	e := store.Retrieve(id)

	var res transport.UserRes
	res.Bind(e)

	router.XMLResponse(w, res, http.StatusOK)
}

func (u *userAPI) Create(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var req transport.UserReq
	err = xml.Unmarshal(body, &req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	e := store.Create(store.User{ID: uuid.NewString(), Name: req.Name})

	var res transport.UserRes
	res.Bind(e)

	router.XMLResponse(w, res, http.StatusCreated)
}

func (u *userAPI) Update(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var req transport.UserReq
	err = xml.Unmarshal(body, &req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	id := strings.Split(r.URL.Path, "/")[3]

	e := store.Retrieve(id)
	e.Name = req.Name
	e = store.Update(e)

	var res transport.UserRes
	res.Bind(e)

	router.XMLResponse(w, res, http.StatusOK)
}

func (u *userAPI) Delete(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/users/")
	store.Delete(id)

	w.WriteHeader(http.StatusNoContent)
}
