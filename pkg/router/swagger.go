package router

import "net/http"

// setupSwaggerRoutes registers OpenAPI JSON, Swagger UI, and root redirect (development only).
func (r *Router) setupSwaggerRoutes() {
	r.mux.HandleFunc("/swagger/guestbook.swagger.json", func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		http.ServeFile(w, req, "gen/guestbook/v1/guestbook.swagger.json")
	})

	r.mux.HandleFunc("/swagger-ui.html", func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		http.ServeFile(w, req, "swagger-ui.html")
	})

	r.mux.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		if req.URL.Path == "/" {
			http.Redirect(w, req, "/swagger-ui.html", http.StatusFound)
			return
		}
		http.NotFound(w, req)
	})
}
