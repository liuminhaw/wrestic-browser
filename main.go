package main

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/joho/godotenv"
	"github.com/liuminhaw/wrestic-brw/controllers"
	"github.com/liuminhaw/wrestic-brw/static"
	"github.com/liuminhaw/wrestic-brw/templates"
	"github.com/liuminhaw/wrestic-brw/views"
)

type config struct {
	Server struct {
		Address string // default localhost:3000
	}
}

// loadEnvConfig loads config setting from .env file
func loadEnvConfig() (config, error) {
	var cfg config
	err := godotenv.Load()
	if err != nil {
		return cfg, err
	}

	// TODO: Read the server value from an ENV variable
	cfg.Server.Address = ":4000"

	return cfg, nil
}

func main() {
	cfg, err := loadEnvConfig()
	if err != nil {
		panic(err)
	}

	r := chi.NewRouter()
	fileServer := http.FileServer(http.FS(static.FS))
	r.Get("/static/*", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.StripPrefix("/static", fileServer).ServeHTTP(w, r)
	}))

	tpl := views.Must(views.ParseFS(templates.FS, "tailwind.gohtml", "home.gohtml"))
	r.Get("/", controllers.StaticHandler(tpl))

	// Start server
	fmt.Printf("Starting the server on %s...", cfg.Server.Address)
	err = http.ListenAndServe(cfg.Server.Address, r)
	if err != nil {
		panic(err)
	}

}

// func TimerMiddleware(h http.HandlerFunc) http.HandlerFunc {
// 	return func(w http.ResponseWriter, r *http.Request) {
// 		start := time.Now()
// 		h(w, r)
// 		fmt.Println("Request time:", time.Since(start))
// 	}
// }
