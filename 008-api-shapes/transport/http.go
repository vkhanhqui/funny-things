package transport

import "net/http"

func NewRouter(pattern string, mux *http.ServeMux) {
	mux.HandleFunc(pattern, func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			List(w, r)
		case http.MethodPost:
			Create(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	mux.HandleFunc(pattern+"/", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			Retrieve(w, r)
		case http.MethodPut:
			Update(w, r)
		case http.MethodDelete:
			Delete(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})
}
