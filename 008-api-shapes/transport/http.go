package transport

import "net/http"

func NewRouter(pattern string, mux *http.ServeMux, userAPI UserAPI) {
	mux.HandleFunc(pattern, func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			userAPI.List(w, r)
		case http.MethodPost:
			userAPI.Create(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	mux.HandleFunc(pattern+"/", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			userAPI.Retrieve(w, r)
		case http.MethodPut:
			userAPI.Update(w, r)
		case http.MethodDelete:
			userAPI.Delete(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})
}
