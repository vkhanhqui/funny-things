package router

import (
	"encoding/json"
	"encoding/xml"
	"log"
	"net/http"
)

func JsonResponse(w http.ResponseWriter, v any, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	err := json.NewEncoder(w).Encode(v)
	if err != nil {
		panic(err)
	}
}

func XMLResponse(w http.ResponseWriter, v any, status int) {
	w.Header().Set("Content-Type", "application/xml")
	w.WriteHeader(status)
	res, err := xml.MarshalIndent(v, "", "  ")
	if err != nil {
		panic(err)
	}
	w.Write(res)
}

func ErrorHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				log.Println(err)
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			}
		}()
		next.ServeHTTP(w, r)
	})
}
